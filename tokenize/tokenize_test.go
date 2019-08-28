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

	type fields struct {
		text              string
		entities          []string
		tokenizedText     []Token
		tokenizedEntities [][]Token
	}
	tests := []struct {
		name    string
		fields  fields
		want    []Token
		wantErr bool
	}{
		{
			name: "six tokens",
			fields: fields{
				text:              "Punchinello was burning to get me",
				entities:          []string{"Punchinello"},
				tokenizedText:     []Token{},
				tokenizedEntities: [][]Token{{}},
			},
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
				tt.fields.text,
				tt.fields.entities,
			)
			if (err != nil) != tt.wantErr {
				t.Errorf("NLP.NewNLP() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			got, err := nlp.tokenize(tt.fields.text)
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
