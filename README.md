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

```
Make(text string, entities []string, tokenizer Tokenizer) (map[string]float64, error)
```

## Usage

```go
import (
    "fmt"

    "github.com/ndabAP/assocentity"
)

func main() {
    text := "The quick brown fox jumps over the lazy dog"
    res, _ := assocentity.Make(text, []string{"fox"}, nil)

    fmt.Println(res)
}

```

## Author

[Julian Claus](https://www.julian-claus.de) and contributors.

## License

MIT