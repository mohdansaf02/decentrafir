package services

import (
	"context"
	"fmt"
	"net/smtp"
	"time"

	"github.com/decntrafir/backend/internal/config"
	"github.com/decntrafir/backend/internal/database"
	"github.com/decntrafir/backend/internal/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type NotificationService struct {
	db  *database.DB
	cfg *config.Config
}

func NewNotificationService(db *database.DB, cfg *config.Config) *NotificationService {
	return &NotificationService{db: db, cfg: cfg}
}

func (s *NotificationService) Create(ctx context.Context, userID primitive.ObjectID, title, message, nType, firID string) error {
	n := models.Notification{
		UserID:    userID,
		Title:     title,
		Message:   message,
		Type:      nType,
		Read:      false,
		FIRID:     firID,
		CreatedAt: time.Now(),
	}
	_, err := s.db.Collection("notifications").InsertOne(ctx, n)
	return err
}

func (s *NotificationService) List(ctx context.Context, userID primitive.ObjectID, limit int64) ([]models.Notification, error) {
	opts := options.Find().SetSort(bson.D{{Key: "createdAt", Value: -1}}).SetLimit(limit)
	cursor, err := s.db.Collection("notifications").Find(ctx, bson.M{"userId": userID}, opts)
	if err != nil {
		return nil, err
	}
	var items []models.Notification
	if err := cursor.All(ctx, &items); err != nil {
		return nil, err
	}
	return items, nil
}

func (s *NotificationService) SendEmail(to, subject, body string) error {
	if s.cfg.SMTPHost == "" || s.cfg.SMTPUser == "" {
		return nil
	}
	addr := fmt.Sprintf("%s:%s", s.cfg.SMTPHost, s.cfg.SMTPPort)
	msg := []byte(fmt.Sprintf("To: %s\r\nSubject: %s\r\n\r\n%s", to, subject, body))
	auth := smtp.PlainAuth("", s.cfg.SMTPUser, s.cfg.SMTPPass, s.cfg.SMTPHost)
	return smtp.SendMail(addr, auth, s.cfg.SMTPFrom, []string{to}, msg)
}
