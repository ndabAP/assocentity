# assocentity

Package assocentity returns the average distance from words to a given entity. **Important**: At the moment, it's not recommend to use special characters for entities.

## Features

- Interfere at every step
- pass aliases to entity
- provides a default tokenizer
- powered by Googles Cloud Natural Language API

## Installation

```bash
$ go get github.com/ndabAP/assocentity/v8
```

## Prerequisites

Sign-up for a Cloud Natural Language API service account key and download the generated JSON file. This equals the `credentialsFile` at the example below. You should never commit that file.

## Usage

```go
import (
	"fmt"
	"log"

	"github.com/ndabAP/assocentity/v8/tokenize"
	"github.com/ndabAP/assocentity/v8"
)

const credentialsFile = "google_nlp_service_account.json"

func main() {
	text := "Punchinello wanted Payne? He'd see the pain."
	entities := []string{"Punchinello", "Payne"}

	// Create a NLP instance
	nlp, err := tokenize.NewNLP(credentialsFile, tokenize.AutoLang)
	if err != nil {
		log.Fatal(err)
	}

	// Allow any part of speech
	psd := tokenize.NewPoSDetermer(tokenize.ANY)

    	// Do calculates the average distances
	assocEntities, err := assocentity.Do(nlp, psd, text, entities)
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

## In-Depth

Section procedure explains the process from a non-technical perspective and API helps to interfere the applications process.

### Procedure

The process is split into three parts. Two of them belong to the tokenization and one calculates the average distance between words and entities.

1. **Tokenization**. Splits the words and assigns part of speech to the token
2. **Part of speech determination**. Keeps only the wanted part of speeches
3. **Calculating the average**. Main function that does the actual work

#### Tokenization

Googles Cloud Natural Language API is the default tokenizer and will split the words and after that this library assigns the part of speech to the tokens. No additional checking should be done here. For this step, it's nessecary to sign-up for a service account key. 

A simpler, offline solution would be using Gos native `strings.Fields` method.

#### Part of speech determination

It's possible to only allow certain part of speeches, e. g. only nouns and verbs. Also the entities should stay included. Therefore, this step is separated so it could be more optimized.

#### Calculating the average

Finally, the average distances get calculated with the given predecessors.

### API

There are two steps to interfere the tokenization process. To interfere, the interfaces have to be implemented. The last takes interfaces from the other steps. For a non-technical explanation, read the procedure section.

#### Tokenization

Interface to implement:

```go
// Tokenizer tokenizes a text
type Tokenizer interface {
	Tokenize(text string) ([]Token, error)
}
```

While `Token` is of type:

```go
// Token represents a tokenized text unit
type Token struct {
	PoS   int
	Token string
}
```

So, for example given this text:

```go
text := "Punchinello was burning to get me"
```

The result from `Tokenize` could be:

```go
res := []Token{
	{
		Token: "Punchinello",
		PoS:   tokenize.NOUN,
	},
	{
		Token: "was",
		PoS:   tokenize.VERB,
	},
	{
		Token: "burning",
		PoS:   tokenize.VERB,
	},
	{
		Token: "to",
		PoS:   tokenize.PRT,
	},
	{
		Token: "get",
		PoS:   tokenize.VERB,
	},
	{
		Token: "me",
		PoS:   tokenize.PRON,
	},
}
```

#### Part of speech determination

Interface to implement:

```go
// PoSDetermer determinates if part of speech tags should be kept
type PoSDetermer interface {
	Determ(tokenizedText []Token, tokenizedEntities [][]Token) ([]Token, error)
}
```

We want to preserve the part of speech information. Therefore, we return `Token` here instead of `string`. This makes it possible to keep the real distances between tokens.

#### Calculating the average

This step can't be changed. It takes a `Tokenizer`, `PoSDetermer`, text as `string` and entities in a form of `[][]string`. The method will call all the necessary implemented methods automatically and returns a `map` with the tokens and distances.

## Projects using assocentity

- [entityscrape](https://github.com/ndabAP/entityscrape) - Distance between word types (default: adjectives) in news articles and persons

## Author

[Julian Claus](https://www.julian-claus.de) and contributors.

## License

MIT
