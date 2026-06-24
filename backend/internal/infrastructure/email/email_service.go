package email

import (
	"bytes"
	"fmt"
	"html/template"

	"github.com/Kal-el21/backend/configs"
	"gopkg.in/gomail.v2"
)

type EmailService interface {
	SendOTP(to string, recipientName string, otp string) error
	SendNotification(to string, subject string, htmlBody string) error
}

type emailService struct {
	cfg *configs.Config
}

func NewEmailService(cfg *configs.Config) EmailService {
	return &emailService{cfg: cfg}
}

func (s *emailService) newDialer() *gomail.Dialer {
	d := gomail.NewDialer(s.cfg.SMTPHost, s.cfg.SMTPPort, s.cfg.SMTPUser, s.cfg.SMTPPassword)
	return d
}

// SendOTP mengirimkan kode OTP 2FA untuk login.
func (s *emailService) SendOTP(to string, recipientName string, otp string) error {
	subject := "Kode Verifikasi Login PPMS"
	body := s.renderOTPEmail(recipientName, otp)
	return s.sendHTML(to, subject, body)
}

// SendNotification mengirimkan notifikasi sistem via email
// (dipakai oleh notification service saat channel=EMAIL dan enabled=true).
func (s *emailService) SendNotification(to string, subject string, htmlBody string) error {
	return s.sendHTML(to, subject, htmlBody)
}

func (s *emailService) sendHTML(to, subject, htmlBody string) error {
	if s.cfg.SMTPUser == "" || s.cfg.SMTPPassword == "" {
		// Email tidak dikonfigurasi — log warning tapi tidak error
		// agar fitur lain tidak terblokir saat SMTP belum di-setup
		return fmt.Errorf("email service not configured: SMTP_USER or SMTP_PASSWORD is empty")
	}

	m := gomail.NewMessage()
	m.SetHeader("From", s.cfg.SMTPFrom)
	m.SetHeader("To", to)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", htmlBody)

	return s.newDialer().DialAndSend(m)
}

// renderOTPEmail menghasilkan HTML email OTP dengan template sederhana
// yang konsisten dengan design system PPMS (biru + merah).
func (s *emailService) renderOTPEmail(name, otp string) string {
	const tmpl = `<!DOCTYPE html>
<html>
<head><meta charset="UTF-8"></head>
<body style="margin:0;padding:0;background:#F1F5F9;font-family:'Inter',Arial,sans-serif;">
  <table width="100%" cellpadding="0" cellspacing="0" style="padding:40px 20px;">
    <tr><td align="center">
      <table width="480" style="background:#FFFFFF;border-radius:12px;border:1px solid #E2E8F0;overflow:hidden;">
        <tr>
          <td style="background:linear-gradient(135deg,#2563EB,#DC2626);padding:28px 32px;">
            <p style="margin:0;color:white;font-size:20px;font-weight:700;letter-spacing:-0.02em;">PPMS</p>
            <p style="margin:6px 0 0;color:rgba(255,255,255,0.8);font-size:13px;">Project Portfolio Management System</p>
          </td>
        </tr>
        <tr>
          <td style="padding:32px;">
            <p style="margin:0 0 8px;font-size:18px;font-weight:600;color:#0F172A;">Kode Verifikasi Login</p>
            <p style="margin:0 0 24px;font-size:14px;color:#64748B;line-height:1.6;">
              Halo {{.Name}}, gunakan kode berikut untuk menyelesaikan login ke PPMS.
              Kode berlaku selama <strong>10 menit</strong>.
            </p>
            <div style="background:#EFF6FF;border:1px solid #BFDBFE;border-radius:8px;padding:20px;text-align:center;margin-bottom:24px;">
              <p style="margin:0;font-size:36px;font-weight:700;letter-spacing:0.15em;color:#1D4ED8;">{{.OTP}}</p>
            </div>
            <p style="margin:0;font-size:12px;color:#94A3B8;line-height:1.5;">
              Jika Anda tidak mencoba login, abaikan email ini.
              Jangan bagikan kode ini kepada siapa pun.
            </p>
          </td>
        </tr>
        <tr>
          <td style="background:#F8FAFC;padding:16px 32px;border-top:1px solid #E2E8F0;">
            <p style="margin:0;font-size:11px;color:#94A3B8;">© 2026 PPMS Internal. Email otomatis, tidak perlu dibalas.</p>
          </td>
        </tr>
      </table>
    </td></tr>
  </table>
</body>
</html>`

	t, err := template.New("otp").Parse(tmpl)
	if err != nil {
		return fmt.Sprintf("<p>Kode OTP Anda: <strong>%s</strong></p>", otp)
	}

	var buf bytes.Buffer
	if err := t.Execute(&buf, map[string]string{"Name": name, "OTP": otp}); err != nil {
		return fmt.Sprintf("<p>Kode OTP Anda: <strong>%s</strong></p>", otp)
	}

	return buf.String()
}
