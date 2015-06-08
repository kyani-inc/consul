package env_test

import (
	"strings"
	"testing"

	"github.com/kyani-inc/consul/env"
)

func testSet(T *testing.T) {
	cli, _ := env.New(env.DefaultConfig())

	err := cli.Set("key", "value")
	if err != nil {
		T.Errorf("Expected no error but received - %s\n", err.Error())
	}
}

func testList(T *testing.T) {
	cli, _ := env.New(env.DefaultConfig())

	list := cli.List()
	if len(list) == 0 {
		T.Error("Expected to at least have one element in the list but did not receive any.")
	}
}

func testListWithNamespace(T *testing.T) {
	ns := "test/namespace"
	cli, _ := env.New(env.DefaultConfig())
	cli = cli.SetNamespace(ns)

	list := cli.List()
	for _, item := range list {
		if strings.Contains(item, ns) {
			T.Errorf("Namespace not properly stripped from list item: %s\n", item)
		}
	}
}

func testGet(T *testing.T) {
	tests := map[string]string{
		"key":         "value",
		"key/key":     "value",
		"key/key/key": "value",
		"key/key-key": "value",
		"key_key":     "value",
		"key.key":     "value",
	}

	cli, _ := env.New(env.DefaultConfig())

	for key, value := range tests {
		cli.Set(key, value)
		_value := cli.Get(key)
		if _value != value {
			T.Errorf("Expected '%s' but received '%s'\n", value, _value)
		}
	}
}
