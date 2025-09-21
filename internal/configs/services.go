package configs

import "os"

type ServicesConfig struct {
	AuthServiceAddr        string
	UserServiceAddr        string
	ContentServiceAddr     string
	InteractionServiceAddr string
	FeedServiceAddr        string
	SearchServiceAddr      string
}

func NewServiceConfig() *ServicesConfig {
	return &ServicesConfig{
		AuthServiceAddr:        os.Getenv("AUTH_SERVICE_ADDR"),
		UserServiceAddr:        os.Getenv("USER_SERVICE_ADDR"),
		ContentServiceAddr:     os.Getenv("CONTENT_SERVICE_ADDR"),
		InteractionServiceAddr: os.Getenv("INTERACTION_SERVICE_ADDR"),
		FeedServiceAddr:        os.Getenv("FEED_SERVICE_ADDR"),
		SearchServiceAddr:      os.Getenv("SEARCH_SERVICE_ADDR"),
	}
}
