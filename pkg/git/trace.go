package git

import (
	"io"
	"log"

	jaeger "github.com/uber/jaeger-client-go"

	jaegercfg "github.com/uber/jaeger-client-go/config"
)

var (
	config jaegercfg.Configuration
)

func init() {
	config = jaegercfg.Configuration{
		Sampler: &jaegercfg.SamplerConfig{
			Type:  jaeger.SamplerTypeConst,
			Param: 1,
		},
		Reporter: &jaegercfg.ReporterConfig{
			LogSpans:          true,
			CollectorEndpoint: "https://jaeger-api.isteven.cn/api/traces", // 替换host
		},
	}
}

func StarTracer() io.Closer {
	closer, err := config.InitGlobalTracer("github.com/iineva/git-nfs")
	if err != nil {
		log.Printf("Could not initialize jaeger tracer: %s", err.Error())
		return nil
	}
	return closer
}
