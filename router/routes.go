package router

import (
	"errors"
	"fmt"
	"regexp"
)

/*
 * Query route
 *
 */

type queryRoute struct {
	pattern *regexp.Regexp
	handler QueryHandlerFn
}

func createQueryRoute(patternAsString string, handler QueryHandlerFn) (QueryRoute, error) {
	pattern, patternErr := createPattern(patternAsString)
	if patternErr != nil {
		return nil, patternErr
	}

	out := queryRoute{
		pattern: pattern,
		handler: handler,
	}

	return &out, nil
}

// Pattern returns the regex pattern
func (obj *queryRoute) Pattern() *regexp.Regexp {
	return obj.pattern
}

// Handler returns the query pattern
func (obj *queryRoute) Handler() QueryHandlerFn {
	return obj.handler
}

/*
 * Trx route
 *
 */

type trxRoute struct {
	pattern *regexp.Regexp
	handler TrxHandlerFn
}

func createTrxRoute(patternAsString string, handler TrxHandlerFn) (TrxRoute, error) {
	pattern, patternErr := createPattern(patternAsString)
	if patternErr != nil {
		return nil, patternErr
	}

	out := trxRoute{
		pattern: pattern,
		handler: handler,
	}

	return &out, nil
}

// Pattern returns the regex pattern
func (obj *trxRoute) Pattern() *regexp.Regexp {
	return obj.pattern
}

// Handler returns the handler
func (obj *trxRoute) Handler() TrxHandlerFn {
	return obj.handler
}

/*
 * Helpers func
 *
 */

func createPattern(patternAsString string) (*regexp.Regexp, error) {
	pattern, patternErr := regexp.Compile(patternAsString)
	if patternErr != nil {
		str := fmt.Sprintf("there was an error while compiling the pattern (%s): %s", patternAsString, patternErr.Error())
		return nil, errors.New(str)
	}

	return pattern, nil
}
