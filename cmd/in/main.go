package main

import (
	"encoding/json"
	"fmt"
	"github.com/cah-tylerrasor/pact-verification-resource/pkg/broker"
	"github.com/cah-tylerrasor/pact-verification-resource/pkg/concourse"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	var request concourse.InRequest
	populateRequest(&request)

	client := broker.NewClient(request.Source.BrokerURL)

	if request.Source.Username != nil && request.Source.Password != nil {
		broker.WithBasicAuth(*request.Source.Username, *request.Source.Password)(client)
	}

	if len(os.Args) != 2 {
		concourse.FailTask("first argument must be a directory")
	}

	dir := os.Args[1]

	bytes, err := client.GetValidationRaw(request.Source.Consumer, request.Version.Provider, request.Version.PactVersion)
	if err != nil {
		concourse.FailTask("could not read bytes: %s", err)
	}

	pactPath := fmt.Sprintf("%s-%s-%s.json",  request.Source.Consumer, request.Version.Provider, request.Version.PactVersion)
	pactPath = strings.ReplaceAll(pactPath, " ", "-")

	file, err := os.Create(filepath.Join(dir, pactPath))
	if err != nil {
		concourse.FailTask("could not open file: %s", err)
	}

	_, err = file.Write(bytes)
	defer file.Close()
	if err != nil {
		concourse.FailTask("could not write to file: %s", err)
	}

	resp := concourse.InResponse{Version: request.Version, Metadata: concourse.Metadata{{Name: "pactVerification", Value: pactPath}}}
	if err := json.NewEncoder(os.Stdout).Encode(resp); err != nil {
		concourse.FailTask("could not encode response: %s", err)
	}
}

func populateRequest(req *concourse.InRequest) {
	if err := json.NewDecoder(os.Stdin).Decode(req); err != nil {
		concourse.FailTask("could not decode request: %s", err)
	}
}
