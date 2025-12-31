package models

import "time"

type ServiceGraph struct {
	Name      string   `json:"name"`
	Health    float32  `json:"health"`
	DependsOn []string `json:"depends_on"`
}

type Payload struct {
	Root        string    `json:"root"`
	FailedNodes []string  `json:"failednodes"`
	Time        time.Time `json:"failed_time"`

	BlastRadius int `json:"blast_radius"`
}
