// File for storing, reading, and deleting data for the app
package common

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
)

var AppDataManager = NewDataManger()

type AppData struct {
	UserId string `json:"userId"`
}

type DataManager struct {
	// Path to file storing app data
	path string
}

func NewDataManger() DataManager {
	return DataManager{
		path: filepath.Join(GetDataPath(), "data.json"),
	}
}

func (dm DataManager) FileCheck() {
	if _, err := os.Stat(dm.path); errors.Is(err, os.ErrNotExist) {
		dm.SetData(AppData{})
	}
}

func (dm DataManager) GetData() (AppData, error) {
	dm.FileCheck()

	bytes, err := os.ReadFile(dm.path)
	if err != nil {
		return AppData{}, err
	}

	var data AppData
	err = json.Unmarshal([]byte(bytes), &data)

	if err != nil {
		return AppData{}, err
	}

	return data, nil
}

func (dm DataManager) SetData(data AppData) error {
	dataBytes, err := json.Marshal(data)
	if err != nil {
		return err
	}

	os.WriteFile(dm.path, dataBytes, 0644)
	return nil
}

// - Property Methods -
func (dm DataManager) GetUserId() (string, error) {
	data, err := dm.GetData()
	if err != nil {
		return "", err
	}

	return data.UserId, nil
}

func (dm DataManager) SetUserId(id string) error {
	data, err := dm.GetData()
	if err != nil {
		return err
	}

	data.UserId = id

	return dm.SetData(data)
}
