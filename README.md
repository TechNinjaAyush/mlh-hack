# 📁 cascade failure
### *Neural-Engineered Incident Response & Dependency Mapping*

**Basalt Sentinel** is an advanced SRE (Site Reliability Engineering) observability framework designed to visualize microservice dependencies and diagnose system-wide collapses in real-time. By combining a **high-concurrency Golang graph engine** with **Google Gemini 1.5 Flash**, it provides automated, sub-second Root Cause Analysis (RCA) for cascading failures.

---

## 🏗 System Architecture

The platform operates on a **Detection → Traversal → Intelligence** pipeline, ensuring that "Alert Fatigue" is replaced by actionable forensic data.



### 1. Ingestion & Simulation Layer (Go)
The backend maintains a persistent state of the system's topology using an adjacency list. When a node failure is detected, the engine triggers a **concurrency-safe DFS (Depth First Search)** to determine the impact.

### 2. Neural Forensics Layer (Gemini AI)
Instead of static log parsing, the Sentinel sends a "Topology Context" to Gemini 1.5 Flash. This allows the AI to understand the *relationship* between services, identifying "Patient Zero" rather than just reporting symptoms.

### 3. Visual Command Center (React)
A D3-driven force-directed graph provides a 60FPS real-time view of the "Blast Radius," using physics-based motion to highlight system stress.

---

## 🧠 Technical Deep-Dive

### Recursive Blast Radius Calculation
The engine identifies the set of affected services $S$ using a directed dependency graph $G = (V, E)$. For any failing "Patient Zero" node $v$, the blast radius is defined as:

$$S = \{u \in V \mid \exists \text{ path from } v \text{ to } u \text{ in } G\}$$

The Go backend processes this traversal in $O(V + E)$ time, ensuring that even with 1,000+ nodes, the impact is identified in microseconds.



### Forensic Payload Schema
The Sentinel streams a unified payload via Server-Sent Events (SSE) to the frontend. The Go struct is mapped as follows:

```go
type Payload struct {
    IncidentID  string    `json:"incidentID"`  // Unique event identifier
    Root        string    `json:"Root"`        // The originating failure node
    FailedNodes []string  `json:"FailedNodes"` // Array of impacted dependents
    Time        time.Time `json:"Time"`        // RFC3339 Timestamp
    BlastRadius int       `json:"BlastRadius"` // Total count of affected nodes
    RCA         string    `json:"rca"`         // Neural Forensic Analysis
}
🛠 Advanced Features

    Atomic State Management: React-based health maps ensure that UI updates are memoized, preventing expensive re-renders during high-frequency failure events.

    Physics-Based Visualization: Utilizes D3 force-simulation to physically shift the graph layout and apply "gravity" to the root cause node during failure.

    Contextual Remediation: Gemini analysis includes specific mitigation steps, such as suggested shell commands or circuit-breaker adjustments.

    Backpressure-Safe SSE: The Go server implements channel-based buffering to handle high-frequency event bursts without overwhelming the client main thread.

    Heuristic Decay: Automatic "health recovery" visual markers that clear after 8 seconds of system stability.

📊 Technical Stack
Component	Technology	Rationale
Backend	Go 1.21+	Optimized for high-concurrency goroutines and memory-safe graph traversal.
AI Engine	Gemini 1.5 Flash	Lowest latency-to-context ratio for real-time forensic generation.
Streaming	SSE (Server-Sent Events)	Native browser support for unidirectional, low-overhead telemetry streams.
Visualization	React-Force-Graph	Canvas-based rendering capable of handling 500+ node meshes at 60FPS.
UI/UX	Tailwind CSS	Cyberpunk-inspired dark mode (Basalt theme) for high-contrast SRE visibility.
🔧 Installation
1. Prerequisites

    Go (Golang) 1.21 or higher.

    Node.js 18.x and npm/yarn.

    Google AI (Gemini) API Key.

2. Backend Setup
Bash

# Clone the repository
git clone [https://github.com/your-username/basalt-sentinel.git](https://github.com/your-username/basalt-sentinel.git)
cd basalt-sentinel/backend

# Set your Gemini API Key
export GEMINI_API_KEY='your_api_key_here'

# Install Go dependencies
go mod tidy

# Start the Go engine
go run main.go

3. Frontend Setup
Bash

# Navigate to frontend directory
cd ../frontend

# Install packages
npm install

# Start the dashboard
npm start

🚢 Deployment
Docker Deployment

The most reliable way to deploy the Basalt Sentinel is via the integrated container:
Bash

# Build the integrated image
docker build -t basalt-sentinel:latest .

# Run the container with environment variables
docker run -p 8080:8080 -e GEMINI_API_KEY='your_key' basalt-sentinel:latest
