# assocentity

Package assocentity returns the average distance from words to a given entity.

## Features

- Interfere at every step
- pass aliases to entity
- provides a default tokenzier

## Installation

```bash
$ go get github.com/ndabAP/assocentity/v4
```

## Prerequisites

Sign-up for a Cloud Natural Language API service account key and download the generated JSON file. This equals the `credentialsFile` at the example below. You can also create your own tokenizer.

## Usage

```go
import (
	"fmt"
	"log"

	"github.com/ndabAP/assocentity/v4/tokenize"
	"github.com/ndabAP/assocentity/v4"
)

const (
	credentialsFile = "google_nlp_service_account.json"
	sep             = " "
)

func main() {
	text := "Punchinello wanted Payne? He'd see the pain."
	entities := []string{"Punchinello", "Payne"}

    	// Create a NLP instance
	nlp, err := tokenize.NewNLP(credentialsFile, text, entities)
	if err != nil {
		log.Fatal(err)
	}

    	// Join merges the entities with a simple algorithm
	dj := tokenize.NewDefaultJoin(sep)
	if err = dj.Join(nlp); err != nil {
		log.Fatal(err)
	}

    	// Assoc calculates the average distances
	assocentities, err := assocentity.Assoc(dj, nlp, entities)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(assocentities) // map[wanted:1 ?:2 He:3 'd:4 see:5 the:6 pain:7 .:8]
}
```

## Author

[Julian Claus](https://www.julian-claus.de) and contributors.

## License

MIT
