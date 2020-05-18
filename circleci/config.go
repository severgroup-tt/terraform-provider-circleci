package circleci

import (
	circleci "github.com/samanthaq/terraform-provider-circleci/circleci/circleci-go"
)

type Config struct {
	AuthToken string
	BaseURL   string
}

func (c *Config) NewClient() *circleci.Client {
	return circleci.NewClient(c.AuthToken, c.BaseURL)
}
