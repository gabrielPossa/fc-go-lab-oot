package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"

	"github.com/gabrielPossa/fc-go-lab-oot/pkg/webserver"
	"go.opentelemetry.io/otel"

	otelProvider "github.com/gabrielPossa/fc-go-lab-oot/pkg/otel"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	shutdown, err := otelProvider.InitProvider("cep")
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err := shutdown(ctx); err != nil {
			log.Fatal("failed to shutdown TracerProvider: %w", err)
		}
	}()

	tracer := otel.Tracer("microservice-tracer")

	httpServer := webserver.NewWebServer(":8080")
	httpServer.AddHandler("/weatherByCEP", webserver.POST, GetWeatherByCEP(tracer))
	fmt.Println("Starting web server on port", "8080")
	httpServer.Start()
}
