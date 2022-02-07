package repo

import (
	"assignment02/internal/assignment02/Videos/entity"
	"assignment02/rpc"
	"context"
	"github.com/kamva/mgm/v3"
	"github.com/kamva/mgm/v3/operator"
	"go.mongodb.org/mongo-driver/bson"
)

type IRepo interface {
	Get(ctx context.Context, req *rpc.GetVideosRequest) (*[]entity.Video, error)
	Search(ctx context.Context, req *rpc.SearchVideosRequest) (*[]entity.Video, error)
	Create(video *entity.Video) error
}

type Repo struct {
	coll *mgm.Collection
}

func NewRepo(collection *mgm.Collection) IRepo {
	return &Repo{
		coll: collection,
	}
}

func (r Repo) Get(ctx context.Context, req *rpc.GetVideosRequest) (*[]entity.Video, error) {
	var videos []entity.Video

	err := r.coll.SimpleAggregate(&videos,
		bson.M{operator.Sort: bson.M{"created_at": -1}}, bson.M{operator.Skip: req.Skip}, bson.M{operator.Limit: req.Count})
	if err != nil {
		return nil, err
	}
	return &videos, nil
}

func (r Repo) Search(ctx context.Context, req *rpc.SearchVideosRequest) (*[]entity.Video, error) {
	var videos []entity.Video
	err := r.coll.SimpleAggregate(&videos,
		bson.M{operator.Match: bson.M{operator.Text: bson.M{"$search": req.Query}}}, bson.M{operator.Sort: bson.M{"created_at": -1}}, bson.M{operator.Skip: req.Skip}, bson.M{operator.Limit: req.Count})
	if err != nil {
		return nil, err
	}
	return &videos, nil
}

func (r Repo) Create(video *entity.Video) error {
	return r.coll.Create(video)
}
