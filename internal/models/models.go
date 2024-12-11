package models

import (
	"time"

	pb "github.com/IskanderSh/tages-task/proto"
)

type MetaInfo struct {
	ID        int64     `json:"id"`
	FileName  string    `json:"fileName"`
	Path      string    `json:"path"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type UploadFileRequest struct {
	FileName string
	Content  []byte
}

func FromProtoToDomain(req *pb.UploadFileRequest) *UploadFileRequest {
	return &UploadFileRequest{
		FileName: req.FileName,
		Content:  req.Content,
	}
}
