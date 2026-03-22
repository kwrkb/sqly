package snippet

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

type Snippet struct {
	Name  string `yaml:"name"`
	Query string `yaml:"query"`
}

func configDir() (string, error) {
	if d := os.Getenv("XDG_CONFIG_HOME"); d != "" {
		return d, nil
	}
	return os.UserConfigDir()
}

func snippetPath() (string, error) {
	dir, err := configDir()
	if err != nil {
		return "", fmt.Errorf("finding user config dir: %w", err)
	}
	return filepath.Join(dir, "asql", "snippets.yaml"), nil
}

func Load() ([]Snippet, error) {
	path, err := snippetPath()
	if err != nil {
		return nil, err
	}

	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, fmt.Errorf("reading snippets: %w", err)
	}

	if len(data) == 0 {
		return nil, nil
	}

	var snippets []Snippet
	if err := yaml.Unmarshal(data, &snippets); err != nil {
		return nil, fmt.Errorf("parsing snippets: %w", err)
	}
	return snippets, nil
}

func Save(snippets []Snippet) error {
	path, err := snippetPath()
	if err != nil {
		return err
	}

	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return fmt.Errorf("creating config dir: %w", err)
	}

	data, err := yaml.Marshal(snippets)
	if err != nil {
		return fmt.Errorf("marshaling snippets: %w", err)
	}

	if err := atomicWrite(path, data, 0o600); err != nil {
		return fmt.Errorf("writing snippets: %w", err)
	}
	return nil
}

func atomicWrite(path string, data []byte, perm fs.FileMode) error {
	tmp, err := os.CreateTemp(filepath.Dir(path), ".tmp-*")
	if err != nil {
		return err
	}
	tmpName := tmp.Name()
	defer os.Remove(tmpName)

	if _, err := tmp.Write(data); err != nil {
		tmp.Close()
		return err
	}
	if err := tmp.Chmod(perm); err != nil {
		tmp.Close()
		return err
	}
	if err := tmp.Close(); err != nil {
		return err
	}
	return os.Rename(tmpName, path)
}
