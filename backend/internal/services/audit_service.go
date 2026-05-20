package services

import (
	"context"
	"time"

	"github.com/decntrafir/backend/internal/database"
	"github.com/decntrafir/backend/internal/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type AuditService struct {
	db *database.DB
}

func NewAuditService(db *database.DB) *AuditService {
	return &AuditService{db: db}
}

func (s *AuditService) Log(ctx context.Context, userID primitive.ObjectID, action, resource, resourceID, ip, txHash string, details map[string]interface{}) error {
	log := models.AuditLog{
		UserID:     userID,
		Action:     action,
		Resource:   resource,
		ResourceID: resourceID,
		Details:    details,
		IPAddress:  ip,
		TxHash:     txHash,
		CreatedAt:  time.Now(),
	}
	_, err := s.db.Collection("audit_logs").InsertOne(ctx, log)
	return err
}
