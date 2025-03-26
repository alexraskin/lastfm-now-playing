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

func ExtractImageUrls(images []struct {
	Size string `xml:"size,attr"`
	Url  string `xml:",chardata"`
}) []string {
	urls := make([]string, len(images))
	for i, img := range images {
		urls[i] = img.Url
	}
	return urls
}

func GetLargestImageURL(images []struct {
	Size string `xml:"size,attr"`
	Url  string `xml:",chardata"`
}) string {
	for i := len(images) - 1; i >= 0; i-- {
		if images[i].Url != "" {
			return images[i].Url
		}
	}
	return "" // fallback if none found
}