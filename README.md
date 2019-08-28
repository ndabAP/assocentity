# assocentity

Package assocentity returns the average distance from words to a given entity.

## Features

- Interfere at every step
- pass aliases to entity
- provides a default tokenzier

## Installation

```bash
$ go get github.com/ndabAP/assocentity/v6
```

## Prerequisites

Sign-up for a Cloud Natural Language API service account key and download the generated JSON file. This equals the `credentialsFile` at the example below. Don't commit that file.

## Usage

```go
import (
	"fmt"
	"log"

	"github.com/ndabAP/assocentity/v6/tokenize"
	"github.com/ndabAP/assocentity/v6"
)

const (
	credentialsFile = "google_nlp_service_account.json"
)

func main() {
	text := "Punchinello wanted Payne? He'd see the pain."
	entities := []string{"Punchinello", "Payne"}

	// Create a NLP instance
	nlp, err := tokenize.NewNLP(credentialsFile, text, entities)
	if err != nil {
		log.Fatal(err)
	}

	// Allow any part of speech
	dps := tokenize.NewPoSDetermer(tokenize.ANY)

    	// Do calculates the average distances
	assocEntities, err := assocentity.Do(nlp, dps, entities)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(assocEntities) 
	// map[tokenize.Token]float64{
	//	tokenize.Token{PoS: tokenize.VERB, Token: "wanted"}: 1,
	//	tokenize.Token{PoS: tokenize.PUNCT, Token: "?"}:     2,
	//	tokenize.Token{PoS: tokenize.PRON, Token: "He"}:     3,
	//	tokenize.Token{PoS: tokenize.VERB, Token: "'d"}:     4,
	//	tokenize.Token{PoS: tokenize.VERB, Token: "see"}:    5,
	//	tokenize.Token{PoS: tokenize.DET, Token: "the"}:     6,
	//	tokenize.Token{PoS: tokenize.NOUN, Token: "pain"}:   7,
	//	tokenize.Token{PoS: tokenize.PUNCT, Token: "."}:     8,
	// }
}
```

## Author

[Julian Claus](https://www.julian-claus.de) and contributors.

## License

MIT
