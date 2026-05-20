package controllers

import (
	"io"
	"net/http"

	"github.com/decntrafir/backend/internal/middleware"
	"github.com/decntrafir/backend/internal/services"
	"github.com/decntrafir/backend/pkg/jwtutil"
	"github.com/gin-gonic/gin"
)

type EvidenceController struct {
	evidence *services.EvidenceService
}

func NewEvidenceController(evidence *services.EvidenceService) *EvidenceController {
	return &EvidenceController{evidence: evidence}
}

func (c *EvidenceController) Upload(ctx *gin.Context) {
	claims := ctx.MustGet(middleware.ContextUserKey).(*jwtutil.Claims)
	firID := ctx.PostForm("firId")
	if firID == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "firId required"})
		return
	}
	file, header, err := ctx.Request.FormFile("file")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "file required"})
		return
	}
	defer file.Close()
	data, err := io.ReadAll(file)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ev, err := c.evidence.Upload(ctx.Request.Context(), firID, claims.UserID, header.Filename, data, ctx.ClientIP())
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusCreated, ev)
}
