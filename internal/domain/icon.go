package domain

import "encoding/json"

// Icon contains icon information
type Icon struct {
	ID  int64  `json:"id"`
	URL string `json:"url"`
}

// CvtJSONToIcons converts JSON string to Icon list
func CvtJSONToIcons(jsonStr string) ([]Icon, error) {
	var icons []Icon
	if err := json.Unmarshal([]byte(jsonStr), &icons); err != nil {
		return nil, err
	}
	return icons, nil
}

// CvtIconsToJSON converts Icon list to JSON string
func CvtIconsToJSON(icons []Icon) (string, error) {
	iconJSON, err := json.Marshal(icons)
	if err != nil {
		return "", err
	}
	return string(iconJSON), nil
}
