package banks

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"path/filepath"

	"github.com/gorilla/mux"
)

const (
	showBanksRoute        = "/banks"
	showBanksTemplateFile = "banks.html"
)

func showBanks(rter *mux.Router, templateDir string) *mux.Route {
	return rter.HandleFunc(showBanksRoute, func(w http.ResponseWriter, r *http.Request) {
		// retrieve the html page:
		content, contentErr := ioutil.ReadFile(filepath.Join(templateDir, showBanksTemplateFile))
		if contentErr != nil {
			w.WriteHeader(http.StatusInternalServerError)
			str := fmt.Sprintf("the template could not be rendered: %s", contentErr.Error())
			w.Write([]byte(str))
		}

		w.WriteHeader(http.StatusOK)
		w.Write(content)
	})
}
