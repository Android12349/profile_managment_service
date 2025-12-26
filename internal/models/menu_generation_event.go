package models

type MenuGenerationRequestEvent struct {
	RequestID   string              `json:"request_id"`
	UserID      int32               `json:"user_id"`
	Preferences MenuGenerationPrefs `json:"preferences"`
	Timestamp   string              `json:"timestamp"`
}

type MenuGenerationPrefs struct {
	BJU      *BJU     `json:"bju,omitempty"`
	Budget   *int32   `json:"budget,omitempty"`
	Products []string `json:"products"`
}
