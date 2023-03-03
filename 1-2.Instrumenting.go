package OpenTelemetry_Go
Instrumenting
OpenTelemetry 如何促进应用程序的自动和手动检测。
为了使系统可观察，必须对其进行检测：也就是说，来自系统组件的代码必须发出跟踪、指标和日志。
无需修改源代码，您就可以使用自动检测从应用程序收集遥测数据。如果您之前使用 APM 代理从您的应用程序中提取遥测数据，Automatic Instrumentation 将为您提供类似的开箱即用体验。
为了进一步促进应用程序的检测，您可以通过针对 OpenTelemetry API 进行编码来手动检测您的应用程序。
为此，您不需要检测应用程序中使用的所有依赖项：
	通过直接调用 OpenTelemetry API 本身，可以开箱即用地观察到您的一些库。这些库有时被称为本地检测。
	对于没有这种集成的库，OpenTelemetry 项目提供特定于语言的 Instrumentation Libraries
请注意，对于大多数语言，可以同时使用手动和自动检测：自动检测将使您能够快速深入了解您的应用程序，而手动检测将使您能够将细粒度的可观察性嵌入到您的代码中。
手动和自动检测的确切安装机制因您开发的语言而异，但以下部分介绍了一些相似之处。

Automatic Instrumentation
如果适用，OpenTelemetry 的特定语言实现将提供一种在不接触源代码的情况下检测应用程序的方法。虽然底层机制取决于语言，但至少这会将 OpenTelemetry
API 和 SDK 功能添加到您的应用程序中。此外，他们可能会添加一组仪器库和导出器依赖项。

配置可通过环境变量和可能的语言特定方式（例如 Java 中的系统属性）获得。至少，必须配置服务名称以标识正在检测的服务。还有各种其他配置选项可用，可能包括：
	数据源特定配置
	导出器配置
	传播子配置
	资源配置

Manual Instrumentation
Import the OpenTelemetry API and SDK
您首先需要将 OpenTelemetry 导入您的服务代码。如果您正在开发旨在由可运行二进制文件使用的库或其他组件，那么您只需依赖 API。如果您的工件是一个
独立的流程或服务，那么您将依赖于 API 和 SDK。有关 OpenTelemetry API 和 SDK 的更多信息，请参阅规范。

Configure the OpenTelemetry API
为了创建跟踪或指标，您需要首先创建一个 tracer 和/或 meter provider。一般来说，我们建议 SDK 应该为这些对象提供一个默认的 provider。然后，
您将从该 provider 处获得一个 tracer 或 meter instance，并为其指定名称和版本。您在此处选择的名称应标识所检测的确切内容——例如，如果您正在
编写一个库，则应以您的库命名它（例如 com.legitimatebusiness.myLibrary），因为此名称将为所有 spans 或 metric events produced。还建议
您提供与您的库或服务的当前版本相对应的版本字符串（即 semver:1.0.0）。

Configure the OpenTelemetry SDK
如果您正在构建服务流程，您还需要使用适当的选项配置 SDK，以便将遥测数据导出到某些分析后端。我们建议通过配置文件或其他机制以编程方式处理此配置。
您可能希望利用每种语言的调整选项。

Create Telemetry Data
配置完 API 和 SDK 后，您就可以通过从 provider 处获得的 tracer 和 meter 对象自由地创建 traces 和 metric events。为您的依赖项使用
Instrumentation Libraries——查看 Registry 或您的语言的存储库以获取更多信息。

Export Data
创建遥测数据后，您将希望将其发送到某个地方。 OpenTelemetry 支持两种将数据从流程导出到分析后端的主要方法，直接从流程或通过 OpenTelemetry Collector 代理它。

进程内导出要求您导入并依赖一个或多个 exporters, 将 OpenTelemetry 的内存 span 和 metric 对象转换为适合遥测分析工具（如 Jaeger 或 Prometheus）
的格式的库。此外，OpenTelemetry 支持称为 OTLP 的有线协议，所有 OpenTelemetry SDK 都支持该协议。此协议可用于将数据发送到 OpenTelemetry Collector
，这是一个独立的二进制进程，可以作为服务实例的代理或 sidecar 运行，也可以在单独的主机上运行。然后可以将收集器配置为将此数据转发和导出到您选择
的分析工具。

除了 Jaeger 或 Prometheus 等开源工具之外，越来越多的公司支持从 OpenTelemetry 获取遥测数据。有关详细信息，请参阅供应商。
