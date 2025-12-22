package services

import (
	"book-management-system/models"
	"book-management-system/repositories"
	"book-management-system/utils"
	"errors"
)

type AuthService interface {
	Register(username, password, email string) (*models.User, error)
	Login(username, password string) (*models.User, string, error)
	GetUserByID(id uint) (*models.User, error)
}

type authService struct {
	userRepo repositories.UserRepository
}

func NewAuthService(userRepo repositories.UserRepository) AuthService {
	return &authService{userRepo: userRepo}
}

func (s *authService) Register(username, password, email string) (*models.User, error) {
	if _, err := s.userRepo.FindByUsername(username); err == nil {
		return nil, errors.New("用户名已存在")
	}

	if _, err := s.userRepo.FindByEmail(email); err == nil {
		return nil, errors.New("邮箱已存在")
	}

	user := &models.User{
		Username: username,
		Password: password,
		Email:    email,
		Role:     models.RoleUser,
	}

	if err := s.userRepo.Create(user); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *authService) Login(username, password string) (*models.User, string, error) {
	user, err := s.userRepo.FindByUsername(username)
	if err != nil {
		return nil, "", errors.New("用户名或密码错误")
	}

	if user.Password != password {
		return nil, "", errors.New("用户名或密码错误")
	}

	token, err := utils.GenerateToken(user)
	if err != nil {
		return nil, "", err
	}

	return user, token, nil
}

func (s *authService) GetUserByID(id uint) (*models.User, error) {
	return s.userRepo.FindByID(id)
}
