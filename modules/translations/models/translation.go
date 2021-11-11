package models

// TranslationInput holds necessary information to add a translated value of a language-agnostic key of a specific language
type TranslationInput struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

// Translation holds the value of a translated message referenced by a language-agnostic key
type Translation struct {
	Lang  string `json:"lang" gorm:"primaryKey"`
	Key   string `json:"key" gorm:"primaryKey"`
	Value string `json:"value"`
}
