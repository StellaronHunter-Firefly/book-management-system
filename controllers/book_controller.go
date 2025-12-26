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

// CreateBookRequest 创建图书请求
type CreateBookRequest struct {
	Title       string `json:"title" binding:"required"`
	Author      string `json:"author" binding:"required"`
	TotalCopies int    `json:"total_copies" binding:"required,min=1"`
}

// UpdateBookRequest 更新图书请求
type UpdateBookRequest struct {
	Title       string `json:"title"`
	Author      string `json:"author"`
	TotalCopies int    `json:"total_copies" binding:"min=1"`
}

// BorrowRequest 借书请求
type BorrowRequest struct {
	BookID uint `json:"book_id" binding:"required"`
}

// DeleteBookRequest 删除图书请求
type DeleteBookRequest struct {
	Confirm bool `json:"confirm" example:"true"`
}

// DeleteBookResponse 删除图书响应
type DeleteBookResponse struct {
	Message string       `json:"message" example:"图书删除成功"`
	Book    *models.Book `json:"deleted_book"`
}

// CreateBook godoc
// @Summary      创建图书
// @Description  管理员创建新图书
// @Tags         图书管理
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        request  body  CreateBookRequest  true  "图书信息"
// @Success      201      {object}  models.Book
// @Failure      400      {object}  ErrorResponse
// @Failure      401      {object}  ErrorResponse
// @Failure      403      {object}  ErrorResponse
// @Router       /admin/books [post]
func (c *BookController) CreateBook(ctx *gin.Context) {
	var req CreateBookRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error1": err.Error()})
		return
	}

	book := &models.Book{
		Title:       req.Title,
		Author:      req.Author,
		TotalCopies: req.TotalCopies,
		Available:   req.TotalCopies,
	}

	if err := c.bookService.CreateBook(book); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error2": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, book)
}

// GetAllBooks godoc
// @Summary      获取所有图书
// @Description  获取所有图书列表
// @Tags         图书
// @Accept       json
// @Produce      json
// @Success      200  {array}  models.Book
// @Failure      500  {object}  ErrorResponse
// @Router       /books [get]
func (c *BookController) GetAllBooks(ctx *gin.Context) {
	books, err := c.bookService.GetAllBooks()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, books)
}

// GetBookByID godoc
// @Summary      获取图书详情
// @Description  根据ID获取图书详情
// @Tags         图书
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "图书ID"
// @Success      200  {object}  models.Book
// @Failure      400  {object}  ErrorResponse
// @Failure      404  {object}  ErrorResponse
// @Router       /books/{id} [get]
func (c *BookController) GetBookByID(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "无效的图书ID"})
		return // 添加return语句
	}

	book, err := c.bookService.GetBookByID(uint(id))
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "图书不存在"}) // 改为404
		return
	}

	ctx.JSON(http.StatusOK, book)
}

// UpdateBook godoc
// @Summary      更新图书
// @Description  管理员更新图书信息
// @Tags         图书管理
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id   path      int               true  "图书ID"
// @Param        request  body  UpdateBookRequest  true  "图书信息"
// @Success      200  {object}  SuccessResponse
// @Failure      400  {object}  ErrorResponse
// @Failure      401  {object}  ErrorResponse
// @Failure      403  {object}  ErrorResponse
// @Router       /admin/books/{id} [put]
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

	// 更新成功后，重新获取图书信息
	updatedBook, err := c.bookService.GetBookByID(uint(id))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "获取更新后的图书信息失败"})
		return
	}

	// 返回更新后的图书信息
	ctx.JSON(http.StatusOK, updatedBook)
}

// DeleteBook godoc
// @Summary      删除图书
// @Description  管理员删除图书
// @Tags         图书管理
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id   path      int  true  "图书ID"
// @Success      200  {object}  DeleteBookResponse
// @Failure      400  {object}  ErrorResponse
// @Failure      401  {object}  ErrorResponse
// @Failure      403  {object}  ErrorResponse
// @Router       /admin/books/{id} [delete]
func (c *BookController) DeleteBook(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "无效的图书ID"}) // JSON大写
		return
	}

	// 先获取图书信息
	book, err := c.bookService.GetBookByID(uint(id))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "图书不存在"}) // JSON大写
		return
	}

	// 删除图书
	if err := c.bookService.DeleteBook(uint(id)); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()}) // JSON大写
		return
	}

	// 返回删除成功的响应，包含被删除的图书信息
	ctx.JSON(http.StatusOK, DeleteBookResponse{ // JSON大写
		Message: "图书删除成功",
		Book:    book,
	})
}

// SearchBooks godoc
// @Summary      搜索图书
// @Description  根据关键词搜索图书
// @Tags         图书
// @Accept       json
// @Produce      json
// @Param        q    query     string  true  "搜索关键词"
// @Success      200  {array}   models.Book
// @Failure      400  {object}  ErrorResponse
// @Failure      500  {object}  ErrorResponse
// @Router       /books/search [get]
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

// BorrowBook godoc
// @Summary      借书
// @Description  借阅图书
// @Tags         借阅
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        request  body  BorrowRequest  true  "借书信息"
// @Success      200  {object}  SuccessResponse
// @Failure      400  {object}  ErrorResponse
// @Failure      401  {object}  ErrorResponse
// @Router       /books/borrow [post]
func (c *BookController) BorrowBook(ctx *gin.Context) {
	var req BorrowRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, _ := ctx.Get("userID")

	// 获取图书信息（用于返回书名）
	book, err := c.bookService.GetBookByID(req.BookID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "图书不存在"})
		return
	}

	// 执行借书操作
	if err := c.bookService.BorrowBook(userID.(uint), req.BookID); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message":    "借书成功",
		"book_title": book.Title,
		"book_id":    book.ID,
		"author":     book.Author,
		"available":  book.Available - 1, // 借阅后的可用库存
	})
}

// ReturnBook godoc
// @Summary      还书
// @Description  归还图书
// @Tags         借阅
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        request  body  BorrowRequest  true  "还书信息"
// @Success      200  {object}  SuccessResponse
// @Failure      400  {object}  ErrorResponse
// @Failure      401  {object}  ErrorResponse
// @Router       /books/return [post]
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

// GetMyBorrowedBooks godoc
// @Summary      获取已借图书
// @Description  获取当前用户已借的图书列表
// @Tags         借阅
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Success      200  {array}  models.Book
// @Failure      401  {object}  ErrorResponse
// @Failure      500  {object}  ErrorResponse
// @Router       /books/my-borrowed [get]
func (c *BookController) GetMyBorrowedBooks(ctx *gin.Context) {
	userID, _ := ctx.Get("userID")

	books, err := c.bookService.GetBorrowedBooks(userID.(uint))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, books)
}

// GetMyBorrowRecords godoc
// @Summary      获取借阅记录
// @Description  获取当前用户的借阅记录
// @Tags         借阅
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Success      200  {array}  models.BorrowRecord
// @Failure      401  {object}  ErrorResponse
// @Failure      500  {object}  ErrorResponse
// @Router       /books/my-records [get]
func (c *BookController) GetMyBorrowRecords(ctx *gin.Context) {
	userID, _ := ctx.Get("userID")

	records, err := c.bookService.GetBorrowRecords(userID.(uint))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, records)
}

// GetAllBorrowRecords godoc
// @Summary      获取所有借阅记录
// @Description  管理员获取所有用户的借阅记录
// @Tags         借阅管理
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Success      200  {array}  models.BorrowRecord
// @Failure      401  {object}  ErrorResponse
// @Failure      403  {object}  ErrorResponse
// @Failure      500  {object}  ErrorResponse
// @Router       /admin/borrow-records [get]
func (c *BookController) GetAllBorrowRecords(ctx *gin.Context) {
	records, err := c.bookService.GetAllBorrowRecords()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, records)
}

// CheckAvailability godoc
// @Summary      检查图书可用性
// @Description  检查图书是否可借
// @Tags         图书
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "图书ID"
// @Success      200  {object}  map[string]bool
// @Failure      400  {object}  ErrorResponse
// @Failure      404  {object}  ErrorResponse
// @Router       /books/{id}/availability [get]
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
