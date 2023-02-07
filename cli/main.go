package main

import (
	"context"
	"encoding/csv"
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"github.com/ndabAP/assocentity/v12"
	"github.com/ndabAP/assocentity/v12/nlp"
	"github.com/ndabAP/assocentity/v12/tokenize"
)

var logger = log.Default()

func init() {
	log.SetFlags(0)
	logger.SetOutput(os.Stderr)
	flag.Parse()
}

var (
	entitiesF = flag.String(
		"entities",
		"",
		"Define entities to be searched within input, example: -entities=\"Max Payne,Payne\"",
	)
	gogSvcLocF = flag.String(
		"gog-svc-loc",
		"",
		"Google Clouds NLP JSON service account file, example: -gog-svc-loc=\"~/gog-svc-loc.json\"",
	)
	opF = flag.String(
		"op",
		"mean",
		"Operation to execute",
	)
	posF = flag.String(
		"pos",
		"any",
		"Defines part of speeches to be included, example: -pos=noun,verb,pron",
	)
)

func main() {
	if len(*gogSvcLocF) == 0 {
		printHelpAndFail(errors.New("missing google service account file"))
	}

	// Read text as stdin
	textBytes, err := io.ReadAll(os.Stdin)
	if err != nil {
		printHelpAndFail(err)
	}
	if len(textBytes) == 0 {
		printHelpAndFail(errors.New("empty text"))
	}

	credentialsFilename := *gogSvcLocF
	nlpTok := nlp.NewNLPTokenizer(credentialsFilename, nlp.AutoLang)

	// Set part of speech
	posArr := strings.Split(*posF, ",")
	if len(posArr) == 0 {
		printHelpAndFail(errors.New("missing pos"))
	}
	// Parse part of speech flag and use PoS type
	poS := parsePoS(posArr)

	// Prepare text and entities
	text := string(textBytes)
	entities := strings.Split(*entitiesF, ",")
	if len(entities) == 0 {
		printHelpAndFail(errors.New("missing entities"))
	}

	// Recover to provide an unified API response
	defer func() {
		if r := recover(); r != nil {
			printHelpAndFail(r)
		}
	}()

	// Should we set a timeout?
	var ctx = context.Background()

	switch *opF {
	case "mean":
		s := assocentity.Source{
			Entities: entities,
			Texts:    []string{text},
		}
		dists, err := assocentity.Dists(ctx, nlpTok, poS, s)
		if err != nil {
			printHelpAndFail(err)
		}
		mean := assocentity.Mean(dists)

		// Write CSV to stdout
		csvwr := csv.NewWriter(os.Stdout)
		defer csvwr.Flush()
		for tok, dist := range mean {
			record := []string{
				// Text
				tok[0],
				// Part of speech
				tok[1],
				// Distance
				fmt.Sprintf("%f", dist),
			}
			if err := csvwr.Write(record); err != nil {
				printHelpAndFail(err)
			}
		}

	default:
		printHelpAndFail(errors.New("unknown operation"))
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

func printHelpAndFail(reason any) {
	logger.Println(reason)
	logger.Println()
	logger.Println("Usage:")
	logger.Println()
	flag.PrintDefaults()

	os.Exit(1)
}
