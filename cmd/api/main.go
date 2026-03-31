package main

import (
	"context"
	"embed"
	"encoding/json"
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

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

	db, err := database.New(cfg.Database)
	if err != nil {
		log.Fatalf("connect database: %v", err)
	}
	defer db.Close()

	futaClient := futa.NewClient(cfg.Futa)

	mux := http.NewServeMux()

	// === API Routes ===

	// Settings
	mux.HandleFunc("/api/settings", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			s, err := db.GetSettings(r.Context())
			if err != nil {
				jsonError(w, err.Error(), 500)
				return
			}
			// Mask refresh token for security
			masked := *s
			if len(masked.RefreshToken) > 10 {
				masked.RefreshToken = masked.RefreshToken[:6] + "..." + masked.RefreshToken[len(masked.RefreshToken)-4:]
			}
			jsonOK(w, masked)

		case http.MethodPut:
			var s database.AppSettings
			if err := json.NewDecoder(r.Body).Decode(&s); err != nil {
				jsonError(w, "invalid JSON", 400)
				return
			}
			if err := db.UpdateSettings(r.Context(), &s); err != nil {
				jsonError(w, err.Error(), 500)
				return
			}
			jsonOK(w, map[string]string{"message": "settings updated"})

		default:
			http.Error(w, "method not allowed", 405)
		}
	})

	// Validate refresh token & extract user info
	mux.HandleFunc("/api/auth/validate-token", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "method not allowed", 405)
			return
		}
		var body struct {
			RefreshToken string `json:"refresh_token"`
		}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			jsonError(w, "invalid JSON", 400)
			return
		}
		if body.RefreshToken == "" {
			jsonError(w, "refresh_token is required", 400)
			return
		}

		tokenResp, err := futa.ExchangeRefreshToken(r.Context(), body.RefreshToken)
		if err != nil {
			jsonError(w, "Token không hợp lệ: "+err.Error(), 400)
			return
		}

		claims, err := futa.ParseIDToken(tokenResp.IDToken)
		if err != nil {
			jsonError(w, "Không thể đọc thông tin từ token: "+err.Error(), 400)
			return
		}

		// Normalize phone number
		phone := claims.PhoneNumber
		if strings.HasPrefix(phone, "+84") {
			phone = "0" + phone[3:]
		}

		jsonOK(w, map[string]any{
			"valid":     true,
			"full_name": claims.FullName,
			"phone":     phone,
			"user_id":   tokenResp.UserID,
		})
	})

	// Stats
	mux.HandleFunc("/api/stats", func(w http.ResponseWriter, r *http.Request) {
		stats, err := db.GetStats(r.Context())
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
			schedules, err := db.ListSchedules(r.Context(), filter)
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
			if s.OriginKeyword == "" || s.DestKeyword == "" {
				jsonError(w, "origin_keyword and dest_keyword are required", 400)
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
		schedules, err := db.GetRecentSchedules(r.Context(), 5)
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
		fileServer.ServeHTTP(w, r)
	})

	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", cfg.Server.Port),
		Handler: corsMiddleware(mux),
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
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
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
