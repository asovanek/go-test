package auth

type SignUpRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
}

// SignInRequest body for signin.
type SignInRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

// AuthTokensResponse is returned after signup/signin.
type AuthTokensResponse struct {
	Token string           `json:"token"`
	User  AuthUserCompact `json:"user"`
}

// AuthUserCompact is nested user summary in auth responses.
type AuthUserCompact struct {
	ID    string `json:"id"`
	Email string `json:"email"`
}
