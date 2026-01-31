package controllers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"servicedependency/models"

	"github.com/gin-gonic/gin"
	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"
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

	ctx := c.Request.Context()
	healthMap := make(map[string]float32)
	c.Writer.Header().Set("Content-Type", "text/event-stream")
	c.Writer.Header().Set("Cache-Control", "no-cache")
	c.Writer.Header().Set("Connection", "keep-alive")
	c.Writer.Header().Set("X-Accel-Buffering", "no")

	ch1 := make(chan models.Payload, 10)
	filepath := "./sample/service.json"

	file, err := os.ReadFile(filepath)

	if err != nil {
		log.Fatalf("Error in opening file\n")
	}

	data := string(file)
	fmt.Printf("data is %v\n", data)

	var deps []models.ServiceGraph
	if err := json.Unmarshal(file, &deps); err != nil {
		log.Printf("ERROR:Failed to unmarshal json:%v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"Error": err.Error(),
		})
		return
	}

	reverseDependency := make(map[string][]string)

	for _, svc := range deps {
		for _, dep := range svc.DependsOn {
			reverseDependency[dep] = append(reverseDependency[dep], svc.Name)

		}
	}

	for _, svc := range deps {

		healthMap[svc.Name] = svc.Health

	}

	fmt.Printf("dependecy is %v", deps)

	for dep, dependents := range reverseDependency {
		fmt.Printf("%s is depended on by %v\n", dep, dependents)
	}

	graphPayload := BuildGraph(reverseDependency)
	graphBytes, _ := json.Marshal(graphPayload)
	fmt.Fprintf(c.Writer, "data: %s\n\n", graphBytes)
	c.Writer.Flush()

	fmt.Println("📊 Sent initial graph structure")

	apiKey := "AIzaSyBPv8vtdFJxApayRmP2_hNQikunnMdfE4c"
	client, err := genai.NewClient(ctx, option.WithAPIKey(apiKey))

	if err != nil {
		log.Printf("failed to create geini client:%v", err)
	}

	defer client.Close()

	go DFS(client, reverseDependency, healthMap, ch1, ctx)

	for {

		select {
		case <-ctx.Done():
			fmt.Println("Client disconnected stop sse")
			return

		case payload, ok := <-ch1:

			if !ok {

				fmt.Println("Channel closed- stop SSE")
				return
			}

			bytes, err := json.Marshal(payload)

			if err != nil {
				continue
			}

			_, err = fmt.Fprintf(c.Writer, "data:%s\n\n", bytes)

			if err != nil {
				return
			}

			c.Writer.Flush()

		}
	}

}
