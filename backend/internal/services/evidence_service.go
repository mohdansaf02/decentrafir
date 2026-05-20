package services

import (
	"context"
	"fmt"
	"path/filepath"
	"strings"
	"time"

	"github.com/decntrafir/backend/internal/database"
	"github.com/decntrafir/backend/internal/models"
	"github.com/decntrafir/backend/pkg/ipfs"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var allowedExtensions = map[string]bool{
	".pdf": true, ".jpg": true, ".jpeg": true, ".png": true, ".mp4": true, ".doc": true, ".docx": true,
}

type EvidenceService struct {
	db     *database.DB
	ipfs   *ipfs.PinataClient
	audit  *AuditService
	maxMB  int64
}

func NewEvidenceService(db *database.DB, ipfsClient *ipfs.PinataClient, audit *AuditService, maxMB int64) *EvidenceService {
	return &EvidenceService{db: db, ipfs: ipfsClient, audit: audit, maxMB: maxMB}
}

func (s *EvidenceService) Upload(ctx context.Context, firID string, userID primitive.ObjectID, filename string, data []byte, ip string) (*models.Evidence, error) {
	ext := strings.ToLower(filepath.Ext(filename))
	if !allowedExtensions[ext] {
		return nil, fmt.Errorf("file type not allowed")
	}
	maxBytes := s.maxMB * 1024 * 1024
	if int64(len(data)) > maxBytes {
		return nil, fmt.Errorf("file exceeds %dMB limit", s.maxMB)
	}

	cid, url, err := s.ipfs.UploadFile(filename, data)
	if err != nil {
		return nil, err
	}

	firMongoID := primitive.NilObjectID
	evidence := models.Evidence{
		ID:         primitive.NewObjectID(),
		FIRID:      firID,
		FIRMongoID: firMongoID,
		UploadedBy: userID,
		FileName:   filename,
		FileType:   ext,
		FileSize:   int64(len(data)),
		IPFSCID:    cid,
		IPFSURL:    url,
		CreatedAt:  time.Now(),
	}
	_, err = s.db.Collection("evidence").InsertOne(ctx, evidence)
	if err != nil {
		return nil, err
	}
	_ = s.audit.Log(ctx, userID, "EVIDENCE_UPLOAD", "evidence", evidence.ID.Hex(), ip, "", map[string]interface{}{"cid": cid, "firId": firID})
	return &evidence, nil
}
