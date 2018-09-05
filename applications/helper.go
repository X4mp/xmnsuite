package applications

import (
	"errors"
	"fmt"
	"regexp"
	"strings"

	"github.com/tendermint/tendermint/crypto"
)

/*
 * Helpers
 */

func createResourceHash(res interface{}) []byte {
	js, jsErr := cdc.MarshalJSON(res)
	if jsErr != nil {
		panic(jsErr)
	}

	return crypto.Sha256(js)
}

func isCodeValid(code int) bool {

	validCodes := []int{
		IsSuccessful,
		NotFound,
		ServerError,
		IsUnAuthorized,
		IsUnAuthenticated,
		RouteNotFound,
		InvalidRoute,
		InvalidRequest,
	}

	for _, oneValidCode := range validCodes {
		if code == oneValidCode {
			return true
		}
	}

	return false
}

func createPattern(patternAsString string) (*regexp.Regexp, error) {
	pattern, patternErr := regexp.Compile(patternAsString)
	if patternErr != nil {
		str := fmt.Sprintf("there was an error while compiling the pattern (%s): %s", patternAsString, patternErr.Error())
		return nil, errors.New(str)
	}

	return pattern, nil
}

func fromURLPatternToRegex(urlPattern string) (*regexp.Regexp, []string, error) {
	//variables:
	delimiters := "|"
	openEl := "<"
	closeEl := ">"

	//define the patterns:
	varNameWithPatternAsString := fmt.Sprintf("%s[^%s]+%s", openEl, closeEl, closeEl)
	varNamePatternAsString := "[a-z_]+"
	regexPatternAsString := fmt.Sprintf("[^%s]+", closeEl)

	// variable name with regex pattern:
	varNameWithPattern, varNameWithPatternErr := regexp.Compile(varNameWithPatternAsString)
	if varNameWithPatternErr != nil {
		return nil, nil, varNameWithPatternErr
	}

	// variable name:
	varNamePattern, varNamePatternErr := regexp.Compile(varNamePatternAsString)
	if varNamePatternErr != nil {
		return nil, nil, varNamePatternErr
	}

	//regex pattern:
	regexPattern, regexPatternErr := regexp.Compile(regexPatternAsString)
	if regexPatternErr != nil {
		return nil, nil, regexPatternErr
	}

	// find the variable name with its regex:
	updatedURLPattern := urlPattern
	varNames := []string{}
	variablesInBrackets := varNameWithPattern.FindAllString(urlPattern, -1)
	for _, oneVariableInBracket := range variablesInBrackets {

		split := strings.Split(oneVariableInBracket, delimiters)
		if len(split) != 2 {
			fmt.Printf("\n\n%s\n\n", oneVariableInBracket)
			str := fmt.Sprintf("there should only be 1 delimiter (%s) per bracket pair", delimiters)
			return nil, nil, errors.New(str)
		}

		// create the map:
		varNameAsString := varNamePattern.FindString(split[0])
		regexAsString := regexPattern.FindString(split[1])

		// add the var names in the slice:
		varNames = append(varNames, varNameAsString)

		//modify the url:
		regexAsStringWithParentheses := fmt.Sprintf("(%s)", regexAsString)
		updatedURLPattern = strings.Replace(updatedURLPattern, oneVariableInBracket, regexAsStringWithParentheses, -1)
	}

	out, outErr := createPattern(updatedURLPattern)
	if outErr != nil {
		return nil, nil, outErr
	}

	return out, varNames, outErr
}
