package models

import "time"

type ServiceGraph struct {
	Name      string   `json:"name"`
	Health    float32  `json:"health"`
	DependsOn []string `json:"depends_on"`
}

type Payload struct {
	Root        string    `json:"Root"`
	FailedNodes []string  `json:"FailedNodes"`
	Time        time.Time `json:"Time"`
	BlastRadius int       `json:"BlastRadius"`
}
