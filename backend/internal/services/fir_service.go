package services

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/decntrafir/backend/internal/database"
	"github.com/decntrafir/backend/internal/models"
	"github.com/decntrafir/backend/pkg/blockchain"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type FIRService struct {
	db           *database.DB
	chain        *blockchain.Client
	audit        *AuditService
	notification *NotificationService
}

type FIRListQuery struct {
	Status    string
	CrimeType string
	FIRID     string
	FromDate  *time.Time
	ToDate    *time.Time
	OfficerID string
	Page      int64
	Limit     int64
}

type FIRAnalytics struct {
	TotalFIRs    int64            `json:"totalFirs"`
	PendingFIRs  int64            `json:"pendingFirs"`
	ApprovedFIRs int64            `json:"approvedFirs"`
	ByStatus     map[string]int64 `json:"byStatus"`
	ByCrimeType  map[string]int64 `json:"byCrimeType"`
	Recent       []models.FIR     `json:"recent"`
}

func NewFIRService(db *database.DB, chain *blockchain.Client, audit *AuditService, notification *NotificationService) *FIRService {
	return &FIRService{db: db, chain: chain, audit: audit, notification: notification}
}

func (s *FIRService) Create(ctx context.Context, citizenID primitive.ObjectID, req models.CreateFIRRequest, ip string) (*models.FIR, error) {
	now := time.Now()
	firID := fmt.Sprintf("FIR-%s", uuid.New().String()[:8])
	fir := models.FIR{
		ID:          primitive.NewObjectID(),
		FIRID:       firID,
		Title:       req.Title,
		Description: req.Description,
		CrimeType:   req.CrimeType,
		CitizenID:   citizenID,
		Status:      models.FIRStatusPending,
		Location:    req.Location,
		IPFSHash:    req.IPFSHash,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
	meta, _ := json.Marshal(map[string]interface{}{
		"firId": firID, "title": req.Title, "description": req.Description,
		"crimeType": req.CrimeType, "location": req.Location, "ipfs": req.IPFSHash,
	})
	fir.MetadataHash = blockchain.MetadataHash(string(meta))
	txHash, blockNum, err := s.chain.CreateFIROnChain(ctx, firID, string(meta))
	if err != nil {
		return nil, err
	}
	fir.TransactionHash = txHash
	fir.BlockNumber = blockNum

	if _, err := s.db.Collection("firs").InsertOne(ctx, fir); err != nil {
		return nil, err
	}
	_ = s.audit.Log(ctx, citizenID, "FIR_CREATE", "fir", fir.FIRID, ip, txHash, map[string]interface{}{"status": fir.Status})
	_ = s.notification.Create(ctx, citizenID, "FIR Registered", "Your FIR has been submitted.", "fir_created", fir.FIRID)
	return &fir, nil
}

func (s *FIRService) GetByID(ctx context.Context, id string) (*models.FIR, error) {
	var fir models.FIR
	filter := bson.M{}
	if oid, err := primitive.ObjectIDFromHex(id); err == nil {
		filter["_id"] = oid
	} else {
		filter["firId"] = id
	}
	err := s.db.Collection("firs").FindOne(ctx, filter).Decode(&fir)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return nil, errors.New("fir not found")
	}
	return &fir, err
}

func (s *FIRService) Update(ctx context.Context, id string, userID primitive.ObjectID, role models.Role, req models.UpdateFIRRequest, ip string) (*models.FIR, error) {
	fir, err := s.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	update := bson.M{"updatedAt": time.Now()}
	if req.Title != nil {
		update["title"] = *req.Title
	}
	if req.Location != nil {
		update["location"] = *req.Location
	}
	if req.Status != nil {
		if role != models.RolePolice && role != models.RoleAdmin {
			return nil, errors.New("only police or admin can update status")
		}
		update["status"] = *req.Status
		txHash, err := s.chain.UpdateStatusOnChain(ctx, fir.FIRID, blockchain.StatusToChain(string(*req.Status)))
		if err != nil {
			return nil, err
		}
		_ = s.audit.Log(ctx, userID, "FIR_STATUS_UPDATE", "fir", fir.FIRID, ip, txHash, map[string]interface{}{"status": *req.Status})
		_ = s.notification.Create(ctx, fir.CitizenID, "FIR Status Updated", fmt.Sprintf("Status changed to %s", *req.Status), "fir_updated", fir.FIRID)
	}
	if req.OfficerID != nil {
		if role != models.RolePolice && role != models.RoleAdmin {
			return nil, errors.New("only police or admin can assign officer")
		}
		oid, err := primitive.ObjectIDFromHex(*req.OfficerID)
		if err != nil {
			return nil, err
		}
		update["officerId"] = oid
	}
	_, err = s.db.Collection("firs").UpdateOne(ctx, bson.M{"_id": fir.ID}, bson.M{"$set": update})
	if err != nil {
		return nil, err
	}
	return s.GetByID(ctx, fir.ID.Hex())
}

func (s *FIRService) List(ctx context.Context, q FIRListQuery, citizenFilter *primitive.ObjectID) ([]models.FIR, int64, error) {
	filter := bson.M{}
	if citizenFilter != nil {
		filter["citizenId"] = *citizenFilter
	}
	if q.Status != "" {
		filter["status"] = q.Status
	}
	if q.CrimeType != "" {
		filter["crimeType"] = q.CrimeType
	}
	if q.FIRID != "" {
		filter["firId"] = bson.M{"$regex": q.FIRID, "$options": "i"}
	}
	if q.OfficerID != "" {
		if oid, err := primitive.ObjectIDFromHex(q.OfficerID); err == nil {
			filter["officerId"] = oid
		}
	}
	if q.FromDate != nil || q.ToDate != nil {
		dateFilter := bson.M{}
		if q.FromDate != nil {
			dateFilter["$gte"] = *q.FromDate
		}
		if q.ToDate != nil {
			dateFilter["$lte"] = *q.ToDate
		}
		filter["createdAt"] = dateFilter
	}
	if q.Limit <= 0 {
		q.Limit = 20
	}
	if q.Page <= 0 {
		q.Page = 1
	}
	skip := (q.Page - 1) * q.Limit
	total, err := s.db.Collection("firs").CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, err
	}
	opts := options.Find().SetSort(bson.D{{Key: "createdAt", Value: -1}}).SetSkip(skip).SetLimit(q.Limit)
	cursor, err := s.db.Collection("firs").Find(ctx, filter, opts)
	if err != nil {
		return nil, 0, err
	}
	var firs []models.FIR
	if err := cursor.All(ctx, &firs); err != nil {
		return nil, 0, err
	}
	return firs, total, nil
}

func (s *FIRService) Analytics(ctx context.Context) (*FIRAnalytics, error) {
	total, _ := s.db.Collection("firs").CountDocuments(ctx, bson.M{})
	pending, _ := s.db.Collection("firs").CountDocuments(ctx, bson.M{"status": models.FIRStatusPending})
	approved, _ := s.db.Collection("firs").CountDocuments(ctx, bson.M{"status": models.FIRStatusApproved})

	byStatus := map[string]int64{}
	statuses := []models.FIRStatus{models.FIRStatusPending, models.FIRStatusUnderReview, models.FIRStatusApproved, models.FIRStatusRejected, models.FIRStatusClosed}
	for _, st := range statuses {
		c, _ := s.db.Collection("firs").CountDocuments(ctx, bson.M{"status": st})
		byStatus[string(st)] = c
	}

	pipeline := mongo.Pipeline{
		{{Key: "$group", Value: bson.M{"_id": "$crimeType", "count": bson.M{"$sum": 1}}}},
	}
	cursor, _ := s.db.Collection("firs").Aggregate(ctx, pipeline)
	type agg struct {
		ID    string `bson:"_id"`
		Count int64  `bson:"count"`
	}
	var crimeAgg []agg
	_ = cursor.All(ctx, &crimeAgg)
	byCrime := map[string]int64{}
	for _, a := range crimeAgg {
		byCrime[a.ID] = a.Count
	}

	recent, _, _ := s.List(ctx, FIRListQuery{Limit: 5}, nil)
	return &FIRAnalytics{
		TotalFIRs: total, PendingFIRs: pending, ApprovedFIRs: approved,
		ByStatus: byStatus, ByCrimeType: byCrime, Recent: recent,
	}, nil
}
