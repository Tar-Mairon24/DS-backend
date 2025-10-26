package services

import (
	"crypto/rand"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"math/big"
	"os"
	"time"

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

func (s *EmailService) SendVerificationEmail(toEmail string, motivo string) error {
	verificationCode := s.generateVerificationCode()

	if err := s.saveVerificationCode(toEmail, verificationCode, motivo); err != nil {
		return err
	}

	if err := s.sendEmail(toEmail, verificationCode); err != nil {
		return err	
	}

	return nil
}

func (s *EmailService) VerifyEmail(verificacionData models.EmailVerification) (bool, error) {
	if verificacionData.Code == "" || verificacionData.Email == "" {
		return false, errors.New("verification code and email must be provided")
	}

	var userID int
	var usado int
	query := "SELECT id_usuario, usado FROM Tokens_Verificacion WHERE token = ? AND id_usuario = (SELECT id_usuario FROM Usuarios WHERE usuario = ? ) AND fecha_expiracion > ?"
	err := s.DB.QueryRow(query, verificacionData.Code, verificacionData.Email, time.Now()).Scan(&userID, &usado)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, errors.New("invalid verification code or email")
		}
		log.Println("Error fetching user by verification code:", err)
		return false, err
	}
	if usado == 1 {
		return false, errors.New("verification code has already been used")
	}

	updateQuery := "UPDATE Usuarios SET verificado = 1 WHERE id_usuario = ?"
	_, err = s.DB.Exec(updateQuery, userID)
	if err != nil {
		log.Println("Error updating user verification status:", err)
		return false, err
	}
	updateQuery = "UPDATE Tokens_Verificacion SET usado = 1, fecha_uso = ? WHERE token = ? AND id_usuario = ?"
	_, err = s.DB.Exec(updateQuery, time.Now(), verificacionData.Code, userID)
	if err != nil {
		log.Println("Error updating token status:", err)
		return false, err
	}

	return true, nil
}

func (s *EmailService) ResendVerificationEmail(toEmail string) error {
	if toEmail == "" {
		return errors.New("email must be provided")
	}

	userID, error := s.getIdFromEmail(toEmail)
	if error != nil {
		return error
	}

	query := "SELECT num_renvios, token FROM Tokens_Verificacion WHERE id_usuario = ? ORDER BY fecha_creacion DESC LIMIT 1"

	var reenviado int
	var token string
	err := s.DB.QueryRow(query, userID).Scan(&reenviado, &token)
	if err != nil {
		if err == sql.ErrNoRows {
			reenviado = 0
		} else {
			log.Println("Error fetching resend count:", err)
			return err
		}
	}

	log.Printf("Resend count for %s: %d", toEmail, reenviado)

	if reenviado >= 3 {
		return errors.New("maximum number of resends reached")
	}

	verificationCode := s.generateVerificationCode()

	if err := s.sendEmail(toEmail, verificationCode); err != nil {
		return err
	}

	expirationDate := time.Now().Add(48 * time.Hour)

	query = "UPDATE Tokens_Verificacion SET num_renvios = num_renvios + 1, token = ?, fecha_modificacion = ?, fecha_expiracion = ? WHERE id_usuario = ? AND token = ? ORDER BY fecha_creacion DESC LIMIT 1"
	_, err = s.DB.Exec(query, verificationCode, time.Now(), expirationDate, userID, token)
	if err != nil {
		log.Println("Error updating resend count:", err)
		return err
	}
	log.Printf("Resend count updated for %s", toEmail)

	return nil
}

func (s *EmailService) saveVerificationCode(toEmail string, code string, motivo string) error {
	expirationDate := time.Now().Add(48 * time.Hour)

	query := "INSERT INTO Tokens_Verificacion (token, id_usuario, fecha_expiracion, fecha_creacion, usado, motivo) VALUES (?, (SELECT id_usuario FROM Usuarios WHERE usuario = ?), ?, ?, ?, ?)"
	_, err := s.DB.Exec(query, code, toEmail, expirationDate, time.Now(), 0, motivo)
	if err != nil {
		log.Println("Error inserting verification code in database:", err)
		return err
	}
	return nil
}

func (s *EmailService) getIdFromEmail(toEmail string) (int, error) {
	var userID int
	query := "SELECT id_usuario FROM Usuarios WHERE usuario = ?"
	err := s.DB.QueryRow(query, toEmail).Scan(&userID)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0, errors.New("user not found")
		}
		log.Println("Error fetching user ID:", err)
		return 0, err
	}
	return userID, nil
}	

func (s *EmailService) sendEmail(toEmail string, verificationCode string) error {
	smtpUser := os.Getenv("SMTP_USER")
	smtpPass := os.Getenv("SMTP_PASS")
	apiPort := os.Getenv("API_PORT")
	if apiPort == "" {
		apiPort = "3000"
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
		<p>Te enviamos este correo para que verificar que eres tú.</p>
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

	log.Printf("Verification email sent to %s", toEmail)
	return nil
}

func (s *EmailService) generateVerificationCode() string {
	n, err := rand.Int(rand.Reader, big.NewInt(1000000)) // 0..999999
	if err != nil {
		log.Printf("failed to generate verification code: %v", err)
		return "000000"
	}
	return fmt.Sprintf("%06d", n.Int64())
}