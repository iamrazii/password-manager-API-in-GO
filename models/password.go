package models

type Password struct {
	Id            int    `json:"id"`
	Userid        int    `json:"userid"`
	Platform_Name string `json:"platform_name"`
	Password      string `json:"password"`
	Account_Email string `json:"account_email"`
}

type UpdatePassword struct {
	Password string `json:"password" binding:"required"`
}
