package Videos

import (
	"assignment02/internal/assignment02/Videos/entity"
	"assignment02/internal/assignment02/Videos/service"
	"assignment02/internal/assignment02/dto"
	"assignment02/internal/assignment02/validator"
	"assignment02/internal/boot/app"
	"assignment02/rpc"
	"context"
	"encoding/json"
	"github.com/kamva/mgm/v3"
	"github.com/segmentio/kafka-go"
	"github.com/twitchtv/twirp"
	"strconv"
	"time"
)

type Server struct {
	service service.IService
}

func NewServer() rpc.FetchVideos {
	collection := mgm.Coll(&entity.Video{})
	return &Server{
		service: service.NewService(collection),
	}
}

func (s Server) Get(ctx context.Context, request *rpc.GetVideosRequest) (*rpc.GetVideosResponse, error) {
	err := validator.ValidateGetRequest(request)
	if err != nil {
		return nil, twirp.NewError(twirp.InvalidArgument, err.Error())
	}
	videos, err := s.service.Get(ctx, request)
	if err != nil {
		return nil, twirp.NewError(twirp.Internal, err.Error())
	}
	var videosResponse []*rpc.Video

	for _, v := range *videos {
		video, err := v.ToProto()
		if err != nil {
			return nil, twirp.NewError(twirp.Internal, err.Error())
		}
		videosResponse = append(videosResponse, video)
	}
	return &rpc.GetVideosResponse{Videos: videosResponse}, nil
}

func (s Server) Search(ctx context.Context, request *rpc.SearchVideosRequest) (*rpc.SearchVideosResponse, error) {
	err := validator.ValidateSearchRequest(request)
	if err != nil {
		return nil, twirp.NewError(twirp.InvalidArgument, err.Error())
	}
	videos, err := s.service.Search(ctx, request)
	if err != nil {
		return nil, twirp.NewError(twirp.Internal, err.Error())
	}
	var videosResponse []*rpc.Video

	for _, v := range *videos {
		video, err := v.ToProto()
		if err != nil {
			return nil, twirp.NewError(twirp.Internal, err.Error())
		}
		videosResponse = append(videosResponse, video)
	}
	return &rpc.SearchVideosResponse{Videos: videosResponse}, nil
}

func (s Server) CronAPI(ctx context.Context, request *rpc.CronRequest) (*rpc.CronResponse, error) {
	// create start time as current time minus 20 and endtime as current time minus 10
	startTime := strconv.FormatInt(time.Now().Unix()-20, 10)
	endTime := strconv.FormatInt(time.Now().Unix()-10, 10)
	message := dto.KafkaMessage{
		StartTime: startTime,
		EndTime:   endTime,
	}
	m, _ := json.Marshal(message)
	err := app.GetAppContext().GetKafkaWriter().WriteMessages(ctx, kafka.Message{Value: m})
	if err != nil {
		return nil, err
	}
	return &rpc.CronResponse{Success: "True"}, nil
}
