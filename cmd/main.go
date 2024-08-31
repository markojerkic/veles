package main

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"os/exec"

	"github.com/markojerkic/veles/internal"
)

func contains[T any](s []T, e T) bool {
	for _, a := range s {
		if &a == &e {
			return true
		}
	}
	return false
}

type BuildGraphNode struct {
	node    Build
	parents []BuildGraphNode
}

func (self *BuildGraphNode) build() {
	cmd := exec.Command(self.node.Command)
	output, err := cmd.Output()

	if err == nil {
		slog.Error("Error running command", slog.Any("error message", err))
		return
	}

	fmt.Print(output)
}

func (self *BuildGraphNode) run() {
}

type Build struct {
	Watch     []string `json:"watch"`
	Exclude   []string `json:"exclude"`
	Command   string   `json:"command"`
	Propagate bool     `json:"propagate"`
	DependsOn []string `json:"dependsOn"`
}

type Config struct {
	Builds     map[string]Build `json:"builds"`
	RunCommand string           `json:"runCommand"`

	dependencyGraph []BuildGraphNode
}

func (c *Config) Load(parse string) error {

	file, err := os.ReadFile(parse)
	if err != nil {
		return fmt.Errorf("failed to read file: %w", err)
	}

	err = json.Unmarshal(file, c)
	if err != nil {
		return fmt.Errorf("failed to unmarshal json: %w", err)
	}

	return nil
}

func findRoots(graph map[string]BuildGraphNode) []BuildGraphNode {
	roots := make([]BuildGraphNode, 0, len(graph))

	for _, build := range graph {
		if build.parents == nil {
			roots = append(roots, build)
		}

	}

	return roots
}

func (c *Config) calculateDependencyGraph() []BuildGraphNode {
	buildConfigs := make(map[string]BuildGraphNode)

	for name, build := range c.Builds {
		dependency := BuildGraphNode{
			node: build,
		}

		buildConfigs[name] = dependency
	}

	for name, build := range c.Builds {
		config := buildConfigs[name]
		for _, depName := range build.DependsOn {
			if depName == name {
				parent := buildConfigs[depName]
				parent.parents = append(parent.parents, config)
			}
		}
	}

	return findRoots(buildConfigs)

}

func (c *Config) run() {
	c.dependencyGraph = c.calculateDependencyGraph()
}

func main() {
	// configFile := flag.String("config", "build.json", "config file")
	// flag.Parse()
	//
	// var config Config

	internal.WatchFiles(".*go")

}
