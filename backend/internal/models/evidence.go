package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Evidence struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	FIRID       string             `bson:"firId" json:"firId"`
	FIRMongoID  primitive.ObjectID `bson:"firMongoId" json:"firMongoId"`
	UploadedBy  primitive.ObjectID `bson:"uploadedBy" json:"uploadedBy"`
	FileName    string             `bson:"fileName" json:"fileName"`
	FileType    string             `bson:"fileType" json:"fileType"`
	FileSize    int64              `bson:"fileSize" json:"fileSize"`
	IPFSCID     string             `bson:"ipfsCid" json:"ipfsCid"`
	IPFSURL     string             `bson:"ipfsUrl" json:"ipfsUrl"`
	CreatedAt   time.Time          `bson:"createdAt" json:"createdAt"`
}
