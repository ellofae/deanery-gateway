package dto

type UserRegistered struct {
	UserName string
	Email    string
	Phone    string
}

type ProfileInformation struct {
	Username   string `json:"user_name"`
	RecordCode int    `json:"record_code"`
}
