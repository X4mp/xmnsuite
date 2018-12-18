package web

import (
	"log"
	"net/http"

	"github.com/xmnservices/xmnsuite/configs"
)

func getConfigsFromCookie(name string, r *http.Request) configs.Configs {
	defer func() {
		if r := recover(); r != nil {
			// log:
			log.Printf("the cookie could not be found or was invalid, redirect: %s", r)
		}
	}()

	// read the cookie:
	cookieConfigs, cookieConfigsErr := r.Cookie(name)
	if cookieConfigsErr != nil {
		panic(cookieConfigsErr)
	}

	log.Printf("\n in cookie: %s \n", cookieConfigs.Value)

	// convert the string to a config instance:
	conf := configs.SDKFunc.Create(configs.CreateParams{
		Encoded: cookieConfigs.Value,
	})

	return conf
}
