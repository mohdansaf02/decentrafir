package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type AuditLog struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	UserID    primitive.ObjectID `bson:"userId" json:"userId"`
	Action    string             `bson:"action" json:"action"`
	Resource  string             `bson:"resource" json:"resource"`
	ResourceID string            `bson:"resourceId,omitempty" json:"resourceId,omitempty"`
	Details   map[string]interface{} `bson:"details,omitempty" json:"details,omitempty"`
	IPAddress string             `bson:"ipAddress,omitempty" json:"ipAddress,omitempty"`
	TxHash    string             `bson:"txHash,omitempty" json:"txHash,omitempty"`
	CreatedAt time.Time          `bson:"createdAt" json:"createdAt"`
}
