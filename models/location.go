package models

type AutocompleteInput struct {
	Input                string   `json:"input"`
	IncludedPrimaryTypes []string `json:"includedPrimaryTypes"`
	SessionToken         string   `json:"sessionToken"`
}

type GooglePlacePrediction struct {
	PlaceID string `json:"placeId"`
	Text    struct {
		Text string `json:"text"`
	} `json:"text"`
}

type GoogleAutocompleteResponse struct {
	Suggestions []struct {
		PlacePrediction GooglePlacePrediction `json:"placePrediction"`
	} `json:"suggestions"`
}

type UIAutocompleteSuggestion struct {
	Text    string `json:"text"`
	PlaceID string `json:"placeId"`
}

type UIAutocompleteResponse struct {
	Suggestions []UIAutocompleteSuggestion `json:"suggestions"`
}
