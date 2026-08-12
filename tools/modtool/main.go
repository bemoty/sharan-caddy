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
	"encoding/json"
	"flag"
	"fmt"
	"os"
)

const defaultManifestPath = "modules.json"

func main() {
	if err := run(os.Args[1:], os.Stdout); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run(args []string, stdout *os.File) error {
	fs := flag.NewFlagSet("modtool", flag.ContinueOnError)
	manifestPath := fs.String("manifest", defaultManifestPath, "path to modules.json")
	if err := fs.Parse(args); err != nil {
		return err
	}

	rest := fs.Args()
	if len(rest) == 0 {
		return fmt.Errorf("usage: modtool <xcaddy-args|ids|bundles> [bundle-name]")
	}
	cmd := rest[0]

	bundleName := "latest"
	if len(rest) > 1 {
		bundleName = rest[1]
	}

	data, err := os.ReadFile(*manifestPath)
	if err != nil {
		return fmt.Errorf("read manifest: %w", err)
	}
	m, err := parseManifest(data)
	if err != nil {
		return err
	}

	if cmd == "bundles" {
		out, err := json.Marshal(m.bundleNames())
		if err != nil {
			return err
		}
		_, _ = fmt.Fprintln(stdout, string(out))
		return nil
	}

	b, err := m.findBundle(bundleName)
	if err != nil {
		return err
	}
	if err := b.validate(); err != nil {
		return err
	}

	switch cmd {
	case "xcaddy-args":
		_, _ = fmt.Fprintln(stdout, b.xcaddyArgs())
	case "ids":
		for _, id := range b.ids() {
			_, _ = fmt.Fprintln(stdout, id)
		}
	default:
		return fmt.Errorf("unknown command %q", cmd)
	}
	return nil
}
