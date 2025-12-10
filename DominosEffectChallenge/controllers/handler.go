package controllers

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"time"

	"servicedependecygraph/DominosEffectChallenge/models"
)

func BuildReverseAdjacencyGraph(services []models.ServiceGraph) map[string][]models.ServiceGraph {
	reverse := make(map[string][]models.ServiceGraph)

	for _, s := range services {
		if _, ok := reverse[s.Name]; !ok {
			reverse[s.Name] = []models.ServiceGraph{}
		}
	}

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
		impacted_services[Service_Name][s.Name] = true

		DFS_TRAVERSAL(s.Name, service, visited_services, impacted_services)

		for dep := range impacted_services[s.Name] {
			impacted_services[Service_Name][dep] = true

		}
	}
}

func ApplyGlitch(service []models.ServiceGraph) []string {
	if len(service) == 0 {
		return nil
	}

	rand.Seed(time.Now().UnixNano())

	totalGlitches := 15
	glitchedRoots := make([]string, 0, totalGlitches)

	// Build list of all service names
	names := make([]string, 0, len(service))
	for _, s := range service {
		names = append(names, s.Name)
	}

	seen := make(map[string]bool)

	for i := 0; i < totalGlitches; i++ {

		idx := rand.Intn(len(names))
		chosen := names[idx]

		min := float32(0.2)
		max := float32(0.5)
		drop := min + rand.Float32()*(max-min)

		for j := range service {
			if service[j].Name == chosen {
				service[j].Health -= drop
				if service[j].Health < 0 {
					service[j].Health = 0
				}

				if service[j].Health < 0.7 && !seen[chosen] {
					glitchedRoots = append(glitchedRoots, chosen)
					seen[chosen] = true
				}

				break
			}
		}
	}

	return glitchedRoots
}

func DFS(service map[string][]models.ServiceGraph, glitch_service []string) {
	Impacted_Services := make(map[string]map[string]bool)

	visited_services := make(map[string]bool)
	for k := range service {
		visited_services[k] = false
	}

	uniqueRoots := make([]string, 0, len(glitch_service))
	seen := make(map[string]bool)
	for _, g := range glitch_service {
		if !seen[g] {
			seen[g] = true
			uniqueRoots = append(uniqueRoots, g)
		}
	}
	for _, glitch := range uniqueRoots {
		if _, ok := service[glitch]; !ok {
			service[glitch] = []models.ServiceGraph{}
			visited_services[glitch] = false
		}

		for k := range visited_services {
			visited_services[k] = false
		}

		DFS_TRAVERSAL(glitch, service, visited_services, Impacted_Services)
	}

	for serviceName, deps := range Impacted_Services {
		time.Sleep(2 * time.Second)
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

func JsonMarhslling(w http.ResponseWriter, r *http.Request) {
	fmt.Print("started json marhaling.....\n")

	fileBytes, err := os.ReadFile("sample/service.json")
	if err != nil {
		log.Printf("failed to read the file %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	var services []models.ServiceGraph
	if err := json.Unmarshal(fileBytes, &services); err != nil {
		log.Printf("ERROR: failed to unmarshal json :%v", err)
		http.Error(w, "Internal server Error", http.StatusInternalServerError)
		return
	}

	for _, service := range services {
		fmt.Printf(" - Service: %s, Health: %.2f, Depends On: %v\n", service.Name, service.Health, service.DependsOn)
	}

	glitchService := ApplyGlitch(services)

	for _, k := range glitchService {

		fmt.Printf("glitch service is %s :\n", k)

	}

	Data := BuildReverseAdjacencyGraph(services)

	DFS(Data, glitchService)

	fmt.Println("\nReverse Dependency Graph:")
	for dep, dependents := range Data {
		fmt.Printf("\nService: %s is depended on by:\n", dep)
		for _, s := range dependents {
			fmt.Printf("   - %s (Health: %.2f, DependsOn: %v)\n", s.Name, s.Health, s.DependsOn)
		}
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(services); err != nil {
		log.Printf("ERROR: failed to encode response %v", err)
	}
}
