package auth

import (
	"URLShorter/internal/user"
	"URLShorter/pkg/di"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"strings"
)

type Service struct {
	userRepository *user.Repository
}

func NewService(userRepository *user.Repository) *Service {
	return &Service{
		userRepository: userRepository,
	}
}

func (s *Service) Register(email, login, password, firstname, lastname string) (*user.User, error) {
	type checkResult struct {
		taken bool
		err   error
	}
	emailCh := make(chan checkResult)
	loginCh := make(chan checkResult)
	go func() {
		defer close(emailCh)
		existedUser, err := s.userRepository.FindByEmail(email)
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			emailCh <- checkResult{err: fmt.Errorf("error checking email: %w", err)}
			return
		}
		emailCh <- checkResult{taken: existedUser != nil}
	}()
	go func() {
		defer close(loginCh)
		existedUser, err := s.userRepository.FindByLogin(login)
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			loginCh <- checkResult{err: fmt.Errorf("error checking login: %w", err)}
			return
		}
		loginCh <- checkResult{taken: existedUser != nil}
	}()
	emailCheckResult := <-emailCh
	loginCheckResult := <-loginCh
	if emailCheckResult.err != nil {
		return nil, emailCheckResult.err
	}
	if loginCheckResult.err != nil {
		return nil, loginCheckResult.err
	}
	var validationErrors []string
	if emailCheckResult.taken {
		validationErrors = append(validationErrors, fmt.Sprintf("%s: %s", ErrUserExistsEmail, email))
	}
	if loginCheckResult.taken {
		validationErrors = append(validationErrors, fmt.Sprintf("%s: %s", ErrUserExistsLogin, login))
	}
	if len(validationErrors) > 0 {
		return nil, errors.New(strings.Join(validationErrors, "; "))
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	verificationToken := uuid.NewString()

	newUser := &user.User{
		Email:             email,
		Login:             login,
		PasswordHash:      string(hashedPassword),
		FirstName:         firstname,
		LastName:          lastname,
		Role:              di.RoleUser,
		IsBanned:          false,
		EmailVerified:     false,
		VerificationToken: verificationToken,
	}
	registeredUser, err := s.userRepository.Create(newUser)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", ErrUserNotCreated, err)
	}

	return registeredUser, nil
}
