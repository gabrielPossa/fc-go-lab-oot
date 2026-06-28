package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"

	"github.com/go-chi/chi/v5"
	"go.opentelemetry.io/otel/trace"

	"github.com/gabrielPossa/fc-go-lab-oot/internal/cep"
	"github.com/gabrielPossa/fc-go-lab-oot/internal/weather"
	"github.com/gabrielPossa/fc-go-lab-oot/pkg/utils"
)

var digitCheck = regexp.MustCompile(`^\d{8}$`)

type response struct {
	City  string  `json:"city"`
	TempC float64 `json:"temp_C"`
	TempF float64 `json:"temp_F"`
	TempK float64 `json:"temp_K"`
}

func GetWeatherByCEP(tracer trace.Tracer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, span := tracer.Start(r.Context(), "GetWeatherByCEP")
		defer span.End()

		cepString := chi.URLParam(r, "CEP")

		if !digitCheck.MatchString(cepString) {
			http.Error(w, "invalid zipcode", http.StatusUnprocessableEntity)
			return
		}

		ctxCep, spanCep := tracer.Start(ctx, "FetchCEPData")
		CEP, err := cep.FetchCEPData(ctxCep, cepString)
		spanCep.End()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if CEP.Erro == "true" {
			http.Error(w, "can not find zipcode", http.StatusNotFound)
			return
		}

		weatherQ := fmt.Sprintf("%s,%s", CEP.Localidade, CEP.Estado)

		ctxWeather, spanWeather := tracer.Start(ctx, "GetWeatherData")
		weatherData, err := weather.GetWeatherData(ctxWeather, weatherQ)
		spanWeather.End()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		resp := response{
			City:  CEP.Localidade,
			TempC: weatherData.Current.TempC,
			TempF: utils.CelciusToFahrenheit(weatherData.Current.TempC),
			TempK: utils.CelciusToKelvin(weatherData.Current.TempC),
		}

		w.Header().Set("Content-Type", "application/json")
		err = json.NewEncoder(w).Encode(resp)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}
