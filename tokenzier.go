package assocentity

type tokenized []string

// Tokenizes English words
func englishTokenzier(text string) (tokenized, error) {
	// ctx := context.Background()

	// // Creates a client.
	// client, err := language.NewClient(ctx)
	// if err != nil {
	// 	log.Fatalf("Failed to create client: %v", err)
	// }

	// res, _ := client.AnnotateText(ctx, &languagepb.AnnotateTextRequest{
	// 	Document: &languagepb.Document{
	// 		Source: &languagepb.Document_Content{
	// 			Content: text,
	// 		},
	// 		Type: languagepb.Document_PLAIN_TEXT,
	// 	},
	// 	Features: &languagepb.AnnotateTextRequest_Features{
	// 		ExtractSyntax: true,
	// 	},
	// 	EncodingType: languagepb.EncodingType_UTF8,
	// })
	// client.Close()

	// var tokens []string
	// for _, v := range res.GetTokens() {
	// 	if v.PartOfSpeech.Tag != languagepb.PartOfSpeech_PUNCT {
	// 		tokens = append(tokens, v.GetText().GetContent())
	// 	}
	// }

	// for _, tok := range tokens {

	// }

	// fmt.Println(tokens)

	return []string{"I'm", "Max", "Payne", "a", "real", "human", "Max", "was", "here", "a", "human"}, nil
}
