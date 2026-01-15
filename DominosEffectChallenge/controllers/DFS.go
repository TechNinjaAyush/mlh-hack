package controllers

import (
	"context"
	"fmt"
	"math/rand"
	"servicedependency/models"
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

func DFS(reverseDepend map[string][]string, healthMap map[string]float32, ch chan models.Payload, ctx context.Context) {
	defer close(ch)

	rand.Seed(time.Now().UnixNano())
	var pairs []string
	for k := range healthMap {
		pairs = append(pairs, k)
	}

	// Run for more iterations or indefinitely
	for iter := 0; iter < 15; iter++ {

		select {

		case <-ctx.Done():
			fmt.Println("DFS  stopped:client disconnected")
			return

		default:

		}
		idx := rand.Intn(len(pairs))
		pick_service := pairs[idx]

		fmt.Printf(" service picked: %s (iteration %d/50)\n", pick_service, iter+1)

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
				if len(v) > 0 {
					p := models.Payload{
						Root:        k,
						FailedNodes: v,
						Time:        time.Now(),
						BlastRadius: len(v),
					}
					fmt.Printf("📤 Sending: %s → impacted: %v (blast radius: %d)\n", k, v, len(v))
					select {
					case <-ctx.Done():
						return
					case ch <- p:
					}
				}
			}
		}

		select {
		case <-ctx.Done():
			return
		case <-time.After(4 * time.Second):
		}
	}

	fmt.Println("⏹️  DFS completed, keeping channel open...")
	time.Sleep(10 * time.Second)
	fmt.Println("📪 Channel closed")

}
