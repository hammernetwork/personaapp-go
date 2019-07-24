package forms

type UserSignup struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}
