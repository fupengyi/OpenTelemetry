package OpenTelemetry_Go
OpenTelemetry Exporters
一旦 OpenTelemetry SDK 创建并处理了遥测数据，就需要将其导出。这个包包含用于此目的的出口商。

Exporter Packages
以下导出程序包随以下 OpenTelemetry 信号支持一起提供。
Exporter Package												Metrics		Traces
go.opentelemetry.io/otel/exporters/jaeger									✓
go.opentelemetry.io/otel/exporters/otlp/otlpmetric				✓
go.opentelemetry.io/otel/exporters/otlp/otlptrace							✓
go.opentelemetry.io/otel/exporters/prometheus					✓
go.opentelemetry.io/otel/exporters/stdout/stdoutmetric			✓
go.opentelemetry.io/otel/exporters/stdout/stdouttrace						✓
go.opentelemetry.io/otel/exporters/zipkin									✓
请参阅与此项目兼容的第 3 部分导出器的 OpenTelemetry 注册表。 给予反馈
