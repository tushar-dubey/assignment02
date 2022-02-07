package entity

import (
	"assignment02/rpc"
	"encoding/json"
	"fmt"
	"github.com/kamva/mgm/v3"
)

type Video struct {
	mgm.DefaultModel `bson:",inline"`
	ChannelId        string           `json:"channel_id" bson:"channel_id"`
	ChannelTitle     string           `json:"channel_title" bson:"channel_title"`
	Description      string           `json:"description" bson:"description"`
	PublishedAt      string           `json:"published_at" bson:"published_at"`
	ThumbnailDetails ThumbnailDetails `json:"thumbnailDetails" bson:"thumbnailDetails"`
	Title            string           `json:"title" bson:"title"`
}

type ThumbnailDetails struct {
	Default  Thumbnail `json:"default,omitempty" bson:"default"`
	MaxRes   Thumbnail `json:"maxRes,omitempty" bson:"maxRes"`
	High     Thumbnail `json:"high,omitempty" bson:"high"`
	Medium   Thumbnail `json:"medium,omitempty" bson:"medium"`
	Standard Thumbnail `json:"standard,omitempty" bson:"standard"`
}

type Thumbnail struct {
	Height int64  `json:"height" bson:"height"`
	Width  int64  `json:"width" bson:"width"`
	URL    string `json:"url" bson:"url"`
}

func (v *Video) ToProto() (*rpc.Video, error) {
	e := &rpc.Video{}

	e.Id = fmt.Sprintf("%v", v.GetID())
	e.ChannelTitle = v.ChannelTitle
	e.ChannelId = v.ChannelId
	e.Description = v.Description
	e.Title = v.Title
	e.PublishedAt = v.PublishedAt

	data, _ := json.Marshal(v.ThumbnailDetails)

	err := json.Unmarshal(data, &e.ThumbnailDetails)
	if err != nil {
		return nil, err
	}
	return e, nil
}
