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

package nodes

import (
	"fmt"
	"strings"

	"github.com/pkg/errors"
	"sigs.k8s.io/kind/pkg/cluster/consts"
	"sigs.k8s.io/kind/pkg/container"

	"sigs.k8s.io/kind/pkg/exec"
)

// Delete deletes nodes by name / ID (see Node.String())
func Delete(nodes ...Node) error {
	if len(nodes) == 0 {
		return nil
	}
	ids := []string{}
	for _, node := range nodes {
		ids = append(ids, node.nameOrID)
	}
	cmd := exec.Command(
		container.Engine,
		append(
			[]string{
				"rm",
				"-f", // force the container to be delete now
				"-v", // delete volumes
			},
			ids...,
		)...,
	)
	return cmd.Run()
}

// List returns the list of container IDs for the kind "nodes", optionally
// filtered by container engine ps filters
// https://docs.docker.com/engine/reference/commandline/ps/#filtering
func List(filters ...string) ([]Node, error) {
	res := []Node{}
	visit := func(cluster string, node *Node) {
		res = append(res, *node)
	}
	return res, list(visit, filters...)
}

// ListByCluster returns a list of nodes by the kind cluster name
func ListByCluster(filters ...string) (map[string][]Node, error) {
	res := make(map[string][]Node)
	visit := func(cluster string, node *Node) {
		res[cluster] = append(res[cluster], *node)
	}
	return res, list(visit, filters...)
}

func list(visit func(string, *Node), filters ...string) error {
	args := []string{
		"ps",
		"-q",         // quiet output for parsing
		"-a",         // show stopped nodes
		"--no-trunc", // don't truncate
		// filter for nodes with the cluster label
		"--filter", "label=" + consts.ClusterLabelKey,
		// format to include friendly name and the cluster name
		"--format", fmt.Sprintf(`{{.Names}}\t{{.Label "%s"}}`, consts.ClusterLabelKey),
	}
	for _, filter := range filters {
		args = append(args, "--filter", filter)
	}
	cmd := exec.Command(container.Engine, args...)
	lines, err := exec.CombinedOutputLines(cmd)
	if err != nil {
		return errors.Wrap(err, "failed to list nodes")
	}
	for _, line := range lines {
		parts := strings.Split(line, "\t")
		if len(parts) != 2 {
			return fmt.Errorf("invalid output when listing nodes: %s", line)
		}
		names := strings.Split(parts[0], ",")
		cluster := parts[1]
		visit(cluster, FromID(names[0]))
	}
	return nil
}
