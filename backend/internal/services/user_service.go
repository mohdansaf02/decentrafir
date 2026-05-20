package services

import (
	"context"
	"time"

	"github.com/decntrafir/backend/internal/database"
	"github.com/decntrafir/backend/internal/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	db    *database.DB
	audit *AuditService
}

func NewUserService(db *database.DB, audit *AuditService) *UserService {
	return &UserService{db: db, audit: audit}
}

func (s *UserService) List(ctx context.Context, role string, page, limit int64) ([]models.User, int64, error) {
	filter := bson.M{}
	if role != "" {
		filter["role"] = role
	}
	if limit <= 0 {
		limit = 20
	}
	if page <= 0 {
		page = 1
	}
	total, err := s.db.Collection("users").CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, err
	}
	opts := options.Find().SetSkip((page - 1) * limit).SetLimit(limit).SetSort(bson.D{{Key: "createdAt", Value: -1}})
	cursor, err := s.db.Collection("users").Find(ctx, filter, opts)
	if err != nil {
		return nil, 0, err
	}
	var users []models.User
	if err := cursor.All(ctx, &users); err != nil {
		return nil, 0, err
	}
	for i := range users {
		users[i].Password = ""
	}
	return users, total, nil
}

func (s *UserService) CreateByAdmin(ctx context.Context, adminID primitive.ObjectID, user models.User, password, ip string) (*models.User, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	now := time.Now()
	user.ID = primitive.NewObjectID()
	user.Password = string(hash)
	user.CreatedAt = now
	user.UpdatedAt = now
	_, err = s.db.Collection("users").InsertOne(ctx, user)
	if err != nil {
		return nil, err
	}
	user.Password = ""
	_ = s.audit.Log(ctx, adminID, "USER_CREATE", "user", user.ID.Hex(), ip, "", map[string]interface{}{"role": user.Role})
	return &user, nil
}

func (s *UserService) Delete(ctx context.Context, adminID, targetID primitive.ObjectID, ip string) error {
	_, err := s.db.Collection("users").DeleteOne(ctx, bson.M{"_id": targetID})
	if err == nil {
		_ = s.audit.Log(ctx, adminID, "USER_DELETE", "user", targetID.Hex(), ip, "", nil)
	}
	return err
}
