package OpenTelemetry_Go
OpenTelemetry-Go
OpenTelemetry-Go 是 OpenTelemetry 的 Go 实现。它提供了一组 API 来直接测量软件的性能和行为，并将这些数据发送到可观察性平台。

Project Status
Signal		Status		Project
Traces		Stable		N/A
Metrics		Alpha		N/A
Logs		Frozen[1]	N/A
[1]：在我们开发 Traces 和 Metrics 时，该项目的日志信号开发已停止。目前没有接受日志拉取请求。
在我们的本地项目板和里程碑中跟踪特定于此存储库的进度和状态。
可以在版本控制文档中找到项目版本控制信息和稳定性保证。

Compatibility
OpenTelemetry-Go 确保与当前支持的 Go 语言版本兼容：
支持每个主要的 Go 版本，直到有两个更新的主要版本。例如，Go 1.5 一直支持到 Go 1.7 版本，Go 1.6 一直支持到 Go 1.8 版本。
对于上游不再支持的 Go 版本，opentelemetry-go 将停止以下列方式确保与这些版本的兼容性：
1.将发布 opentelemetry-go 的次要版本以添加对新支持的 Go 版本的支持。
2.以下 opentelemetry-go 的次要版本将删除对最旧版本（现已归档上游）的 Go 的兼容性测试。这个和未来的 opentelemetry-go 版本可能包含仅受当前支持的 Go 版本支持的功能。
目前，该项目支持以下环境。
...
虽然这个项目应该适用于其他系统，但目前没有为这些系统提供兼容性保证。

Getting Started
您可以在 opentelemetry.io 上找到入门指南。
OpenTelemetry 的目标是提供一组 API 来从您的应用程序捕获分布式跟踪和指标并将它们发送到可观察性平台。该项目允许您为用 Go 编写的应用程序执行
此操作。此过程有两个步骤：检测您的应用程序和配置导出器。

Instrumentation
要开始从您的应用程序中捕获分布式跟踪和指标事件，首先需要对其进行检测。最简单的方法是为您的代码使用检测库。请务必查看官方支持的检测库-the officially supported instrumentation libraries.
如果您需要扩展仪器库提供的遥测，或者想直接为您的应用程序构建您自己的仪器，您将需要使用 Go otel 包。包含的示例是查看此过程的一些实际用途的好方法。

Export
现在您的应用程序已被检测以收集遥测数据，它需要一个导出管道来将该遥测数据发送到可观察性平台。
OpenTelemetry 项目的所有官方支持的导出器都包含在导出器目录中。
Exporter		Metrics			Traces
Jaeger							✓
OTLP			✓				✓
Prometheus		✓
stdout			✓				✓
Zipkin							✓

Contributing
请参阅贡献文档。

Documentation
Overview ¶
包 otel 提供对 OpenTelemetry API 的全局访问。 otel 包的子包提供了 OpenTelemetry API 的实现。

提供的 API 用于检测代码并测量有关该代码的性能和操作的数据。默认情况下，测量数据不会在任何地方进行处理或传输。 OpenTelemetry SDK 的实现，如
默认 SDK 实现 (go.opentelemetry.io/otel/sdk)，以及关联的导出器用于处理和传输此数据。

要阅读入门指南，请参阅 https://opentelemetry.io/docs/go/getting-started/。
要阅读有关跟踪的更多信息，请参阅 go.opentelemetry.io/otel/trace。
要了解有关指标的更多信息，请参阅 go.opentelemetry.io/otel/metric。
要阅读有关传播的更多信息，请参阅 go.opentelemetry.io/otel/propagation 和 go.opentelemetry.io/otel/baggage。
