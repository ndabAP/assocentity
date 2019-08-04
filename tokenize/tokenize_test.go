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
	}
	tests := []struct {
		name    string
		fields  fields
		want    []string
		wantErr bool
	}{
		{
			name: "punctation",
			fields: fields{
				text:     "No Payne, No Gain.",
				entities: []string{"Payne"},
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
	if testing.Short() {
		t.SkipNow()
	}

	type fields struct {
		text     string
		entities []string
	}
	tests := []struct {
		name    string
		fields  fields
		want    [][]string
		wantErr bool
	}{
		{
			name: "one entitiy",
			fields: fields{
				text:     "You're in a computer game, Max.",
				entities: []string{"Max"},
			},
			want:    [][]string{{"Max"}},
			wantErr: false,
		},
		{
			name: "two entities",
			fields: fields{
				text:     "Mona Sax. Lisa's evil twin.",
				entities: []string{"Mona Sax", "Mona"},
			},
			want:    [][]string{{"Mona", "Sax"}, {"Mona"}},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			nlp := &NLP{
				text:     tt.fields.text,
				entities: tt.fields.entities,
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
