package email

import (
	"fmt"
	"log"
	"net/url"

	"github.com/resend/resend-go/v3"
	"github.com/tuandoquoc/futa-ticket-hunter/internal/config"
)

type Sender struct {
	client *resend.Client
	cfg    config.EmailConfig
}

func NewSender(cfg config.EmailConfig) *Sender {
	if cfg.ResendAPIKey == "" {
		return &Sender{cfg: cfg}
	}
	return &Sender{
		client: resend.NewClient(cfg.ResendAPIKey),
		cfg:    cfg,
	}
}

type PaymentInfo struct {
	BookingID   string
	BookingCode string
	PhoneNumber string
	ToEmail     string
	ToName      string
	OriginName  string
	DestName    string
	TravelDate  string
	SeatName    string
	TicketPrice int
	RouteName   string
}

func (s *Sender) SendPaymentLink(info PaymentInfo) error {
	if s.client == nil {
		log.Printf("Email not configured, skipping payment link for %s", info.ToEmail)
		return nil
	}

	paymentURL := fmt.Sprintf("https://futabus.vn/thanh-toan?bookingId=%s&bookingCode=%s&phoneNumber=%s",
		url.QueryEscape(info.BookingID),
		url.QueryEscape(info.BookingCode),
		url.QueryEscape(info.PhoneNumber),
	)

	from := "FutaHunter <noreply@resend.dev>"
	if s.cfg.FromAddress != "" {
		name := s.cfg.FromName
		if name == "" {
			name = "FutaHunter"
		}
		from = fmt.Sprintf("%s <%s>", name, s.cfg.FromAddress)
	}

	priceStr := fmt.Sprintf("%d", info.TicketPrice)
	if info.TicketPrice > 0 {
		priceStr = formatVND(info.TicketPrice)
	}

	html := fmt.Sprintf(`
<div style="font-family:'Segoe UI',Arial,sans-serif;max-width:600px;margin:0 auto;background:#fff">
  <div style="background:linear-gradient(135deg,#E8431A,#FF6B35);padding:24px 32px;border-radius:12px 12px 0 0">
    <h1 style="color:#fff;margin:0;font-size:20px">FutaHunter</h1>
    <p style="color:rgba(255,255,255,0.85);margin:4px 0 0;font-size:13px">Dat ve thanh cong!</p>
  </div>
  <div style="padding:24px 32px;border:1px solid #eee;border-top:none;border-radius:0 0 12px 12px">
    <p style="font-size:15px;color:#333">Xin chao <strong>%s</strong>,</p>
    <p style="font-size:14px;color:#555;margin:12px 0">Ve xe cua ban da duoc dat thanh cong. Vui long thanh toan de hoan tat.</p>
    <div style="background:#f9f9f9;border-radius:8px;padding:16px;margin:16px 0">
      <table style="width:100%%;font-size:13px;color:#333">
        <tr><td style="padding:4px 0;color:#999">Tuyen</td><td style="padding:4px 0;font-weight:600;text-align:right">%s</td></tr>
        <tr><td style="padding:4px 0;color:#999">Ngay di</td><td style="padding:4px 0;font-weight:600;text-align:right">%s</td></tr>
        <tr><td style="padding:4px 0;color:#999">Ghe</td><td style="padding:4px 0;font-weight:600;text-align:right">%s</td></tr>
        <tr><td style="padding:4px 0;color:#999">Ma ve</td><td style="padding:4px 0;font-weight:600;text-align:right">%s</td></tr>
        <tr><td style="padding:4px 0;color:#999">Gia ve</td><td style="padding:4px 0;font-weight:700;color:#E8431A;text-align:right">%s</td></tr>
      </table>
    </div>
    <a href="%s" style="display:block;text-align:center;background:#E8431A;color:#fff;padding:14px 24px;border-radius:8px;text-decoration:none;font-size:15px;font-weight:700;margin:20px 0">Thanh toan ngay</a>
    <p style="font-size:12px;color:#999;margin-top:16px;text-align:center">Neu nut khong hoat dong, copy link sau vao trinh duyet:<br><a href="%s" style="color:#E8431A;word-break:break-all">%s</a></p>
  </div>
</div>`,
		esc(info.ToName),
		esc(info.RouteName),
		esc(info.TravelDate),
		esc(info.SeatName),
		esc(info.BookingCode),
		priceStr,
		paymentURL, paymentURL, paymentURL,
	)

	params := &resend.SendEmailRequest{
		From:    from,
		To:      []string{info.ToEmail},
		Subject: fmt.Sprintf("FutaHunter - Thanh toan ve %s → %s (%s)", info.OriginName, info.DestName, info.BookingCode),
		Html:    html,
	}

	sent, err := s.client.Emails.Send(params)
	if err != nil {
		return fmt.Errorf("send email to %s: %w", info.ToEmail, err)
	}

	log.Printf("Payment email sent to %s (id: %s)", info.ToEmail, sent.Id)
	return nil
}

func formatVND(amount int) string {
	s := fmt.Sprintf("%d", amount)
	n := len(s)
	if n <= 3 {
		return s + " VND"
	}
	var result []byte
	for i, c := range s {
		if i > 0 && (n-i)%3 == 0 {
			result = append(result, '.')
		}
		result = append(result, byte(c))
	}
	return string(result) + " VND"
}

func esc(s string) string {
	// Basic HTML escaping for email content
	r := ""
	for _, c := range s {
		switch c {
		case '<':
			r += "&lt;"
		case '>':
			r += "&gt;"
		case '&':
			r += "&amp;"
		case '"':
			r += "&quot;"
		default:
			r += string(c)
		}
	}
	return r
}
