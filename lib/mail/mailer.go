package mail

import (
	"fmt"
	"health-tech/lib"
	"log"
	"os"

	"gopkg.in/gomail.v2"
)

const (
	smtpHost = "smtp.gmail.com"
	smtpPort = 587

	// Base URL for the API endpoint that will handle the verification request
	scheme  = "http"
	baseURL = "%s://%s/health-tech/confirm"
)

// SendConfirmationEmail sends a verification link containing the JWT.
func SendConfirmationEmail(recipientEmail, jwtToken string) error {

	// 1. CONSTRUCT THE FULL VERIFICATION URL
	verificationURL := fmt.Sprintf("%s?token=%s", fmt.Sprintf(baseURL, scheme, lib.GetServiceHostWithPort("user")), jwtToken)

	// 2. BUILD THE EMAIL CONTENT
	msg := gomail.NewMessage()
	msg.SetHeader("From", "no-reply@yourdomain.com")
	msg.SetHeader("To", recipientEmail)
	msg.SetHeader("Subject", "Action Required: Confirm Your Email Address")

	// HTML body:
	htmlBody := fmt.Sprintf(`
		<html>
		<body>
			<p>Welcome! Your HealthTech registration is almost complete.</p>
			<p>Please click the button below to **verify your email address** and activate your account:</p>
			
			<div style="margin-top: 20px;">
				<a href="%s" 
				   style="padding: 10px 20px; background-color: #007bff; color: white; text-decoration: none; border-radius: 5px; font-weight: bold;">
					Verify My Account
				</a>
			</div>
			
			<p style="margin-top: 20px; font-size: 10px; color: #888;">
				If the button doesn't work, copy and paste this link into your browser:<br>
				<a href="%s">%s</a>
			</p>
			<p>The link will expire in 1 hour.</p>
		</body>
		</html>
	`, verificationURL, verificationURL, verificationURL)

	msg.SetBody("text/html", htmlBody)

	// get the smtp pass from the environment variable
	smtpPass := os.Getenv("SMTP_PASS")
	smtpUser := os.Getenv("SMTP_USER")
	if smtpPass == "" || smtpUser == "" {
		return fmt.Errorf("SMTP_PASS or SMTP_USER is not set")
	}

	dial := gomail.NewDialer(smtpHost, smtpPort, smtpUser, smtpPass)

	if err := dial.DialAndSend(msg); err != nil {
		return fmt.Errorf("failed to send confirmation email: %w", err)
	}

	log.Printf("Verification email successfully dispatched to %s.", recipientEmail)
	return nil
}
