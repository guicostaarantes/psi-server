package models

// TranslationInput holds necessary information to add a translated value of a language-agnostic key of a specific language
type TranslationInput struct {
	Key   string `json:"key" bson:"key"`
	Value string `json:"value" bson:"value"`
}

// Translation holds the value of a translated message referenced by a language-agnostic key
type Translation struct {
	Lang  string `json:"lang" bson:"lang"`
	Key   string `json:"key" bson:"key"`
	Value string `json:"value" bson:"value"`
}
