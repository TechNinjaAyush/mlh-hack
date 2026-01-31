package controllers

import (
	"context"
	"fmt"
	"math/rand"
	"servicedependency/models"
	"time"

	"github.com/google/generative-ai-go/genai"
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

func GenerateRCAReport(ctx context.Context, client *genai.Client, root string, impacted []string, meta map[string]models.ServiceGraph) string {
	if client == nil {
		return "AI Analysis unavailable: Client not initialized"
	}

	model := client.GenerativeModel("gemini-2.5-flash")

	// Fetch the specific metadata from your metaLookup map
	serviceInfo, exists := meta[root]
	reason := "Unknown technical failure"
	if exists {
		reason = serviceInfo.Reason
	}

	// Construct the SRE-style prompt
	prompt := fmt.Sprintf(`
        You are an SRE bot. 
        Service '%s' failed. 
        Predefined Root Cause: %s. 
        Cascading Impact: %v. 
        Provide a concise, technical RCA explaining the propagation you have to provide also the solution of impact in terms of technical aspect  in precise way you have to just pick any  dummy random users and just tell the  any dummy  number of users impacted  also the the dummy risk score.
    `, root, reason, impacted)

	resp, err := model.GenerateContent(ctx, genai.Text(prompt))

	if err != nil {
		return fmt.Sprintf("RCA Error: %s (Check: %s)", err.Error(), reason)
	}

	if len(resp.Candidates) > 0 && resp.Candidates[0].Content != nil {
		// Return the text from the first part of the first candidate
		return fmt.Sprintf("%v", resp.Candidates[0].Content.Parts[0])
	}

	return "No specific RCA generated."
}
func DFS(client *genai.Client, reverseDepend map[string][]string, healthMap map[string]float32, ch chan models.Payload, ctx context.Context) {
	defer close(ch)

	rand.Seed(time.Now().UnixNano())
	var pairs []string

	metaLookup := make(map[string]models.ServiceGraph)
	for k := range healthMap {
		pairs = append(pairs, k)
	}

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
					rcaReport := GenerateRCAReport(ctx, client, k, v, metaLookup)

					p := models.Payload{

						Root:        k,
						FailedNodes: v,
						Time:        time.Now(),
						BlastRadius: len(v),
						RCA:         rcaReport,
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
