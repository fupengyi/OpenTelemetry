package OpenTelemetry_Go

import (
	"context"
	"time"
)

github.com/GoogleCloudPlatform/opentelemetry-operations-go
OpenTelemetry Google Cloud Trace Exporter 允许用户将收集的跟踪和跨度发送到 Google Cloud。
Google Cloud Trace 是一个分布式跟踪后端系统。它可以帮助开发人员收集解决微服务和单体架构中的延迟问题所需的时序数据。它管理收集的跟踪数据的收集和查找。
此导出程序包假定您的应用程序已使用 OpenTelemetry SDK 进行检测。准备好导出 OpenTelemetry 数据后，您可以将此导出器添加到您的应用程序中。


Setup
Google Cloud Trace 是 Google Cloud Platform 提供的托管服务。官方 GCP 文档中提供了 OpenTelemetry 的端到端设置指南，因此该文档通过导出器设置。


Usage
导入跟踪导出程序包后，创建并安装新的导出管道，然后就可以开始跟踪了。如果您在 GCP 环境中运行，导出器将使用环境的服务帐户自动进行身份验证。如果
没有，您将需要按照身份验证中的说明进行操作。
package main

import (
	"context"
	"log"
	
	texporter "github.com/GoogleCloudPlatform/opentelemetry-operations-go/exporter/trace"
	
	"go.opentelemetry.io/otel"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

func main() {
	exporter, err := texporter.New()
	if err != nil {
		log.Fatalf("unable to set up tracing: %v", err)
	}
	tp := sdktrace.NewTracerProvider(sdktrace.WithBatcher(exporter))
	defer tp.Shutdown(context.Background())
	
	otel.SetTracerProvider(tp)
	
	tracer := tp.Tracer("example.com/trace")
	ctx, span := tracer.Start(ctx, "foo")
	defer span.End()
	
	// Do some work.
}



Authentication
Google Cloud Trace 导出器依赖于 google.FindDefaultCredentials，因此默认情况下会自动检测服务帐户，但也可以在特定条件下检测自定义凭证文
件（所谓的 service_account_key.json）。引用 google.FindDefaultCredentials 的文档：
	一个 JSON 文件，其路径由 GOOGLE_APPLICATION_CREDENTIALS 环境变量指定。
	位于 gcloud 命令行工具已知位置的 JSON 文件。在 Windows 上，这是 %APPDATA%/gcloud/application_default_credentials.json。在其他系统上，$HOME/.config/gcloud/application_default_credentials.json。

在本地运行代码时，除了 GOOGLE_APPLICATION_CREDENTIALS 之外，您可能还需要指定 Google 项目 ID。这最好使用环境变量（例如 GOOGLE_CLOUD_PROJECT）
和 WithProjectID 方法来完成，例如：
projectID := os.Getenv("GOOGLE_CLOUD_PROJECT")
_, shutdown, err := texporter.InstallNewPipeline(
	[]texporter.Option {
		texporter.WithProjectID(projectID),
		// other optional exporter options
	},
	...
)



Useful links
有关 OpenTelemetry 的更多信息，请访问：https://opentelemetry.io/
有关 OpenTelemetry Go 的更多信息，请访问：https://github.com/open-telemetry/opentelemetry-go
在 https://cloud.google.com/trace 上了解有关 Google Cloud Trace 的更多信息



Documentation
Functions ¶
1.func Version() string														// Version 是正在使用的 OpenTelemetry Operations Trace Exporter 的当前发行版本。
2.func WithAttributeMapping(mapping AttributeMapping) func(o *options)		// WithAttributeMapping 配置如何将 OpenTelemetry 跨度属性映射到谷歌云跟踪跨度属性。默认情况下，它映射到在跟踪 UI 中突出使用的属性。
3.func WithContext(ctx context.Context) func(o *options)					// WithContext 设置跟踪导出器和度量导出器所依赖的上下文。
4.func WithDestinationProjectQuota() func(o *options)						// WithDestinationProjectQuota 允许按请求使用目标项目的配额。例如，设置 gcp.project.id 资源属性时。
5.func WithErrorHandler(handler otel.ErrorHandler) func(o *options)			// WithErrorHandler 设置在将跨度数据上传到 Stackdriver 时发生错误时调用的挂钩。如果未设置自定义挂钩，则会记录错误。
6.func WithProjectID(projectID string) func(o *options)						// WithProjectID 将 Google Cloud Platform 项目设置为 projectID。如果不使用此选项，它会自动从默认凭据检测过程中检测项目 ID。请在文档中找到默认凭证检测过程的详细顺序：https://godoc.org/golang.org/x/oauth2/google#FindDefaultCredentials
7.func WithTimeout(t time.Duration) func(o *options)						// WithTimeout 设置跟踪导出器和度量导出器的超时 如果未设置，则默认为 12 秒超时。
8.func WithTraceClientOptions(opts []option.ClientOption) func(o *options)	// WithTraceClientOptions 设置用于跟踪的附加客户端选项。


Types ¶
type AttributeMapping func(attribute.Key) attribute.Key						// AttributeMapping 确定如何从 OpenTelemetry 跨度属性键映射到云跟踪属性键。


type Exporter struct {																				// Exporter 是将数据上传到 Stackdriver 的跟踪导出器。
	// contains filtered or unexported fields														// TODO(yoshifumi)：一旦规范定义过程和采样器实现完成，就添加一个指标导出器。
}
1.func New(opts ...Option) (*Exporter, error)														// New 创建一个新的 Exporter，它实现了 trace.Exporter。
2.func (e *Exporter) ExportSpans(ctx context.Context, spanData []sdktrace.ReadOnlySpan) error		// ExportSpans 将 ReadOnlySpan 导出到 Stackdriver Trace。
3.func (e *Exporter) Shutdown(ctx context.Context) error											// Shutdown 等待导出的数据上传。出于我们的目的，它关闭了客户端。


type Option func(*options)													// Option 是传递给导出器初始化函数的函数类型。
