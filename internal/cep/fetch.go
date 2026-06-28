package cep

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

const baseUrl = "https://viacep.com.br/ws/"

// client com transport otel gera um span HTTP filho para a chamada ao ViaCEP.
var client = http.Client{Transport: otelhttp.NewTransport(http.DefaultTransport)}

func FetchCEPData(ctx context.Context, cep string) (*CEP, error) {

	url := fmt.Sprintf("%s%s/json/", baseUrl, cep)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	var cepData CEP
	err = json.NewDecoder(res.Body).Decode(&cepData)
	if err != nil {
		return nil, err
	}

	return &cepData, nil
}
