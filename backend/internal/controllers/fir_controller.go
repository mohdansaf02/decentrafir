package controllers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/decntrafir/backend/internal/middleware"
	"github.com/decntrafir/backend/internal/models"
	"github.com/decntrafir/backend/internal/services"
	"github.com/decntrafir/backend/pkg/jwtutil"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type FIRController struct {
	fir *services.FIRService
}

func NewFIRController(fir *services.FIRService) *FIRController {
	return &FIRController{fir: fir}
}

func (c *FIRController) Create(ctx *gin.Context) {
	claims := ctx.MustGet(middleware.ContextUserKey).(*jwtutil.Claims)
	var req models.CreateFIRRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	fir, err := c.fir.Create(ctx.Request.Context(), claims.UserID, req, ctx.ClientIP())
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusCreated, fir)
}

func (c *FIRController) Get(ctx *gin.Context) {
	fir, err := c.fir.GetByID(ctx.Request.Context(), ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, fir)
}

func (c *FIRController) Update(ctx *gin.Context) {
	claims := ctx.MustGet(middleware.ContextUserKey).(*jwtutil.Claims)
	var req models.UpdateFIRRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	fir, err := c.fir.Update(ctx.Request.Context(), ctx.Param("id"), claims.UserID, models.Role(claims.Role), req, ctx.ClientIP())
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, fir)
}

func (c *FIRController) List(ctx *gin.Context) {
	claims := ctx.MustGet(middleware.ContextUserKey).(*jwtutil.Claims)
	q := parseFIRQuery(ctx)
	var citizenFilter *primitive.ObjectID
	if claims.Role == string(models.RoleCitizen) {
		citizenFilter = &claims.UserID
	}
	firs, total, err := c.fir.List(ctx.Request.Context(), q, citizenFilter)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"data": firs, "total": total, "page": q.Page, "limit": q.Limit})
}

func (c *FIRController) Analytics(ctx *gin.Context) {
	data, err := c.fir.Analytics(ctx.Request.Context())
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, data)
}

func parseFIRQuery(ctx *gin.Context) services.FIRListQuery {
	page, _ := strconv.ParseInt(ctx.DefaultQuery("page", "1"), 10, 64)
	limit, _ := strconv.ParseInt(ctx.DefaultQuery("limit", "20"), 10, 64)
	q := services.FIRListQuery{
		Status: ctx.Query("status"), CrimeType: ctx.Query("crimeType"),
		FIRID: ctx.Query("firId"), OfficerID: ctx.Query("officerId"),
		Page: page, Limit: limit,
	}
	if from := ctx.Query("fromDate"); from != "" {
		if t, err := time.Parse(time.RFC3339, from); err == nil {
			q.FromDate = &t
		}
	}
	if to := ctx.Query("toDate"); to != "" {
		if t, err := time.Parse(time.RFC3339, to); err == nil {
			q.ToDate = &t
		}
	}
	return q
}
