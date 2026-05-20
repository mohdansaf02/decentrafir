package controllers

import (
	"net/http"

	"github.com/decntrafir/backend/internal/middleware"
	"github.com/decntrafir/backend/internal/services"
	"github.com/gin-gonic/gin"
)

type AuthController struct {
	auth *services.AuthService
}

func NewAuthController(auth *services.AuthService) *AuthController {
	return &AuthController{auth: auth}
}

// Register godoc
// @Summary Register citizen
// @Tags auth
// @Accept json
// @Produce json
// @Param body body services.RegisterInput true "Register"
// @Success 201 {object} services.AuthTokens
// @Router /register [post]
func (c *AuthController) Register(ctx *gin.Context) {
	var in services.RegisterInput
	if err := ctx.ShouldBindJSON(&in); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	tokens, err := c.auth.Register(ctx.Request.Context(), in, ctx.ClientIP())
	if err != nil {
		status := http.StatusInternalServerError
		if err == services.ErrUserExists {
			status = http.StatusConflict
		}
		ctx.JSON(status, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusCreated, tokens)
}

// Login godoc
// @Summary Login
// @Tags auth
// @Accept json
// @Produce json
// @Router /login [post]
func (c *AuthController) Login(ctx *gin.Context) {
	var body struct {
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required"`
	}
	if err := ctx.ShouldBindJSON(&body); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	tokens, err := c.auth.Login(ctx.Request.Context(), body.Email, body.Password, ctx.ClientIP())
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}
	ctx.JSON(http.StatusOK, tokens)
}

func (c *AuthController) Refresh(ctx *gin.Context) {
	var body struct {
		RefreshToken string `json:"refreshToken" binding:"required"`
	}
	if err := ctx.ShouldBindJSON(&body); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	tokens, err := c.auth.Refresh(ctx.Request.Context(), body.RefreshToken)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, tokens)
}

func (c *AuthController) Me(ctx *gin.Context) {
	claims, _ := ctx.Get(middleware.ContextUserKey)
	ctx.JSON(http.StatusOK, claims)
}
