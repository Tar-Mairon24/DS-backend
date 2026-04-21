package service

import (
	"context"
	"crypto/rand"
	"errors"
	"fmt"
	"math/big"
	"os"

	"github.com/sirupsen/logrus"
	"github.com/wneessen/go-mail"

	"ds-backend/internal/domain/models"
	"ds-backend/internal/domain/ports"
)

type EmailService struct {
	emailRepo ports.EmailRepository
}

func NewEmailService(emailRepo ports.EmailRepository) ports.EmailService {
	return &EmailService{
		emailRepo: emailRepo,
	}
}

func (s *EmailService) SendVerificationEmail(ctx context.Context, toEmail string, motivo string) error {
	verificationCode := s.generateVerificationCode()

	if err := s.emailRepo.SaveVerificationCode(ctx, toEmail, verificationCode, motivo); err != nil {
		return err
	}

	if err := s.sendEmail(toEmail, verificationCode); err != nil {
		return err
	}

	return nil
}

func (s *EmailService) VerifyEmail(ctx context.Context, verificacionData models.EmailVerification) (bool, error) {
	if verificacionData.Code == "" || verificacionData.Email == "" {
		return false, errors.New("verification code and email must be provided")
	}

	userID, isUsed, err := s.emailRepo.VerifyTokenAndUser(ctx, verificacionData.Code, verificacionData.Email)
	if err != nil {
		return false, err
	}

	if isUsed {
		return false, errors.New("verification code has already been used")
	}

	if err := s.emailRepo.UpdateUserVerificationStatus(ctx, userID); err != nil {
		return false, err
	}

	if err := s.emailRepo.UpdateTokenAsUsed(ctx, verificacionData.Code, userID); err != nil {
		return false, err
	}

	return true, nil
}

func (s *EmailService) ResendVerificationEmail(ctx context.Context, toEmail string) error {
	if toEmail == "" {
		return errors.New("email must be provided")
	}

	userID, err := s.emailRepo.GetIDFromEmail(ctx, toEmail)
	if err != nil {
		return err
	}

	usado, motivo, err := s.emailRepo.GetTokenVerificationStatus(ctx, userID)
	if err != nil {
		return err
	}

	if usado == true {
		logrus.Println("Code already verified for user:", toEmail)
		motivo = "Reintento de verificacion por el motivo: " + motivo
		return s.SendVerificationEmail(ctx, toEmail, motivo)
	}

	reenviado, token, err := s.emailRepo.GetLatestTokenInfo(ctx, userID)
	if err != nil {
		return err
	}

	logrus.Printf("Resend count for %s: %d", toEmail, reenviado)

	if reenviado >= 3 {
		return errors.New("maximum number of resends reached")
	}

	verificationCode := s.generateVerificationCode()

	if err := s.sendEmail(toEmail, verificationCode); err != nil {
		return err
	}

	if err := s.emailRepo.UpdateTokenResendInfo(ctx, userID, token, verificationCode); err != nil {
		return err
	}

	logrus.Printf("Resend count updated for %s", toEmail)

	return nil
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
		logrus.Errorf("failed to set From address: %s", err)
		return err
	}
	if err := message.To(toEmail); err != nil {
		logrus.Errorf("failed to set To address: %s", err)
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
		logrus.Errorf("failed to create mail client: %s", err)
		return err
	}
	if err := client.DialAndSend(message); err != nil {
		logrus.Errorf("failed to send mail: %s", err)
		return err
	}

	logrus.Printf("Verification email sent to %s", toEmail)
	return nil
}

func (s *EmailService) generateVerificationCode() string {
	n, err := rand.Int(rand.Reader, big.NewInt(1000000)) // 0..999999
	if err != nil {
		logrus.Errorf("failed to generate verification code: %v", err)
		return "000000"
	}
	return fmt.Sprintf("%06d", n.Int64())
}
