package main

import (
	"context"
	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"github.com/ndabAP/assocentity/v11"
	"github.com/ndabAP/assocentity/v11/nlp"
	"github.com/ndabAP/assocentity/v11/tokenize"
)

var logger = log.Default()

func init() {
	log.SetFlags(0)

	logger.SetOutput(os.Stderr)

	flag.Parse()
}

var (
	poSMap = map[string]tokenize.PoS{
		"any":     tokenize.ANY,
		"adj":     tokenize.ADJ,
		"adv":     tokenize.ADV,
		"affix":   tokenize.AFFIX,
		"conj":    tokenize.CONJ,
		"det":     tokenize.DET,
		"noun":    tokenize.NOUN,
		"num":     tokenize.NUM,
		"pron":    tokenize.PRON,
		"prt":     tokenize.PRT,
		"punct":   tokenize.PUNCT,
		"unknown": tokenize.UNKN,
		"verb":    tokenize.VERB,
		"x":       tokenize.X,
	}
	poSMapIds = map[tokenize.PoS]string{
		tokenize.UNKN:  "UNKNOWN",
		tokenize.ADJ:   "ADJ",
		tokenize.ADP:   "ADP",
		tokenize.ADV:   "ADV",
		tokenize.CONJ:  "CONJ",
		tokenize.DET:   "DET",
		tokenize.NOUN:  "NOUN",
		tokenize.NUM:   "NUM",
		tokenize.PRON:  "PRON",
		tokenize.PRT:   "PRT",
		tokenize.PUNCT: "PUNCT",
		tokenize.VERB:  "VERB",
		tokenize.X:     "X",
		tokenize.AFFIX: "AFFIX",
	}
)

var (
	gogSvcLocF = flag.String(
		"gog-svc-loc",
		"",
		"Google Clouds NLP JSON service account file, example: -gog-svc-loc=\"~/gog-svc-loc.json\"",
	)
	posF = flag.String(
		"pos",
		"any",
		"Defines part of speeches to keep, example: -pos=noun,verb,pron",
	)
	entitiesF = flag.String(
		"entities",
		"",
		"Define entities to be searched within input, example: -entities=\"Max Payne, Payne\"",
	)
)

// Should we set a timeout?
var ctx = context.Background()

func main() {
	// Read text as stdin
	textBytes, err := io.ReadAll(os.Stdin)
	if err != nil {
		printHelpAndFail(err)
	}

	// assocentity
	credentialsFilename := *gogSvcLocF
	nlpTok := nlp.NewNLPTokenizer(credentialsFilename, nlp.AutoLang)

	// Set part of speech
	posArr := strings.Split(*posF, ",")
	// Parse part of speech flag and use PoS type
	pos := parsePoS(posArr)
	posDeterm := nlp.NewNLPPoSDetermer(pos)

	// Prepare text and entities
	text := string(textBytes)
	entities := strings.Split(*entitiesF, ",")

	// Recover to provide an unified API response
	defer func() {
		if r := recover(); r != nil {
			printHelpAndFail(r)
		}
	}()
	assocEntities, err := assocentity.Do(ctx, nlpTok, posDeterm, text, entities)
	if err != nil {
		printHelpAndFail(err)
	}

	// Write CSV to stdout
	w := csv.NewWriter(os.Stdout)
	defer w.Flush()
	for token, dist := range assocEntities {
		record := []string{
			// Text, part of speech, distance
			token.Text, poSMapIds[token.PoS], fmt.Sprintf("%v", dist),
		}
		if err := w.Write(record); err != nil {
			printHelpAndFail(err)
		}
	}
}

// ["1", "3", "2", "5"] -> 11
func parsePoS(posArr []string) (pos tokenize.PoS) {
	for _, p := range posArr {
		if p, ok := poSMap[p]; ok {
			// Add bits
			pos += p
		}
	}
	return
}

func printHelpAndFail(reason any) {
	logger.Println(reason)
	logger.Println()

	logger.Println("Usage:")
	logger.Println()
	flag.PrintDefaults()
	os.Exit(1)
}
