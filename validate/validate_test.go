package validate_test

import (
	"os"
	"path/filepath"
	"testing"

	"gopkg.in/yaml.v3"
)

// These types mirror the gomnibus software definition schema.
// Kept local so this repo has no binary dependency on gomnibus itself.

type source struct {
	URL    string `yaml:"url"`
	Git    string `yaml:"git"`
	SHA256 string `yaml:"sha256"`
	SHA512 string `yaml:"sha512"`
	SHA1   string `yaml:"sha1"`
	MD5    string `yaml:"md5"`
}

type versionBlock struct {
	Version string  `yaml:"version"`
	Source  *source `yaml:"source,omitempty"`
}

type definition struct {
	Name           string         `yaml:"name"`
	DefaultVersion string         `yaml:"default_version"`
	License        string         `yaml:"license"`
	Source         *source        `yaml:"source,omitempty"`
	RelativePath   string         `yaml:"relative_path,omitempty"`
	Dependencies   []string       `yaml:"dependencies,omitempty"`
	Versions       []versionBlock `yaml:"versions,omitempty"`
}

const softwareDir = "../config/software"

func loadAll(t *testing.T) map[string]definition {
	t.Helper()
	entries, err := os.ReadDir(softwareDir)
	if err != nil {
		t.Fatalf("reading %s: %v", softwareDir, err)
	}
	defs := map[string]definition{}
	for _, e := range entries {
		if e.IsDir() || filepath.Ext(e.Name()) != ".yaml" {
			continue
		}
		data, err := os.ReadFile(filepath.Join(softwareDir, e.Name()))
		if err != nil {
			t.Fatalf("reading %s: %v", e.Name(), err)
		}
		var def definition
		if err := yaml.Unmarshal(data, &def); err != nil {
			t.Fatalf("parsing %s: %v", e.Name(), err)
		}
		defs[def.Name] = def
	}
	return defs
}

// TestAllDefinitionsParse verifies every YAML file is valid and has required fields.
func TestAllDefinitionsParse(t *testing.T) {
	defs := loadAll(t)
	if len(defs) == 0 {
		t.Fatal("no definitions found in " + softwareDir)
	}
	for _, def := range defs {
		t.Run(def.Name, func(t *testing.T) {
			if def.Name == "" {
				t.Error("missing name")
			}
			if def.DefaultVersion == "" {
				t.Errorf("%s: missing default_version", def.Name)
			}
			if def.License == "" {
				t.Errorf("%s: missing license", def.Name)
			}
		})
	}
}

// TestSourcesHaveChecksums verifies that every URL source has at least one checksum.
func TestSourcesHaveChecksums(t *testing.T) {
	defs := loadAll(t)
	hasChecksum := func(s *source) bool {
		return s.SHA256 != "" || s.SHA512 != "" || s.SHA1 != "" || s.MD5 != ""
	}
	for _, def := range defs {
		t.Run(def.Name, func(t *testing.T) {
			if def.Source != nil && def.Source.URL != "" && !hasChecksum(def.Source) {
				t.Errorf("%s: default source URL has no checksum", def.Name)
			}
			for _, vb := range def.Versions {
				if vb.Source != nil && vb.Source.URL != "" && !hasChecksum(vb.Source) {
					t.Errorf("%s version %s: source URL has no checksum", def.Name, vb.Version)
				}
			}
		})
	}
}

// TestDependenciesAreKnown verifies all declared dependencies exist in the library.
func TestDependenciesAreKnown(t *testing.T) {
	defs := loadAll(t)
	names := map[string]bool{}
	for name := range defs {
		names[name] = true
	}
	for _, def := range defs {
		t.Run(def.Name, func(t *testing.T) {
			for _, dep := range def.Dependencies {
				if !names[dep] {
					t.Errorf("%s: dependency %q not found in library", def.Name, dep)
				}
			}
		})
	}
}

// TestVersionBlocksHaveVersionField verifies every versions[] entry has a version string.
func TestVersionBlocksHaveVersionField(t *testing.T) {
	defs := loadAll(t)
	for _, def := range defs {
		t.Run(def.Name, func(t *testing.T) {
			for i, vb := range def.Versions {
				if vb.Version == "" {
					t.Errorf("%s: versions[%d] missing version field", def.Name, i)
				}
			}
		})
	}
}
