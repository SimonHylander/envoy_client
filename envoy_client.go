package envoy_client

import (
	"fmt"
	"net/http"
	"os"
)

/**
https://www.envoyproxy.io/docs/envoy/latest/configuration/http/http_filters/router_filter#x-envoy-retry-on
*/

type RetryOn string

const (
	ServerError          RetryOn = "5xx"
	GatewayError         RetryOn = "gateway-error"
	Reset                RetryOn = "reset"
	ConnectFailure       RetryOn = "connect-failure"
	EnvoyRateLimited     RetryOn = "envoy-ratelimited"
	Retriable4xx         RetryOn = "retriable-4xx"
	RefusedStream        RetryOn = "refused-stream"
	RetriableStatusCodes RetryOn = "retriable-status-codes"
	RetriableHeaders     RetryOn = "retriable-headers"
	Http3PostConnectFail RetryOn = "http3-post-connect-failure"
)

type EnvoyClient struct {
	ClientProxy string
	ServiceName string
	URI         string
	HTTPMethod  string
	Data        string
	Headers     map[string]string
	HTTPClient  *http.Client
}

func NewEnvoyClient() *EnvoyClient {
	return &EnvoyClient{
		Headers: make(map[string]string),
	}
}

func (ec *EnvoyClient) Service(serviceName string) *EnvoyClient {
	ec.ServiceName = serviceName

	return ec
}

/**
https://www.envoyproxy.io/docs/envoy/latest/configuration/http/http_filters/router_filter#x-envoy-max-retries
*/
func (ec *EnvoyClient) Retries(retries int) *EnvoyClient {
	ec.Header("x-envoy-max-retries", fmt.Sprintf("%d", retries))
	return ec
}

func (ec *EnvoyClient) RetryOn(retryOnHeader RetryOn) *EnvoyClient {
	ec.Header("x-envoy-retry-on", fmt.Sprintf("%s", retryOnHeader))
	return ec
}

func (ec *EnvoyClient) Header(headerName, headerValue string) *EnvoyClient {
	ec.Headers[headerName] = headerValue
	return ec
}

func (ec *EnvoyClient) Get(path string) *EnvoyClient {
	host := getEnv("CLIENT_EGRESS", "http://envoy:9000")
	ec.URI = host + "/" + ec.ServiceName + path
	ec.HTTPMethod = http.MethodGet

	return ec
}

func (ec *EnvoyClient) Post(path string) *EnvoyClient {
	host := getEnv("CLIENT_EGRESS", "http://envoy:9000")
	ec.URI = host + "/" + ec.ServiceName + path
	ec.HTTPMethod = http.MethodPost

	return ec
}

func (ec *EnvoyClient) Put(path string) *EnvoyClient {
	host := getEnv("CLIENT_EGRESS", "http://envoy:9000")
	ec.URI = host + "/" + ec.ServiceName + path
	ec.HTTPMethod = http.MethodPut

	return ec
}

func (ec *EnvoyClient) Delete(path string) *EnvoyClient {
	host := getEnv("CLIENT_EGRESS", "http://envoy:9000")
	ec.URI = host + "/" + ec.ServiceName + path
	ec.HTTPMethod = http.MethodDelete

	return ec
}

func (ec *EnvoyClient) Request() *http.Request {
	/* var body *bytes.Reader

	if ec.Data != nil && ec.Data != "" {
		body = bytes.NewReader([]byte(ec.Data))
	} */

	req, err := http.NewRequest(ec.HTTPMethod, ec.URI, nil)

	if err != nil {
		panic(err)
	}

	if ec.Headers != nil {
		for key, value := range ec.Headers {
			req.Header.Set(key, value)
		}
	}

	return req
}

func (ec *EnvoyClient) Call() (*http.Response, error) {
	if ec.HTTPClient == nil {
		ec.HTTPClient = &http.Client{}
	}

	req := ec.Request()

	response, err := ec.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func getEnv(env, defaultValue string) string {
	value := os.Getenv(env)
	if value == "" {
		value = defaultValue
	}
	return value
}
