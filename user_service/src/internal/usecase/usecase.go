package usecase

import (
	"errors"
	"fmt"
	"user_service/src/internal/config"
	session "user_service/src/internal/core/sessions"
	"user_service/src/internal/core/users"
	"user_service/src/pkg/generatejwt"
	"user_service/src/pkg/hashpassword"
)

type UserService struct {
	userRepo    users.UserRepoImpl
	sessionRepo session.SessionRepoImpl
	jwtKey      string
}

func NewUserService(userRepo users.UserRepoImpl, sessionRepo session.SessionRepoImpl, jwtKey string) users.UserServiceImpl {
	return UserService{userRepo: userRepo, sessionRepo: sessionRepo, jwtKey: jwtKey}
}

func (us UserService) CreateUser(user users.User) (users.User, error) {
	createdUser, err := us.userRepo.CreateUser(user)
	if err != nil {
		// return createdUser, err
		return createdUser, errors.New("user already exists. Username should be unique")
	}
	return createdUser, nil
}

func (us UserService) LoginUser(user users.User, config *config.Config) (users.LoginResponse, error) {
	loginResponse := users.LoginResponse{}

	foundUser, err := us.userRepo.FindUserByUsername(user.Username)
	if err != nil {
		return loginResponse, fmt.Errorf("invalid USERNAME or password")
	}

	if err := hashpassword.CheckPassword(foundUser.Password, user.Password); err != nil {
		return loginResponse, fmt.Errorf("invalid username or PASSWORD")
	}

	// Generate access token
	tokenString, tokenExpiration, err := generatejwt.GenerateJWT(foundUser.Id, config.JWT_SECRET)
	if err != nil {
		return loginResponse, fmt.Errorf("failed to create access token")
	}

	// Generate Session token

	session, err := generatejwt.GenerateSession(foundUser.Id)
	if err != nil {
		return loginResponse, fmt.Errorf("failed to generate session token")
	}

	err = us.sessionRepo.CreateSession(session)
	if err != nil {
		return loginResponse, fmt.Errorf("failed to create session")
	}

	loginResponse.FoundUser = foundUser
	loginResponse.TokenString = tokenString
	loginResponse.TokenExpire = tokenExpiration
	loginResponse.Session = session

	return loginResponse, nil
}

func (us UserService) GetUserById(userId int) (users.UserProfile, error) {
	user, err := us.userRepo.FindUserById(userId)
	if err != nil {
		return user, fmt.Errorf("can't get user try again later")
	}
	return user, nil
}

func (us UserService) LogoutUser(userId int) error {
	err := us.sessionRepo.DeleteSession(userId)
	if err != nil {
		return fmt.Errorf("failed to logout user, please try again later")
	}

	return nil
}
