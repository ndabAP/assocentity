# assocentity

Package assocentity returns the average distance from words to a given entity.

## Features

- Accepts a custom tokenizer
- pass aliases to entity
- provides a default tokenzier

## Installation

```bash
$ go get github.com/ndabAP/assocentity/v2
```

## API

```go
Make(text string, entities []string, tokenizer func(string) ([]string, error)) (map[string]float64, error)
```

## Usage

```go
import (
    "fmt"

    "github.com/ndabAP/assocentity/v2"
)

func main() {
    text := "The quick brown fox jumps over the lazy dog"
    res, _ := assocentity.Make(text, []string{"fox"}, nil)

    fmt.Println(res) // map[The:3 brown:1 dog:5 jumps:1 lazy:4 over:2 quick:2 the:3]
}

```

## Author

[Julian Claus](https://www.julian-claus.de) and contributors.

## License

MIT
