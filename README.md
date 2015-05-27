# consul
A collection of abstractions that Kyäni uses to ease use with consul

## Env

Package Env utilizes the build flag `production`, if the flag is not present the package will run in "dev mode" which utilizes the `os` package to read environmental variables so a `env.Get()` in the background will run `os.Getenv()`. However, if the build flag `production` is set the package will abstract and grossly smiplify the consul K/V API.

**Note**: This package no where near leverages the power of the Consul K/V API. It is designed specifically to be a very stupid `set`, `get`, `list` so backwards compatibility can be maintained with the `os` package.

#### Env Installation

```
go get github.com/kyani-inc/consul/env
```

### Ideals

##### Namespaces

The way we leverage Consul is each application has their own folder structure that we call Namespaces to one versed in Consul K/V this is just a basic folder structure like `app/settings`. The main reason to break things into their own namespace is to fully support the backing dev version which relies on the `os`. As such the Namespace functions don't do anything when downgraded to the dev package.

Additionally key's with a `/` are converted to `.` in the dev package.

As such in dev you can have the following variables:

```
DB.USER=root
DB.PASSWORD=root
```

And you can reference these as `env.Get("DB/USER")` and `env.Get("DB/PASSWORD")` respectively. So you don't have to dumb down your production environment at all.

##### Live Environments

One of the reasons we use Consul is it gives us the option to change environment variables on the fly in realtime, whether that be to enable logging, change a debug level, or cycle database credentials on the fly.

As such the recommended method for utilizing this package (Your choice ultimately) is to not read env variables into global scopes but re-read the variable everytime you use it. Due to each of our servers running the Consul Agent we have determined that the overhead (very small) to re-reading in the variables each time outweighs the benefits of being able to change variables on the fly.

#### Basic Usage

```
package main

import (
    "fmt"

    "github.com/kyani-inc/consul/env"
)

func main() {
    e, err := env.New(env.DefaultConfig()) 
    if err != nil {
        // Handle the error
    }

    // Put something in the env
    e.Set("DB/USER", "root")

    // Read it back
    fmt.Println(e.Get("DB/USER"))

    // Change `e` to a new Namespace
    e = e.SetNamespace("app/myapp")
    fmt.Println(e.Namespace)

    // Add something in our new namespace
    e.Set("GITHUB/API_KEY", "12345")

    // Get something from our old namespace
    // This is only a temporary namespace change.
    fmt.Println(e.SetNamespace("").Get("DB/USER"))

    // Show everything in our environment
    fmt.Printf("%q\n", e.List())
}
```

This example will completely work in both a consul environment and local with os environment variables.


##### Running your application when using env

As mentioned previously, env utilizes build flags to switch between the dev environments as such the following methods will need to be adapted into your build process:

**development** (No change):

```
go run main.go
```

**production**:

```
go run -tags production main.go
```

##### Gotchas

- In Dev, Namespaces are not supported, so if your application uses the same env name in two different namespaces only one will be utilized.

##### Todo

- Tests
- Benchmarks