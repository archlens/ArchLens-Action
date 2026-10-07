// A generated module for ArchLensAction functions
//
// This module has been generated via dagger init and serves as a reference to
// basic module structure as you get started with Dagger.
//
// Two functions have been pre-created. You can modify, delete, or add to them,
// as needed. They demonstrate usage of arguments and return types using simple
// echo and grep commands. The functions can be called from the dagger CLI or
// from one of the SDKs.
//
// The first line in this comment block is a short description line and the
// rest is a long description with more detail on the module's purpose or usage,
// if appropriate. All modules should have a short description.

package main

import (
	"bytes"
	"context"
	"dagger/arch-lens-action/internal/dagger"
	"encoding/json"
	"fmt"
	"path"
)

type ArchLensAction struct {
	Src *dagger.Directory
}

func New(src *dagger.Directory) *ArchLensAction {
	return &ArchLensAction{Src: src}
}

type View struct {
	Include []string `json:"include"`
	Exclude []string `json:"exclude"`
}

type Github struct {
	Url    string `json:"url"`
	Branch string `json:"branch"`
}

type Input struct {
	Name         string          `json:"name"`
	Github       Github          `json:"github"`
	RootFolder   string          `json:"rootFolder"`
	Views        map[string]View `json:"views"`
	SaveLocation string          `json:"saveLocation"`
	RunCommand 	 string 		 `json:"runCommand"`
}

// firstViewName returns the first key of the "views" object in declaration order.
func firstViewName(raw []byte) (string, error) {
	dec := json.NewDecoder(bytes.NewReader(raw))
	if _, err := dec.Token(); err != nil { // opening {
		return "", err
	}
	for dec.More() {
		keyTok, err := dec.Token()
		if err != nil {
			return "", err
		}
		key, _ := keyTok.(string)
		if key == "views" {
			if _, err := dec.Token(); err != nil { // opening { of views
				return "", err
			}
			if !dec.More() {
				return "", fmt.Errorf("archlens.json defines no views")
			}
			nameTok, err := dec.Token()
			if err != nil {
				return "", err
			}
			name, ok := nameTok.(string)
			if !ok {
				return "", fmt.Errorf("unexpected token in views")
			}
			return name, nil
		}
		var skip json.RawMessage // not "views": skip its value
		if err := dec.Decode(&skip); err != nil {
			return "", err
		}
	}
	return "", fmt.Errorf(`"views" not found in archlens.json`)
}

// Returns a container that echoes whatever string argument is provided
func (m *ArchLensAction) Container() *dagger.Container {
	return dag.Container().
		From("debian:bookworm").
		WithExec([]string{"apt", "update"}).
		WithExec([]string{"apt", "install", "curl", "git", "tar", "-y"}).
		WithDirectory("/proj", m.Src).
		WithWorkdir("/proj").
		WithExec([]string{"sh", "-c",
			"curl -fsSL https://github.com/archlens/ArchLensGo/releases/latest/download/archlens-linux-amd64.tar.gz | tar -xz"}).
		// WithExec([]string{"sh", "-c", "tar -xvf ./archlens.tar.gz"}).
		WithExec([]string{"git", "config", "--global", "--add", "safe.directory", "/proj"})
	}


func (m *ArchLensAction) RenderDiff(ctx context.Context, baseRef, headRef string) (string, error) {
	raw, err := m.Src.File("archlens.json").Contents(ctx)
	if err != nil {
		return "", fmt.Errorf("archlens.json not found in the repository root: %w", err)
	}

	var cfg Input
	if err := json.Unmarshal([]byte(raw), &cfg); err != nil {
		return "", fmt.Errorf("parsing archlens.json: %w", err)
	}
	if cfg.SaveLocation == "" {
		return "", fmt.Errorf("archlens.json has no saveLocation")
	}

	viewName, err := firstViewName([]byte(raw))
	if err != nil {
		return "", err
	}

	ctr := m.Container().
		WithExec([]string{"./archlens", "renderDiff", baseRef, headRef})

	outPath := path.Join("/proj", cfg.SaveLocation, viewName+".md")
	return ctr.File(outPath).Contents(ctx)	
}
