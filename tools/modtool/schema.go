// Copyright 2026 Joshua Winkler and The Sharan Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"bytes"
	"encoding/json"
	"fmt"
)

type module struct {
	ID      string `json:"id"`
	Import  string `json:"import"`
	Version string `json:"version"`
}

type bundle struct {
	Name        string   `json:"name"`
	Description string   `json:"description"`
	BaseImage   string   `json:"base_image"`
	Modules     []module `json:"modules"`
}

type manifest struct {
	Bundles []bundle `json:"bundles"`
}

func parseManifest(data []byte) (*manifest, error) {
	dec := json.NewDecoder(bytes.NewReader(data))
	dec.DisallowUnknownFields()
	var m manifest
	if err := dec.Decode(&m); err != nil {
		return nil, fmt.Errorf("parse manifest: %w", err)
	}
	return &m, nil
}

func (m *manifest) findBundle(name string) (*bundle, error) {
	for i := range m.Bundles {
		if m.Bundles[i].Name == name {
			return &m.Bundles[i], nil
		}
	}
	return nil, fmt.Errorf("bundle %q not found", name)
}

func (b *bundle) validate() error {
	if len(b.Modules) == 0 {
		return fmt.Errorf("bundle %q has zero modules", b.Name)
	}
	seen := make(map[string]bool, len(b.Modules))
	for _, mod := range b.Modules {
		if mod.ID == "" {
			return fmt.Errorf("bundle %q: module has empty id", b.Name)
		}
		if mod.Import == "" {
			return fmt.Errorf("bundle %q: module %q has empty import", b.Name, mod.ID)
		}
		if mod.Version == "" {
			return fmt.Errorf("bundle %q: module %q has empty version", b.Name, mod.ID)
		}
		if seen[mod.ID] {
			return fmt.Errorf("bundle %q: duplicate module id %q", b.Name, mod.ID)
		}
		seen[mod.ID] = true
	}
	return nil
}

func (b *bundle) xcaddyArgs() string {
	var buf bytes.Buffer
	for i, mod := range b.Modules {
		if i > 0 {
			buf.WriteByte(' ')
		}
		fmt.Fprintf(&buf, "--with %s@%s", mod.Import, mod.Version)
	}
	return buf.String()
}

func (b *bundle) ids() []string {
	ids := make([]string, len(b.Modules))
	for i, mod := range b.Modules {
		ids[i] = mod.ID
	}
	return ids
}

func (m *manifest) bundleNames() []string {
	names := make([]string, len(m.Bundles))
	for i, b := range m.Bundles {
		names[i] = b.Name
	}
	return names
}
