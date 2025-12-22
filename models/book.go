package models

import (
	"time"

	"gorm.io/gorm"
)

type Book struct {
	ID          uint           `gorm:"primarykey" json:"id"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
	Title       string         `gorm:"not null;index" json:"title"`
	Author      string         `gorm:"not null" json:"author"`
	TotalCopies int            `gorm:"not null;default:1" json:"total_copies"`
	Available   int            `gorm:"not null" json:"available"`
	BorrowedBy  []User         `gorm:"many2many:user_borrowed_books;" json:"-"`
}

type BorrowRecord struct {
	ID         uint       `gorm:"primarykey" json:"id"`
	UserID     uint       `gorm:"not null;index" json:"user_id"`
	BookID     uint       `gorm:"not null;index" json:"book_id"`
	BorrowedAt time.Time  `gorm:"not null" json:"borrowed_at"`
	ReturnedAt *time.Time `json:"returned_at"`
	DueDate    time.Time  `gorm:"not null" json:""`
	Book       Book       `gorm:"foreignKey:BookID" json:"book"`
	User       User       `gorm:"foreignKey:UserID" json:"user"`
}
