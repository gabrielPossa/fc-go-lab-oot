package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"regexp"
	"strings"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel/trace"

	"github.com/gabrielPossa/fc-go-lab-oot/internal/cep"
)

var digitCheck = regexp.MustCompile(`^\d{8}$`)

var weatherClient = http.Client{Transport: otelhttp.NewTransport(http.DefaultTransport)}

func GetWeatherByCEP(tracer trace.Tracer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()

		ctx, span := tracer.Start(r.Context(), "GetWeatherByCEP")
		defer span.End()

		var rCep cep.RequestCEP
		err := json.NewDecoder(r.Body).Decode(&rCep)
		if err != nil {
			log.Println("JSON Invalido")
			http.Error(w, "invalid zipcode", http.StatusUnprocessableEntity)
			return
		}

		cepString := strings.Replace(rCep.Cep, "-", "", -1)

		if !digitCheck.MatchString(cepString) {
			log.Println("CEP invalido,CEP deve ser composto por 8 números. Formatos aceitos: 12345678 ou 12345-678")
			http.Error(w, "invalid zipcode", http.StatusUnprocessableEntity)
			return
		}

		callWeatherAPI(ctx, tracer, cepString, w)
	}
}

func callWeatherAPI(ctx context.Context, tracer trace.Tracer, cep string, w http.ResponseWriter) {
	ctx, span := tracer.Start(ctx, "callWeatherAPI")
	defer span.End()

	url := fmt.Sprintf("%s/%s/weather", "http://weather:8080/cep", cep)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		log.Printf("Falha ao criar request para o MS Weather: %s", err)
		http.Error(w, "Internal ERROR", http.StatusInternalServerError)
		return
	}
	res, err := weatherClient.Do(req)
	if err != nil {
		log.Printf("Falha ao chamar MS Weather: %s", err)
		http.Error(w, "Internal ERROR", http.StatusInternalServerError)
		return
	}

	defer res.Body.Close()
	copyHeader(w.Header(), res.Header)
	w.WriteHeader(res.StatusCode)
	io.Copy(w, res.Body)
}

func copyHeader(dst, src http.Header) {
	for k, vv := range src {
		for _, v := range vv {
			dst.Add(k, v)
		}
	}
}
