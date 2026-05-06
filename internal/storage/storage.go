package storage

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"

	"github.com/bugsy6/recurlit/internal/models"
)

func configDir() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	dir := filepath.Join(home, ".config", "recurlit", "requests")
	if err := os.MkdirAll(dir, 0755); err != nil {
		return "", err
	}
	return dir, nil
}

func Save(req models.SavedRequest) error {
	dir, err := configDir()
	if err != nil {
		return err
	}
	name := sanitizeName(req.Name)
	path := filepath.Join(dir, name+".json")
	data, err := json.MarshalIndent(req, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0644)
}

func LoadAll() ([]models.SavedRequest, error) {
	dir, err := configDir()
	if err != nil {
		return nil, err
	}
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}
	var results []models.SavedRequest
	for _, e := range entries {
		if e.IsDir() || !strings.HasSuffix(e.Name(), ".json") {
			continue
		}
		path := filepath.Join(dir, e.Name())
		data, err := os.ReadFile(path)
		if err != nil {
			continue
		}
		var sr models.SavedRequest
		if err := json.Unmarshal(data, &sr); err != nil {
			continue
		}
		results = append(results, sr)
	}
	return results, nil
}

func sanitizeName(name string) string {
	r := strings.NewReplacer(
		"/", "_",
		"\\", "_",
		":", "_",
		"*", "_",
		"?", "_",
		`"`, "_",
		"<", "_",
		">", "_",
		"|", "_",
		" ", "_",
	)
	return r.Replace(name)
}
