package messages

import (
	"encoding/json"
	"fmt"
	"log"
)

const (
	SCANSTARTED    string = "ScanStarted"
	SCANCOMPLETE   string = "ScanCompleted"
	SCANPROGRESSED string = "ScanProgressed"
	SCANFAILED     string = "ScanFailed"
	SCANCANCELLED  string = "ScanCancelled"
)

var AllTopics = []string{
	SCANSTARTED,
	SCANCOMPLETE,
	SCANPROGRESSED,
	SCANFAILED,
	SCANCANCELLED,
}

// The raw message sent from redis.
type BaseEvent struct {
	Event string          `json:"event"`
	Data  json.RawMessage `json:"data"`
}

type OsseEvent interface {
	GetType() string
}

// The individual message types from Osse

type ScanStarted struct {
	Directories int `json:"directories"`
}

func (s ScanStarted) GetType() string {
	return SCANSTARTED
}

type ScanProgressed struct {
	DirectoryName       string `json:"directoryName"`
	FilesScanned        int    `json:"filesScanned"`
	FilesSkipped        int    `json:"filesSkipped"`
	TotalDirectories    int    `json:"totalDirectories"`
	ScannedDirectories  int    `json:"scannedDirectories"`
	NextDirectoryToScan string `json:"nextDirectoryToScan"` // May be null (empty string in go)
}

func (s ScanProgressed) GetType() string {
	return SCANPROGRESSED
}

type ScanCompleted struct {
	DirectoryCount int `json:"directoryCount"`
}

func (s ScanCompleted) GetType() string {
	return SCANCOMPLETE
}

type ScanFailed struct {
	Reason string `json:"message"`
}

func (s ScanFailed) GetType() string {
	return SCANFAILED
}

type ScanCancelled struct {
	DirectoriesScannedBeforeCancellation int `json:"directoriesScannedBeforeCancellation"`
}

func (s ScanCancelled) GetType() string {
	return SCANCANCELLED
}

func GetEventFromMessage(message string) (OsseEvent, error) {
	var base BaseEvent
	err := json.Unmarshal([]byte(message), &base)
	if err != nil {
		log.Println("Error parsing event: ", err)
		return nil, err
	}

	// Now that its valid json, we determine the event type
	switch base.Event {
	case "App\\Events\\ScanStarted":
		var data ScanStarted
		err = json.Unmarshal(base.Data, &data)
		if err == nil {
			return data, nil
		}

	case "App\\Events\\ScanProgressed":
		var data ScanProgressed
		err = json.Unmarshal(base.Data, &data)
		if err == nil {
			return data, nil
		}

	case "App\\Events\\ScanCompleted":
		var data ScanCompleted
		err = json.Unmarshal(base.Data, &data)
		if err == nil {
			return data, nil
		}

	case "App\\Events\\ScanFailed":
		var data ScanFailed
		err = json.Unmarshal(base.Data, &data)
		if err == nil {
			return data, nil
		}

	case "App\\Events\\ScanCancelled":
		var data ScanCancelled
		err = json.Unmarshal(base.Data, &data)
		if err == nil {
			return data, nil
		}
	}

	return nil, fmt.Errorf("Unknown event: %s", base.Event)
}

func GetJsonOfEvent(event OsseEvent) (string, error) {
	jsonData, err := json.Marshal(event)
	if err != nil {
		log.Println("Error with converting osse event to json.")
		return "", err
	}

	return string(jsonData), nil
}
