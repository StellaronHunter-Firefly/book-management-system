package repositories

import (
	"book-management-system/config"
	"book-management-system/models"
	"fmt"

	"gorm.io/gorm"
)

type UserRepository interface {
	Create(user *models.User) error
	FindByID(id uint) (*models.User, error)
	FindByUsername(username string) (*models.User, error)
	FindByEmail(email string) (*models.User, error)
	Update(user *models.User) error
	Delete(id uint) error
	FindAll() ([]models.User, error)
	Count() (int64, error)
	CountByRole(role models.UserRole) (int64, error)
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository() UserRepository {
	return &userRepository{db: config.DB}
}

func (r *userRepository) Create(user *models.User) error {
	if user.Username == "" {
		return fmt.Errorf("用户名不能为空")
	}

	var existing models.User
	if err := r.db.Where("username = ?", user.Username).First(&existing).Error; err == nil {
		return fmt.Errorf("用户名已存在")
	}

	if user.Email != "" {
		if err := r.db.Where("email = ?", user.Email).First(&existing).Error; err == nil {
			return fmt.Errorf("邮箱已存在")
		}
	}

	return r.db.Create(user).Error
}

func (r *userRepository) FindByID(id uint) (*models.User, error) {
	if id == 0 {
		return nil, fmt.Errorf("无效的用户ID")
	}

	var user models.User
	err := r.db.First(&user, id).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("用户不存在")
		}
		return nil, fmt.Errorf("查询用户失败: %w", err)
	}

	return &user, nil
}

func (r *userRepository) FindByUsername(username string) (*models.User, error) {
	if username == "" {
		return nil, fmt.Errorf("用户名不能为空")
	}

	var user models.User
	err := r.db.Where("username = ?", username).First(&user).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("用户不存在")
		}
		return nil, fmt.Errorf("查询用户失败: %w", err)
	}

	return &user, nil
}

func (r *userRepository) FindByEmail(email string) (*models.User, error) {
	if email == "" {
		return nil, fmt.Errorf("邮箱不能为空")
	}

	var user models.User
	err := r.db.Where("email = ?", email).First(&user).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("用户不存在")
		}
		return nil, fmt.Errorf("查询用户失败: %w", err)
	}

	return &user, nil
}

func (r *userRepository) Update(user *models.User) error {
	var existing models.User
	if err := r.db.First(&existing, user.ID).Error; err != nil {
		return fmt.Errorf("用户不存在")
	}

	if user.Username != "" && user.Username != existing.Username {
		var count int64
		if err := r.db.Model(&models.User{}).
			Where("username = ? AND id != ?", user.Username, user.ID).
			Count(&count).Error; err != nil {
			return fmt.Errorf("检查用户名失败: %w", err)
		}
		if count > 0 {
			return fmt.Errorf("用户名已存在")
		}
	}

	if user.Email != "" && user.Email != existing.Email {
		var count int64
		if err := r.db.Model(&models.User{}).
			Where("email = ? AND id != ?", user.Email, user.ID).
			Count(&count).Error; err != nil {
			return fmt.Errorf("检查邮箱失败: %w", err)
		}
		if count > 0 {
			return fmt.Errorf("邮箱已存在")
		}
	}

	return r.db.Save(user).Error
}

func (r *userRepository) Delete(id uint) error {
	var user models.User
	if err := r.db.First(&user, id).Error; err != nil {
		return fmt.Errorf("用户不存在")
	}

	var borrowedCount int64
	if err := r.db.Model(&models.BorrowRecord{}).
		Where("user_id = ? AND returned_at IS NULL", id).
		Count(&borrowedCount).Error; err != nil {
		return fmt.Errorf("检查借阅记录失败: %w", err)
	}

	if borrowedCount > 0 {
		return fmt.Errorf("用户有未归还的图书，无法删除")
	}

	return r.db.Delete(&models.User{}, id).Error
}

func (r *userRepository) FindAll() ([]models.User, error) {
	var users []models.User

	if err := r.db.Model(&models.User{}).
		Select("id, username, email, role, created_at, updated_at").
		Order("created_at DESC").
		Find(&users).Error; err != nil {
		return nil, fmt.Errorf("查询用户列表失败: %w", err)
	}

	return users, nil
}

func (r *userRepository) Count() (int64, error) {
	var count int64
	if err := r.db.Model(&models.User{}).Count(&count).Error; err != nil {
		return 0, fmt.Errorf("统计用户失败: %w", err)
	}
	return count, nil
}

func (r *userRepository) CountByRole(role models.UserRole) (int64, error) {
	var count int64
	if err := r.db.Model(&models.User{}).
		Where("role = ?", role).
		Count(&count).Error; err != nil {
		return 0, fmt.Errorf("按角色统计失败: %w", err)
	}
	return count, nil
}
