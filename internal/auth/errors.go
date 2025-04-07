package auth

const (
	ErrUserExistsEmail = "user with that email already exists"
	ErrUserExistsLogin = "user with that login already exists"
	ErrUserNotFound    = "user not found"
	ErrUserNotCreated  = "user not created"
	ErrWrongCredential = "wrong email or password"
)
