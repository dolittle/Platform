package start

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
)

const endpoint string = "/api/backup/start"

type payload struct {
	DumpFilename string `json:"dumpFilename"`
	Tenant       string `json:"tenant"`
	Environment  string `json:"environment"`
	EventSource  string `json:"eventSource"`
	Application  string `json:"application"`
}

func createPayload(dumpFilename string, tenant string, environment string, eventSource string, application string) payload {
	return payload{
		DumpFilename: dumpFilename,
		Tenant:       tenant,
		Environment:  environment,
		EventSource:  eventSource,
		Application:  application,
	}
}

type backupStarted struct {
	apiURL *url.URL
}

func CreateBackupStarted(backend string) (*backupStarted, error) {
	postURL, err := url.Parse(fmt.Sprintf("%s%s", backend, endpoint))
	if err != nil {
		return nil, err
	}
	return &backupStarted{
		apiURL: postURL,
	}, nil
}

func (b *backupStarted) Commit(dumpFilename string, tenant string, environment string, eventSource string, application string) error {
	requestPayload := createPayload(dumpFilename, tenant, environment, eventSource, application)
	jsonPayload, err := json.Marshal(requestPayload)
	if err != nil {
		return err
	}

	log.Printf("Committing BackupStarted event with payload %s to endpoint %s", requestPayload, b.apiURL.String())
	err = b.sendPayload(jsonPayload)
	if err != nil {
		return err
	}
	return nil
}

func (b *backupStarted) sendPayload(jsonPayload []byte) error {
	response, err := http.Post(b.apiURL.String(), "application/json", bytes.NewBuffer(jsonPayload))
	if err != nil {
		return err
	}

	log.Printf("Received response %s", response.Status)
	if response.StatusCode != http.StatusOK {
		var responseJSON map[string]interface{}
		json.NewDecoder(response.Body).Decode(&responseJSON)
		return fmt.Errorf("Received non-ok response %s with body %s", response.Status, responseJSON)
	}
	return nil
}
