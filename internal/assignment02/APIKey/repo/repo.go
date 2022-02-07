package repo

import (
	"assignment02/internal/assignment02/APIKey/entity"
	"context"
	"github.com/kamva/mgm/v3"
	"github.com/kamva/mgm/v3/operator"
	"go.mongodb.org/mongo-driver/bson"
	"strconv"
	"time"
)

type IRepo interface {
	Create(ctx context.Context, apiKey *entity.APIKey) error
	GetValidKey(apiKey *entity.APIKey) error
	Update(apiKey *entity.APIKey) error
}

type Repo struct {
	coll *mgm.Collection
}

func NewRepo() IRepo {
	return &Repo{coll: mgm.Coll(&entity.APIKey{})}
}

func (r Repo) Create(ctx context.Context, apiKey *entity.APIKey) error {
	err := r.coll.CreateWithCtx(ctx, apiKey)
	if err != nil {
		return err
	}
	return nil
}

func (r Repo) GetValidKey(apiKey *entity.APIKey) error {
	t := strconv.FormatInt(time.Now().Unix(), 10)
	return r.coll.First(bson.M{"enableAt": bson.M{operator.Lt: t}}, apiKey)
}

func (r Repo) Update(apiKey *entity.APIKey) error {
	return r.coll.Update(apiKey)
}
