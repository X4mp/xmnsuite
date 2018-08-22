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

// Matches returns true if the route matches, false otherwise
func (obj *queryRoute) Matches(req Request) bool {
	return true
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
	handler TrxHandlerFn
}

func createTrxRoute(handler TrxHandlerFn) (TrxRoute, error) {
	out := trxRoute{
		handler: handler,
	}

	return &out, nil
}

// Matches returns true if the route matches, false otherwise
func (obj *trxRoute) Matches(req Request) bool {
	return true
}

// Handler returns the handler
func (obj *trxRoute) Handler() TrxHandlerFn {
	return obj.handler
}

/*
 * TrxChk route
 *
 */

type trxChkRoute struct {
	handler TrxChkHandlerFn
}

func createTrxChkRoute(handler TrxChkHandlerFn) (TrxChkRoute, error) {
	out := trxChkRoute{
		handler: handler,
	}

	return &out, nil
}

// Matches returns true if the route matches, false otherwise
func (obj *trxChkRoute) Matches(req TrxChkRequest) bool {
	return true
}

// Handler returns the handler
func (obj *trxChkRoute) Handler() TrxChkHandlerFn {
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
