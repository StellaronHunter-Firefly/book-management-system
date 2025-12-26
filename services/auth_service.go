package services

import (
	"book-management-system/models"
	"book-management-system/repositories"
	"book-management-system/utils"
	"fmt"
	"errors"
)

type AuthService interface {
	Register(username, password, email string) (*models.User, error)
	Login(username, password string) (*models.User, string, error)
	GetUserByID(id uint) (*models.User, error)
	ChangeUsername(userID uint, newUsername string) error
	ChangePassword(userID uint, oldPassword, newPassword string) error
	GetAllUsers() ([]models.User, error)
	VerifyPassword(userID uint, password string) (bool, error)
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

// ChangeUsername 更改用户名
func (s *authService) ChangeUsername(userID uint, newUsername string) error {
	if newUsername == "" {
		return errors.New("新用户名不能为空")
	}
	
	if len(newUsername) < 3 {
		return errors.New("用户名长度至少3个字符")
	}
	
	if len(newUsername) > 50 {
		return errors.New("用户名长度不能超过50个字符")
	}
	
	return s.userRepo.UpdateUsername(userID, newUsername)
}

// VerifyPassword 验证原密码
func (s *authService) VerifyPassword(userID uint, password string) (bool, error) {
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return false, errors.New("用户不存在")
	}
	
	// 由于移除了密码加密，直接比较
	return user.Password == password, nil
}

// ChangePassword 更改密码
func (s *authService) ChangePassword(userID uint, oldPassword, newPassword string) error {
	if newPassword == "" {
		return errors.New("新密码不能为空")
	}
	
	if len(newPassword) < 6 {
		return errors.New("密码长度至少6个字符")
	}
	
	// 验证原密码
	valid, err := s.VerifyPassword(userID, oldPassword)
	if err != nil {
		return err
	}
	if !valid {
		return errors.New("原密码错误")
	}
	
	return s.userRepo.UpdatePassword(userID, newPassword)
}

// GetAllUsers 获取所有用户（仅管理员使用）
func (s *authService) GetAllUsers() ([]models.User, error) {
	users, err := s.userRepo.FindAll()
	if err != nil {
		return nil, fmt.Errorf("获取用户列表失败: %w", err)
	}
	
	// 移除密码字段，保护隐私
	for i := range users {
		users[i].Password = ""
	}
	
	return users, nil
}