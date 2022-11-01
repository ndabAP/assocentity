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

	"github.com/ndabAP/assocentity/v9"
	"github.com/ndabAP/assocentity/v9/nlp"
	"github.com/ndabAP/assocentity/v9/tokenize"
)

var logger = log.Default()

func init() {
	log.SetFlags(0)

	logger.SetOutput(os.Stderr)
}

var poSMap = map[string]tokenize.PoS{
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

func main() {
	// Should we set a timeout?
	ctx := context.TODO()

	flag.Parse()

	textBytes, err := io.ReadAll(os.Stdin)
	if err != nil {
		printHelpAndFail(err)
	}

	// assocentity
	// Create a NLP instance
	credentialsFilename := *gogSvcLocF
	nlpTok := nlp.NewNLPTokenizer(credentialsFilename, nlp.AutoLang)

	// Set part of speech
	posArr := strings.Split(*posF, ",")
	pos := parsePoS(posArr)
	posDeterm := nlp.NewNLPPoSDetermer(pos)

	// Prepare text and entities
	text := string(textBytes)
	entities := strings.Split(*entitiesF, ",")

	// Recover to provide an unified response
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
	for token, dist := range assocEntities {
		record := []string{
			token, fmt.Sprintf("%v", dist),
		}
		w.Write(record)
	}
	w.Flush()
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
