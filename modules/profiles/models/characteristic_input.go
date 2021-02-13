package models

// SetCharacteristicInput is the schema for information needed to create a characteristic and its possible values
type SetCharacteristicInput struct {
	Name           string               `json:"name" bson:"name"`
	Type           CharacteristicType   `json:"type" bson:"type"`
	Target         CharacteristicTarget `json:"target" bson:"target"`
	PossibleValues []string             `json:"possibleValues" bson:"possibleValues"`
}

// SetCharacteristicChoiceInput is the schema for information needed to assign a characteristic to a profile
type SetCharacteristicChoiceInput struct {
	CharacteristicName string   `json:"characteristicName" bson:"characteristicName"`
	SelectedValues     []string `json:"selectedValues" bson:"selectedValues"`
}

// SetPreferenceInput is the schema for information needed to set the preferences of a patient
type SetPreferenceInput struct {
	CharacteristicName string `json:"characteristicName" bson:"characteristicName"`
	SelectedValue      string `json:"selectedValue" bson:"selectedValue"`
	Weight             int64  `json:"weight" bson:"weight"`
}
