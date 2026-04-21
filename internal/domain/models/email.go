package models

type EmailVerification struct {
	Email string `json:"email"`
	Code  string `json:"code"`
}

type EmailRequest struct {
	Email 	string `json:"email"`
	Reason 	string `json:"reason,omitempty"`
}

type EmailResendRequest struct {
	Email string `json:"email"`
}