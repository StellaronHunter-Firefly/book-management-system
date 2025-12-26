package repositories

import (
	"book-management-system/config"
	"book-management-system/models"
	"fmt"

	"gorm.io/gorm"
)

type BookRepository interface {
	Create(book *models.Book) error
	FindByID(id uint) (*models.Book, error)
	Update(book *models.Book) error
	Delete(id uint) error
	FindAll() ([]models.Book, error)
	FindAvailable() ([]models.Book, error)
	Search(query string) ([]models.Book, error)
	CheckAvailability(bookID uint) (bool, error)
	ExistsByTitleAndAuthor(title, author string) (bool, error)
}

type bookRepository struct {
	db *gorm.DB
}

func NewBookRepository() BookRepository {
	return &bookRepository{db: config.DB}
}

func (r *bookRepository) Create(book *models.Book) error {
	if book.Title == "" {
		return fmt.Errorf("书名不能为空")
	}
	if book.Author == "" {
		return fmt.Errorf("作者不能为空")
	}
	if book.TotalCopies <= 0 {
		return fmt.Errorf("总库存必须大于0")
	}

	if book.Available == 0 {
		book.Available = book.TotalCopies
	}

	return r.db.Create(book).Error
}

func (r *bookRepository) FindByID(id uint) (*models.Book, error) {
	if id == 0 {
		return nil, fmt.Errorf("无效的图书ID")
	}

	var book models.Book
	err := r.db.First(&book, id).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("图书不存在")
		}
		return nil, fmt.Errorf("查询图书失败: %w", err) // 修正格式字符串
	}

	return &book, nil
}

func (r *bookRepository) Update(book *models.Book) error {
	var existing models.Book
	if err := r.db.First(&existing, book.ID).Error; err != nil {
		return fmt.Errorf("图书不存在")
	}

	if book.TotalCopies != existing.TotalCopies {
		diff := book.TotalCopies - existing.TotalCopies
		book.Available = existing.Available + diff
		if book.Available < 0 {
			book.Available = 0
		}
	} else {
		book.Available = existing.Available
	}

	return r.db.Save(book).Error
}

func (r *bookRepository) Delete(id uint) error {
	var book models.Book
	if err := r.db.First(&book, id).Error; err != nil {
		return fmt.Errorf("图书不存在")
	}

	if book.Available != book.TotalCopies {
		return fmt.Errorf("图书有未归还记录，无法删除")
	}

	return r.db.Delete(&models.Book{}, id).Error
}

func (r *bookRepository) FindAll() ([]models.Book, error) {
	var books []models.Book

	if err := r.db.Model(&models.Book{}).
		Select("id, title, author, total_copies, available, created_at, updated_at").
		Order("created_at DESC").
		Find(&books).Error; err != nil {
		return nil, fmt.Errorf("查询图书列表失败: %w", err)
	}

	return books, nil
}

func (r *bookRepository) FindAvailable() ([]models.Book, error) {
	var books []models.Book

	if err := r.db.Model(&models.Book{}).
		Select("id, title, author, total_copies, available, created_at").
		Where("available > 0").
		Order("created_at DESC").
		Find(&books).Error; err != nil {
		return nil, fmt.Errorf("查询可借图书失败: %w", err)
	}

	return books, nil
}

func (r *bookRepository) Search(query string) ([]models.Book, error) {
	var books []models.Book

	if query == "" {
		return r.FindAll()
	}
	searchPattern := "%" + query + "%"

	searchQuery := r.db.Model(&models.Book{}).
		Where("title LIKE ? OR author LIKE ?",
			searchPattern, searchPattern) // 传入3个参数

	if err := searchQuery.
		Select("id, title, author, total_copies, available, created_at").
		Order("created_at DESC").
		Find(&books).Error; err != nil {
		return nil, fmt.Errorf("搜索图书失败: %w", err)
	}

	return books, nil
}

func (r *bookRepository) CheckAvailability(bookID uint) (bool, error) {
	var book models.Book
	if err := r.db.Select("available").First(&book, bookID).Error; err != nil {
		return false, fmt.Errorf("查询图书失败: %w", err)
	}

	return book.Available > 0, nil
}

func (r *bookRepository) ExistsByTitleAndAuthor(title, author string) (bool, error) {
	if title == "" || author == "" {
		return false, fmt.Errorf("书名和作者不能为空")
	}

	var count int64
	err := r.db.Model(&models.Book{}).
		Where("title = ? AND author = ?", title, author).
		Count(&count).Error

	if err != nil {
		return false, fmt.Errorf("查询图书存在性失败: %w", err)
	}

	return count > 0, nil
}
