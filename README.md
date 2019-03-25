# assocentity

## Features

- Calculates the distance between an entity and words
- accepts a custom tokenizer
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