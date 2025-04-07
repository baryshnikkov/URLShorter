package auth

type RegisterReq struct {
	Email     string `json:"email" validate:"required,email,max=255"`
	Login     string `json:"login" validate:"required,min=3,max=50"`
	Password  string `json:"password" validate:"required,min=8,max=64"`
	FirstName string `json:"first_name" validate:"omitempty,max=100"`
	LastName  string `json:"last_name" validate:"omitempty,max=100"`
}

type RegisterRes struct {
	Email        string `json:"email"`
	Login        string `json:"login"`
	FirstName    string `json:"first_name"`
	LastName     string `json:"last_name"`
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type LoginReq struct {
	Email    string `json:"email" validate:"email,max=255"`
	Password string `json:"password" validate:"required,min=8,max=64"`
}

type LoginRes struct {
	Email        string `json:"email"`
	Login        string `json:"login"`
	FirstName    string `json:"first_name"`
	LastName     string `json:"last_name"`
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}
