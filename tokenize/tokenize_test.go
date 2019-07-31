package tokenize

import (
	"flag"
	"reflect"
	"testing"
)

const credentialsFile = "../configs/google_nlp_service_account.json"

var api = flag.Bool("api", false, "call google api")

func TestNLP_Tokenize(t *testing.T) {
	flag.Parse()
	// Call API only when flag is given
	if !*api {
		t.SkipNow()
	}

	type fields struct {
		text     string
		entities []string
		punct    bool
	}
	tests := []struct {
		name    string
		fields  fields
		want    []string
		wantErr bool
	}{
		{
			name: "no punctation",
			fields: fields{
				text:     "No Payne, No Gain.",
				entities: []string{"Payne"},
				punct:    false,
			},
			want:    []string{"No", "Payne", "No", "Gain"},
			wantErr: false,
		},
		{
			name: "punctation",
			fields: fields{
				text:     "No Payne, No Gain.",
				entities: []string{"Payne"},
				punct:    true,
			},
			want:    []string{"No", "Payne", ",", "No", "Gain", "."},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			nlp, err := NewNLP(
				credentialsFile,
				tt.fields.text,
				tt.fields.entities,
				tt.fields.punct,
			)
			if err != nil {
				t.Errorf("NLP.Tokenize() error = %v", err)

				return
			}
			got, err := nlp.Tokenize()
			if (err != nil) != tt.wantErr {
				t.Errorf("NLP.Tokenize() error = %v, wantErr %v", err, tt.wantErr)

				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NLP.Tokenize() = %v, want %v", got, tt.want)
			}
		})
	}
}
