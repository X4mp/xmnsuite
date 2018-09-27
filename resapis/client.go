package restapis

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/url"

	"github.com/go-resty/resty"
	crypto "github.com/xmnservices/xmnsuite/crypto"
)

type client struct {
	server *url.URL
	pk     crypto.PrivateKey
}

func createClient(pk crypto.PrivateKey, server *url.URL) Client {
	out := client{
		pk:     pk,
		server: server,
	}

	return &out
}

// CreateAccount creates an account
func (app *client) CreateAccount(name string, seedwords []string) error {
	ur, urErr := url.Parse(fmt.Sprintf("%s/%s", app.server.String(), "accounts"))
	if urErr != nil {
		return urErr
	}

	jsSeedWords, jsErr := json.Marshal(seedwords)
	if jsErr != nil {
		return jsErr
	}

	formData := map[string]string{
		"name":      name,
		"seedwords": string(jsSeedWords),
	}

	token := createTokenWithFormData(resty.MethodPost, ur.RequestURI(), formData).Hash()
	resp, respErr := resty.R().SetHeader("X-Session-Token", app.pk.Sign(token).String()).SetFormData(formData).Post(ur.String())
	if respErr != nil {
		return respErr
	}

	code := resp.StatusCode()
	if code < 200 || code >= 300 {
		str := fmt.Sprintf("there was an error while executing the http request.  Returned http code: %d, Returned message: %s", code, resp.Body())
		return errors.New(str)
	}

	return nil
}

// RetrieveAccounts retrieves the accounts
func (app *client) RetrieveAccounts() ([]Account, error) {
	ur, urErr := url.Parse(fmt.Sprintf("%s/%s", app.server.String(), "accounts"))
	if urErr != nil {
		return nil, urErr
	}

	token := createTokenWithoutData(resty.MethodGet, ur.RequestURI()).Hash()
	resp, respErr := resty.R().SetHeader("X-Session-Token", app.pk.Sign(token).String()).Get(ur.String())
	if respErr != nil {
		return nil, respErr
	}

	code := resp.StatusCode()
	if code < 200 || code >= 300 {
		str := fmt.Sprintf("there was an error while executing the http request.  Returned http code: %d, Returned message: %s", code, resp.Body())
		return nil, errors.New(str)
	}

	accs := new([]*account)
	jsErr := json.Unmarshal(resp.Body(), accs)
	if jsErr != nil {
		return nil, jsErr
	}

	out := []Account{}
	for _, oneAccount := range *accs {
		out = append(out, oneAccount)
	}

	return out, nil
}

// RetrieveAccountByName retrieves the account by name
func (app *client) RetrieveAccountByName(name string) (Account, error) {
	ur, urErr := url.Parse(fmt.Sprintf("%s/%s/%s", app.server.String(), "accounts", name))
	if urErr != nil {
		return nil, urErr
	}

	token := createTokenWithoutData(resty.MethodGet, ur.RequestURI()).Hash()
	resp, respErr := resty.R().SetHeader("X-Session-Token", app.pk.Sign(token).String()).Get(ur.String())
	if respErr != nil {
		return nil, respErr
	}

	code := resp.StatusCode()
	if code < 200 || code >= 300 {
		str := fmt.Sprintf("there was an error while executing the http request.  Returned http code: %d, Returned message: %s", code, resp.Body())
		return nil, errors.New(str)
	}

	out := new(account)
	jsErr := json.Unmarshal(resp.Body(), out)
	if jsErr != nil {
		return nil, jsErr
	}

	return out, nil
}
