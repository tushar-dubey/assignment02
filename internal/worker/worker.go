package worker

import (
	"assignment02/internal/assignment02/APIKey/entity"
	repo2 "assignment02/internal/assignment02/APIKey/repo"
	entity2 "assignment02/internal/assignment02/Videos/entity"
	repo3 "assignment02/internal/assignment02/Videos/repo"
	"assignment02/internal/assignment02/dto"
	"assignment02/internal/boot/app"
	"context"
	"encoding/json"
	"github.com/golang-module/carbon"
	"github.com/kamva/mgm/v3"
	"github.com/segmentio/kafka-go"
	"google.golang.org/api/option"
	"google.golang.org/api/youtube/v3"
	"google.golang.org/appengine/log"
	log2 "log"
	"strconv"
	"time"
)

func RunWorker(ctx context.Context) {
	defer func(KafkaReader *kafka.Reader) {
		err := KafkaReader.Close()
		if err != nil {
			log2.Fatalf("Cannot Close Kafka Reader Connection: %v", err.Error())
		}
	}(app.GetAppContext().GetKafkaReader())
	for {
		m, err := app.GetAppContext().GetKafkaReader().FetchMessage(ctx)
		if err != nil {
			log.Infof(ctx, "error reading message from kafka: %v", err.Error())
			continue
		}
		message := &dto.KafkaMessage{}
		err = json.Unmarshal(m.Value, message)
		if err != nil {
			panic(err)
		}
		repo := repo2.NewRepo()
		apiKey := &entity.APIKey{}
		err = repo.GetValidKey(apiKey)
		if err != nil {
			log2.Fatalf("Cannot Get Valid API Key: %v", err.Error())
		}
		youtubeClient, err := GetYoutubeClient(ctx, apiKey)
		if err != nil {
			log2.Fatalf("Cannot Create Valid Youtube Client: %v", err.Error())
		}
		part := []string{"snippet"}
		publishedAfterInt, _ := strconv.ParseInt(message.StartTime, 10, 64)
		publishedAfter := carbon.CreateFromTimestamp(publishedAfterInt).ToRfc3339String()
		publishedBeforeInt, _ := strconv.ParseInt(message.EndTime, 10, 64)
		publishedBefore := carbon.CreateFromTimestamp(publishedBeforeInt).ToRfc3339String()
		response, err := youtubeClient.Search.List(part).Type("video").
			Order("date").PublishedAfter(publishedAfter).PublishedBefore(publishedBefore).Do()
		if err != nil {
			// Handle the case of request limit being breached for the API key
			log.Infof(ctx, "Cannot Get Valid Youtube Response: %v", err.Error())
			markAPIKeyAsExpired(repo, apiKey)
			// Continue to the next loop iteration without committing the message, the message will remain in Kafka to be picked up again
			continue
		}
		err = parseAndSaveResponseToDB(response)
		if err != nil {
			log2.Fatalf("Error Parsing response from Youtube: %v", err.Error())
		}
		err = app.GetAppContext().GetKafkaReader().CommitMessages(ctx, m)
		if err != nil {
			log2.Fatalf("Error Commiting Message to Kafka: %v", err.Error())
		}
	}
}

func markAPIKeyAsExpired(repo repo2.IRepo, key *entity.APIKey) {
	// mark the key as disabled till 24 hours from now
	enableAt := strconv.FormatInt(time.Now().Local().Add(time.Hour*time.Duration(24)).Unix(), 10)
	key.EnableAt = enableAt
	_ = repo.Update(key)
}

func parseAndSaveResponseToDB(response *youtube.SearchListResponse) error {
	repo := repo3.NewRepo(mgm.Coll(&entity2.Video{}))
	for _, item := range response.Items {
		var defaultThumbnail entity2.Thumbnail
		var maxresThumbnail entity2.Thumbnail
		var standardThumbnail entity2.Thumbnail
		var highThumbnail entity2.Thumbnail
		var mediumThumbnail entity2.Thumbnail
		if item.Snippet.Thumbnails.Default != nil {
			defaultThumbnail = entity2.Thumbnail{
				Height: item.Snippet.Thumbnails.Default.Height,
				Width:  item.Snippet.Thumbnails.Default.Width,
				URL:    item.Snippet.Thumbnails.Default.Url,
			}
		}
		if item.Snippet.Thumbnails.Maxres != nil {
			maxresThumbnail = entity2.Thumbnail{
				Height: item.Snippet.Thumbnails.Maxres.Height,
				Width:  item.Snippet.Thumbnails.Maxres.Width,
				URL:    item.Snippet.Thumbnails.Maxres.Url,
			}
		}
		if item.Snippet.Thumbnails.Standard != nil {
			standardThumbnail = entity2.Thumbnail{
				Height: item.Snippet.Thumbnails.Standard.Height,
				Width:  item.Snippet.Thumbnails.Standard.Width,
				URL:    item.Snippet.Thumbnails.Standard.Url,
			}
		}
		if item.Snippet.Thumbnails.High != nil {
			highThumbnail = entity2.Thumbnail{
				Height: item.Snippet.Thumbnails.High.Height,
				Width:  item.Snippet.Thumbnails.High.Width,
				URL:    item.Snippet.Thumbnails.High.Url,
			}
		}
		if item.Snippet.Thumbnails.Medium != nil {
			mediumThumbnail = entity2.Thumbnail{
				Height: item.Snippet.Thumbnails.Medium.Height,
				Width:  item.Snippet.Thumbnails.Medium.Width,
				URL:    item.Snippet.Thumbnails.Medium.Url,
			}
		}
		thubnailDetails := entity2.ThumbnailDetails{
			Default:  defaultThumbnail,
			MaxRes:   maxresThumbnail,
			High:     highThumbnail,
			Medium:   mediumThumbnail,
			Standard: standardThumbnail,
		}
		video := &entity2.Video{
			ChannelId:        item.Snippet.ChannelId,
			ChannelTitle:     item.Snippet.ChannelTitle,
			Description:      item.Snippet.Description,
			PublishedAt:      item.Snippet.PublishedAt,
			ThumbnailDetails: thubnailDetails,
			Title:            item.Snippet.Title,
		}
		err := repo.Create(video)
		if err != nil {
			return err
		}
	}
	return nil
}

func GetYoutubeClient(ctx context.Context, apiKey *entity.APIKey) (*youtube.Service, error) {
	clientOption := option.WithAPIKey(apiKey.Key)
	client, err := youtube.NewService(ctx, clientOption)
	if err != nil {
		return nil, err
	}
	return client, nil
}
