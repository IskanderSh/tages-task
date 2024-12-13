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
	GetFiles(ctx context.Context) ([]models.MetaInfo, error)
	DownloadFile(name string, counter int, buffer []byte) (bytesRead int, err error)
}

func (h *Handler) UploadFile(stream pb.FileProvider_UploadFileServer) error {
	var (
		name    string
		counter int
	)

	for {
		req, err := stream.Recv()
		if err == io.EOF {
			h.MustClose(name)

			return stream.SendAndClose(&pb.UploadFileResponse{
				FileName: name,
			})
		}
		if err != nil {
			h.log.Error(err.Error())
			return status.Error(codes.InvalidArgument, fmt.Sprintf("%s: %v", errorlist.ErrInvalidValues, err))
		}

		uploadFileReq := models.FromProtoToDomain(req)
		// set saved file name
		if name != "" {
			uploadFileReq.FileName = name
		}

		name, err = h.service.UploadFile(context.Background(), *uploadFileReq, counter)
		if err != nil {
			h.log.Error(err.Error())
			return status.Error(codes.Internal, err.Error())
		}

		counter++
	}
}

func (h *Handler) DownloadFile(req *pb.DownloadFileRequest, stream pb.FileProvider_DownloadFileServer) error {
	buffer := make([]byte, 1024*1024)
	counter := 0

	for {
		// DownloadFile is writing values into buffer, this is not go way,
		// but it helps to not initialize new slice every function call
		bytesRead, err := h.service.DownloadFile(req.FileName, counter, buffer)
		if err != nil {
			if err == io.EOF {
				h.MustClose(req.FileName)
				break
			}

			h.log.Error("error downloading file: ", err)
			return status.Error(codes.Internal, err.Error())
		}

		err = stream.Send(&pb.DownloadFileResponse{
			Content: buffer[:bytesRead],
		})
		if err != nil {
			h.log.Error("error sending chunk to client: ", err)
			return status.Error(codes.Internal, err.Error())
		}

		counter++
	}

	h.log.Info("stream sending is ended successfully")
	return nil
}

func (h *Handler) FetchFiles(ctx context.Context, empty *emptypb.Empty) (*pb.FetchFilesResponse, error) {
	values, err := h.service.GetFiles(ctx)
	if err != nil {
		h.log.Error(err.Error())
		return nil, status.Error(codes.Internal, err.Error())
	}

	list := make([]*pb.File, 0, len(values))
	for _, value := range values {
		list = append(list, value.ToProto())
	}

	return &pb.FetchFilesResponse{
		Data: list,
	}, nil
}

func (h *Handler) MustClose(name string) {
	if err := h.service.FinishUpload(name); err != nil {
		h.log.Error("error closing file: ", err)
	}
}
