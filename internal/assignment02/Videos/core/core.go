package core

import (
	"assignment02/internal/assignment02/Videos/entity"
	"assignment02/internal/assignment02/Videos/repo"
	"assignment02/rpc"
	"context"
	"github.com/kamva/mgm/v3"
)

type ICore interface {
	Get(ctx context.Context, req *rpc.GetVideosRequest) (*[]entity.Video, error)
	Search(ctx context.Context, req *rpc.SearchVideosRequest) (*[]entity.Video, error)
}

type Core struct {
	repo repo.IRepo
}

func NewCore(collection *mgm.Collection) ICore {
	return &Core{
		repo: repo.NewRepo(collection),
	}
}

func (c Core) Get(ctx context.Context, req *rpc.GetVideosRequest) (*[]entity.Video, error) {
	return c.repo.Get(ctx, req)
}

func (c Core) Search(ctx context.Context, req *rpc.SearchVideosRequest) (*[]entity.Video, error) {
	return c.repo.Search(ctx, req)
}
