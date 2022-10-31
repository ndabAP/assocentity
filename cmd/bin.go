package main

import (
	"flag"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"time"

	"github.com/ndabAP/assocentity/v9"
	"github.com/ndabAP/assocentity/v9/nlp"
	"github.com/ndabAP/assocentity/v9/tokenize"
	"golang.org/x/net/context"
)

var logger = log.Default()

func init() {
	log.SetFlags(0)
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

const default_timeout = 60

func main() {
	// Flags
	_ = flag.Int(
		"timeout",
		60,
		"Max allowed timeout per request",
	)
	gogSvcLocF := flag.String(
		"gog-svc-loc",
		"",
		"Google Clouds NLP JSON service account file, example: -gog-svc-loc=\"~/gog-svc-loc.json\"",
	)
	posF := flag.String(
		"pos",
		"any",
		"Defines part of speeches to keep, example: -pos=noun,verb,pron",
	)
	entitiesF := flag.String(
		"entities",
		"",
		"Define entities to be searched within input, example: -entities=\"Donald Trump, Trump\"",
	)
	flag.Parse()

	// Stdin
	textBytes, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		printUsageAndPanic(err)
	}

	// assocentity
	// TODO!: Fix timeout
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*default_timeout)
	defer cancel()

	// Create a NLP instance
	credentialsFilename := *gogSvcLocF
	nlpTok := nlp.NewNLPTokenizer(credentialsFilename, nlp.AutoLang)

	// Allow any part of speech
	posArr := strings.Split(*posF, ",")
	pos := parsePoS(posArr)
	posDeterm := nlp.NewNLPPoSDetermer(pos)

	text := string(textBytes)
	entities := strings.Split(*entitiesF, ",")
	assocEntities, err := assocentity.Do(ctx, nlpTok, posDeterm, text, entities)
	if err != nil {
		printUsageAndPanic(err)
	}

	logger.Println(assocEntities)
}

func parsePoS(posArr []string) (pos tokenize.PoS) {
	for _, p := range posArr {
		if p, ok := poSMap[p]; ok {
			pos &= p
		}
	}
	return
}

func printUsageAndPanic(err error) {
	logger.SetOutput(os.Stderr)
	logger.Println(err)

	flag.PrintDefaults()
	os.Exit(1)
}
