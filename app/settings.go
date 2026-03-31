package app

import (
	"encoding/json"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

const (
	settingsJSONCFile  = "settings.jsonc"
	settingsLegacyFile = "settings.json"
)

type Settings struct {
	Terminal string         `json:"terminal"`
	Bento4   Bento4Settings `json:"bento4,omitempty"`
}

type Bento4Settings struct {
	Mp4decryptPath string `json:"mp4decryptPath,omitempty"`
	BinDir         string `json:"binDir,omitempty"`
}

func (a *App) GetOS() string {
	return runtime.GOOS
}

func (a *App) GetSettings() (*Settings, error) {
	appDataDir, err := a.GetAppDataDir()
	if err != nil {
		return &Settings{}, err
	}
	paths := []string{
		filepath.Join(appDataDir, settingsJSONCFile),
		filepath.Join(appDataDir, settingsLegacyFile),
	}

	var data []byte
	for _, p := range paths {
		data, err = os.ReadFile(p)
		if err == nil {
			break
		}
		if !os.IsNotExist(err) {
			return &Settings{}, err
		}
	}

	if data == nil {
		return &Settings{}, nil
	}

	s := &Settings{}
	if err := json.Unmarshal(stripJSONCComments(data), s); err != nil {
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
	content := strings.Join([]string{
		"// settings.json can be deleted, this is the new settings file",
		"// NOTE: Bento4 removal only applies to AMDecrypt-gui-managed installs.",
		"// If mp4decrypt path differs from the recorded managed path, the uninstall is skipped.",
		"// Default Bento4 install dir: ~/.local/bin/",
		string(data),
		"",
	}, "\n")
	return os.WriteFile(filepath.Join(appDataDir, settingsJSONCFile), []byte(content), 0644)
}

func stripJSONCComments(input []byte) []byte {
	var b strings.Builder
	b.Grow(len(input))

	inString := false
	escaped := false
	inLineComment := false
	inBlockComment := false

	for i := 0; i < len(input); i++ {
		c := input[i]

		if inLineComment {
			if c == '\n' {
				inLineComment = false
				b.WriteByte(c)
			}
			continue
		}

		if inBlockComment {
			if c == '*' && i+1 < len(input) && input[i+1] == '/' {
				inBlockComment = false
				i++
			}
			continue
		}

		if inString {
			b.WriteByte(c)
			if escaped {
				escaped = false
				continue
			}
			if c == '\\' {
				escaped = true
				continue
			}
			if c == '"' {
				inString = false
			}
			continue
		}

		if c == '/' && i+1 < len(input) {
			next := input[i+1]
			if next == '/' {
				inLineComment = true
				i++
				continue
			}
			if next == '*' {
				inBlockComment = true
				i++
				continue
			}
		}

		if c == '"' {
			inString = true
		}

		b.WriteByte(c)
	}

	return []byte(b.String())
}
