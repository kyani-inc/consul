# consul
A collection of abstractions that Kyäni uses to ease use with consul

Links

- [Discovery](consul/tree/master/discovery)

## Env

Package Env allows developers to use the [consul](https://github.com/hashicorp/consul) api as a storage for environment variables. This package utilizes build flags to allow developers to degrade their application to a dev environment that relies on OS environment variables, without the need to change their source code.

At the production level this package abstracts out the consul.KV() api and at the dev level it abstracts out the core `os` package.

**Note**: This package is designed to be simplistic enough to work with both the `os` package and the [consul](https://github.com/hashicorp/consul) api.

#### Env Installation

```
go get github.com/kyani-inc/consul/env
```

### Ideals

##### Namespaces

The way we leverage Consul is each application has their own folder structure that we call Namespaces. In normal Consul ideals a "key" is a full folder path such as `app/settings/db_user`. The way we would view a namespace is `app/settings` with the key being `db_user`.

We have chosen to break things into namespaces in order to support the fallback to OS environment variables where folders might not be supported. 

###### In Dev

- Namespaces in development are silently ignored. See Gotcha below.
- Folders are supported by converting `/` to `.`

**Example**:
In dev you can have the following variables:

```
DB.USER=root
DB.PASSWORD=root
```

You can then reference these as `env.Get("DB/USER")` and `env.Get("DB/PASSWORD")` respectively. This should allow you to best match your production setup.


##### Live Environments

One of the reasons we use Consul is it gives us the option to change environment variables in realtime. Examples include; enabling logging, changing a debug level, or cycling database credentials on the fly.

The way we recommend using this package is to not read your environment variables into a persistent scope, but to instead re-read the variables everytime you need them. 

Our internal tests show that the overhead of re-reading these variables is very small when compared to the benefits of not needing to restart a program to refresh the variables.

#### Basic Usage

```
package main

import (
    "fmt"

    "github.com/kyani-inc/consul/env"
)

func main() {
    // e is returned with a namespace of ""
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

The above example will work both in development using the `os` package and in production using the `consul` api.

**Note**: When a client is created using `env.New()` the namespace is empty, this is true in both dev and production.


##### Running your application when using env

As mentioned previously, env utilizes build flags to switch between the dev environments. As such the following methods will need to be adapted into your build process. Note that running in dev mode is the default option:

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