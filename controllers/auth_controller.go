package controllers

import (
	"book-management-system/models"
	"book-management-system/services"
	"book-management-system/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AuthController struct {
	authService services.AuthService
}

func NewAuthController(authService services.AuthService) *AuthController {
	return &AuthController{authService: authService}
}

// RegisterRequest 注册请求体
type RegisterRequest struct {
	Username string `json:"username" binding:"required,min=3,max=50"`
	Password string `json:"password" binding:"required,min=6"`
	Email    string `json:"email" binding:"required,email"`
}

// LoginRequest 登录请求体
type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// AuthResponse 认证响应
type AuthResponse struct {
	User  *models.User `json:"user"`
	Token string       `json:"token"`
}

// SuccessResponse 通用成功响应
type SuccessResponse struct {
	Message string `json:"message" example:"操作成功"`
	Data    any    `json:"data,omitempty"`
}

// ErrorResponse 错误响应
type ErrorResponse struct {
	Error string `json:"error" example:"错误信息"`
}

type ChangeUsernameRequest struct {
	NewUsername string `json:"new_username" binding:"required,min=3,max=50" example:"newuser123"`
}

type ChangePasswordRequest struct {
	OldPassword string `json:"old_password" binding:"required" example:"oldpassword123"`
	NewPassword string `json:"new_password" binding:"required,min=6" example:"newpassword456"`
}

type UserInfo struct {
	ID        uint            `json:"id"`
	Username  string          `json:"username"`
	Email     string          `json:"email"`
	Role      models.UserRole `json:"role"`
	CreatedAt string          `json:"created_at"`
}

// Register godoc
// @Summary      用户注册
// @Description  新用户注册
// @Tags         认证
// @Accept       json
// @Produce      json
// @Param        request  body  RegisterRequest  true  "注册信息"
// @Success      201      {object}  AuthResponse
// @Failure      400      {object}  ErrorResponse
// @Failure      500      {object}  ErrorResponse
// @Router       /auth/register [post]
func (c *AuthController) Register(ctx *gin.Context) {
	var req RegisterRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := c.authService.Register(req.Username, req.Password, req.Email)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	token, err := utils.GenerateToken(user)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "生成令牌失败"})
		return
	}

	response := AuthResponse{
		User:  user,
		Token: token,
	}

	ctx.JSON(http.StatusCreated, response)
}

// Login godoc
// @Summary      用户登录
// @Description  用户登录获取JWT令牌
// @Tags         认证
// @Accept       json
// @Produce      json
// @Param        request  body  LoginRequest  true  "登录凭证"
// @Success      200      {object}  AuthResponse
// @Failure      400      {object}  ErrorResponse
// @Failure      401      {object}  ErrorResponse
// @Router       /auth/login [post]
func (c *AuthController) Login(ctx *gin.Context) {
	var req LoginRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, token, err := c.authService.Login(req.Username, req.Password)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	response := AuthResponse{
		User:  user,
		Token: token,
	}

	ctx.JSON(http.StatusOK, response)
}

// GetProfile godoc
// @Summary      获取用户资料
// @Description  获取当前登录用户的资料
// @Tags         用户
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Success      200  {object}  models.User
// @Failure      401  {object}  ErrorResponse
// @Failure      404  {object}  ErrorResponse
// @Router       /users/profile [get]
func (c *AuthController) GetProfile(ctx *gin.Context) {
	userID, exists := ctx.Get("userID")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "未认证"})
		return
	}

	user, err := c.authService.GetUserByID(userID.(uint))
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "用户不存在"})
		return
	}

	ctx.JSON(http.StatusOK, user)
}

// ChangeUsername godoc
// @Summary      更改用户名
// @Description  用户更改自己的用户名
// @Tags         用户
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        request  body  ChangeUsernameRequest  true  "用户名更改信息"
// @Success      200  {object}  SuccessResponse
// @Failure      400  {object}  ErrorResponse
// @Failure      401  {object}  ErrorResponse
// @Router       /users/username [put]
func (c *AuthController) ChangeUsername(ctx *gin.Context) {
	// 获取用户ID
	userID, exists := ctx.Get("userID")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "未认证"})
		return
	}

	var req ChangeUsernameRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := c.authService.ChangeUsername(userID.(uint), req.NewUsername); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message":      "用户名修改成功",
		"new_username": req.NewUsername,
	})
}

// ChangePassword godoc
// @Summary      更改密码
// @Description  用户更改自己的密码（需要输入原密码）
// @Tags         用户
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        request  body  ChangePasswordRequest  true  "密码更改信息"
// @Success      200  {object}  SuccessResponse
// @Failure      400  {object}  ErrorResponse
// @Failure      401  {object}  ErrorResponse
// @Router       /users/password [put]
func (c *AuthController) ChangePassword(ctx *gin.Context) {
	// 获取用户ID
	userID, exists := ctx.Get("userID")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "未认证"})
		return
	}

	var req ChangePasswordRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := c.authService.ChangePassword(userID.(uint), req.OldPassword, req.NewPassword); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "密码修改成功"})
}

// GetAllUsers godoc
// @Summary      获取所有用户列表
// @Description  管理员获取所有用户列表
// @Tags         用户管理
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Success      200  {array}   UserInfo
// @Failure      401  {object}  ErrorResponse
// @Failure      403  {object}  ErrorResponse
// @Failure      500  {object}  ErrorResponse
// @Router       /admin/users [get]
func (c *AuthController) GetAllUsers(ctx *gin.Context) {
	users, err := c.authService.GetAllUsers()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 转换为UserInfo格式
	var userInfos []UserInfo
	for _, user := range users {
		userInfos = append(userInfos, UserInfo{
			ID:        user.ID,
			Username:  user.Username,
			Email:     user.Email,
			Role:      user.Role,
			CreatedAt: user.CreatedAt.Format("2006-01-02 15:04:05"),
		})
	}

	ctx.JSON(http.StatusOK, userInfos)
}

// Logout godoc
// @Summary      用户退出登录
// @Description  用户退出登录，前端需要删除本地存储的token
// @Tags         认证
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Success      200  {object}  SuccessResponse
// @Failure      401  {object}  ErrorResponse
// @Router       /auth/logout [post]
func (c *AuthController) Logout(ctx *gin.Context) {
	// JWT是无状态的，退出只需要前端删除token
	// 这里只需返回成功信息
	ctx.JSON(http.StatusOK, gin.H{
		"message": "退出登录成功",
	})
}
