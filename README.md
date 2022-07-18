# Ctx Echo
Returns a map of all key-values stored inside a context. It only works with `context.Context` package.

## Installation
Using go get:
```sh
go get github.com/AmirSoleimani/ctxecho@v1.0.0
```

## Usage
It's very easy to use this pkg, You only need to call `Inspect` function.
```golang

import (
    // ...
    "github.com/AmirSoleimani/ctxecho"
    // ...
)

func MyServiceHandler(ctx context.Context) error {
    kvMap := ctxecho.Inspect(ctx)
    fmt.Printf("%+v\n", kvMap)

    // ...
    return nil
}
```


## License

This code is licensed under the MIT license.
