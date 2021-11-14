package characteristics_models

// CharacteristicResponse is the schema for a characteristic and its possible values to be returned to the user
type CharacteristicResponse struct {
	Name           string             `json:"name"`
	Type           CharacteristicType `json:"type"`
	PossibleValues []string           `json:"possibleValues"`
}

// CharacteristicChoiceResponse is the schema for a characteristic and its possible values to be returned to the user
type CharacteristicChoiceResponse struct {
	Name           string             `json:"name"`
	Type           CharacteristicType `json:"type"`
	SelectedValues []string           `json:"selectedValues"`
	PossibleValues []string           `json:"possibleValues"`
}

// PreferenceResponse is the schema for the preferences to be returned to the user
type PreferenceResponse struct {
	CharacteristicName string `json:"characteristicName"`
	SelectedValue      string `json:"selectedValue"`
	Weight             int64  `json:"weight"`
}
