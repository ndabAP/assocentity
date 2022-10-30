# assocentity

Package assocentity returns the average distance from tokens to given entities.
**Important**: If you use the provided NLP tokenizer, you can't use special
characters in entities due its nature of tokenization.

## Features

- Interfere at every step
- Pass aliases to entity
- Provides a default NLP tokenizer

## Installation

```bash
$ go get github.com/ndabAP/assocentity/v9
```

## Prerequisites

Sign-up for a Cloud Natural Language API service account key and download the
generated JSON file. This equals the `credentialsFile` at the example below.
You should never commit that file.

## Usage

```go
import (
	"context"
	"fmt"
	"log"

	"github.com/ndabAP/assocentity/v9"
	"github.com/ndabAP/assocentity/v9/nlp"
	"github.com/ndabAP/assocentity/v9/tokenize"
)

const credentialsFile = "google_nlp_service_account.json"

func main() {
	text := "Punchinello wanted Payne? He'd see the pain."
	entities := []string{"Punchinello", "Payne"}

	// Create a NLP instance
	nlpTok := nlp.NewNLPTokenizer(credentialsFile, nlp.AutoLang)

	// Allow any part of speech
	posDeterm := nlp.NewNLPPoSDetermer(tokenize.ANY)

    	// Do calculates the average distances
	ctx := context.Background()
	assocEntities, err := assocentity.Do(ctx, nlpTok, posDeterm, text, entities)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(assocEntities)
	// map[string]float64{
	// 	"wanted": 1,
	// 	"?":      2,
	// 	"He":     3,
	// 	"'d":     4,
	// 	"see":    5,
	// 	"the":    6,
	// 	"pain":   7,
	// 	".":      8,
	// }
}
```

## In-Depth

Section procedure explains the process from a non-technical perspective and API
helps to interfere the applications process.

### Procedure

The process is split into three parts. Two of them belong to the tokenization
and one calculates the average distance between words and entities.

1. **Tokenization**. Splits the tokens and assigns part of speech
2. **Part of speech determination**. Keeps only the wanted part of speeches
3. **Calculating the average**. Main function that does the actual work

#### Tokenization

Googles Cloud Natural Language API is the default tokenizer and will split the
tokens, and after that this it assigns the part of speech to the tokens. No
additional checking should be done here. For this step, it's nessecary to
sign-up for a service account key.

A simpler, offline solution would be using Gos native `strings.Fields` method.

#### Part of speech determination

It's possible to only allow certain part of speeches, e. g. only nouns and
verbs. Also the entities must stay included. Therefore, this step is separated
so it could be more optimized.

#### Calculating the average

Finally, the average distances get calculated with the given predecessors.

### API

There are two steps to interfere into the tokenization process. You just need to
implement the interfaces. `Do` takes interfaces all calls their methods. For a
non-technical explanation, read the procedure section.

#### Tokenization

Interface to implement:

```go
type Tokenizer interface {
	Tokenize(ctx context.Context, text string) ([]Token, error)
}
```

While `Token` is of type:

```go
type Token struct {
	PoS  PoS    // Part of speech
	Text string // Text
}
```

So, for example given this text:

```go
text := "Punchinello was burning to get me"
```

The result from `Tokenize` would be:

```go
[]Token{
	{
		Text: "Punchinello",
		PoS:  tokenize.NOUN,
	},
	{
		Text: "was",
		PoS:  tokenize.VERB,
	},
	{
		Text: "burning",
		PoS:  tokenize.VERB,
	},
	{
		Text: "to",
		PoS:  tokenize.PRT,
	},
	{
		Text: "get",
		PoS:  tokenize.VERB,
	},
	{
		Text: "me",
		PoS:  tokenize.PRON,
	},
}
```

#### Part of speech determination

Interface to implement:

```go
type PoSDetermer interface {
	DetermPoS(textTokens []Token, entitiesTokens [][]Token) []Token
}
```

We want to preserve the part of speech information. Therefore, we return `Token`
here instead of `string`. This makes it possible to keep the real distances
between tokens, e. g.

#### Calculating the average

This step can't be changed. It takes a `Tokenizer`, `PoSDetermer`, text as
`string` and entities in a form of `[][]string`. The method will call all the
interface methods and returns a `map` with the tokens and distances.

## Projects using assocentity

- [entityscrape](https://github.com/ndabAP/entityscrape) - Distance between word
  types (default: adjectives) in news articles and persons

## Author

[Julian Claus](https://www.julian-claus.de) and contributors.

## License

MIT
