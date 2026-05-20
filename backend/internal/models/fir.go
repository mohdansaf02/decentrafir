package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type FIRStatus string

const (
	FIRStatusPending     FIRStatus = "pending"
	FIRStatusUnderReview FIRStatus = "under_review"
	FIRStatusApproved    FIRStatus = "approved"
	FIRStatusRejected    FIRStatus = "rejected"
	FIRStatusClosed      FIRStatus = "closed"
)

type FIR struct {
	ID              primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	FIRID           string             `bson:"firId" json:"firId"`
	Title           string             `bson:"title" json:"title"`
	Description     string             `bson:"description" json:"description"`
	CrimeType       string             `bson:"crimeType" json:"crimeType"`
	CitizenID       primitive.ObjectID `bson:"citizenId" json:"citizenId"`
	OfficerID       *primitive.ObjectID `bson:"officerId,omitempty" json:"officerId,omitempty"`
	Status          FIRStatus          `bson:"status" json:"status"`
	Location        string             `bson:"location" json:"location"`
	EvidenceFiles   []string           `bson:"evidenceFiles,omitempty" json:"evidenceFiles,omitempty"`
	IPFSHash        string             `bson:"ipfsHash,omitempty" json:"ipfsHash,omitempty"`
	MetadataHash    string             `bson:"metadataHash,omitempty" json:"metadataHash,omitempty"`
	TransactionHash string             `bson:"transactionHash,omitempty" json:"transactionHash,omitempty"`
	BlockNumber     uint64             `bson:"blockNumber,omitempty" json:"blockNumber,omitempty"`
	CreatedAt       time.Time          `bson:"createdAt" json:"createdAt"`
	UpdatedAt       time.Time          `bson:"updatedAt" json:"updatedAt"`
}

type CreateFIRRequest struct {
	Title       string `json:"title" binding:"required,min=3"`
	Description string `json:"description" binding:"required,min=10"`
	CrimeType   string `json:"crimeType" binding:"required"`
	Location    string `json:"location" binding:"required"`
	IPFSHash    string `json:"ipfsHash,omitempty"`
}

type UpdateFIRRequest struct {
	Status    *FIRStatus `json:"status,omitempty"`
	OfficerID *string    `json:"officerId,omitempty"`
	Title     *string    `json:"title,omitempty"`
	Location  *string    `json:"location,omitempty"`
}
