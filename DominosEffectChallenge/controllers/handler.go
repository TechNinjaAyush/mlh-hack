package controllers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"servicedependency/models"

	"runtime"

	"github.com/gin-gonic/gin"
)

type ServiceDependency struct {
	Name      string   `json:"name"`
	DependsOn []string `json:"depends_on"`
	Health    float32  `json:"health"`
}

func JsonMarshalling(c *gin.Context) {
	// Set SSE headers

	ctx := c.Request.Context()
	c.Writer.Header().Set("Content-Type", "text/event-stream")
	c.Writer.Header().Set("Cache-Control", "no-cache")
	c.Writer.Header().Set("Connection", "keep-alive")
	c.Writer.Header().Set("X-Accel-Buffering", "no")   
	

	ch1 := make(chan models.Payload, 10)
	filepath := "sample/service.json"

	file, err := os.ReadFile(filepath)
	if err != nil {
		log.Fatalf("Error in opening file\n")
	}

	data := string(file)
	fmt.Printf("data is %v\n", data)

	services := []ServiceDependency{}

	if err := json.Unmarshal(file, &services); err != nil {
		log.Printf("ERROR:Failed to unmarshal json:%v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"Error": err.Error(),
		})
		return
	}

	healthMap := make(map[string]float32)
	for _, s := range services {
		healthMap[s.Name] = s.Health
	}

	reverseDependency := make(map[string][]string)
	for _, s := range services {
		for _, d := range s.DependsOn {
			reverseDependency[d] = append(reverseDependency[d], s.Name)
		}
	}

	for k, s := range reverseDependency {
		fmt.Printf("%s depends on these services %v\n", k, s)
	}

	graphPayload := map[string]interface{}{
		"type": "reverse_dependency",
		"data": reverseDependency,
	}
	graphBytes, _ := json.Marshal(graphPayload)
	fmt.Fprintf(c.Writer, "data: %s\n\n", string(graphBytes))
	c.Writer.Flush()
	fmt.Println("📊 Sent initial graph structure")

	go DFS(reverseDependency, healthMap, ch1, ctx)

	n := runtime.NumGoroutine()
	fmt.Printf("total goroutines running are %d", n)

	for {

		select {

		case <-ctx.Done():
			fmt.Println("Client disconnected stop sse")
			return

		case payload, ok := <-ch1:

			if !ok {
				fmt.Println("Channel closed → stop SSE")
				return
			}

			bytes, err := json.Marshal(payload)
			if err != nil {
				continue
			}

			_, err = fmt.Fprintf(c.Writer, "data: %s\n\n", bytes)
			if err != nil {
				return
			}
			c.Writer.Flush()

		}
	}

}
