package backups

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
)

const (
	api                        string = "/api/backup"
	backupStoredMethodEndpoint string = "stored"
)

type requestPayload struct {
	BackupFileName string `json:"backupFileName"`
	Tenant         string `json:"tenant"`
	Environment    string `json:"environment"`
	Application    string `json:"application"`
	ShareName      string `json:"shareName"`
}
type storedPayload struct {
	requestPayload
	DurationInSeconds uint `json:"durationInSeconds"`
}

func createPayload(b *Backups, backupFileName string) requestPayload {
	return requestPayload{
		BackupFileName: backupFileName,
		Tenant:         b.tenant,
		Environment:    b.environment,
		Application:    b.application,
		ShareName:      b.shareName,
	}
}
func createStoredPayload(b *Backups, backupFileName string, backupDurationInSeconds uint) storedPayload {
	return storedPayload{
		requestPayload:    createPayload(b, backupFileName),
		DurationInSeconds: backupDurationInSeconds,
	}
}

type Backups struct {
	apiURL      *url.URL
	tenant      string
	environment string
	application string
	shareName   string
}

func CreateBackups(
	host string,
	port int,
	tenant string,
	environment string,
	application string,
	shareName string) (*Backups, error) {
	apiURL, err := url.Parse(fmt.Sprintf("%s:%d%s", host, port, api))
	if err != nil {
		return nil, err
	}
	return &Backups{
		apiURL:      apiURL,
		tenant:      tenant,
		environment: environment,
		application: application,
		shareName:   shareName,
	}, nil
}

func (b *Backups) NotifyStored(backupFileName string, backupDurationInSeconds uint) error {
	payload := createStoredPayload(b, backupFileName, backupDurationInSeconds)
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	log.Printf("Notifying Backups microservice that backup has been successfully stored with payload %v", payload)
	err = b.sendPayload(jsonPayload, backupStoredMethodEndpoint)
	if err != nil {
		return err
	}
	return nil
}

func (b *Backups) sendPayload(jsonPayload []byte, apiMethod string) error {
	apiMethodEndpointURL, err := url.Parse(fmt.Sprintf("%s/%s", b.apiURL.String(), apiMethod))
	if err != nil {
		return err
	}
	log.Printf("Sending payload to endpoint %s\n", apiMethodEndpointURL.String())
	response, err := http.Post(apiMethodEndpointURL.String(), "application/json", bytes.NewBuffer(jsonPayload))
	if err != nil {
		if response != nil && response.StatusCode != http.StatusOK {
			var responseJSON map[string]interface{}
			json.NewDecoder(response.Body).Decode(&responseJSON)
			return fmt.Errorf("Received non-ok response %s with body %s", response.Status, responseJSON)
		}
		return err
	}
	if response.StatusCode != http.StatusOK {
		var responseJSON map[string]interface{}
		json.NewDecoder(response.Body).Decode(&responseJSON)
		return fmt.Errorf("Received non-ok response %s with body %s", response.Status, responseJSON)
	}
	return nil
}
