# consul
A collection of abstractions that Kyäni uses to ease use with consul

Links

- [Discovery](https://github.com/kyani-inc/consul/tree/master/discovery)
- [Env](#Env)

## Env

Package Env allows developers to use the [consul](https://github.com/hashicorp/consul) api as a storage for environment variables. This package attempts to use an available consul connection or silently falls back to OS Environment Variables.

**Note**: This package is designed to be simplistic enough to work with both the `os` package and the [consul](https://github.com/hashicorp/consul) api.

#### Env Installation

```
go get github.com/kyani-inc/consul/env
```

### Ideals

##### Namespaces

The way we leverage Consul is each application has their own folder structure that we call Namespaces. In normal Consul ideals a "key" is a full folder path such as `app/settings/db_user`. The way we would view a namespace is `app/settings` with the key being `db_user`.

We have chosen to break things into namespaces in order to support the fallback to OS environment variables where folders might not be supported. 

Another use case we have is to still leverage folders but doing so inside of a namespace as such:

```
namespace: app/app1
key: db1/user
key: db1/password
```

When this falls back to the OS level these are read as `db1.user` (See section below.)

###### OS Support

- Namespaces while using the OS Package are silently ignored. See Gotcha below.
- Folders are supported by converting `/` to `.`

**Example**:
When using the os you can have the following variables:

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

**Note**: When a client is created using `env.New()` the namespace is empty, this is true when using both the OS package and the Consul API.

##### Build Flags

This package has a build flag you can choose to enable that when set will notify you if the package is falling back to the OS package. This can be useful in troubleshooting if `env` is actually talking to your Consul cluster.

Usage is as follows:

```
[#] go run -tags production main.go
[#] go build -tags production .
```

##### Gotchas

- At the OS level, Namespaces are not supported, so if your application uses the same env name in two different namespaces only one will be utilized.

##### Todo

- ~~Tests~~
- Benchmarks