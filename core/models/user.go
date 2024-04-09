package models

type UserRegistration struct {
	UserName   string `json:"user_name" validate:"required,min=1"`
	Email      string `json:"email" validate:"required,min=1,email"`
	Phone      string `json:"phone" validate:"required,e164"`
	UserStatus string `json:"user_status" validate:"required"`
}

type UserLogin struct {
	RecordCode int    `json:"record_code"`
	Password   string `json:"password"`
}
