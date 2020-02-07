package tokenize

import (
	"log"
	"os"
	"reflect"
	"testing"

	"github.com/joho/godotenv"
)

var credentialsFile string

func TestNLP_tokenize(t *testing.T) {
	if testing.Short() {
		t.SkipNow()
	}

	if err := godotenv.Load("../.env"); err != nil {
		log.Fatal(err)
	}

	credentialsFile = os.Getenv("GOOGLE_NLP_SERVICE_ACCOUNT_FILE_LOCATION")

	tests := []struct {
		name    string
		text    string
		want    []Token
		wantErr bool
	}{
		{
			name: "six tokens",
			text: "Punchinello was burning to get me",
			want: []Token{
				{
					Token: "Punchinello",
					PoS:   NOUN,
				},
				{
					Token: "was",
					PoS:   VERB,
				},
				{
					Token: "burning",
					PoS:   VERB,
				},
				{
					Token: "to",
					PoS:   PRT,
				},
				{
					Token: "get",
					PoS:   VERB,
				},
				{
					Token: "me",
					PoS:   PRON,
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			nlp, err := NewNLP(
				credentialsFile,
				AutoLang,
			)
			if (err != nil) != tt.wantErr {
				t.Errorf("NLP.NewNLP() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			got, err := nlp.tokenize(tt.text)
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
