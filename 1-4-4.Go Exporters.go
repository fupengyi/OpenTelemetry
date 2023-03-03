package OpenTelemetry_Go

Exporters
为了可视化和分析您的跟踪和指标，您需要将它们导出到后端。


OTLP Exporter
go.opentelemetry.io/otel/exporters/otlp/otlptrace 和 go.opentelemetry.io/otel/exporters/otlp/otlpmetric 包中提供了 OpenTelemetry 协议 (OTLP) 导出。
请在 GitHub 上找到更多文档


1.Jaeger Exporter
Jaeger 导出在 go.opentelemetry.io/otel/exporters/jaeger 包中可用。

OpenTelemetry-Go Jaeger Exporter
Jaeger 实现的 OpenTelemetry 跨度导出器。

Installation
	go get -u go.opentelemetry.io/otel/exporters/jaeger

Example
See ../../example/jaeger.

Configuration
导出器可用于将跨度发送到：
	Jaeger 代理通过 WithAgentEndpoint 选项使用 jaeger.thrift over compact thrift 协议。
	Jaeger 收集器通过 WithCollectorEndpoint 选项在 HTTP 上使用 jaeger.thrift。

Environment Variables
可以使用以下环境变量（而不是选项对象）来覆盖默认配置。
Environment variable					Option			Default value
OTEL_EXPORTER_JAEGER_AGENT_HOST			WithAgentHost	localhost
OTEL_EXPORTER_JAEGER_AGENT_PORT			WithAgentPort	6831
OTEL_EXPORTER_JAEGER_ENDPOINT			WithEndpoint	http://localhost:14268/api/traces
OTEL_EXPORTER_JAEGER_USER				WithUsername
OTEL_EXPORTER_JAEGER_PASSWORD			WithPassword
使用选项的配置优先于环境变量。

Contributing
此导出器在自定义导入路径中使用 Apache Thrift 库 (v0.14.1) 的销售副本。将来重新生成 Thrift 代码时，请根据需要调整导入路径。

References
	1-4-4-4-1.Jaeger
	1-4-4-4-2.OpenTelemetry to Jaeger Transformation
	1-4-4-4-3.OpenTelemetry Environment Variable Specification
	
2.
Prometheus Exporter
Prometheus 导出在 go.opentelemetry.io/otel/exporters/prometheus 包中可用。
请在 GitHub 上找到更多文档
