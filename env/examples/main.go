package main

import (
	"fmt"

	"github.com/kyani-inc/consul/env"
	"github.com/subosito/gotenv" // Only needed for this demo
)

var logging bool

func init() {
	gotenv.Load() // Read the .env file into environment variables.
}

func main() {
	e, err := env.New(env.DefaultConfig())
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	logging = e.Get("ENABLE_LOGGING") == "on"

	if logging {
		fmt.Println("Connecting to database: ")
		fmt.Printf("%s:%s@tcp(%s:%s)/%s\n", e.Get("DB/USER"), e.Get("DB/PASS"), e.Get("DB/HOST"), e.Get("DB/PORT"), e.Get("DB/NAME"))
	}
}
