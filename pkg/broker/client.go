package broker

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/cah-tylerrasor/pact-verification-resource/pkg/concourse"
	"io/ioutil"
	"net/http"
	"time"
)

type (
	Client struct {
		baseURL  string
		client   *http.Client
		token    string
		username string
		password string
	}

	Option func(*Client)
)

func WithBasicAuth(username, password string) Option {
	return func(broker *Client) {
		broker.username = username
		broker.password = password
	}
}

func WithClient(client *http.Client) Option {
	return func(broker *Client) {
		broker.client = client
	}
}

func NewClient(brokerURL string, opts ...Option) *Client {
	broker := Client{baseURL: brokerURL}
	for _, o := range opts {
		o(&broker)
	}

	if broker.client == nil {
		broker.client = &http.Client{Timeout: time.Second * 5}
	}

	return &broker
}

func (c *Client) get(url string) (resp *http.Response, err error) {
	req, err := http.NewRequest(http.MethodGet, url, nil)

	if len(c.username) > 0 && len(c.password) > 0 {
		auth := base64.StdEncoding.EncodeToString([]byte(c.username + ":" + c.password))
		req.Header.Add("Authorization", "Basic " + auth)
	}

	return c.client.Do(req)
}

func (c *Client) GetValidation(consumer, provider string) HalPactVerification {
	url := fmt.Sprintf("%s/pacts/provider/%s/consumer/%s/verification-results/latest", c.baseURL, provider, consumer)
	return c.GetValidationFromUrl(consumer, provider, url)
}

func (c *Client) GetTaggedValidation(consumer, provider, tag string) HalPactVerification {
	url := fmt.Sprintf("%s/pacts/provider/%s/consumer/%s/latest/%s/verification-results/latest", c.baseURL, provider, consumer, tag)
	return c.GetValidationFromUrl(consumer, provider, url)
}

func (c *Client) GetValidationFromUrl(consumer, provider, url string) HalPactVerification {
	resp, err := c.get(url)
	if err != nil {
		concourse.FailTask("error while requesting information: %s\n", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode > 399 {
		concourse.FailTask("error while requesting information: %d\n", resp.StatusCode)
	}

	var halPactVerification HalPactVerification
	err = json.NewDecoder(resp.Body).Decode(&halPactVerification)
	if err != nil {
		concourse.FailTask("error decoding json: %s\n", err)
	}

	return halPactVerification
}

func (c *Client) GetValidationRaw(consumer, provider, version string) ([]byte, error) {
	url := fmt.Sprintf("%s/pacts/provider/%s/consumer/%s/pact-version/%s/verification-results/latest", c.baseURL, provider, consumer, version)

	resp, err := c.get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode > 399 {
		return nil, fmt.Errorf("error while requesting information: %d", resp.StatusCode)
	}

	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return bytes, nil
}
