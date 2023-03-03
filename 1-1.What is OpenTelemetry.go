package OpenTelemetry_Go
What is OpenTelemetry?
关于 OpenTelemetry 的背景信息。
微服务架构使开发人员能够更快地构建和发布软件，并具有更大的独立性，因为他们不再受制于与单体架构相关的复杂发布流程。
随着这些现在分布式系统的扩展，开发人员越来越难以了解他们自己的服务如何依赖或影响其他服务，尤其是在部署之后或中断期间，速度和准确性至关重要。
可观察性使开发人员和运维人员都有可能获得对其系统的可见性。

So what?
为了使系统可观察，必须对其进行检测。也就是说，代码必须发出跟踪、指标和日志。然后必须将经过检测的数据发送到可观察性后端。有许多可观察性后端，从
自托管开源工具（例如 Jaeger 和 Zipkin）到商业 SaaS 产品。

过去，检测代码的方式会有所不同，因为每个可观察性后端都有自己的检测库和代理，用于向工具发送数据。

这意味着没有用于将数据发送到可观察性后端的标准化数据格式。此外，如果一家公司选择切换 Observability 后端，这意味着他们将不得不重新检测他们的
代码并配置新代理，以便能够将遥测数据发送到新选择的工具。

由于缺乏标准化，最终结果是缺乏数据可移植性和用户维护仪器库的负担。
认识到标准化的必要性后，云社区聚集在一起，诞生了两个开源项目：OpenTracing（一个云原生计算基金会 (CNCF) 项目）和 OpenCensus（一个谷歌开源社区项目）。
OpenTracing 提供了供应商中立的 API，用于将遥测数据发送到可观察性后端；但是，它依赖于开发人员实现自己的库来满足规范。
OpenCensus 提供了一组特定于语言的库，开发人员可以使用这些库来检测他们的代码并发送到他们支持的任何一个后端。

Hello, OpenTelemetry!
为了拥有一个单一的标准，OpenCensus 和 OpenTracing 于 2019 年 5 月合并形成了 OpenTelemetry（简称 OTel）。作为 CNCF 孵化项目，OpenTelemetry
充分利用了两个世界的优点，然后是一些。

OTel 的目标是提供一组标准化的与供应商无关的 SDK、API 和工具，用于摄取、转换数据并将数据发送到可观察性后端（即开源或商业供应商）。

What can OpenTelemetry do for me?
OTel 拥有来自云提供商、供应商和最终用户的广泛行业支持和采用。它为您提供：
	每种语言都有一个独立于供应商的仪器库，支持自动和手动仪器。
	可以以多种方式部署的单个供应商中立的收集器二进制文件。
	生成、发出、收集、处理和导出遥测数据的端到端实现。
	完全控制您的数据，能够通过配置将数据并行发送到多个目的地。
	开放标准语义约定以确保与供应商无关的数据收集
	能够并行支持多种上下文传播格式，以协助随着标准的发展进行迁移。
	无论您在可观察性之旅的哪个阶段，都有一条前进的道路。
由于支持各种开源和商业协议、格式和上下文传播机制以及为 OpenTracing 和 OpenCensus 项目提供填充程序，因此很容易采用 OpenTelemetry。

What OpenTelemetry is not
OpenTelemetry 不是像 Jaeger 或 Prometheus 那样的可观察性后端。相反，它支持将数据导出到各种开源和商业后端。它提供了一个可插入的架构，因此
可以轻松添加额外的技术协议和格式。