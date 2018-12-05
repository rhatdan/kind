/*
Copyright 2018 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

// Package container contains helpers for working with container engines
// This package has no stability guarantees whatsoever!
package container

import (
	"sigs.k8s.io/kind/pkg/exec"
)

//Engine is the container engine to use to launch containers
var Engine = ""

const engines = []string{"podman", "docker"}

func init() {
	for _, engine := range engines {
		if Engine, err := exec.LookPath(engine); err == nil {
			return
		}
	}
	Engine = "docker"
}
