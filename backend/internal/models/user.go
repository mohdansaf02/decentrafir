package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Role string

const (
	RoleCitizen Role = "citizen"
	RolePolice  Role = "police"
	RoleAdmin   Role = "admin"
)

type User struct {
	ID            primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Name          string             `bson:"name" json:"name"`
	Email         string             `bson:"email" json:"email"`
	Password      string             `bson:"password" json:"-"`
	WalletAddress string             `bson:"walletAddress,omitempty" json:"walletAddress,omitempty"`
	Role          Role               `bson:"role" json:"role"`
	Phone         string             `bson:"phone,omitempty" json:"phone,omitempty"`
	Department    string             `bson:"department,omitempty" json:"department,omitempty"`
	BadgeNumber   string             `bson:"badgeNumber,omitempty" json:"badgeNumber,omitempty"`
	RefreshToken  string             `bson:"refreshToken,omitempty" json:"-"`
	CreatedAt     time.Time          `bson:"createdAt" json:"createdAt"`
	UpdatedAt     time.Time          `bson:"updatedAt" json:"updatedAt"`
}
