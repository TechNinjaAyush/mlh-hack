package controllers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"servicedependecygraph/DominosEffectChallenge/models"

	"github.com/gin-gonic/gin"
)

type ServiceDependency struct {
	Name string `json:"name"`

	DependsOn []string `json:"depends_on"`

	Health float32 `json:"health"`
}

func JsonMarshalling(c *gin.Context) {

	c.Writer.Header().Set("Content-Type", "text/event-stream")
	c.Writer.Header().Set("Cache-Control", "no-cache")
	c.Writer.Header().Set("Connection", "keep-ali")

	ch1 := make(chan models.Payload)

	filepath := "sample/service.json"

	file, err := os.ReadFile(filepath)

	if err != nil {
		log.Fatalf("Error in opening file\n")

	}

	data := string(file)

	fmt.Printf("data is %v", data)

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

	//created reversedependecy  map

	reverseDependency := make(map[string][]string)

	for _, s := range services {

		for _, d := range s.DependsOn {

			reverseDependency[d] = append(reverseDependency[d], s.Name)

		}
	}

	for k, s := range reverseDependency {

		fmt.Printf("%s depends on this services %s\n", k, s)

	}

	revBytes, _ := json.Marshal(gin.H{
		"type": "reverse_dependency",
		"data": reverseDependency,
	})

	fmt.Fprintf(c.Writer, "data: %s\n\n", revBytes)
	c.Writer.Flush()  

	go DFS(reverseDependency, healthMap, ch1)

	for payload := range ch1 {

		bytes, _ := json.Marshal(payload)

		fmt.Fprintf(c.Writer, "data: %s\n\n", bytes)
		c.Writer.Flush()

	}

}
