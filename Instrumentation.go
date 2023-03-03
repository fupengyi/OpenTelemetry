package OpenTelemetry_Go
Instrumentation
此目录中包含的代码包含对第 3 方 Go 包和标准库中的一些包的检测。

New Instrumentation
请勿在未阅读以下内容的情况下提交新工具的拉取请求。

该项目致力于使用 OpenTelemetry 促进质量仪器的开发。为了实现这一目标，我们认识到需要使用 OpenTelemetry 的最佳实践来编写检测，而且还需要由
了解他们正在检测的包的开发人员编写。此外，生产的仪器需要维护和发展。

OpenTelemetry Go 开发人员社区的规模不足以支持不断增长的检测量。因此，为了实现我们的目标，我们对仪器包的存放位置提出了以下建议。
	1.原生于检测包
	2.专用的公共存储库
	3.在 opentelemetry-go-contrib 存储库中
如果可能，OpenTelemetry 检测应包含在检测包中。这将确保工具到达所有包用户，并由了解包的开发人员持续维护。

如果检测不能直接包含在它正在检测的包中，则应将其托管在其维护者拥有的专用公共存储库中。这将适当地分配仪器的维护职责，并确保这些维护人员拥有维护
代码所需的特权。

最后一个应该托管仪器的地方是在这个存储库中。在这里维护仪器会阻碍 Go 的 OpenTelemetry 的开发，因此应该避免。当检测不能包含在目标包中并且有充
分的理由不将其托管在单独的专用存储库中时，应提交检测请求。在可以考虑合并仪器的任何拉取请求之前，需要接受仪器请求。

无论仪器托管在何处，都需要能够被发现。存在 OpenTelemetry 注册表以确保可以发现检测。您可以在此处了解如何将工具添加到注册表。

Instrumentation Packages
OpenTelemetry 注册表是发现检测包的最佳位置。它将包括该项目之外的包。
为流行的 Go 包和用例提供了以下检测包。
Instrumentation Package						Metrics		Traces
github.com/astaxie/beego					✓			✓
github.com/aws/aws-sdk-go-v2							✓
...
runtime										✓

Organization
为了确保检测包的可维护性和可发现性，必须遵循以下准则。

Packaging
所有仪器包应该是以下形式：
	go.opentelemetry.io/contrib/instrumentation/{IMPORT_PATH}/otel{PACKAGE_NAME}
{IMPORT_PATH} 和 {PACKAGE_NAME} 是被检测包的标准 Go 标识符。
For example:
	go.opentelemetry.io/contrib/instrumentation/github.com/gorilla/mux/otelmux
	go.opentelemetry.io/contrib/instrumentation/gopkg.in/macaron.v1/otelmacaron
	go.opentelemetry.io/contrib/instrumentation/database/sql/otelsql
此规则存在例外情况。例如，runtime 和 host 检测不检测任何 Go 包，因此不适合这种结构。

Contents
所有仪器包必须遵守项目的贡献指南。此外，还需要遵循以下包装组成指南。
	所有检测包必须是 Go 模块。因此，每个包都需要存在适当配置的 go.mod 和 go.sum。
	为了帮助理解工具，应该包含一个 Go 包文档。如果包中有多个文件，则此文档应该位于专用的 doc.go 文件中。它应该包含有用的信息，例如检测的目的是什么、如何使用它以及可能存在的任何兼容性限制。
	应该包括如何实际使用仪器的示例。
	所有检测包必须提供一个选项，如果它使用 Tracer，则接受 TracerProvider，如果它使用 Meter，则接受 MeterProvider，如果它处理任何上下文传播，则提供 Propagators。此外，如果没有提供可选的，包必须使用全局包提供的默认 TracerProvider、MeterProvider 和 Propagator。
	所有仪器包不得提供接受示踪剂或仪表的选项。
	所有检测包必须创建任何使用的 Tracer 或 Meter，其名称与检测包名称匹配。
	所有检测包必须创建任何使用过的 Tracer 或 Meter，其语义版本与包含检测的模块的版本相对应。
