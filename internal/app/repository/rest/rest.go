package rest

import (
	"github.com/denis-oreshkevich/gopher-mart/internal/app/config"
	"github.com/go-resty/resty/v2"
)

type Repository struct {
	client *resty.Client
	conf   *config.Config
}

func NewRepository(client *resty.Client, conf *config.Config) *Repository {
	return &Repository{
		client: client,
		conf:   conf,
	}
}
