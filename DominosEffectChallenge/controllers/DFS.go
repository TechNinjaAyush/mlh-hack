package controllers

import (
	"fmt"
	"math/rand"
	"servicedependecygraph/DominosEffectChallenge/models"
	"time"
)

func DFStraversal(glitch_service string, current string, reverseDepend map[string][]string, visited_services map[string]bool, impacted_services map[string][]string) {

	visited_services[current] = true

	dependents := reverseDepend[current]

	for _, dep := range dependents {

		impacted_services[glitch_service] = append(impacted_services[glitch_service], dep)

		if !visited_services[dep] {
			DFStraversal(glitch_service, dep, reverseDepend, visited_services, impacted_services)
		}
	}

}

func DFS(reverseDepend map[string][]string, healthMap map[string]float32, ch chan models.Payload) {
	// n random ticks where we take any random service and
	//convert healthmap into slice to pick random service to reduce health

	rand.Seed(time.Now().UnixNano())
	var pairs []string
	for k := range healthMap {
		pairs = append(pairs, k)

	}

	for iter := 0; iter < 10; iter++ {

		idx := rand.Intn(len(pairs))
		pick_service := pairs[idx]

		fmt.Printf("service picked: %s\n", pick_service)

		drop := 0.2 + rand.Float32()*(0.5-0.2)
		healthMap[pick_service] -= drop

		fmt.Printf("current health: %.2f\n", healthMap[pick_service])

		if healthMap[pick_service] < 0.7 && healthMap[pick_service] > 0 {

			fmt.Printf("Applying DFS on %s\n", pick_service)

			visited_services := make(map[string]bool)
			impacted_services := make(map[string][]string)

			DFStraversal(pick_service, pick_service, reverseDepend, visited_services, impacted_services)

			alpha := float32(0.3)

			for root_node, deps := range impacted_services {
				for _, dep := range deps {
					healthMap[dep] = max(0, healthMap[dep]-alpha*(0.7-healthMap[root_node]))
				}

			}

			for k, v := range impacted_services {

				p := models.Payload{
					Root:        k,
					FailedNodes: v,
					Time:        time.Now(),
					BlastRadius: len(v),
				}

				fmt.Printf("failed service %s → impacted: %v\n and blast radius of  this service is %d", k, v, len(v))
				ch <- p

			}

		}

		time.Sleep(4 * time.Second)

	}

	close(ch)

	// DFStraversal("service-I", "service-I", reverseDepend, visited_services, impacted_services)

}
