package models

// SetCharacteristicInput is the schema for information needed to create a characteristic and its possible values
type SetCharacteristicInput struct {
	Name           string               `json:"name"`
	Type           CharacteristicType   `json:"type"`
	Target         CharacteristicTarget `json:"target"`
	PossibleValues []string             `json:"possibleValues"`
}

// SetCharacteristicChoiceInput is the schema for information needed to assign a characteristic to a profile
type SetCharacteristicChoiceInput struct {
	CharacteristicName string   `json:"characteristicName"`
	SelectedValues     []string `json:"selectedValues"`
}

// SetPreferenceInput is the schema for information needed to set the preferences of a patient
type SetPreferenceInput struct {
	CharacteristicName string `json:"characteristicName"`
	SelectedValue      string `json:"selectedValue"`
	Weight             int64  `json:"weight"`
}
