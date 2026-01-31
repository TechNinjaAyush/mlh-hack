package models

import "time"

type ServiceGraph struct {
	Name      string   `json:"name"`
	Health    float32  `json:"health"`
	DependsOn []string `json:"depends_on"`
	Reason    string   `json:"failure_check_reason"`
	
}

type Payload struct {
	IncidentID  string    `json:"incidentID"`
	Root        string    `json:"Root"`
	FailedNodes []string  `json:"FailedNodes"`
	Time        time.Time `json:"Time"`
	BlastRadius int       `json:"BlastRadius"`
	RCA         string    `json:"rca"`
}
