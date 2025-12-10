package users

type UserRequest struct {
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

type ChangePasswordRequest struct {
	OldPassword     string `json:"old_password" binding:"required"`
	NewPassword     string `json:"new_password" binding:"required,min=6"`
	ConfirmPassword string `json:"confirm_password" binding:"required,min=6"`
}
