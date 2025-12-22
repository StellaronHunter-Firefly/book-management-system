package services

import (
	"book-management-system/models"
	"book-management-system/repositories"
	"errors"
)

type BookService interface {
	CreateBook(book *models.Book) error
	GetBookByID(id uint) (*models.Book, error)
	GetAllBooks() ([]models.Book, error)
	UpdateBook(id uint, book *models.Book) error
	DeleteBook(id uint) error
	SearchBooks(query string) ([]models.Book, error)
	BorrowBook(userID, bookID uint) error
	ReturnBook(userID, bookID uint) error
	GetBorrowedBooks(userID uint) ([]models.Book, error)
	GetBorrowRecords(userID uint) ([]models.BorrowRecord, error)
	GetAllBorrowRecords() ([]models.BorrowRecord, error)
	CheckBookAvailability(bookID uint) (bool, error)
}

type bookService struct {
	bookRepo repositories.BookRepositoryWithBorrow
}

func NewBookService(bookRepo repositories.BookRepositoryWithBorrow) BookService {
	return &bookService{bookRepo: bookRepo}
}

func (s *bookService) CreateBook(book *models.Book) error {
	if book.Title == "" {
		return errors.New("书名不能为空")
	}
	if book.Author == "" {
		return errors.New("作者不能为空")
	}
	if book.TotalCopies <= 0 {
		return errors.New("库存数量必须大于0")
	}

	return s.bookRepo.Create(book)
}

func (s *bookService) GetBookByID(id uint) (*models.Book, error) {
	return s.bookRepo.FindByID(id)
}

func (s *bookService) GetAllBooks() ([]models.Book, error) {
	return s.bookRepo.FindAll()
}

func (s *bookService) UpdateBook(id uint, book *models.Book) error {
	existing, err := s.bookRepo.FindByID(id)
	if err != nil {
		return err
	}

	existing.Title = book.Title
	existing.Author = book.Author

	if existing.TotalCopies != book.TotalCopies {
		diff := book.TotalCopies - existing.TotalCopies
		existing.Available += diff
		if existing.Available < 0 {
			existing.Available = 0
		}
		existing.TotalCopies = book.TotalCopies
	}

	return s.bookRepo.Update(existing)
}

func (s *bookService) DeleteBook(id uint) error {
	book, err := s.bookRepo.FindByID(id)
	if err != nil {
		return err
	}

	if book.Available != book.TotalCopies {
		return errors.New("无法删除正在借阅的书籍")
	}

	return s.bookRepo.Delete(id)
}

func (s *bookService) SearchBooks(query string) ([]models.Book, error) {
	return s.bookRepo.Search(query)
}

func (s *bookService) BorrowBook(userID, bookID uint) error {
	book, err := s.bookRepo.FindByID(bookID)
	if err != nil {
		return errors.New("图书不存在")
	}

	if book.Available <= 0 {
		return errors.New("图书已全部借出")
	}

	return s.bookRepo.BorrowBook(userID, bookID)
}

func (s *bookService) ReturnBook(userID, bookID uint) error {
	return s.bookRepo.ReturnBook(userID, bookID)
}

func (s *bookService) GetBorrowedBooks(userID uint) ([]models.Book, error) {
	return s.bookRepo.GetBorrowedBooks(userID)
}

func (s *bookService) GetBorrowRecords(userID uint) ([]models.BorrowRecord, error) {
	return s.bookRepo.GetBorrowRecords(userID)
}

func (s *bookService) GetAllBorrowRecords() ([]models.BorrowRecord, error) {
	return s.bookRepo.GetAllBorrowRecords()
}

func (s *bookService) CheckBookAvailability(bookID uint) (bool, error) {
	return s.bookRepo.CheckAvailability(bookID)
}
