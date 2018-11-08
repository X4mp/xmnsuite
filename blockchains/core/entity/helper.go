package entity

import (
	"fmt"
	"strconv"

	uuid "github.com/satori/go.uuid"
)

func keynameByID(keyname string, id *uuid.UUID) string {
	return fmt.Sprintf("%s:by_id:%s", keyname, id.String())
}

func fetchIntFromParams(params map[string]string, keyname string, def int) int {
	value := fetchFromParams(params, keyname)
	if value == "" {
		return def
	}

	valAsInt, valAsIntErr := strconv.Atoi(value)
	if valAsIntErr != nil {
		return def
	}

	return valAsInt
}

func fetchFromParams(params map[string]string, keyname string) string {
	if value, ok := params[keyname]; ok {
		return value
	}

	return ""
}
