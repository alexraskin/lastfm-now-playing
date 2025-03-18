package utils

import (
	"bytes"
	"encoding/json"
)

// CustomJSONEncoder creates a JSON encoder that doesn't escape HTML
func CustomJSONEncoder(v any) ([]byte, error) {
	buffer := &bytes.Buffer{}
	encoder := json.NewEncoder(buffer)
	encoder.SetEscapeHTML(false)
	err := encoder.Encode(v)
	if err != nil {
		return nil, err
	}
	// Remove the trailing newline that json.Encoder adds
	return bytes.TrimRight(buffer.Bytes(), "\n"), nil
}
