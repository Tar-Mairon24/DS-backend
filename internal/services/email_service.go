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

	sq "github.com/Masterminds/squirrel"
	"github.com/wneessen/go-mail"

	"backend/internal/models"
)

type EmailService struct {
	DB *sql.DB
	sq sq.StatementBuilderType
}

func NewEmailService(db *sql.DB) *EmailService {
	return &EmailService{
		DB: db,
		sq: sq.StatementBuilder.PlaceholderFormat(sq.Question),
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
	query := s.sq.Select("id_usuario", "usado").
		From("Tokens_Verificacion").
		Where(sq.Eq{"token": verificacionData.Code}).
		Where(sq.Expr("id_usuario = (SELECT id_usuario FROM Usuarios WHERE usuario = ?)", verificacionData.Email)).
		Where(sq.Gt{"fecha_expiracion": time.Now()})
	sqlStr, args, err := query.ToSql()
	if err != nil {
		log.Println("Error building SQL query:", err)
		return false, err
	}
	err = s.DB.QueryRow(sqlStr, args...).Scan(&userID, &usado)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Println("No matching verification code or email found")
			return false, errors.New("invalid verification code or email")
		}
		log.Println("Error fetching user by verification code:", err)
		return false, err
	}
	if usado == 1 {
		return false, errors.New("verification code has already been used")
	}

	updateQuery := s.sq.Update("Usuarios").Set("verificado", 1).Where(sq.Eq{"id_usuario": userID})
	updateSQL, updateArgs, err := updateQuery.ToSql()
	if err != nil {
		log.Println("Error building SQL query:", err)
		return false, err
	}
	_, err = s.DB.Exec(updateSQL, updateArgs...)
	if err != nil {
		log.Println("Error updating user verification status:", err)
		return false, err
	}
	updateTokenQuery := s.sq.Update("Tokens_Verificacion").
		Set("usado", 1).
		Set("fecha_uso", time.Now()).
		Where(sq.Eq{"token": verificacionData.Code, "id_usuario": userID})
	updateTokenSQL, updateTokenArgs, err := updateTokenQuery.ToSql()
	if err != nil {
		log.Println("Error building SQL query:", err)
		return false, err
	}
	_, err = s.DB.Exec(updateTokenSQL, updateTokenArgs...)
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

	query := s.sq.Select("usado", "motivo").From("Tokens_Verificacion").Where(sq.Eq{"id_usuario": userID})
	var usado int
	var motivo string
	sqlStr, args, err := query.ToSql()
	if err != nil {
		log.Println("Error building SQL query:", err)
		return err
	}
	err = s.DB.QueryRow(sqlStr, args...).Scan(&usado, &motivo)
	if err != nil {
		log.Println("Error fetching user verification status:", err)
		return err
	}
	if usado == 1 {
		log.Println("Code already verified for user:", toEmail)
		s.SendVerificationEmail(toEmail, "Retry")
		return nil
	}

	query = s.sq.Select("num_renvios", "token").
		From("Tokens_Verificacion").
		Where(sq.Eq{"id_usuario": userID}).
		OrderBy("fecha_creacion DESC").
		Limit(1)

	var reenviado int
	var token string
	sqlStr, args, err = query.ToSql()
	if err != nil {
		log.Println("Error building SQL query:", err)
		return err
	}
	err = s.DB.QueryRow(sqlStr, args...).Scan(&reenviado, &token)
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

	updateLatestQuery := s.sq.Update("Tokens_Verificacion").
		Set("num_renvios", sq.Expr("num_renvios + 1")).
		Set("token", verificationCode).
		Set("fecha_modificacion", time.Now()).
		Set("fecha_expiracion", expirationDate).
		Where(sq.Eq{"id_usuario": userID, "token": token}).
		Suffix("ORDER BY fecha_creacion DESC LIMIT 1")
	updateLatestSQL, updateLatestArgs, err := updateLatestQuery.ToSql()
	if err != nil {
		log.Println("Error building SQL query:", err)
		return err
	}
	_, err = s.DB.Exec(updateLatestSQL, updateLatestArgs...)
	if err != nil {
		log.Println("Error updating resend count:", err)
		return err
	}
	log.Printf("Resend count updated for %s", toEmail)

	return nil
}

func (s *EmailService) saveVerificationCode(toEmail string, code string, motivo string) error {
	expirationDate := time.Now().Add(48 * time.Hour)

	query := s.sq.Insert("Tokens_Verificacion").
		Columns("token", "id_usuario", "fecha_expiracion", "fecha_creacion", "usado", "motivo").
		Values(code, sq.Expr("(SELECT id_usuario FROM Usuarios WHERE usuario = ?)", toEmail), expirationDate, time.Now(), 0, motivo)
	sqlStr, args, err := query.ToSql()
	if err != nil {
		log.Println("Error building SQL query:", err)
		return err
	}
	_, err = s.DB.Exec(sqlStr, args...)
	if err != nil {
		log.Println("Error inserting verification code in database:", err)
		return err
	}
	return nil
}

func (s *EmailService) getIdFromEmail(toEmail string) (int, error) {
	var userID int
	query := s.sq.Select("id_usuario").From("Usuarios").Where(sq.Eq{"usuario": toEmail})
	sqlStr, args, err := query.ToSql()
	if err != nil {
		log.Println("Error building SQL query:", err)
		return 0, err
	}
	err = s.DB.QueryRow(sqlStr, args...).Scan(&userID)
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
