package envfile

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// Profile represents a named environment profile (e.g. "dev", "staging", "prod").
type Profile struct {
	Name    string
	Entries map[string]string
}

// LoadProfile loads a .env file for the given profile name from a base directory.
// It looks for files named <dir>/<profile>.env or <dir>/.env.<profile>.
func LoadProfile(dir, profile string) (*Profile, error) {
	candidates := []string{
		filepath.Join(dir, profile+".env"),
		filepath.Join(dir, ".env."+profile),
	}

	for _, path := range candidates {
		if _, err := os.Stat(path); err == nil {
			entries, err := Parse(path)
			if err != nil {
				return nil, fmt.Errorf("profile %q: %w", profile, err)
			}
			return &Profile{Name: profile, Entries: entries}, nil
		}
	}

	return nil, fmt.Errorf("profile %q not found in %s (tried %s)", profile, dir, strings.Join(candidates, ", "))
}

// ListProfiles returns all profile names discoverable in the given directory.
func ListProfiles(dir string) ([]string, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, fmt.Errorf("list profiles: %w", err)
	}

	seen := map[string]bool{}
	var profiles []string

	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		name := e.Name()
		var profile string
		switch {
		case strings.HasSuffix(name, ".env") && name != ".env":
			profile = strings.TrimSuffix(name, ".env")
		case strings.HasPrefix(name, ".env."):
			profile = strings.TrimPrefix(name, ".env.")
		default:
			continue
		}
		if profile != "" && !seen[profile] {
			seen[profile] = true
			profiles = append(profiles, profile)
		}
	}

	return profiles, nil
}

// DiffProfiles returns the Diff between two profiles.
func DiffProfiles(a, b *Profile) []DiffEntry {
	return Diff(a.Entries, b.Entries)
}
