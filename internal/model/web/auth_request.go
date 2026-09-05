package web

type RegisterRequest struct {
	Name     string `json:"name" validate:"required,min=3,max=255"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
}

type LogOutRequest struct {
	RefreshToken string `cookie:"refresh_token" validate:"required"`
}

type RefreshTokenRequest struct {
	RefreshToken string `cookie:"refresh_token" json:"refresh_token" validate:"required"`
	AccessToken  string `header:"authorization" json:"access_token" validate:"required"`
}
