package trivi

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/textproto"
	"path"

	"github.com/pkg/errors"
)

func (t *triviApi) UploadBankStatement(filename string) (err error, bankStatementResponse BankStatementResponse) {
	body := &bytes.Buffer{}
	err, writer := writeData(body, filename)
	if err != nil {
		err = errors.Wrap(err, "cannot write data")
	}

	req, err := http.NewRequest(http.MethodPost, TRIVI_BANK_STATEMENT_UPLOAD, body)
	if err != nil {
		err = errors.Wrap(err, "cannot create request")
		return
	}

	req.Header.Set("Content-Type", fmt.Sprintf("multipart/related; boundary=%s", writer.Boundary()))
	req.Header.Add("Authorization", fmt.Sprintf("%v %v", t.auth.TokenType, t.auth.AccessToken))

	res, err := t.client.Do(req)
	if err != nil {
		err = errors.Wrap(err, "cannot execute request")
		return
	}

	resp, err := ioutil.ReadAll(res.Body)
	if err != nil {
		err = errors.Wrap(err, "cannot read response body")
		return
	}

	if res.StatusCode == http.StatusBadRequest {
		printRequest(req)
		printResponse(res)
		bankStatementErrorResponse := BankStatementErrorResponse{}
		err = json.Unmarshal(resp, &bankStatementErrorResponse)
		if err != nil {
			return errors.Wrap(err, "cannot unmarshal response"), bankStatementResponse
		}
		return &bankStatementErrorResponse, bankStatementResponse
	}

	err = json.Unmarshal(resp, &bankStatementResponse)
	if err != nil {
		printRequest(req)
		printResponse(res)
		return errors.Wrap(err, "cannot unmarshal response"), bankStatementResponse
	}

	return
}

func writeData(body *bytes.Buffer, filename string) (err error, writer *multipart.Writer) {
	writer = multipart.NewWriter(body)
	err = writeMetaData(writer, filename)
	if err != nil {
		err = errors.Wrap(err, "cannot write metadata")
	}

	err = writeFile(filename, writer)
	if err != nil {
		err = errors.Wrap(err, "cannot write file data")
	}

	err = writer.Close()
	if err != nil {
		err = errors.Wrap(err, "cannot close writer")
		return
	}
	return
}
func writeMetaData(writer *multipart.Writer, filename string) (err error) {
	metadataHeader := textproto.MIMEHeader{}
	metadataHeader.Set("Content-Type", "application/json; charset=UTF-8")
	part, _ := writer.CreatePart(metadataHeader)
	bankStatementRequest := BankStatementRequest{path.Base(filename)}
	encoded, err := json.Marshal(bankStatementRequest)
	if err != nil {
		return
	}
	reader := bytes.NewBuffer(encoded)
	_, err = part.Write(reader.Bytes())
	if err != nil {
		return
	}
	return
}

func writeFile(filename string, writer *multipart.Writer) (err error) { //mediaPart
	mediaData, _ := ioutil.ReadFile(filename)
	mediaHeader := textproto.MIMEHeader{}
	//mediaHeader.Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%v\".", filename))
	//mediaHeader.Set("Content-ID", "media")
	//mediaHeader.Set("Content-Filename", filename)
	mediaHeader.Set("Content-Type", "application/pdf")

	mediaPart, _ := writer.CreatePart(mediaHeader)
	_, err = io.Copy(mediaPart, bytes.NewReader(mediaData))
	if err != nil {
		err = errors.Wrap(err, "cannot copy data into stream")
		return
	}
	return
}

type BankStatementRequest struct {
	Name string `json:"name"`
}

type BankStatementResponse struct {
	ID                int    `json:"id"`
	Name              string `json:"name"`
	AuthorName        string `json:"authorName"`
	Uploaded          string `json:"uploaded"`
	BankAccountID     int    `json:"bankAccountId"`
	ProcessingState   int    `json:"processingState"`
	ProcessingMessage string `json:"processingMessage"`
}

type BankStatementErrorResponse struct {
	Code           string           `json:"code"`
	Message        string           `json:"message"`
	FieldsAffected []FieldsAffected `json:"fieldsAffected"`
}
type FieldsAffected struct {
	Name    string `json:"name"`
	Message string `json:"message"`
	Code    string `json:"code"`
}

func (t *BankStatementErrorResponse) Error() string {
	return fmt.Sprintf("%v: %v", t.Code, t.Message)
}
