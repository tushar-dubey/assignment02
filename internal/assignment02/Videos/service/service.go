package service

import (
	"assignment02/internal/assignment02/Videos/core"
	"assignment02/internal/assignment02/Videos/entity"
	"assignment02/rpc"
	"context"
	"github.com/kamva/mgm/v3"
)

type IService interface {
	Get(ctx context.Context, req *rpc.GetVideosRequest) (*[]entity.Video, error)
	Search(ctx context.Context, req *rpc.SearchVideosRequest) (*[]entity.Video, error)
}

type Service struct {
	core core.ICore
}

func NewService(collection *mgm.Collection) IService {
	return &Service{
		core: core.NewCore(collection),
	}
}

func (s Service) Get(ctx context.Context, req *rpc.GetVideosRequest) (*[]entity.Video, error) {
	return s.core.Get(ctx, req)
}

func (s Service) Search(ctx context.Context, req *rpc.SearchVideosRequest) (*[]entity.Video, error) {
	return s.core.Search(ctx, req)
}
