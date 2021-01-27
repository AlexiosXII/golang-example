package users

type UserRequest struct {
	Email    string `json:"email" validate:"required"`
	Password string `json:"password" validate:"required"`
	Name     string `json:"name" validate:"required"`
	Age      int    `json:"age" validate:"required"`
}

type UserLogin struct {
	Email    string `json:"email" validate:"required"`
	Password string `json:"password" validate:"required"`
}
