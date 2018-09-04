package meander

import (
	"strings"
)

type j struct {
	Name       string
	PlaceTypes []string
}

var Journeys = []interface{}{
	&j{Name: "Romaintic", PlaceTypes: []string{"park", "bar", "movie_theatre", "restaurant"}},
	&j{Name: "Shopping", PlaceTypes: []string{"cafe", "department_store", "clothing_store", "shoe_store"}},
	&j{Name: "NightLife", PlaceTypes: []string{"bar", "casino", "food", "night_club"}},
	&j{Name: "Culture", PlaceTypes: []string{"museum", "cafe", "library", "art_gallery"}},
}

func (j *j) Public() interface{} {
	return map[string]interface{}{
		"name":    j.Name,
		"journey": strings.Join(j.PlaceTypes, "|"),
	}
}
