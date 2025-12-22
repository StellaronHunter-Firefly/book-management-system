package controllers

import (
	"book-management-system/models"
	"book-management-system/services"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type BookController struct {
	bookService services.BookService
}

func NewBookController(bookService services.BookService) *BookController {
	return &BookController{bookService: bookService}
}

type CreateBookRequest struct {
	Title       string `json:"title" binding:"required"`
	Author      string `json:"author" binding:"required"`
	TotalCopies int    `json:"total_copies" binding:"required,min=1"`
}

type UpdateBookRequest struct {
	Title       string `json:"title"`
	Author      string `json:"author"`
	TotalCopies int    `json:"total_copies" binding:"min=1"`
}

type BorrowRequest struct {
	BookID uint `json:"book_id" binding:"required"`
}

func (c *BookController) CreateBook(ctx *gin.Context) {
	var req CreateBookRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	book := &models.Book{
		Title:       req.Title,
		Author:      req.Author,
		TotalCopies: req.TotalCopies,
		Available:   req.TotalCopies,
	}

	if err := c.bookService.CreateBook(book); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, book)
}

func (c *BookController) GetAllBooks(ctx *gin.Context) {
	books, err := c.bookService.GetAllBooks()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, books)
}

func (c *BookController) GetBookByID(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "无效的图书ID"})
		return
	}

	book, err := c.bookService.GetBookByID(uint(id))
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "图书不存在"})
		return
	}

	ctx.JSON(http.StatusOK, book)
}

func (c *BookController) UpdateBook(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "无效的图书ID"})
		return
	}

	var req UpdateBookRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	book := &models.Book{
		Title:       req.Title,
		Author:      req.Author,
		TotalCopies: req.TotalCopies,
	}

	if err := c.bookService.UpdateBook(uint(id), book); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "图书更新成功"})
}

func (c *BookController) DeleteBook(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "无效的图书ID"})
		return
	}

	if err := c.bookService.DeleteBook(uint(id)); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "图书删除成功"})
}

func (c *BookController) SearchBooks(ctx *gin.Context) {
	query := ctx.Query("q")
	if query == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "搜索关键词不能为空"})
		return
	}

	books, err := c.bookService.SearchBooks(query)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, books)
}

func (c *BookController) BorrowBook(ctx *gin.Context) {
	var req BorrowRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, _ := ctx.Get("userID")

	if err := c.bookService.BorrowBook(userID.(uint), req.BookID); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "借书成功"})
}

func (c *BookController) ReturnBook(ctx *gin.Context) {
	var req BorrowRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, _ := ctx.Get("userID")

	if err := c.bookService.ReturnBook(userID.(uint), req.BookID); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "还书成功"})
}

func (c *BookController) GetMyBorrowedBooks(ctx *gin.Context) {
	userID, _ := ctx.Get("userID")

	books, err := c.bookService.GetBorrowedBooks(userID.(uint))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, books)
}

func (c *BookController) GetMyBorrowRecords(ctx *gin.Context) {
	userID, _ := ctx.Get("userID")

	records, err := c.bookService.GetBorrowRecords(userID.(uint))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, records)
}

func (c *BookController) GetAllBorrowRecords(ctx *gin.Context) {
	records, err := c.bookService.GetAllBorrowRecords()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, records)
}

func (c *BookController) CheckAvailability(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "无效的图书ID"})
		return
	}

	available, err := c.bookService.CheckBookAvailability(uint(id))
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "图书不存在"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"available": available})
}
