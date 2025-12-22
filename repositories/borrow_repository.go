package repositories

import (
	"book-management-system/config"
	"book-management-system/models"
	"fmt"
	"time"

	"gorm.io/gorm"
)

type BorrowRepository interface {
	Borrow(record *models.BorrowRecord) error
	Return(recordID uint) error
	FindActiveByUserAndBook(userID, bookID uint) (*models.BorrowRecord, error)
	FindActiveByUser(userID uint) ([]models.BorrowRecord, error)
	FindAll() ([]models.BorrowRecord, error)
	FindByUser(userID uint) ([]models.BorrowRecord, error)
}

type borrowRepository struct {
	db *gorm.DB
}

func NewBorrowRepository() BorrowRepository {
	return &borrowRepository{db: config.DB}
}

func (r *borrowRepository) Borrow(record *models.BorrowRecord) error {
	var user models.User
	if err := r.db.First(&user, record.UserID).Error; err != nil {
		return fmt.Errorf("用户不存在")
	}

	var book models.Book
	if err := r.db.First(&book, record.BookID).Error; err != nil {
		return fmt.Errorf("图书不存在")
	}

	if book.Available <= 0 {
		return fmt.Errorf("图书已全部借出")
	}

	var existingBorrow int64
	if err := r.db.Model(&models.BorrowRecord{}).
		Where("user_id = ? AND book_id = ? AND returned_at IS NULL",
			record.UserID, record.BookID).
		Count(&existingBorrow).Error; err != nil {
		return fmt.Errorf("检查借阅记录失败: %w", err)
	}

	if existingBorrow > 0 {
		return fmt.Errorf("您已借阅此书且尚未归还")
	}

	if record.BorrowedAt.IsZero() {
		record.BorrowedAt = time.Now()
	}
	if record.DueDate.IsZero() {
		record.DueDate = record.BorrowedAt.Add(14 * 24 * time.Hour)
	}

	if err := r.db.Create(record).Error; err != nil {
		return fmt.Errorf("创建借阅记录失败: %w", err)
	}

	book.Available -= 1
	return r.db.Save(&book).Error
}

// func (r *borrowRepository) Return(recordID uint) error {
// 	var record models.BorrowRecord
// 	if err := r.db.Preload("Book").
// 		First(&record, recordID).Error; err != nil {
// 		return fmt.Errorf("借阅记录不存在")
// 	}

// 	if record.ReturnedAt != nil {
// 		return fmt.Errorf("图书已归还")
// 	}

// 	now := time.Now()
// 	record.ReturnedAt = &now

// 	if err := r.db.Save(&record).Error; err != nil {
// 		return fmt.Errorf("更新归还时间失败: %w", err)
// 	}

// 	if record.Book != nil {
// 		record.Book.Available += 1
// 		return r.db.Save(&record.Book).Error
// 	}

//		return nil
//	}
func (r *borrowRepository) Return(recordID uint) error {
	// 1. 获取借阅记录
	var record models.BorrowRecord
	if err := r.db.Preload("Book").
		First(&record, recordID).Error; err != nil {
		return fmt.Errorf("借阅记录不存在")
	}

	// 2. 检查是否已归还
	if record.ReturnedAt != nil {
		return fmt.Errorf("图书已归还")
	}

	// 3. 更新归还时间
	now := time.Now()
	record.ReturnedAt = &now

	if err := r.db.Save(&record).Error; err != nil {
		return fmt.Errorf("更新归还时间失败: %w", err)
	}

	// 4. 增加图书可用数量
	// 修复：检查 Book 的 ID 是否有效
	if record.Book.ID > 0 {
		record.Book.Available += 1
		return r.db.Save(&record.Book).Error
	}

	return nil
}

func (r *borrowRepository) FindActiveByUserAndBook(userID, bookID uint) (*models.BorrowRecord, error) {
	if userID == 0 || bookID == 0 {
		return nil, fmt.Errorf("无效的用户ID或图书ID")
	}

	var record models.BorrowRecord
	err := r.db.Where("user_id = ? AND book_id = ? AND returned_at IS NULL",
		userID, bookID).
		First(&record).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("未找到未归还的借阅记录")
		}
		return nil, fmt.Errorf("查询借阅记录失败: %w", err)
	}

	return &record, nil
}

func (r *borrowRepository) FindActiveByUser(userID uint) ([]models.BorrowRecord, error) {
	if userID == 0 {
		return nil, fmt.Errorf("无效的用户ID")
	}

	var records []models.BorrowRecord
	err := r.db.Preload("Book").
		Where("user_id = ? AND returned_at IS NULL", userID).
		Order("borrowed_at DESC").
		Find(&records).Error

	if err != nil {
		return nil, fmt.Errorf("查询借阅记录失败: %w", err)
	}

	return records, nil
}

func (r *borrowRepository) FindAll() ([]models.BorrowRecord, error) {
	var records []models.BorrowRecord
	err := r.db.Preload("User").Preload("Book").
		Order("borrowed_at DESC").
		Find(&records).Error

	if err != nil {
		return nil, fmt.Errorf("查询借阅记录失败: %w", err)
	}

	return records, nil
}

func (r *borrowRepository) FindByUser(userID uint) ([]models.BorrowRecord, error) {
	if userID == 0 {
		return nil, fmt.Errorf("无效的用户ID")
	}

	var records []models.BorrowRecord
	err := r.db.Preload("Book").
		Where("user_id = ?", userID).
		Order("borrowed_at DESC").
		Find(&records).Error

	if err != nil {
		return nil, fmt.Errorf("查询借阅记录失败: %w", err)
	}

	return records, nil
}
