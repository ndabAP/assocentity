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
	gogSvcLocF = flag.String(
		"gog-svc-loc",
		"",
		"Google Clouds NLP JSON service account file, example: -gog-svc-loc=\"~/gog-svc-loc.json\"",
	)
	opF = flag.String(
		"op",
		"mean",
		"Operation",
	)
	posF = flag.String(
		"pos",
		"any",
		"Defines part of speeches to keep, example: -pos=noun,verb,pron",
	)
	entitiesF = flag.String(
		"entities",
		"",
		"Define entities to be searched within input, example: -entities=\"Max Payne,Payne\"",
	)
)

// Should we set a timeout?
var ctx = context.Background()

func main() {
	// Read text as stdin
	textBytes, err := io.ReadAll(os.Stdin)
	if err != nil {
		printUsageAndFail(err)
	}

	// assocentity
	credentialsFilename := *gogSvcLocF
	nlpTok := nlp.NewNLPTokenizer(credentialsFilename, nlp.AutoLang)

	// Set part of speech
	posArr := strings.Split(*posF, ",")
	// Parse part of speech flag and use PoS type
	poS := parsePoS(posArr)

	// Prepare text and entities
	text := string(textBytes)
	entities := strings.Split(*entitiesF, ",")

	// Recover to provide an unified API response
	defer func() {
		if r := recover(); r != nil {
			printUsageAndFail(r)
		}
	}()

	switch *opF {
	case "mean":
		mean, err := assocentity.Mean(ctx, nlpTok, poS, text, entities)
		if err != nil {
			printUsageAndFail(err)
		}

		// Write CSV to stdout
		w := csv.NewWriter(os.Stdout)
		defer w.Flush()
		for token, dist := range mean {
			record := []string{
				// Text, part of speech, distance
				token.Text, tokenize.PoSMapIds[token.PoS], fmt.Sprintf("%v", dist),
			}
			if err := w.Write(record); err != nil {
				printUsageAndFail(err)
			}
		}
	}
}

// ["1", "3", "2", "5"] -> 11
func parsePoS(posArr []string) (pos tokenize.PoS) {
	for _, p := range posArr {
		if p, ok := tokenize.PoSMap[p]; ok {
			// Add bits
			pos += p
		}
	}
	return
}

func printUsageAndFail(reason any) {
	logger.Println(reason)
	logger.Println()

	logger.Println("Usage:")
	logger.Println()
	flag.PrintDefaults()
	os.Exit(1)
}
