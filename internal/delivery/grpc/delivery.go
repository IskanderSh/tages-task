package delivery

import (
	"context"
	"fmt"
	"io"
	"log/slog"

	"github.com/IskanderSh/tages-task/internal/models"
	"github.com/IskanderSh/tages-task/pkg/errorlist"
	pb "github.com/IskanderSh/tages-task/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

type Handler struct {
	pb.UnimplementedFileProviderServer
	log     *slog.Logger
	service ServiceProvider
}

func NewHandler(log *slog.Logger, service ServiceProvider) *Handler {
	return &Handler{
		log:     log,
		service: service,
	}
}

type ServiceProvider interface {
	UploadFile(ctx context.Context, req models.UploadFileRequest, counter int) (fileName string, err error)
	FinishUpload(name string) error
}

func (h *Handler) UploadFile(stream pb.FileProvider_UploadFileServer) error {
	var (
		name    string
		counter int
	)

	for {
		req, err := stream.Recv()
		if err == io.EOF {
			closeErr := h.service.FinishUpload(name)
			if closeErr != nil {
				h.log.Error("error closing file after upload finish:", closeErr)
			}

			return stream.SendAndClose(&pb.UploadFileResponse{
				FileName: name,
			})
		}
		if err != nil {
			return status.Error(codes.InvalidArgument, fmt.Sprintf("%s: %v", errorlist.ErrInvalidValues, err))
		}

		uploadFileReq := models.FromProtoToDomain(req)
		if name != "" {
			uploadFileReq.FileName = name
		}

		name, err = h.service.UploadFile(context.Background(), *uploadFileReq, counter)
		if err != nil {
			return status.Error(codes.Internal, err.Error())
		}

		counter++
	}
}

func (h *Handler) DownloadFile(req *pb.DownloadFileRequest, stream pb.FileProvider_DownloadFileServer) error {
	return nil
}

func (h *Handler) FetchFiles(ctx context.Context, empty *emptypb.Empty) (*pb.FetchFilesResponse, error) {
	return nil, nil
}
