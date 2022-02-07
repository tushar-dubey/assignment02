package config

type Config struct {
	App          App
	MongoAddress string
	APIKeys      []string
	Kafka        Kafka
}

type App struct {
	Env           string
	Port          string
	ShutdownDelay string
	Loglevel      string
	Query         string
}

type Kafka struct {
	Topic   string
	Address string
	Port    string
}
