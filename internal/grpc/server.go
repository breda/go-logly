package grpc

import (
	"context"

	logv1 "github.com/breda/logly/api/v1"
	"github.com/breda/logly/internal/logly"
	"google.golang.org/grpc"
)

type GrpcServer struct {
	logv1.UnimplementedLoglyServer
	logly *logly.Logly
}

func newGrpcServer(logly *logly.Logly) *GrpcServer {
	return &GrpcServer{
		logly: logly,
	}
}

func NewGrpcServer(logly *logly.Logly) *grpc.Server {
	gsrv := grpc.NewServer()
	srv := newGrpcServer(logly)

	logv1.RegisterLoglyServer(gsrv, srv)
	return gsrv
}

func (s *GrpcServer) Append(ctx context.Context, req *logv1.AppendRequest) (*logv1.AppendResponse, error) {
	id, err := s.logly.Append(req.GetData())
	if err != nil {
		return nil, err
	}

	return &logv1.AppendResponse{
		Record: &logv1.Record{
			Id:   int64(id),
			Data: req.GetData(),
		},
	}, nil
}

func (s *GrpcServer) Fetch(ctx context.Context, req *logv1.FetchRequest) (*logv1.FetchResponse, error) {
	record, err := s.logly.Fetch(uint64(req.GetId()))
	if err != nil {
		return nil, err
	}

	return &logv1.FetchResponse{
		Record: record,
	}, nil
}

// stream here is a bi-directional
func (s *GrpcServer) AppendStream(stream logv1.Logly_AppendStreamServer) error {
	for {
		req, err := stream.Recv()
		if err != nil {
			return err
		}

		res, err := s.Append(stream.Context(), req)
		if err != nil {
			return err
		}

		if err = stream.Send(res); err != nil {
			return err
		}
	}
}

// stream here is a response stream
func (s *GrpcServer) FetchStream(req *logv1.FetchRequest, stream logv1.Logly_FetchStreamServer) error {
	for {
		select {
		case <-stream.Context().Done():
			return nil

		default:
			res, err := s.Fetch(stream.Context(), req)
			if err != nil {
				return err
			}

			if err = stream.Send(res); err != nil {
				return err
			}
		}
	}
}
