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
	"strings"
	"testing"
)

const validManifest = `{
  "bundles": [
    {
      "name": "latest",
      "description": "desc",
      "base_image": "caddy:2-alpine",
      "modules": [
        { "id": "a", "import": "example.com/a", "version": "v1.0.0" },
        { "id": "b", "import": "example.com/b", "version": "v2.0.0" }
      ]
    }
  ]
}`

func TestParseManifest_UnknownField(t *testing.T) {
	data := `{"bundles": [{"name": "latest", "bogus": true, "modules": []}]}`
	if _, err := parseManifest([]byte(data)); err == nil {
		t.Fatal("expected error for unknown field, got nil")
	}
}

func TestParseManifest_Valid(t *testing.T) {
	m, err := parseManifest([]byte(validManifest))
	if err != nil {
		t.Fatalf("parseManifest: %v", err)
	}
	if len(m.Bundles) != 1 {
		t.Fatalf("got %d bundles, want 1", len(m.Bundles))
	}
}

func TestFindBundle_NotFound(t *testing.T) {
	m, err := parseManifest([]byte(validManifest))
	if err != nil {
		t.Fatalf("parseManifest: %v", err)
	}
	if _, err := m.findBundle("nonexistent"); err == nil {
		t.Fatal("expected error for unknown bundle, got nil")
	}
}

func TestFindBundle_Found(t *testing.T) {
	m, err := parseManifest([]byte(validManifest))
	if err != nil {
		t.Fatalf("parseManifest: %v", err)
	}
	b, err := m.findBundle("latest")
	if err != nil {
		t.Fatalf("findBundle: %v", err)
	}
	if b.Name != "latest" {
		t.Fatalf("got name %q, want latest", b.Name)
	}
}

func TestBundleValidate(t *testing.T) {
	tests := []struct {
		name    string
		bundle  bundle
		wantErr string
	}{
		{
			name:    "zero modules",
			bundle:  bundle{Name: "b", Modules: nil},
			wantErr: "zero modules",
		},
		{
			name: "empty id",
			bundle: bundle{Name: "b", Modules: []module{
				{ID: "", Import: "example.com/a", Version: "v1.0.0"},
			}},
			wantErr: "empty id",
		},
		{
			name: "empty import",
			bundle: bundle{Name: "b", Modules: []module{
				{ID: "a", Import: "", Version: "v1.0.0"},
			}},
			wantErr: "empty import",
		},
		{
			name: "empty version",
			bundle: bundle{Name: "b", Modules: []module{
				{ID: "a", Import: "example.com/a", Version: ""},
			}},
			wantErr: "empty version",
		},
		{
			name: "duplicate id",
			bundle: bundle{Name: "b", Modules: []module{
				{ID: "a", Import: "example.com/a", Version: "v1.0.0"},
				{ID: "a", Import: "example.com/a2", Version: "v2.0.0"},
			}},
			wantErr: "duplicate module id",
		},
		{
			name: "happy path",
			bundle: bundle{Name: "b", Modules: []module{
				{ID: "a", Import: "example.com/a", Version: "v1.0.0"},
				{ID: "b", Import: "example.com/b", Version: "v2.0.0"},
			}},
			wantErr: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.bundle.validate()
			if tt.wantErr == "" {
				if err != nil {
					t.Fatalf("validate() = %v, want nil", err)
				}
				return
			}
			if err == nil {
				t.Fatalf("validate() = nil, want error containing %q", tt.wantErr)
			}
			if !strings.Contains(err.Error(), tt.wantErr) {
				t.Fatalf("validate() = %v, want error containing %q", err, tt.wantErr)
			}
		})
	}
}

func TestXcaddyArgs(t *testing.T) {
	b := bundle{Modules: []module{
		{ID: "a", Import: "example.com/a", Version: "v1.0.0"},
		{ID: "b", Import: "example.com/b", Version: "v2.0.0"},
	}}
	want := "--with example.com/a@v1.0.0 --with example.com/b@v2.0.0"
	if got := b.xcaddyArgs(); got != want {
		t.Fatalf("xcaddyArgs() = %q, want %q", got, want)
	}
}

func TestIds(t *testing.T) {
	b := bundle{Modules: []module{
		{ID: "a"}, {ID: "b"}, {ID: "c"},
	}}
	got := b.ids()
	want := []string{"a", "b", "c"}
	if len(got) != len(want) {
		t.Fatalf("ids() = %v, want %v", got, want)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("ids()[%d] = %q, want %q", i, got[i], want[i])
		}
	}
}

func TestBundleNames(t *testing.T) {
	m, err := parseManifest([]byte(validManifest))
	if err != nil {
		t.Fatalf("parseManifest: %v", err)
	}
	names := m.bundleNames()
	if len(names) != 1 || names[0] != "latest" {
		t.Fatalf("bundleNames() = %v, want [latest]", names)
	}
}
