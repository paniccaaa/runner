package extracterr

import "regexp"

func ExtractSyntaxError(errMsg string) string {
	re := regexp.MustCompile(`syntax error:.*`)
	match := re.FindString(errMsg)
	if match != "" {
		return match
	}
	return errMsg // return the original error message if no match found
}
