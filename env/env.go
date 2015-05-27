package env

import (
	consul "github.com/hashicorp/consul/api"
)

type Environmenter interface {
	Set(string, string) error
	Get(string) string
	List() []string

	// Set or change the namespace. Doesn't do anything in Dev.
	SetNamespace(string) Environmenter
}

// DefaultConfig() abstracts out the api.DefaultConfig call
// entirely. This is designed so in default usecases you will
// only ever need to import this one package.
// Obviously if you require more power you have the ability to import
// github.com/hashicorp/consul/api and configure your client as needed.
func DefaultConfig() *consul.Config {
	return consul.DefaultConfig()
}
