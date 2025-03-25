package searchValidationUtil

import "regexp"

var searchRegex = regexp.MustCompile(`^[a-zA-Z0-9]+$`)

func IsValidSearch(query string) bool {
	return searchRegex.MatchString(query)
}
