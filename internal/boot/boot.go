package boot

import (
	"assignment02/internal/assignment02/APIKey/entity"
	repo2 "assignment02/internal/assignment02/APIKey/repo"
	"assignment02/internal/assignment02/server/Videos"
	"assignment02/internal/boot/app"
	"assignment02/internal/config"
	"assignment02/rpc"
	"context"
	"github.com/kamva/mgm/v3"
	"github.com/segmentio/kafka-go"
	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"path"
	"runtime"
	"syscall"
)

var (
	Config config.Config
)

// Default options for configuration loading.
const (
	DefaultConfigType     = "toml"
	DefaultConfigDir      = "./config"
	DefaultConfigFileName = "default"
	WorkDirEnv            = "WORKDIR"
)

func NewContext(ctx context.Context) context.Context {
	if ctx == nil {
		ctx = context.Background()
	}
	return ctx
}

func GetEnv() string {
	env := os.Getenv("APP_ENV")
	if env == "" {
		env = "default"
	}

	return env
}

func InitConfigAndDB(ctx context.Context, config string) error {
	err := initConfig(&Config)
	if err != nil {
		return err
	}
	err = initDb(Config.MongoAddress)
	if err != nil {
		return err
	}
	appContext := app.NewAppContext(ctx)
	appContext.SetConfig(Config)
	return nil
}

func initConfig(config interface{}) error {
	var configPath string
	vi := viper.New()
	vi.SetConfigName(DefaultConfigFileName)
	vi.SetConfigType(DefaultConfigType)
	workDir := os.Getenv(WorkDirEnv)
	if workDir != "" {
		configPath = path.Join(workDir, DefaultConfigDir)
	} else {
		_, thisFile, _, _ := runtime.Caller(1)
		configPath = path.Join(path.Dir(thisFile), "../../"+DefaultConfigDir)
	}
	vi.AddConfigPath(configPath)
	vi.AutomaticEnv()
	if err := vi.ReadInConfig(); err != nil {
		return err
	}
	return vi.Unmarshal(config)
}

func initDb(connectionURI string) error {
	err := mgm.SetDefaultConfig(nil, "mgm_lab", options.Client().ApplyURI(connectionURI))
	if err != nil {
		return err
	}
	return nil
}

// RegisterAndServe Register the http server on the port defined in the Config and Shutdown Gracefully
func RegisterAndServe(ctx context.Context) {
	server := Videos.NewServer()
	handler := rpc.NewFetchVideosServer(server)

	mux := http.NewServeMux()

	mux.Handle(rpc.FetchVideosPathPrefix, handler)

	appServer := &http.Server{
		Handler: mux,
	}
	go func() {
		listener, err := net.Listen("tcp4", Config.App.Port)
		if err != nil {
			panic(err)
		}
		if err := appServer.Serve(listener); err != nil {
			panic(err)
		}
	}()

	c := make(chan os.Signal, 1)

	signal.Notify(c, syscall.SIGTERM, syscall.SIGINT)

	<-c

	err := appServer.Shutdown(ctx)
	if err != nil {
		log.Fatal(err)
	}
}

func CreateKafkaTopic() error {
	// create topic with a controller connection and close it on exit of this function
	var controllerConn *kafka.Conn
	controllerConn, err := kafka.Dial("tcp", net.JoinHostPort(Config.Kafka.Address, Config.Kafka.Port))
	if err != nil {
		return err
	}
	defer controllerConn.Close()

	topicConfigs := []kafka.TopicConfig{
		{
			Topic:             Config.Kafka.Topic,
			NumPartitions:     1,
			ReplicationFactor: 1,
		},
	}

	err = controllerConn.CreateTopics(topicConfigs...)
	if err != nil {
		return err
	}
	return nil
}

func InitKafkaReaderConnection(ctx context.Context) error {
	// create connection and save it in global var
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:   []string{net.JoinHostPort(Config.Kafka.Address, Config.Kafka.Port)},
		GroupID:   "group-A",
		Topic:     Config.Kafka.Topic,
		Partition: 0,
		MinBytes:  10e3, // 10KB
		MaxBytes:  10e6, // 10MB
	})
	app.GetAppContext().SetKafkaReader(r)
	return nil
}

func InitKafkaWriterConnection(ctx context.Context) error {
	// make a writer that produces to topic-A, using the least-bytes distribution
	w := &kafka.Writer{
		Addr:     kafka.TCP(net.JoinHostPort(Config.Kafka.Address, Config.Kafka.Port)),
		Topic:    Config.Kafka.Topic,
		Balancer: &kafka.LeastBytes{},
	}
	app.GetAppContext().SetKafkaWriter(w)
	return nil
}

func WriteAPIKeysToDB(ctx context.Context) error {
	repo := repo2.NewRepo()
	for _, key := range Config.APIKeys {
		apiKey := &entity.APIKey{Key: key}
		err := repo.Create(ctx, apiKey)
		if err != nil {
			return err
		}
	}
	return nil
}
