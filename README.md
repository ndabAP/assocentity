# assocentity

Package assocentity returns the average distance from words to a given entity.

## Features

- Interfere at every step
- pass aliases to entity
- provides a default tokenzier

## Installation

```bash
$ go get github.com/ndabAP/assocentity/v5
```

## Prerequisites

Sign-up for a Cloud Natural Language API service account key and download the generated JSON file. This equals the `credentialsFile` at the example below. Don't commit that file.

## Usage

```go
import (
	"fmt"
	"log"

	"github.com/ndabAP/assocentity/v5/tokenize"
	"github.com/ndabAP/assocentity/v5"
)

const credentialsFile = "google_nlp_service_account.json"

func main() {
	text := "Punchinello wanted Payne? He'd see the pain."
	entities := []string{"Punchinello", "Payne"}

	// Create a NLP instance
	nlp, err := tokenize.NewNLP(credentialsFile, text, entities)
	if err != nil {
		log.Fatal(err)
	}

	dps := tokenize.NewPoSDetermer(tokenize.ANY)
	dj := tokenize.NewJoin(tokenize.Whitespace)

    	// Do calculates the average distances
	assocEntities, err := assocentity.Do(dj, nlp, entities)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(assocEntities) // map[wanted:1 ?:2 He:3 'd:4 see:5 the:6 pain:7 .:8]
}
```

## Author

[Julian Claus](https://www.julian-claus.de) and contributors.

## License

MIT
