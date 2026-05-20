package services

import (
	"context"
	"errors"
	"time"

	"github.com/decntrafir/backend/internal/config"
	"github.com/decntrafir/backend/internal/database"
	"github.com/decntrafir/backend/internal/models"
	"github.com/decntrafir/backend/pkg/blockchain"
	"github.com/decntrafir/backend/pkg/jwtutil"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrUserExists   = errors.New("user already exists")
	ErrInvalidCreds = errors.New("invalid credentials")
)

type AuthService struct {
	db    *database.DB
	cfg   *config.Config
	audit *AuditService
}

type AuthTokens struct {
	AccessToken  string       `json:"accessToken"`
	RefreshToken string       `json:"refreshToken"`
	User         models.User  `json:"user"`
}

type RegisterInput struct {
	Name          string      `json:"name"`
	Email         string      `json:"email"`
	Password      string      `json:"password"`
	Role          models.Role `json:"role"`
	Phone         string      `json:"phone"`
	Department    string      `json:"department"`
	BadgeNumber   string      `json:"badgeNumber"`
	WalletAddress string      `json:"walletAddress"`
	WalletSig     string      `json:"walletSignature"`
	WalletMessage string      `json:"walletMessage"`
}

func NewAuthService(db *database.DB, cfg *config.Config, audit *AuditService) *AuthService {
	return &AuthService{db: db, cfg: cfg, audit: audit}
}

func (s *AuthService) Register(ctx context.Context, in RegisterInput, ip string) (*AuthTokens, error) {
	count, err := s.db.Collection("users").CountDocuments(ctx, bson.M{"email": in.Email})
	if err != nil {
		return nil, err
	}
	if count > 0 {
		return nil, ErrUserExists
	}
	if in.WalletAddress != "" && in.WalletSig != "" {
		if !blockchain.VerifyWalletSignature(in.WalletAddress, in.WalletMessage, in.WalletSig) {
			return nil, errors.New("invalid wallet signature")
		}
	}
	role := in.Role
	if role == "" {
		role = models.RoleCitizen
	}
	if role != models.RoleCitizen {
		return nil, errors.New("only citizen self-registration allowed")
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(in.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	now := time.Now()
	user := models.User{
		ID:            primitive.NewObjectID(),
		Name:          in.Name,
		Email:         in.Email,
		Password:      string(hash),
		WalletAddress: in.WalletAddress,
		Role:          role,
		Phone:         in.Phone,
		CreatedAt:     now,
		UpdatedAt:     now,
	}
	if _, err := s.db.Collection("users").InsertOne(ctx, user); err != nil {
		return nil, err
	}
	_ = s.audit.Log(ctx, user.ID, "REGISTER", "user", user.ID.Hex(), ip, "", nil)
	return s.issueTokens(ctx, user)
}

func (s *AuthService) Login(ctx context.Context, email, password, ip string) (*AuthTokens, error) {
	var user models.User
	err := s.db.Collection("users").FindOne(ctx, bson.M{"email": email}).Decode(&user)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, ErrInvalidCreds
		}
		return nil, err
	}
	if bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)) != nil {
		return nil, ErrInvalidCreds
	}
	_ = s.audit.Log(ctx, user.ID, "LOGIN", "user", user.ID.Hex(), ip, "", nil)
	return s.issueTokens(ctx, user)
}

func (s *AuthService) Refresh(ctx context.Context, refreshToken string) (*AuthTokens, error) {
	userID, err := jwtutil.ParseRefreshToken(refreshToken, s.cfg.JWTSecret)
	if err != nil {
		return nil, err
	}
	var user models.User
	if err := s.db.Collection("users").FindOne(ctx, bson.M{"_id": userID, "refreshToken": refreshToken}).Decode(&user); err != nil {
		return nil, ErrInvalidCreds
	}
	return s.issueTokens(ctx, user)
}

func (s *AuthService) issueTokens(ctx context.Context, user models.User) (*AuthTokens, error) {
	access, err := jwtutil.GenerateAccessToken(user.ID, user.Email, string(user.Role), s.cfg.JWTSecret, s.cfg.JWTExpiry)
	if err != nil {
		return nil, err
	}
	refresh, err := jwtutil.GenerateRefreshToken(user.ID, s.cfg.JWTSecret, s.cfg.RefreshExpiry)
	if err != nil {
		return nil, err
	}
	_, err = s.db.Collection("users").UpdateOne(ctx, bson.M{"_id": user.ID}, bson.M{"$set": bson.M{"refreshToken": refresh, "updatedAt": time.Now()}})
	if err != nil {
		return nil, err
	}
	user.Password = ""
	user.RefreshToken = ""
	return &AuthTokens{AccessToken: access, RefreshToken: refresh, User: user}, nil
}
