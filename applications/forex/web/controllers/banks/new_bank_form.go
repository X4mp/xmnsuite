package banks

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"path/filepath"

	"github.com/gorilla/mux"
	"github.com/xmnservices/xmnsuite/applications/forex/objects/currency"
)

const (
	newBankFormRoute        = "/banks"
	newBankFormTemplateFile = "banks_new.html"
)

func newBankForm(rter *mux.Router, templateDir string, currencyRepository currency.Repository) *mux.Route {
	return rter.HandleFunc(newBankFormRoute, func(w http.ResponseWriter, r *http.Request) {
		// retrieve the currencies:
		currPS, currPSErr := currencyRepository.RetrieveSet(0, -1)
		if currPSErr != nil {
			w.WriteHeader(http.StatusInternalServerError)
			str := fmt.Sprintf("there was an error while retrieving Currency partial set: %s", currPSErr.Error())
			w.Write([]byte(str))
		}

		// retrieve the html page:
		content, contentErr := ioutil.ReadFile(filepath.Join(templateDir, newBankFormTemplateFile))
		if contentErr != nil {
			w.WriteHeader(http.StatusInternalServerError)
			str := fmt.Sprintf("the template could not be rendered: %s", contentErr.Error())
			w.Write([]byte(str))
		}

		tmpl, err := template.New("banks_new").Parse(string(content))
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			str := fmt.Sprintf("the template could not be rendered: %s", err.Error())
			w.Write([]byte(str))
		}

		w.WriteHeader(http.StatusOK)
		tmpl.Execute(w, currPS)
	})
}
