package start

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
	startBackupMethodEndpoint  string = "start"
	backupStoredMethodEndpoint string = "stored"
)

type requestPayload struct {
	BackupFileName  string `json:"backupFileName"`
	Tenant          string `json:"tenant"`
	Environment     string `json:"environment"`
	EventSource     string `json:"eventSource"`
	Application     string `json:"application"`
	ApplicationName string `json:"applicationName"`
	ShareName       string `json:"shareName"`
}

func createRequestPayload(backupFileName string, tenant string, environment string, eventSource string, application string, applicationName string, shareName string) requestPayload {
	return requestPayload{
		BackupFileName:  backupFileName,
		Tenant:          tenant,
		Environment:     environment,
		EventSource:     eventSource,
		Application:     application,
		ApplicationName: applicationName,
		ShareName:       shareName,
	}
}

type backup struct {
	apiURL          *url.URL
	backupFileName  string
	tenant          string
	environment     string
	eventSource     string
	application     string
	applicationName string
	shareName       string
}

func CreateBackup(
	backend string,
	backupFileName string,
	tenant string,
	environment string,
	eventSource string,
	application string,
	applicationName string,
	shareName string) (*backup, error) {
	postURL, err := url.Parse(fmt.Sprintf("%s%s", backend, api))
	if err != nil {
		return nil, err
	}
	return &backup{
		apiURL:          postURL,
		backupFileName:  backupFileName,
		tenant:          tenant,
		environment:     environment,
		eventSource:     eventSource,
		application:     application,
		applicationName: applicationName,
		shareName:       shareName,
	}, nil
}

func (b *backup) NotifyStart() error {
	payload := createRequestPayload(
		b.backupFileName,
		b.tenant,
		b.environment,
		b.eventSource,
		b.application,
		b.applicationName,
		b.shareName)
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	log.Printf("Notifying Backup microservice that backup has started with payload %s", payload)
	err = b.sendPayload(jsonPayload, startBackupMethodEndpoint)
	if err != nil {
		return err
	}
	return nil
}

func (b *backup) NotifyStored() error {
	payload := createRequestPayload(
		b.backupFileName,
		b.tenant,
		b.environment,
		b.eventSource,
		b.application,
		b.applicationName,
		b.shareName)
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	log.Printf("Notifying Backup microservice that backup has been successfully stored with payload %s", payload)
	err = b.sendPayload(jsonPayload, backupStoredMethodEndpoint)
	if err != nil {
		return err
	}
	return nil
}

func (b *backup) sendPayload(jsonPayload []byte, apiMethod string) error {
	apiMethodEndpointURL, err := url.Parse(fmt.Sprintf("%s/%s", b.apiURL.String(), apiMethod))
	if err != nil {
		return err
	}
	log.Printf("Sending payload to endpoint %s\n", apiMethodEndpointURL.String())
	response, err := http.Post(apiMethodEndpointURL.String(), "application/json", bytes.NewBuffer(jsonPayload))
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
