package OpenTelemetry_Go

import (
	"context"
	"net/http"
	"time"
)

OpenTelemetry-Go Jaeger Exporter
OpenTelemetry Jaeger 导出器


Installation
	go get -u go.opentelemetry.io/otel/exporters/trace/jaeger


Maintenance
此导出器在自定义导入路径中使用 Apache Thrift 库 (v0.14.1) 的销售副本。将来重新生成 Thrift 代码时，请根据需要调整导入路径。


Documentation
Overview ¶
包 jaeger 包含一个用于 Jaeger 的 OpenTelemetry 跟踪导出器。
该软件包目前处于 pre-GA 阶段。在我们努力跟踪不断发展的 OpenTelemetry 规范和用户反馈时，可能会在后续的次要版本发布中引入向后不兼容的更改。


Functions ¶
func InstallNewPipeline(endpointOption EndpointOption) (*sdktrace.TracerProvider, error)		// InstallNewPipeline 使用推荐的配置实例化一个 NewExportPipeline 并在全局注册它。
func NewExportPipeline(endpointOption EndpointOption) (*sdktrace.TracerProvider, error)			// NewExportPipeline 使用推荐的跟踪提供程序设置设置完整的导出管道


Types ¶
type AgentEndpointOption func(o *AgentEndpointOptions)
type AgentEndpointOptions struct {
	// contains filtered or unexported fields
}
1.func WithAgentHost(host string) AgentEndpointOption											// WithAgentHost 设置要在代理客户端端点中使用的主机。
																								// 此选项会覆盖为 OTEL_EXPORTER_JAEGER_AGENT_HOST 环境变量设置的任何值。
																								// 如果未传递此选项且未设置环境变量，则默认使用“localhost”。
2.func WithAgentPort(port string) AgentEndpointOption											// WithAgentPort 设置要在代理客户端端点中使用的端口。
																								// 此选项会覆盖为 OTEL_EXPORTER_JAEGER_AGENT_PORT 环境变量设置的任何值。
																								// 如果未传递此选项且未设置环境变量，则默认使用“6832”。
3.func WithAttemptReconnectingInterval(interval time.Duration) AgentEndpointOption				// WithAttemptReconnectingInterval 设置尝试重新解析代理端点之间的间隔。
4.func WithDisableAttemptReconnecting() AgentEndpointOption										// WithDisableAttemptReconnecting 设置选项以禁用重新连接 udp 客户端。
5.func WithLogger(logger *log.Logger) AgentEndpointOption										// WithLogger 设置代理客户端使用的记录器。






type CollectorEndpointOption func(o *CollectorEndpointOptions)
type CollectorEndpointOptions struct {
	// contains filtered or unexported fields
}
1.func WithEndpoint(endpoint string) CollectorEndpointOption									// WithEndpoint 是发送到 span 的 Jaeger 收集器的 URL。
																								// 此选项会覆盖为 OTEL_EXPORTER_JAEGER_ENDPOINT 环境变量设置的任何值。
																								// 如果不传递此选项且未设置环境变量，则默认使用“http://localhost:14250”。
2.func WithHTTPClient(client *http.Client) CollectorEndpointOption								// WithHTTPClient 设置用于向收集器端点发出请求的 http 客户端。
3.func WithPassword(password string) CollectorEndpointOption									// WithPassword 设置用于发送给收集器的所有请求的授权标头中的密码。
																								// 此选项覆盖为 OTEL_EXPORTER_JAEGER_PASSWORD 环境变量设置的任何值。
																								// 如果不传递此选项且未设置环境变量，则不会设置密码。
4.func WithUsername(username string) CollectorEndpointOption									// WithUsername 设置要在发送给收集器的所有请求的授权标头中使用的用户名。
																								// 此选项会覆盖为 OTEL_EXPORTER_JAEGER_USER 环境变量设置的任何值。
																								// 如果未传递此选项且未设置环境变量，则不会设置用户名。






type EndpointOption func() (batchUploader, error)
1.func WithAgentEndpoint(options ...AgentEndpointOption) EndpointOption							// WithAgentEndpoint 将 Jaeger 导出器配置为将 span 发送到 jaeger-agent。
																								// 如果没有提供明确的选项，这将使用以下环境变量进行配置：
																								// - OTEL_EXPORTER_JAEGER_AGENT_HOST 用于代理地址主机
																								// - OTEL_EXPORTER_JAEGER_AGENT_PORT 用于代理地址端口
																								// 传递的选项将优先于任何环境变量，如果没有提供，将使用默认值。
2.func WithCollectorEndpoint(options ...CollectorEndpointOption) EndpointOption					// WithCollectorEndpoint 定义 Jaeger HTTP Thrift 收集器的完整 URL。
																								// 如果没有提供明确的选项，这将使用以下环境变量进行配置：
																								// - OTEL_EXPORTER_JAEGER_ENDPOINT 是用于将跨度直接发送到收集器的 HTTP 端点。
																								// - OTEL_EXPORTER_JAEGER_USER 是要作为身份验证发送到收集器端点的用户名。
																								// - OTEL_EXPORTER_JAEGER_PASSWORD 是作为身份验证发送到收集器端点的密码。
																								// 传递的选项将优先于任何环境变量。
																								// 如果没有为端点提供任何值，将使用默认值“http://localhost:14250”。
																								// 如果没有为用户名或密码提供值，则不会设置它们，因为没有默认值。




type Exporter struct {																			// Exporter 将 OpenTelemetry 跨度导出到 Jaeger 代理或收集器。
	// contains filtered or unexported fields
}
1.func NewRawExporter(endpointOption EndpointOption) (*Exporter, error)							// NewRawExporter 返回一个 OTel Exporter 实现，它将收集的 span 导出到 Jaeger。
2.func (e *Exporter) ExportSpans(ctx context.Context, spans []*sdktrace.SpanSnapshot) error		// ExportSpans 将 OpenTelemetry 跨度转换并导出到 Jaeger。
3.func (e *Exporter) Shutdown(ctx context.Context) error										// Shutdown 停止出口商。这将关闭所有连接并释放导出器持有的所有资源。
