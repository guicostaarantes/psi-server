package models

// CharacteristicResponse is the schema for a characteristic and its possible values to be returned to the user
type CharacteristicResponse struct {
	Name           string             `json:"name" bson:"name"`
	Type           CharacteristicType `json:"type" bson:"type"`
	PossibleValues []string           `json:"possibleValues" bson:"possibleValues"`
}

// CharacteristicChoiceResponse is the schema for a characteristic and its possible values to be returned to the user
type CharacteristicChoiceResponse struct {
	Name           string             `json:"name" bson:"name"`
	Type           CharacteristicType `json:"type" bson:"type"`
	SelectedValues []string           `json:"selectedValues" bson:"selectedValues"`
	PossibleValues []string           `json:"possibleValues" bson:"possibleValues"`
}

// PreferenceResponse is the schema for the preferences to be returned to the user
type PreferenceResponse struct {
	CharacteristicName string `json:"characteristicName" bson:"characteristicName"`
	SelectedValue      string `json:"selectedValue" bson:"selectedValue"`
	Weight             int64  `json:"weight" bson:"weight"`
}
