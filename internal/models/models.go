package models

import (
	"time"

	pb "github.com/IskanderSh/tages-task/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type MetaInfo struct {
	ID        int64     `json:"id"`
	FileName  string    `json:"fileName"`
	Path      string    `json:"path"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

func (m *MetaInfo) ToProto() *pb.File {
	return &pb.File{
		FileName:  m.FileName,
		CreatedAt: timestamppb.New(m.CreatedAt),
		UpdatedAt: timestamppb.New(m.UpdatedAt),
	}
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

type FetchFilesResponse struct {
	ID        int64
	FileName  string
	CreatedAt time.Time
	UpdatedAt time.Time
}
