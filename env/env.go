// Package Env allows developers to use the
// [consul](https://github.com/hashicorp/consul) api
// as a storage for environment variables.
//
// At the production level this package abstracts out the
// consul.KV() api and at the dev level it abstracts out
// the core `os` package.
package env // import "github.com/kyani-inc/consul/env"

import (
	"fmt"
	"os"
	"strings"

	consul "github.com/hashicorp/consul/api"
)

// DefaultConfig() abstracts out the api.DefaultConfig call
// entirely. This is so in default usecases you will
// only ever need to import this one package.
// Obviously if you require more power you have the ability to import
// github.com/hashicorp/consul/api and configure your client as needed.
func DefaultConfig() *consul.Config {
	return consul.DefaultConfig()
}

type Env struct {
	namespace string
	kv        *consul.KV
}

// New is used to return a new instance of a KV() consul client.
// It is provided a config to create the client.
// No default namespace is used.
func New(config *consul.Config) (Env, error) {
	client, err := consul.NewClient(config)
	if err != nil {
		return Env{}, err
	}

	return Env{
		kv: client.KV(),
	}, nil
}

// osCleanKey is designed to clean a key for fallback purposes
func (env Env) osCleanKey(key string) string {
	// key shouldn't have the namespace but we'll remove it just in case
	key = strings.Replace(key, env.namespace, "", -1)

	// Convert '/' to '___'
	key = strings.Replace(key, "/", "___", -1)

	// Trim "/"
	key = strings.Trim(key, "/")
	return key
}

func (env Env) Get(key string) string {
	key = env.Namespace() + key

	pair, _, err := env.kv.Get(key, &consul.QueryOptions{
		// AllowStale: true is a set as an optimization technique,
		// allowing for us to query the consul agent and potentially
		// receive stale data.
		AllowStale: true,
	})

	if err != nil {
		return env.osGet(key)
	}

	if pair == nil {
		return ""
	}

	return string(pair.Value)
}

func (env Env) osGet(key string) string {
	debug("[env] Falling back to os.Getenv().")

	return os.Getenv(env.osCleanKey(key))
}

// Set calls the consul.Put() method to save a value.
func (env Env) Set(key, value string) error {
	key = env.Namespace() + key

	_, err := env.kv.Put(&consul.KVPair{
		Key:   key,
		Value: []byte(value),
	}, &consul.WriteOptions{})

	if err != nil {
		err = env.osSet(key, value)
	}

	return err
}

// osSet is the fallback for Set() which
// uses os.Setenv()
func (env Env) osSet(key, value string) error {
	debug("[env] Falling back to os.Setenv().")

	return os.Setenv(env.osCleanKey(key), value)
}

// List calls the consul.List() method to retrieve all values
// from a particular namespace. For backwards compatability,
// each KV Pair is formatted to match a standard k=v as per
// os.Environ.
func (env Env) List() []string {
	pairs, _, err := env.kv.List(env.Namespace(), &consul.QueryOptions{
		// AllowStale: false forces List() to query the consul servers directly
		AllowStale: false,
	})

	if err != nil {
		return env.osList()
	}

	fmtPair := func(kvPair *consul.KVPair) string {
		key := strings.Replace(kvPair.Key, env.Namespace(), "", -1)
		return fmt.Sprintf("%s=%s", key, kvPair.Value)
	}

	// Iterate over the pairs and fmt them like os.Environ does (k=v)
	var p []string
	for _, pair := range pairs {
		if pair.Key == env.Namespace() {
			continue
		}

		p = append(p, fmtPair(pair))
	}

	return p
}

// osList is the fallback for List() which uses os.Environ()
func (Env) osList() []string {
	debug("[env] Falling back to os.Envrion().")

	return os.Environ()
}

// Namespace allows you to set and change the namespace.
func (env Env) SetNamespace(ns string) Env {
	env.namespace = strings.TrimRight(ns, "/") + "/"
	return env
}

// Namespace returns the set namespace
func (env Env) Namespace() string {
	return env.namespace
}
