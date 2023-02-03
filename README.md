# assocentity

Package assocentity is a social science tool to analyze the relative distance
from tokens to entities. The motiviation is to make conclusions based on the
distance from interesting tokens to a certain entity and its synonyms.

## Features

- Provide your own tokenizer
- Define aliases for entities
- Provides a default NLP tokenizer (by Google)
- Provides a multi-OS, language-agnostic CLI version

## Installation

```bash
$ go get github.com/ndabAP/assocentity/v12
```

## Prerequisites

If you want to analyze human readable texts you can use the provided Natural
Language tokenizer (powered by Google). To do so, sign-up for a Cloud Natural
Language API service account key and download the generated JSON file. This
equals the `credentialsFile` at the example below. You should never commit that
file.

A possible offline tokenizer would be a white space tokenizer. You also might
use a parser depending on your purposes.

## Example

We would like to find out which adjectives are how close in average to a certain
public person. Let's take George W. Bush and 1,000 NBC news articles as an
example. "George Bush" is the entity and synonyms are "George Walker Bush" and
"Bush" and so on. The text is each of the 1,000 NBC news articles.

Defining a text source and to set the entity would be first step. Next, we need
to instantiate our tokenizer. In this case, we use the provided Google NLP
tokenizer. Finally, we can calculate our mean distances. Since we have multiple
articles we can use `assocentity.MeanN`, which accepts multiple texts. Notice
how we pass `tokenize.ADJ` to only include adjectives as part of speech.

```go
// Define text source and entity
texts := []string{
	"Former Presidents Barack Obama, Bill Clinton and ...",
	"At the pentagon on the afternoon of 9/11, ...",
	"Tony Blair moved swiftly to place his relationship with ...",
}
entities := []string{
	"Goerge Walker Bush",
	"Goerge Bush",
	"Bush",
}

// Instantiate the NLP tokenizer (powered by Google)
nlpTok := nlp.NewNLPTokenizer(credentialsFile, nlp.AutoLang)

// Get the mean distances to adjectives
ctx := context.TODO()
meanN, err := assocentity.MeanN(ctx, nlpTok, tokenize.ADJ, texts, entities)
if err != nil {
	panic(err)
}
```

### Tokenization

If you provide your own tokenizer you must implement the interface with the
method `Tokenize` and the following signature:

```go
type Tokenizer interface {
	Tokenize(ctx context.Context, text string) ([]Token, error)
}
```

`Token` is of type:

```go
type Token struct {
	PoS  PoS    // Part of speech
	Text string // Text
}

// Part of speech
type PoS int
```

For example, given the text:

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

## CLI

There is also a terminal version available for either Windows, Mac (Darwin) or
Linux (only with 64-bit support) if you don't have Go available. The application
expects the text as stdin and accepts the following flags:

| Flag          | Description                                                                               | Type     | Default |
| ------------- | ----------------------------------------------------------------------------------------- | -------- | ------- |
| `entities`    | Define entities to be searched within input, example: `-entities="Max Payne,Payne"`       | `string` |         |
| `gog-svc-loc` | Google Clouds NLP JSON service account file, example: `-gog-svc-loc="~/gog-svc-loc.json"` | `string` |         |
| `op`          | Operation to excute: `-op="mean"`                                                         | `string` | `mean`  |
| `pos`         | Defines part of speeches to keep, example: `-pos=noun,verb,pron`                          | `string` | `any`   |

Example:

```bash
echo "Relax, Max. You're a nice guy." | ./bin/assocentity_linux_amd64_v12.0.1-0-g948274a-dirty -gog-svc-loc=/home/max/.config/assocentity/google-service.json -entities="Max Payne,Payne,Max"
```

The output is written to stdout in appropoiate formats.

## Projects using assocentity

- [entityscrape](https://github.com/ndabAP/entityscrape) - Distance between word
  types (default: adjectives) in news articles and persons

## Author

[Julian Claus](https://www.julian-claus.de) and contributors.

## License

MIT
