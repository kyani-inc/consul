// +build !production

package env

import (
	"os"
	"strings"

	consul "github.com/hashicorp/consul/api"
)

type devEnv struct {
	Namespace string
}

// New returns an empty instance of devEnv
func New(config *consul.Config) (Environmenter, error) {
	return env, nil
}

// Set will use os.Setenv to set an env variable
func (devEnv) Set(key, value string) error {
	key = strings.Replace(key, "/", ".", -1)

	return os.Setenv(key, value)
}

// Get will use os.Getenv to get an env variable
func (devEnv) Get(key string) string {
	key = strings.Replace(key, "/", ".", -1)

	return os.Getenv(key)
}

// List returns a k=v slice of strings as per os.Environ
func (devEnv) List() []string {
	return os.Environ()
}

// Namespace in dev doesn't do much
func (devEnv) SetNamespace(ns string) Environmenter {
	return devEnv{}
}
