package controllers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"os"
	"servicedependecygraph/DominosEffectChallenge/models"
)

func BuildReverseAdjacencyGraph(services []models.ServiceGraph) map[string][]models.ServiceGraph {

	reverse := make(map[string][]models.ServiceGraph)

	for _, s := range services {
		for _, dep := range s.DependsOn {
			reverse[dep] = append(reverse[dep], s)
		}
	}

	return reverse
}

func DFS_TRAVERSAL(Service_Name string, service map[string][]models.ServiceGraph, visited_services map[string]bool, impacted_services map[string]map[string]bool) {

	if visited_services[Service_Name] {
		return
	}
	visited_services[Service_Name] = true

	if _, exists := impacted_services[Service_Name]; !exists {
		impacted_services[Service_Name] = make(map[string]bool)
	}

	for _, s := range service[Service_Name] {
		// mark direct dependency
		impacted_services[Service_Name][s.Name] = true

		DFS_TRAVERSAL(s.Name, service, visited_services, impacted_services)

		for dep := range impacted_services[s.Name] {
			impacted_services[Service_Name][dep] = true
		}
	}
}

func DFS(service map[string][]models.ServiceGraph) {
	Impacted_Services := make(map[string]map[string]bool)
	visited_services := make(map[string]bool)

	// initialize visited map
	for k := range service {
		visited_services[k] = false
	}

	// traverse from all nodes (or a specific node)

	for k := range service {

		DFS_TRAVERSAL(k, service, visited_services, Impacted_Services)

	}

	// print results
	for serviceName, deps := range Impacted_Services {
		fmt.Printf("Failed service: %s\n", serviceName)
		fmt.Printf("Impacted services: %v\n", keys(deps))
	}
}

func keys(m map[string]bool) []string {
	result := []string{}
	for k := range m {
		result = append(result, k)
	}
	return result
}

func ApplyGlitch(service models.ServiceGraph) {

}

func JsonMarhslling(w http.ResponseWriter, r *http.Request) {

	fmt.Print("started json marhaling.....\n")

	fileBytes, err := os.ReadFile("sample/service.json")

	if err != nil {
		log.Printf("failed to read the file %v", err)

		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	var services []models.ServiceGraph

	err = json.Unmarshal(fileBytes, &services)

	if err != nil {
		log.Printf("ERROR:failed to encode response :%v", err)
		http.Error(w, "Internal server Error", http.StatusInternalServerError)

		return

	}

	// fmt.Println("unmarshalling is done succesfully...")

	w.Header().Set("Content-Type", "application/json")

	err = json.NewEncoder(w).Encode(services)

	if err != nil {
		log.Printf("ERROR: faile to encode response %v", err)
		return
	}

	for _, service := range services {
		fmt.Printf(" - Service: %s, Health: %.2f, Depends On: %v\n", service.Name, service.Health, service.DependsOn)
	}

	Data := BuildReverseAdjacencyGraph(services)

	DFS(Data)

	fmt.Println("\nReverse Dependency Graph:")
	for dep, dependents := range Data {
		fmt.Printf("\nService: %s is depended on by:\n", dep)
		for _, s := range dependents {
			fmt.Printf("   - %s (Health: %.2f, DependsOn: %v)\n", s.Name, s.Health, s.DependsOn)
		}
	}

}
