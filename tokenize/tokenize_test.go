package tokenize

import (
	"reflect"
	"testing"
)

const credentialsFile = "../configs/google_nlp_service_account.json"

func TestNLP_TokenizeText(t *testing.T) {
	if testing.Short() {
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
			got, err := nlp.TokenizeText()
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

func TestNLP_TokenizeEntities(t *testing.T) {
	type fields struct {
		text     string
		entities []string
		punct    bool
	}
	tests := []struct {
		name    string
		fields  fields
		want    [][]string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			nlp := &NLP{
				text:     tt.fields.text,
				entities: tt.fields.entities,
				punct:    tt.fields.punct,
			}
			got, err := nlp.TokenizeEntities()
			if (err != nil) != tt.wantErr {
				t.Errorf("NLP.TokenizeEntities() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NLP.TokenizeEntities() = %v, want %v", got, tt.want)
			}
		})
	}
}
