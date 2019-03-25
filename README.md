# assocentity

## Features

- Calculates the distance between an entity and words
- accepts a custom tokenizer
- provides a default tokenzier

## Installation

```bash
$ go get github.com/ndabAP/assocentity
```

## API

```
Make(text string, entities []string, tokenizer Tokenizer) (map[string]float64, error)
```

## Usage

```go
import "github.com/ndabAP/assocentity"

func main() {
    text := "The quick brown fox jumps over the lazy dog"
    res := assocentity.Make(text, []string{"fox"}, nil)
}

```

## Author

[Julian Claus](https://www.julian-claus.de) and contributors.

## License

MIT