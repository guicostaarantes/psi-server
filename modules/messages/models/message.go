package models

// MessageInput holds necessary information to add a translated value of a language-agnostic key of a specific language
type MessageInput struct {
	Key   string `json:"key" bson:"key"`
	Value string `json:"value" bson:"value"`
}

// Message holds the value of a translated message referenced by a language-agnostic key
type Message struct {
	Lang  string `json:"lang" bson:"lang"`
	Key   string `json:"key" bson:"key"`
	Value string `json:"value" bson:"value"`
}
