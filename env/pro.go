// +build production

package env

import (
	"fmt"
	"strings"

	consul "github.com/hashicorp/consul/api"
)

type proEnv struct {
	namespace string
	consul    *consul.KV
}

// New is used to return a new instance of a KV() consul client.
// It is provided a config to create the client.
// No default namespace is used.
func New(config *consul.Config) (Environmenter, error) {
	client, err := consul.NewClient(config)
	if err != nil {
		return nil, err
	}

	return &proEnv{
		consul: client.KV(),
	}, nil
}

// Set calls the consul.Put() method to save a value.
func (env proEnv) Set(key, value string) error {
	key = env.Namespace() + key

	_, err := env.consul.Put(&consul.KVPair{
		Key:   key,
		Value: []byte(value),
	}, &consul.WriteOptions{})

	return err
}

// Get calls the consul.Get() method to retrieve a value.
func (env proEnv) Get(key string) string {
	key = env.Namespace() + key

	pair, _, err := env.consul.Get(key, &consul.QueryOptions{
		// AllowStale: true is a set as an optimization technique,
		// allowing for us to query the consul agent and potentially
		// receive stale data.
		AllowStale: true,
	})

	if err != nil {
		return ""
	}

	return fmt.Sprintf("%s", pair.Value)
}

// List calls the consul.List() method to retrieve all values
// from a particular namespace. For backwards compatability,
// each KV Pair is formatted to match a standard k=v as per
// os.Environ.
func (env proEnv) List() []string {
	pairs, _, err := env.consul.List(env.Namespace(), &consul.QueryOptions{
		// AllowStale: false forces List() to query the consul servers directly
		AllowStale: false,
	})

	if err != nil {
		return []string{}
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

// Namespace allows you to set and change the namespace.
func (env proEnv) SetNamespace(ns string) Environmenter {
	env.namespace = strings.TrimRight(ns, "/") + "/"
	return env
}

func (env proEnv) Namespace() string {
	return env.namespace
}
