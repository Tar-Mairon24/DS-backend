package services

import (
	"crypto/rand"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"math/big"
	"os"

	"github.com/wneessen/go-mail"

	"backend/internal/models"
)

type EmailService struct {
	DB *sql.DB
}

func NewEmailService(db *sql.DB) *EmailService {
	return &EmailService{
		DB: db,
	}
}

func (s *EmailService) SendVerificationEmail(toEmail string) error {
	smtpUser := os.Getenv("SMTP_USER")
	smtpPass := os.Getenv("SMTP_PASS")
	apiPort := os.Getenv("API_PORT")
	if apiPort == "" {
		apiPort = "3000"
	}
	//url := fmt.Sprintf("http://localhost:%s/verify-email", apiPort)

	verificationCode := s.GenerateVerificationCode()

	if toEmail == "" || smtpUser == "" || smtpPass == "" {
		log.Println("EMAIL_TO, SMTP_USER or SMTP_PASS environment variables are not set. Skipping email sending.")
		return errors.New("missing email configuration")
	}

	if err := s.saveVerificationCode(toEmail, verificationCode); err != nil {
		return err
	}

	message := mail.NewMsg()
	if err := message.From(smtpUser); err != nil {
		log.Fatalf("failed to set From address: %s", err)
		return err
	}
	if err := message.To(toEmail); err != nil {
		log.Fatalf("failed to set To address: %s", err)
		return err
	}
	message.Subject("Verify your email address for Desarrollo Seguro")
	message.SetBodyString(mail.TypeTextHTML, fmt.Sprintf(`
		<html>
		<body>
		<p>Bienvenido a la app para Desarrollo Seguro!</p>
		<p>Por favor verifica tu dirección de correo electrónico haciendo clic en el botón de abajo:</p>
		<p>Tu código de verificación es:</p>
		<h2>%s</h2>
		<p>Si no te has registrado en esta cuenta, por favor ignora este correo electrónico.</p>
		</body>
		</html>
	`, verificationCode))
	client, err := mail.NewClient("smtp.gmail.com", mail.WithSMTPAuth(mail.SMTPAuthAutoDiscover),
		mail.WithUsername(smtpUser), mail.WithPassword(smtpPass))
	if err != nil {
		log.Fatalf("failed to create mail client: %s", err)
		return err
	}
	if err := client.DialAndSend(message); err != nil {
		log.Fatalf("failed to send mail: %s", err)
		return err
	}

	return nil
}

func (s *EmailService) GenerateVerificationCode() string {
    n, err := rand.Int(rand.Reader, big.NewInt(1000000)) // 0..999999
    if err != nil {
        log.Printf("failed to generate verification code: %v", err)
        return "000000"
    }
    return fmt.Sprintf("%06d", n.Int64())
}

func (s *EmailService) saveVerificationCode(toEmail string, code string) error {
	query := "UPDATE Usuarios SET codigo_verificacion = ? WHERE usuario = ?"
	_, err := s.DB.Exec(query, code, toEmail)
	if err != nil {
		log.Println("Error updating verification code in database:", err)
		return err
	}
	return nil
}

func (s *EmailService) VerifyEmail(verificacionData models.EmailVerification) (bool, error) {
	if verificacionData.Code == "" || verificacionData.Email == "" {
		return false, errors.New("verification code and email must be provided")
	}

	var userID int
	query := "SELECT id_usuario FROM Usuarios WHERE codigo_verificacion = ? && usuario = ?"
	err := s.DB.QueryRow(query, verificacionData.Code, verificacionData.Email).Scan(&userID)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, errors.New("invalid verification code or email")
		}
		log.Println("Error fetching user by verification code:", err)
		return false, err
	}

	updateQuery := "UPDATE Usuarios SET verificado = 1, codigo_verificacion = NULL WHERE id_usuario = ?"
	_, err = s.DB.Exec(updateQuery, userID)
	if err != nil {
		log.Println("Error updating user verification status:", err)
		return false, err
	}

	return true, nil
}
