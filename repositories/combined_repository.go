package repositories

import (
	"book-management-system/models"
	"fmt"
	"time"
)

type BookRepositoryWithBorrow interface {
	BookRepository
	BorrowBook(userID, bookID uint) error
	ReturnBook(userID, bookID uint) error
	GetBorrowedBooks(userID uint) ([]models.Book, error)
	GetBorrowRecords(userID uint) ([]models.BorrowRecord, error)
	GetAllBorrowRecords() ([]models.BorrowRecord, error)
	GetActiveBorrowRecord(userID, bookID uint) (*models.BorrowRecord, error)
	ExistsByTitleAndAuthor(title, author string) (bool, error)
}

type combinedBookRepository struct {
	bookRepo   BookRepository
	borrowRepo BorrowRepository
}

func NewCombinedBookRepository() BookRepositoryWithBorrow {
	return &combinedBookRepository{
		bookRepo:   NewBookRepository(),
		borrowRepo: NewBorrowRepository(),
	}
}

func (r *combinedBookRepository) Create(book *models.Book) error {
	return r.bookRepo.Create(book)
}

func (r *combinedBookRepository) FindByID(id uint) (*models.Book, error) {
	return r.bookRepo.FindByID(id)
}

func (r *combinedBookRepository) Update(book *models.Book) error {
	return r.bookRepo.Update(book)
}

func (r *combinedBookRepository) Delete(id uint) error {
	return r.bookRepo.Delete(id)
}

func (r *combinedBookRepository) FindAll() ([]models.Book, error) {
	return r.bookRepo.FindAll()
}

func (r *combinedBookRepository) FindAvailable() ([]models.Book, error) {
	return r.bookRepo.FindAvailable()
}

func (r *combinedBookRepository) Search(query string) ([]models.Book, error) {
	return r.bookRepo.Search(query)
}

func (r *combinedBookRepository) CheckAvailability(bookID uint) (bool, error) {
	return r.bookRepo.CheckAvailability(bookID)
}

func (r *combinedBookRepository) BorrowBook(userID, bookID uint) error {
	record := &models.BorrowRecord{
		UserID:     userID,
		BookID:     bookID,
		BorrowedAt: time.Now(),
		DueDate:    time.Now().Add(14 * 24 * time.Hour),
	}
	return r.borrowRepo.Borrow(record)
}

func (r *combinedBookRepository) ReturnBook(userID, bookID uint) error {
	record, err := r.borrowRepo.FindActiveByUserAndBook(userID, bookID)
	if err != nil {
		return fmt.Errorf("未找到借阅记录: %w", err)
	}

	return r.borrowRepo.Return(record.ID)
}

func (r *combinedBookRepository) GetBorrowedBooks(userID uint) ([]models.Book, error) {
	records, err := r.borrowRepo.FindActiveByUser(userID)
	if err != nil {
		return nil, err
	}

	var books []models.Book
	for _, record := range records {
		book, err := r.bookRepo.FindByID(record.BookID)
		if err == nil {
			books = append(books, *book)
		}
	}

	return books, nil
}

func (r *combinedBookRepository) GetBorrowRecords(userID uint) ([]models.BorrowRecord, error) {
	return r.borrowRepo.FindByUser(userID)
}

func (r *combinedBookRepository) GetAllBorrowRecords() ([]models.BorrowRecord, error) {
	return r.borrowRepo.FindAll()
}

func (r *combinedBookRepository) GetActiveBorrowRecord(userID, bookID uint) (*models.BorrowRecord, error) {
	return r.borrowRepo.FindActiveByUserAndBook(userID, bookID)
}

func (r *combinedBookRepository) ExistsByTitleAndAuthor(title, author string) (bool, error) {
	return r.bookRepo.ExistsByTitleAndAuthor(title, author)
}
