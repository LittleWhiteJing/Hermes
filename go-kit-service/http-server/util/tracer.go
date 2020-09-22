package util

import (
	"github.com/openzipkin/zipkin-go"
	zipkinhttp "github.com/openzipkin/zipkin-go/reporter/http"
)

func GetTracer (zipkinURL *string, hostPort string, serviceName string) (*zipkin.Tracer, error) {
	useNoopTracer := (*zipkinURL == "")
	reporter      := zipkinhttp.NewReporter(*zipkinURL)
	defer reporter.Close()
	zEP, _ := zipkin.NewEndpoint(serviceName, hostPort)
	return zipkin.NewTracer(
			reporter,
			zipkin.WithLocalEndpoint(zEP),
			zipkin.WithNoopTracer(useNoopTracer),
	)
}



