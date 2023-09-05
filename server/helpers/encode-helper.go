package helpers

import "encoding/base64"

type EncodeHelper struct{}

func (h *EncodeHelper) Encode(s string) string {
	data := base64.StdEncoding.EncodeToString(([]byte(s)))
	return string(data)
}

func (h *EncodeHelper) Decode(s string) (string, error) {
	data, err := base64.StdEncoding.DecodeString(s)
	if err != nil {
		return "", err
	}

	return string(data), nil
}
