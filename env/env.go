// Package Env allows developers to use the
// [consul](https://github.com/hashicorp/consul) api
// as a storage for environment variables. This package utilizes
// build flags to allow developers to degrade their application
// to a dev environment that relies on OS environment variables,
// without the need to change their source code.
//
// At the production level this package abstracts out the
// consul.KV() api and at the dev level it abstracts out
// the core `os` package.
package env // import "github.com/kyani-inc/consul/env"

import (
	consul "github.com/hashicorp/consul/api"
)

type Environmenter interface {
	Set(string, string) error
	Get(string) string
	List() []string

	// Set or change the namespace. Doesn't do anything in Dev.
	SetNamespace(string) Environmenter
	Namespace() string
}

// DefaultConfig() abstracts out the api.DefaultConfig call
// entirely. This is designed so in default usecases you will
// only ever need to import this one package.
// Obviously if you require more power you have the ability to import
// github.com/hashicorp/consul/api and configure your client as needed.
func DefaultConfig() *consul.Config {
	return consul.DefaultConfig()
}
