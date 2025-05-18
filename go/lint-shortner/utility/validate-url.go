package utility

import (
	"net/url"
)

// var urlRegEx *regexp.Regexp

func init() {
	// R, _ := regexp.Compile(`/(http[s]?\/\/)(([\w\d]+)\.)*([\w]{2,3})(\/[^ ]*)?/m`)
	// urlRegEx = R
}
func ValidateUrl(urlS string) bool {
	if _, err := url.ParseRequestURI(urlS); err != nil {
		return false
	}
	return true
	// return urlRegEx.MatchString(url)
}
