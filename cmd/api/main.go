package main

import (
	"context"
	"embed"
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/tuandoquoc/futa-ticket-hunter/internal/auth"
	"github.com/tuandoquoc/futa-ticket-hunter/internal/config"
	"github.com/tuandoquoc/futa-ticket-hunter/internal/database"
	"github.com/tuandoquoc/futa-ticket-hunter/internal/futa"
)

//go:embed static
var staticFiles embed.FS

func main() {
	cfgPath := "config.yaml"
	if len(os.Args) > 1 {
		cfgPath = os.Args[1]
	}

	cfg, err := config.Load(cfgPath)
	if err != nil {
		log.Fatalf("load config: %v", err)
	}

	b, _ := json.MarshalIndent(cfg, "", "  ")
	fmt.Println(string(b))

	db, err := database.New(cfg.Database)
	if err != nil {
		log.Fatalf("connect database: %v", err)
	}
	defer db.Close()

	futaClient := futa.NewClient(cfg.Futa)

	if cfg.Clerk.PublishableKey == "" {
		log.Println("WARNING: clerk.publishable_key is not configured — Clerk Sign-In will not work")
	}
	clerkVerifier, err := auth.NewClerkVerifier(cfg.Clerk.Issuer, cfg.Clerk.JWKSURL)
	if err != nil {
		log.Printf("WARNING: clerk auth is not configured: %v (server will start in degraded auth mode)", err)
	}

	mux := http.NewServeMux()

	// === Auth Routes (public) ===

	// Public config (exposes non-secret runtime config to the frontend)
	mux.HandleFunc("/api/config", func(w http.ResponseWriter, r *http.Request) {
		jsonOK(w, map[string]string{
			"clerk_publishable_key": cfg.Clerk.PublishableKey,
			"posthog_api_key":       cfg.Posthog.APIKey,
			"posthog_host":          cfg.Posthog.Host,
		})
	})

	// Get current authenticated user
	mux.HandleFunc("/api/auth/me", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "no-store")
		sess := auth.SessionFromContext(r.Context())
		if sess == nil {
			jsonError(w, "unauthorized", http.StatusUnauthorized)
			return
		}
		jsonOK(w, map[string]string{"email": sess.Email, "name": sess.Name, "picture": sess.Picture})
	})

	// Logout: delete session and clear cookie
	mux.HandleFunc("/api/auth/logout", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		// Clerk handles logout client-side; this endpoint exists for API parity.
		jsonOK(w, map[string]string{"message": "logged out"})
	})

	// === API Routes ===

	// Stats
	mux.HandleFunc("/api/stats", func(w http.ResponseWriter, r *http.Request) {
		emailAddr := r.URL.Query().Get("email")
		var stats *database.Stats
		var err error
		if emailAddr != "" {
			stats, err = db.GetStatsByEmail(r.Context(), emailAddr)
		} else {
			stats, err = db.GetStats(r.Context())
		}
		if err != nil {
			jsonError(w, err.Error(), 500)
			return
		}
		jsonOK(w, stats)
	})

	// Schedules
	mux.HandleFunc("/api/schedules", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			filter := r.URL.Query().Get("status")
			emailAddr := r.URL.Query().Get("email")
			var schedules []database.BookingSchedule
			var err error
			if emailAddr != "" {
				schedules, err = db.ListSchedulesByEmail(r.Context(), emailAddr, filter)
			} else {
				schedules, err = db.ListSchedules(r.Context(), filter)
			}
			if err != nil {
				jsonError(w, err.Error(), 500)
				return
			}
			if schedules == nil {
				schedules = []database.BookingSchedule{}
			}
			jsonOK(w, schedules)

		case http.MethodPost:
			var s database.BookingSchedule
			if err := json.NewDecoder(r.Body).Decode(&s); err != nil {
				jsonError(w, "invalid JSON body", 400)
				return
			}
			if s.PassengerName == "" || s.PassengerPhone == "" {
				jsonError(w, "passenger_name and passenger_phone are required", 400)
				return
			}
			if s.PassengerEmail == "" {
				jsonError(w, "passenger_email is required", 400)
				return
			}
			if s.OriginAreaID == "" || s.OriginName == "" {
				jsonError(w, "origin_area_id and origin_name are required", 400)
				return
			}
			if s.DestAreaID == "" || s.DestName == "" {
				jsonError(w, "dest_area_id and dest_name are required", 400)
				return
			}
			if s.TravelDate == "" {
				jsonError(w, "travel_date is required", 400)
				return
			}
			if s.SeatCount <= 0 {
				s.SeatCount = 1
			}
			if s.TimeFrom == "" {
				s.TimeFrom = "00:00"
			}
			if s.TimeTo == "" {
				s.TimeTo = "23:59"
			}
			if s.SeatType == "" {
				s.SeatType = "any"
			}
			if s.SeatFloor == "" {
				s.SeatFloor = "any"
			}
			switch s.SeatFloor {
			case "any", "up", "down":
			default:
				jsonError(w, "seat_floor must be one of: any, up, down", 400)
				return
			}

			if s.SeatWindow == "" {
				s.SeatWindow = "any"
			}
			switch s.SeatWindow {
			case "any", "window", "non_window":
			default:
				jsonError(w, "seat_window must be one of: any, window, non_window", 400)
				return
			}

			if s.PriorityTopRows < 0 {
				s.PriorityTopRows = 0
			}
			if s.PriorityTopRows > 10 {
				s.PriorityTopRows = 10
			}

			if s.MaxPrice < 0 {
				s.MaxPrice = 0
			}

			if err := db.CreateSchedule(r.Context(), &s); err != nil {
				jsonError(w, err.Error(), 500)
				return
			}
			s.Status = "pending"
			w.WriteHeader(http.StatusCreated)
			jsonOK(w, s)

		default:
			http.Error(w, "method not allowed", 405)
		}
	})

	mux.HandleFunc("/api/schedules/", func(w http.ResponseWriter, r *http.Request) {
		path := strings.TrimPrefix(r.URL.Path, "/api/schedules/")
		parts := strings.SplitN(path, "/", 2)
		id := parts[0]
		if id == "" {
			http.Error(w, "missing id", 400)
			return
		}

		// Handle /api/schedules/:id/cancel
		if len(parts) == 2 && parts[1] == "cancel" && r.Method == http.MethodPost {
			if err := db.CancelSchedule(r.Context(), id); err != nil {
				jsonError(w, err.Error(), 500)
				return
			}
			jsonOK(w, map[string]string{"message": "cancelled"})
			return
		}

		// Handle /api/schedules/:id/check-payment
		if len(parts) == 2 && parts[1] == "check-payment" && r.Method == http.MethodPost {
			s, err := db.GetSchedule(r.Context(), id)
			if err != nil {
				jsonError(w, "schedule not found", 404)
				return
			}
			if s.Status == "paid" {
				jsonOK(w, map[string]any{
					"status":         s.Status,
					"payment_status": "paid",
					"amount":         s.TicketPrice,
				})
				return
			}
			if s.Status != "success" || s.BookingCode == "" {
				jsonOK(w, map[string]any{"status": s.Status, "payment_status": "not_applicable"})
				return
			}

			isPaid, err := futaClient.PaymentStatus(r.Context(), s.BookingCode)
			if err != nil {
				jsonOK(w, map[string]any{"status": s.Status, "payment_status": "unknown", "error": err.Error()})
				return
			}

			if isPaid {
				if err := db.UpdateSchedulePaymentStatus(r.Context(), s.ID, "paid"); err != nil {
					log.Printf("mark schedule paid: %v", err)
				}
				jsonOK(w, map[string]any{
					"status":         "paid",
					"payment_status": "paid",
					"amount":         s.TicketPrice,
				})
				return
			} else if time.Since(s.UpdatedAt) > 5*time.Minute {
				if err := db.UpdateSchedulePaymentStatus(r.Context(), s.ID, "expired"); err != nil {
					log.Printf("mark schedule payment expired: %v", err)
				}
				jsonOK(w, map[string]any{
					"status":         "expired",
					"payment_status": "expired",
					"amount":         s.TicketPrice,
				})
				return
			}

			jsonOK(w, map[string]any{
				"status":         s.Status,
				"payment_status": "pending",
				"booking_code":   s.BookingCode,
			})
			return
		}

		switch r.Method {
		case http.MethodGet:
			s, err := db.GetSchedule(r.Context(), id)
			if err != nil {
				jsonError(w, "schedule not found", 404)
				return
			}
			jsonOK(w, s)
		case http.MethodDelete:
			if err := db.DeleteSchedule(r.Context(), id); err != nil {
				jsonError(w, err.Error(), 500)
				return
			}
			jsonOK(w, map[string]string{"message": "deleted"})
		default:
			http.Error(w, "method not allowed", 405)
		}
	})

	// Recent schedules for dashboard
	mux.HandleFunc("/api/schedules/recent", func(w http.ResponseWriter, r *http.Request) {
		emailAddr := r.URL.Query().Get("email")
		var schedules []database.BookingSchedule
		var err error
		if emailAddr != "" {
			schedules, err = db.GetRecentSchedulesByEmail(r.Context(), emailAddr, 5)
		} else {
			schedules, err = db.GetRecentSchedules(r.Context(), 5)
		}
		if err != nil {
			jsonError(w, err.Error(), 500)
			return
		}
		if schedules == nil {
			schedules = []database.BookingSchedule{}
		}
		jsonOK(w, schedules)
	})

	// Search proxy for FUTA API
	mux.HandleFunc("/api/search/pickup-points", func(w http.ResponseWriter, r *http.Request) {
		keyword := r.URL.Query().Get("keyword")
		if keyword == "" {
			jsonError(w, "keyword is required", 400)
			return
		}
		groups, areas, err := futaClient.SearchPickupPoints(r.Context(), keyword)
		if err != nil {
			jsonError(w, err.Error(), 500)
			return
		}
		jsonOK(w, map[string]any{"groups": groups, "areas": areas})
	})

	// Serve static frontend
	staticSub, err := fs.Sub(staticFiles, "static")
	if err != nil {
		log.Fatalf("static fs: %v", err)
	}
	fileServer := http.FileServer(http.FS(staticSub))
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// SPA: serve index.html for all non-API, non-asset routes
		if !strings.HasPrefix(r.URL.Path, "/api") {
			if r.URL.Path != "/" && !strings.Contains(r.URL.Path, ".") {
				r.URL.Path = "/"
			}
		}
		// Prevent browsers from caching the HTML shell so that a fresh
		// page load always runs checkAuth() with up-to-date JS/HTML.
		if r.URL.Path == "/" || r.URL.Path == "/index.html" {
			w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
			w.Header().Set("Pragma", "no-cache")
			w.Header().Set("Expires", "0")
		}
		fileServer.ServeHTTP(w, r)
	})

	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", cfg.Server.Port),
		Handler: corsMiddleware(authMiddleware(mux, clerkVerifier)),
	}

	go func() {
		log.Printf("API server starting on http://localhost:%d", cfg.Server.Port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("server error: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	server.Shutdown(ctx)
	log.Println("Server stopped")
}

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		next.ServeHTTP(w, r)
	})
}

// authMiddleware protects all /api/* routes except /api/config and /api/auth/logout.
// Valid requests have their session attached to the request context.
func authMiddleware(next http.Handler, clerkVerifier *auth.ClerkVerifier) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path
		if strings.HasPrefix(path, "/api/") &&
			path != "/api/config" &&
			path != "/api/auth/logout" {
			token, err := bearerToken(r)
			if err != nil {
				jsonError(w, "unauthorized", http.StatusUnauthorized)
				return
			}
			if clerkVerifier == nil {
				jsonError(w, "authentication not configured", http.StatusInternalServerError)
				return
			}
			sess, err := clerkVerifier.VerifySession(token)
			if err != nil {
				jsonError(w, "unauthorized", http.StatusUnauthorized)
				return
			}
			r = r.WithContext(context.WithValue(r.Context(), auth.SessionKey, sess))
		}
		next.ServeHTTP(w, r)
	})
}

func jsonOK(w http.ResponseWriter, data any) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}

func jsonError(w http.ResponseWriter, msg string, code int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(map[string]string{"error": msg})
}

func bearerToken(r *http.Request) (string, error) {
	authHeader := strings.TrimSpace(r.Header.Get("Authorization"))
	if authHeader == "" {
		return "", errors.New("missing authorization header")
	}
	parts := strings.SplitN(authHeader, " ", 2)
	if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") || strings.TrimSpace(parts[1]) == "" {
		return "", errors.New("invalid authorization header")
	}
	return strings.TrimSpace(parts[1]), nil
}
