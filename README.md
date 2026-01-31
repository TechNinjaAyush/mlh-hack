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
