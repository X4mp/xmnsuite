package roles

import (
	"errors"
	"fmt"
	"regexp"

	crypto "github.com/xmnservices/xmnsuite/crypto"
	"github.com/xmnservices/xmnsuite/helpers"
	"github.com/xmnservices/xmnsuite/lists"
)

type concreteRoles struct {
	Lst lists.Lists
}

func createConcreteRoles() Roles {
	out := concreteRoles{
		Lst: lists.SDKFunc.CreateSet(),
	}

	return &out
}

// Lists returns the lists
func (app *concreteRoles) Lists() lists.Lists {
	return app.Lst
}

// Copy copies the roles instance
func (app *concreteRoles) Copy() Roles {
	out := concreteRoles{
		Lst: app.Lst.Copy(),
	}

	return &out
}

// Add adds users to a role key and returns the amount of users in that role
func (app *concreteRoles) Add(key string, usrs ...crypto.PublicKey) int {
	lst := app.convertUsers(usrs)
	return app.Lists().Add(key, lst...)
}

// Del deletes users from a role and returns the amount of users deleted
func (app *concreteRoles) Del(key string, usrs ...crypto.PublicKey) int {
	lst := app.convertUsers(usrs)
	return app.Lists().Del(key, lst...)
}

// EnableWriteAccess enables the write access on keys, on a role
func (app *concreteRoles) EnableWriteAccess(key string, keyPatterns ...string) int {
	writeAccessKey := app.writeKey(key)
	lst := []interface{}{}
	for _, onePattern := range keyPatterns {
		_, regErr := regexp.Compile(onePattern)
		if regErr == nil {
			lst = append(lst, onePattern)
		}
	}

	return app.Lists().Add(writeAccessKey, lst...)
}

// DisableWriteAccess disables the write access on keys, on a role
func (app *concreteRoles) DisableWriteAccess(key string, keyPatterns ...string) int {
	writeAccessKey := app.writeKey(key)
	lst := app.convertStrings(keyPatterns)
	return app.Lists().Del(writeAccessKey, lst...)
}

// HasWriteAccess returns the keys we have write access on
func (app *concreteRoles) HasWriteAccess(key string, keys ...string) []string {
	out := []string{}
	writeAccessKey := app.writeKey(key)
	patterns := app.Lists().Retrieve(writeAccessKey, 0, -1)
	for _, onePattern := range patterns {
		reg, regErr := regexp.Compile(onePattern.(string))
		if regErr != nil {
			str := fmt.Sprintf("there was an error while compiling the pattern (%s) stored on role key (%s): %s", onePattern, key, regErr.Error())
			panic(errors.New(str))
		}

		for _, oneKey := range keys {
			found := reg.FindString(oneKey)
			if found != oneKey {
				continue
			}

			out = append(out, oneKey)
		}
	}

	converted := app.convertStrings(out)
	unique := helpers.MakeUnique(converted...)

	uniqueOut := []string{}
	for _, oneUnique := range unique {
		uniqueOut = append(uniqueOut, oneUnique.(string))
	}

	return uniqueOut

}

func (app *concreteRoles) convertUsers(usrs []crypto.PublicKey) []interface{} {
	lst := []interface{}{}
	for _, oneUser := range usrs {
		lst = append(lst, oneUser.String())
	}

	return lst
}

func (app *concreteRoles) convertStrings(strs []string) []interface{} {
	lst := []interface{}{}
	for _, oneStr := range strs {
		lst = append(lst, oneStr)
	}

	return lst
}

func (app *concreteRoles) writeKey(key string) string {
	return fmt.Sprintf("%s:write-access", key)
}
