package models

// GetModels return the model list
func GetModels() []interface{} {
	return []interface{}{
		Fingerprint{},
		User{},
		Dataset{},
		DatasetMessage{},
		SearchHistory{},
		Request{},
		TrialRequest{},
	}
}
