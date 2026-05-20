package controllers

import (
	"net/http"
	"strconv"

	"github.com/decntrafir/backend/internal/middleware"
	"github.com/decntrafir/backend/internal/models"
	"github.com/decntrafir/backend/internal/services"
	"github.com/decntrafir/backend/pkg/jwtutil"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type AdminController struct {
	users         *services.UserService
	audit         *services.AuditService
	notifications *services.NotificationService
}

func NewAdminController(users *services.UserService, audit *services.AuditService, notifications *services.NotificationService) *AdminController {
	return &AdminController{users: users, audit: audit, notifications: notifications}
}

func (c *AdminController) ListUsers(ctx *gin.Context) {
	page, _ := strconv.ParseInt(ctx.DefaultQuery("page", "1"), 10, 64)
	limit, _ := strconv.ParseInt(ctx.DefaultQuery("limit", "20"), 10, 64)
	users, total, err := c.users.List(ctx.Request.Context(), ctx.Query("role"), page, limit)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"data": users, "total": total})
}

func (c *AdminController) CreateUser(ctx *gin.Context) {
	claims := ctx.MustGet(middleware.ContextUserKey).(*jwtutil.Claims)
	var body struct {
		Name        string      `json:"name" binding:"required"`
		Email       string      `json:"email" binding:"required,email"`
		Password    string      `json:"password" binding:"required,min=8"`
		Role        models.Role `json:"role" binding:"required"`
		Phone       string      `json:"phone"`
		Department  string      `json:"department"`
		BadgeNumber string      `json:"badgeNumber"`
	}
	if err := ctx.ShouldBindJSON(&body); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	user := models.User{
		Name: body.Name, Email: body.Email, Role: body.Role,
		Phone: body.Phone, Department: body.Department, BadgeNumber: body.BadgeNumber,
	}
	created, err := c.users.CreateByAdmin(ctx.Request.Context(), claims.UserID, user, body.Password, ctx.ClientIP())
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusCreated, created)
}

func (c *AdminController) DeleteUser(ctx *gin.Context) {
	claims := ctx.MustGet(middleware.ContextUserKey).(*jwtutil.Claims)
	id, err := primitive.ObjectIDFromHex(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	if err := c.users.Delete(ctx.Request.Context(), claims.UserID, id, ctx.ClientIP()); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "deleted"})
}

func (c *AdminController) Notifications(ctx *gin.Context) {
	claims := ctx.MustGet(middleware.ContextUserKey).(*jwtutil.Claims)
	items, err := c.notifications.List(ctx.Request.Context(), claims.UserID, 50)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, items)
}
