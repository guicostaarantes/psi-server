package models

// CreateSlotInput is the schema for information needed to create a new slot
type CreateSlotInput struct {
	Duration int64 `json:"duration" bson:"duration"`
	Price    int64 `json:"price" bson:"price"`
	Interval int64 `json:"interval" bson:"interval"`
}

// UpdateSlotInput is the schema for information needed to update a slot
type UpdateSlotInput struct {
	Duration int64 `json:"duration" bson:"duration"`
	Price    int64 `json:"price" bson:"price"`
	Interval int64 `json:"interval" bson:"interval"`
}
