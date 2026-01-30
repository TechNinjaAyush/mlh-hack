package controllers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/goccy/go-yaml"
)

type DependencyFile struct {
	Services map[string]ServiceDeps `yaml:"services"`
}

type ServiceDeps struct {
	DependsOn []string `yaml:"depends_on"`
}
type Edge struct {
	From string `json:"from"`
	To   string `json:"to"`
}

func BuildGraph(reverseDependency map[string][]string) map[string]interface{} {
	nodesMap := make(map[string]bool)
	edges := []Edge{}

	for dep, dependents := range reverseDependency {
		nodesMap[dep] = true
		for _, svc := range dependents {
			nodesMap[svc] = true
			edges = append(edges, Edge{
				From: svc,
				To:   dep,
			})
		}
	}

	nodes := []map[string]string{}
	for node := range nodesMap {
		nodes = append(nodes, map[string]string{
			"id": node,
		})
	}
	return map[string]interface{}{
		"type":  "initial_graph",
		"nodes": nodes,
		"edges": edges,
	}

}

func JsonMarshalling(c *gin.Context) {
	// Set SSE headers

	// ctx := c.Request.Context()
	c.Writer.Header().Set("Content-Type", "text/event-stream")
	c.Writer.Header().Set("Cache-Control", "no-cache")
	c.Writer.Header().Set("Connection", "keep-alive")
	c.Writer.Header().Set("X-Accel-Buffering", "no")

	// ch1 := make(chan models.Payload, 10)
	filepath := "./dependencies.yaml"

	file, err := os.ReadFile(filepath)
	if err != nil {
		log.Fatalf("Error in opening file\n")
	}

	data := string(file)
	fmt.Printf("data is %v\n", data)

	var deps DependencyFile
	if err := yaml.Unmarshal(file, &deps); err != nil {
		log.Printf("ERROR:Failed to unmarshal yaml:%v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"Error": err.Error(),
		})
		return
	}

	for service, dep := range deps.Services {
		fmt.Printf("service name is %s :\ndepends on:%v\n", service, dep.DependsOn)
	}

	reverseDependency := make(map[string][]string)

	for service, cfg := range deps.Services {
		for _, dep := range cfg.DependsOn {
			reverseDependency[dep] = append(reverseDependency[dep], service)
		}
	}

	for dep, dependents := range reverseDependency {
		fmt.Printf("%s is depended on by %v\n", dep, dependents)
	}

	graphPayload := BuildGraph(reverseDependency)
	graphBytes, _ := json.Marshal(graphPayload)
	fmt.Fprintf(c.Writer, "data: %s\n\n", graphBytes)
	c.Writer.Flush()

	fmt.Println("📊 Sent initial graph structure")

}
