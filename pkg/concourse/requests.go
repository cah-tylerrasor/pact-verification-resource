package concourse

import (
	"time"
)

type (
	// Version holds information about a Concourse resource version
	Version struct {
		Provider  string       `json:"provider"`
		ProviderVersion string `json:"providerVersion"`
		UpdatedAt time.Time    `json:"updated_at"`
		PactVersion   string       `json:"pactVersion"`
	}

	// Source is the configuration of the resource
	Source struct {
		BrokerURL string   `json:"broker_url"`
		Consumer  string   `json:"consumer"`
		Providers []string `json:"providers"`
		Tag       *string  `json:"tag"`
		Username  *string  `json:"username"`
		Password  *string  `json:"password"`
	}

	// Metadata is a key-value attribute pair that describes a version
	Metadata []struct {
		Name  string `json:"name"`
		Value string `json:"value"`
	}

	// CheckRequest is the Concourse input when checking for a new version
	CheckRequest struct {
		Source  Source  `json:"source"`
		Version Version `json:"version"`
	}

	// InRequest is the Concourse input when passing the input
	InRequest struct {
		Source  Source  `json:"source"`
		Version Version `json:"version"`
	}

	// InResponse is the data being sent out to the next job
	InResponse struct {
		Version  Version  `json:"version"`
		Metadata Metadata `json:"metadata"`
	}
)
