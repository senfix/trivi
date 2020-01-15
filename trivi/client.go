package trivi

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httputil"
	"time"

	"log"

	"github.com/pkg/errors"
)

const (
	TRIVI_BANK_ACC_ID           = "489"
	TRIVI_HOST                  = "https://api.trivi.com/v2/"
	TRIVI_AUTH                  = "https://api.trivi.com/auth/token"
	TRIVI_BANK_ACCOUNTS         = TRIVI_HOST + "bankaccounts"
	TRIVI_BANK_STATEMENT_UPLOAD = TRIVI_HOST + "bankaccounts/" + TRIVI_BANK_ACC_ID + "/statements"
)

type TriviApi interface {
	Authenticate(body AuthenticationRequest) (err error, response AuthenticationResponse)
	ListBankAccounts() (err error, bankAccountsResponse BankAccountsResponse)
	GetCurrentAuth() (auth *AuthenticationResponse)
	UploadBankStatement(filename string) (err error, bankStatementResponse BankStatementResponse)
}

type triviApi struct {
	client http.Client
	auth   *AuthenticationResponse
}

func NewTriviApi() TriviApi {
	client := http.Client{
		Timeout: time.Second * 2, // Maximum of 2 secs
	}

	return &triviApi{
		client: client,
	}
}

func (t *triviApi) GetCurrentAuth() (auth *AuthenticationResponse) {
	return t.auth
}

func (t *triviApi) execute(method string, url string, body interface{}, v interface{}) (err error) {
	encoded, err := json.Marshal(body)
	if err != nil {
		return errors.Wrap(err, "cannot marshal request")
	}
	reader := bytes.NewBuffer(encoded)

	req, err := http.NewRequest(method, url, reader)
	if err != nil {
		return errors.Wrap(err, "cannot create request")
	}

	if body != nil {
		req.Header.Add("Content-Type", "application/json")
	}

	if t.auth != nil {
		req.Header.Add("Authorization", fmt.Sprintf("%v %v", t.auth.TokenType, t.auth.AccessToken))
	}

	res, err := t.client.Do(req)
	if err != nil {
		return errors.Wrap(err, "cannot execute request")
	}

	resp, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return errors.Wrap(err, "cannot read response body")
	}

	err = json.Unmarshal(resp, v)
	if err != nil {
		printRequest(req)
		printResponse(res)
		return errors.Wrap(err, "cannot unmarshal response")
	}

	return
}

func printRequest(r *http.Request) {
	data, _ := httputil.DumpRequest(r, true)
	log.Printf("======================\n%s\n", string(data))
}

func printResponse(r *http.Response) {
	data, _ := httputil.DumpResponse(r, true)
	log.Printf("======================\n%s\n", string(data))
}
