package helpers

import "github.com/MarselBisengaliev/go-react-blog/models"

type TagHelper struct{}

func (h *TagHelper) MakeSliceOfInterfacesFromTags(tags []models.Tag) []interface{} {
	slice := make([]interface{}, len(tags))
	for i, v := range tags {
		slice[i] = v
	}

	return slice
}

func (h *TagHelper) ValidateTags(tags []models.Tag) error {
	var validationErr error

	if len(tags) <= 0 {
		return nil
	}

	for _, tag := range tags {
		if validationErr = validate.Struct(tag); validationErr != nil {
			return validationErr
		}
	}

	return validationErr
}
