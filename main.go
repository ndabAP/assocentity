package assocentity

import "fmt"

// Make accepts a text, entities including aliases and a tokenizer which defaults to an English tokenizer.
func Make(text string, entities []string, tokenizer func(string) ([]string, error)) {
	text = "I'm Max Payne a real human. Max was here, a human."
	entities = []string{"Max Payne", "Max"}

	tokenzied, _ := englishTokenzier(text)

	fmt.Println(buildGraph(tokenzied, entities))
}
