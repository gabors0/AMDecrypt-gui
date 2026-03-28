package app

import (
	"encoding/json"
	"os"
	"path/filepath"
	"runtime"
)

type Settings struct {
	Terminal string `json:"terminal"`
}

func (a *App) GetOS() string {
	return runtime.GOOS
}

func (a *App) GetSettings() (*Settings, error) {
	appDataDir, err := a.GetAppDataDir()
	if err != nil {
		return &Settings{}, err
	}
	data, err := os.ReadFile(filepath.Join(appDataDir, "settings.json"))
	if os.IsNotExist(err) {
		return &Settings{}, nil
	}
	if err != nil {
		return &Settings{}, err
	}
	s := &Settings{}
	if err := json.Unmarshal(data, s); err != nil {
		return &Settings{}, err
	}
	return s, nil
}

func (a *App) SaveSettings(s *Settings) error {
	appDataDir, err := a.GetAppDataDir()
	if err != nil {
		return err
	}
	data, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(filepath.Join(appDataDir, "settings.json"), data, 0644)
}
