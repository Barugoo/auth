package tracer

import (
	"io"

	opentracing "github.com/opentracing/opentracing-go"
	jaeger "github.com/uber/jaeger-client-go"
	jaeger_config "github.com/uber/jaeger-client-go/config"
	"github.com/uber/jaeger-lib/metrics"

	"github.com/barugoo/oscillo-auth/config"
)

func NewTracer(config *config.ServiceConfig) (opentracing.Tracer, io.Closer, error) {
	cfg := &jaeger_config.Configuration{
		ServiceName: config.ServiceName,
		Sampler: &jaeger_config.SamplerConfig{
			Type:  "const",
			Param: 1,
		},
		Reporter: &jaeger_config.ReporterConfig{
			LogSpans: true,
		},
	}
	return cfg.NewTracer(
		jaeger_config.Logger(jaeger.StdLogger),
		jaeger_config.Metrics(metrics.NullFactory),
	)
}
