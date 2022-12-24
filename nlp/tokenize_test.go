package nlp_test

import (
	"context"
	"log"
	"os"
	"reflect"
	"testing"

	"github.com/joho/godotenv"
	"github.com/ndabAP/assocentity/v11/nlp"
	"github.com/ndabAP/assocentity/v11/tokenize"
)

func TestTokenize(t *testing.T) {
	if testing.Short() {
		t.SkipNow()
	}

	if err := godotenv.Load("../.env"); err != nil {
		log.Fatal(err)
	}

	credentialsFile := os.Getenv("GOOGLE_NLP_SERVICE_ACCOUNT_FILE_LOCATION")

	tests := []struct {
		name    string
		text    string
		want    []tokenize.Token
		wantErr bool
	}{
		{
			name: "six tokens",
			text: "Punchinello was burning to get me",
			want: []tokenize.Token{
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
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			nlp := nlp.NewNLPTokenizer(
				credentialsFile,
				nlp.AutoLang,
			)
			got, err := nlp.Tokenize(context.Background(), tt.text)
			if (err != nil) != tt.wantErr {
				t.Errorf("NLP.tokenize() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NLP.TokenitokenizezeText() = %v, want %v", got, tt.want)
			}
		})
	}
}
