package OpenTelemetry_Go

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"os/signal"
	"strconv"
)

Getting Started
欢迎使用 OpenTelemetry for Go 入门指南！本指南将引导您完成从 OpenTelemetry 安装、检测、配置和导出数据的基本步骤。在开始之前，请务必安装 Go 1.16 或更新版本。

了解系统在出现故障或出现问题时的运行方式对于解决这些问题至关重要。理解这一点的一种策略是跟踪。本指南展示了如何使用 OpenTelemetry Go 项目来
跟踪示例应用程序。您将从一个为用户计算斐波那契数的应用程序开始，然后您将从那里添加工具以使用 OpenTelemetry Go 生成跟踪遥测数据。

作为参考，您可以在此处找到您将构建的代码的完整示例。
要开始构建应用程序，请创建一个名为 fib 的新目录来存放我们的 Fibonacci 项目。接下来，将以下内容添加到该目录中名为 fib.go 的新文件中。
package main

// Fibonacci returns the n-th fibonacci number.							// Fibonacci 返回第 n 个斐波那契数。
func Fibonacci(n uint) (uint64, error) {
	if n <= 1 {
		return uint64(n), nil
	}
	
	var n2, n1 uint64 = 0, 1
	for i := uint(2); i < n; i++ {
		n2, n1 = n1, n1+n2
	}
	
	return n2 + n1, nil
}
========================================================================================================================

添加核心逻辑后，您现在可以围绕它构建应用程序。使用以下应用程序逻辑添加一个新的 app.go 文件。
package main

import (
	"context"
	"fmt"
	"io"
	"log"
)

// App is a Fibonacci computation application.									// App 是一个斐波那契计算应用程序。
type App struct {
	r io.Reader
	l *log.Logger
}

// NewApp returns a new App.													// NewApp 返回一个新的 App。
func NewApp(r io.Reader, l *log.Logger) *App {
	return &App{r: r, l: l}
}

// Run starts polling users for Fibonacci number requests and writes results.	// Run 开始轮询用户的斐波那契数字请求并写入结果。
func (a *App) Run(ctx context.Context) error {
	for {
		n, err := a.Poll(ctx)
		if err != nil {
			return err
		}
		
		a.Write(ctx, n)
	}
}

// Poll asks a user for input and returns the request.							// Poll 要求用户输入并返回请求。
func (a *App) Poll(ctx context.Context) (uint, error) {
	a.l.Print("What Fibonacci number would you like to know: ")
	
	var n uint
	_, err := fmt.Fscanf(a.r, "%d\n", &n)
	return n, err
}

// Write writes the n-th Fibonacci number back to the user.						// Write 将第 n 个斐波那契数写回给用户。
func (a *App) Write(ctx context.Context, n uint) {
	f, err := Fibonacci(n)
	if err != nil {
		a.l.Printf("Fibonacci(%d): %v\n", n, err)
	} else {
		a.l.Printf("Fibonacci(%d) = %d\n", n, f)
	}
}
========================================================================================================================

在您的应用程序完全组合后，您需要一个 main() 函数来实际运行该应用程序。在新的 main.go 文件中添加以下运行逻辑。
package main

import (
	"context"
	"log"
	"os"
	"os/signal"
)

func main() {
	l := log.New(os.Stdout, "", 0)
	
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt)
	
	errCh := make(chan error)
	app := NewApp(os.Stdin, l)
	go func() {
		errCh <- app.Run(context.Background())
	}()
	
	select {
	case <-sigCh:
		l.Println("\ngoodbye")
		return
	case err := <-errCh:
		if err != nil {
			l.Fatal(err)
		}
	}
}
========================================================================================================================

代码完成后，几乎可以运行应用程序了。在执行此操作之前，您需要将此目录初始化为 Go 模块。从您的终端，在 fib 目录中运行命令 go mod init fib。
这将创建一个 go.mod 文件，Go 使用它来管理导入。现在您应该可以运行该应用程序了！
	$ go run .
	What Fibonacci number would you like to know:
	42
	Fibonacci(42) = 267914296
	What Fibonacci number would you like to know:
	^C
	goodbye
可以使用 Ctrl+C 退出应用程序。您应该会看到与上面类似的输出，如果没有，请确保返回并修复任何错误。
========================================================================================================================


Trace Instrumentation
OpenTelemetry 分为两部分：用于检测代码的 API 和实现 API 的 SDK。要开始将 OpenTelemetry 集成到任何项目中，API 用于定义遥测数据的生成方
式。要在您的应用程序中生成跟踪遥测，您将使用 go.opentelemetry.io/otel/trace 包中的 OpenTelemetry Trace API。

首先，您需要为 Trace API 安装必要的包。在您的工作目录中运行以下命令。
	go get go.opentelemetry.io/otel \
	go.opentelemetry.io/otel/trace

现在已经安装了包，您可以开始使用将在 app.go 文件中使用的导入来更新您的应用程序。
import (
	"context"
	"fmt"
	"io"
	"log"
	"strconv"
	
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)
添加导入后，您可以开始检测。

OpenTelemetry Tracing API 提供了一个 Tracer 来创建跟踪。这些跟踪器旨在与一个检测库相关联。这样，他们产生的遥测就可以理解为来自代码库的那
部分。为了向 Tracer 唯一标识您的应用程序，您将使用在 app.go 中使用包名称创建一个常量。
// name is the Tracer name used to identify this instrumentation library.	// name 是用于标识此检测库的 Tracer 名称。
const name = "fib"

使用完全限定的包名称，对于 Go 包来说应该是唯一的，是识别 Tracer 的标准方法。如果您的示例包名称不同，请务必更新您在此处使用的名称以匹配。
现在一切都应该就绪，可以开始跟踪您的应用程序了。但首先，什么是 trace ？而且，您应该如何为您的应用程序构建它们？
稍微备份一下，跟踪是一种遥测，表示服务正在完成的工作。跟踪是处理交易的参与者之间的连接记录，通常通过客户端/服务器请求处理和其他形式的通信。

服务执行的每个工作部分都在跟踪中由一个跨度表示。这些跨度不仅仅是一个无序的集合。就像我们应用程序的调用堆栈一样，这些跨度是通过彼此之间的关系来
定义的。 “根”跨度是唯一没有父跨度的跨度，它表示服务请求是如何开始的。所有其他跨度与同一跟踪中的另一个跨度具有父关系。

如果关于跨度关系的最后一部分现在没有完全理解，请不要担心。最重要的一点是，你的代码的每一部分，做一些工作，都应该表示为一个跨度。检测代码后，您
将对这些跨度关系有更好的理解，所以让我们开始吧。

首先检测 Run 方法。
// Run starts polling users for Fibonacci number requests and writes results.	// Run 开始轮询用户的斐波那契数字请求并写入结果。
func (a *App) Run(ctx context.Context) error {
	for {
		// Each execution of the run loop, we should get a new "root" span and context.		// 每次运行循环的执行，我们应该得到一个新的“根”跨度和上下文。
		newCtx, span := otel.Tracer(name).Start(ctx, "Run")
		
		n, err := a.Poll(newCtx)
		if err != nil {
			span.End()
			return err
		}
		
		a.Write(newCtx, n)
		span.End()
	}
}
上面的代码为 for 循环的每次迭代创建一个跨度。跨度是使用来自全局 TracerProvider 的 Tracer 创建的。在后面的部分中安装 SDK 时，您将了解有关
TracerProvider 的更多信息并处理设置全局 TracerProvider 的另一面。现在，作为仪器作者，您需要担心的是，当您编写 otel.Tracer(name) 时，您
正在使用来自 TracerProvider 的适当命名的 Tracer。

接下来，检测 Poll 方法。
// Poll asks a user for input and returns the request.			// Poll 要求用户输入并返回请求。
func (a *App) Poll(ctx context.Context) (uint, error) {
	_, span := otel.Tracer(name).Start(ctx, "Poll")
	defer span.End()
	
	a.l.Print("What Fibonacci number would you like to know: ")
	
	var n uint
	_, err := fmt.Fscanf(a.r, "%d\n", &n)
	
	// Store n as a string to not overflow an int64.			// 将 n 存储为字符串以防止 int64 溢出。
	nStr := strconv.FormatUint(uint64(n), 10)
	span.SetAttributes(attribute.String("request.n", nStr))
	
	return n, err
}
与 Run 方法检测类似，这会向方法添加一个跨度以跟踪执行的计算。但是，它还添加了一个属性来注释跨度。当您认为您的应用程序的用户在查看遥测时希望查
看有关运行环境的状态或详细信息时，可以添加此注释。

最后，检测 Write 方法。
// Write writes the n-th Fibonacci number back to the user.		// Write 将第 n 个斐波那契数写回给用户。
func (a *App) Write(ctx context.Context, n uint) {
	var span trace.Span
	ctx, span = otel.Tracer(name).Start(ctx, "Write")
	defer span.End()
	
	f, err := func(ctx context.Context) (uint64, error) {
		_, span := otel.Tracer(name).Start(ctx, "Fibonacci")
		defer span.End()
		return Fibonacci(n)
	}(ctx)
	if err != nil {
		a.l.Printf("Fibonacci(%d): %v\n", n, err)
	} else {
		a.l.Printf("Fibonacci(%d) = %d\n", n, f)
	}
}
此方法使用两个跨度进行检测。一个跟踪 Write 方法本身，另一个跟踪使用 Fibonacci 函数对核心逻辑的调用。你看到上下文是如何通过跨度传递的了吗？
您是否看到这也定义了跨度之间的关系？

在 OpenTelemetry Go 中，跨度关系是使用 context.Context 明确定义的。创建跨度时，上下文会与跨度一起返回。该上下文将包含对创建的跨度的引用
。如果在创建另一个跨度时使用该上下文，则两个跨度将相关。原始 span 将成为新 span 的父级，并且作为推论，新 span 被称为原始 span 的子级。该层
次结构提供了跟踪结构，该结构有助于显示通过系统的计算路径。根据您上面的检测和对跨度关系的理解，您应该期望运行循环的每次执行的跟踪看起来像这样。
Run
├── Poll
└── Write
	└── Fibonacci
Run span 将是 Poll 和 Write span 的父级，而 Write span 将是 Fibonacci span 的父级。
现在你如何真正看到产生的跨度？为此，您需要配置和安装 SDK。


SDK Installation
OpenTelemetry 在其 OpenTelemetry API 的实现中被设计为模块化。 OpenTelemetry Go 项目提供了一个 SDK 包，go.opentelemetry.io/otel/sdk
，它实现了这个 API 并遵守 OpenTelemetry 规范。要开始使用此 SDK，您首先需要创建一个导出器，但在任何事情发生之前，我们需要安装一些包。在 fib
目录中运行以下命令以安装跟踪 STDOUT 导出器和 SDK。
	go get go.opentelemetry.io/otel/sdk \
	go.opentelemetry.io/otel/exporters/stdout/stdouttrace

现在将所需的导入添加到 main.go。
import (
	"context"
	"io"
	"log"
	"os"
	"os/signal"
	
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"
)


Creating a Console Exporter
SDK 将来自 OpenTelemetry API 的遥测数据连接到导出器。导出器是允许将遥测数据发送到某处的包——发送到控制台（这就是我们在这里所做的），或者发
送到远程系统或收集器以进行进一步分析和/或丰富。 OpenTelemetry 通过其生态系统支持各种出口商，包括 Jaeger、Zipkin 和 Prometheus 等流行的开源工具。

要初始化控制台导出器，请将以下函数添加到 main.go 文件中：
// newExporter returns a console exporter.						// newExporter 返回一个控制台导出器。
func newExporter(w io.Writer) (trace.SpanExporter, error) {
	return stdouttrace.New(
		stdouttrace.WithWriter(w),
		// Use human-readable output.							// 使用人类可读的输出。
		stdouttrace.WithPrettyPrint(),
		// Do not print timestamps for the demo.				// 不要为演示打印时间戳。
		stdouttrace.WithoutTimestamps(),
	)
}
这将创建一个带有基本选项的新控制台导出器。稍后您将在配置 SDK 向其发送遥测数据时使用此功能，但首先您需要确保数据是可识别的。



Creating a Resource
遥测数据对于解决服务问题至关重要。问题是，您需要一种方法来识别数据来自哪个服务，甚至是哪个服务实例。 OpenTelemetry 使用资源来表示生成遥测的
实体。将以下函数添加到 main.go 文件中，为应用程序创建适当的资源。
// newResource returns a resource describing this application.
func newResource() *resource.Resource {
	r, _ := resource.Merge(
		resource.Default(),
		resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceName("fib"),
			semconv.ServiceVersion("v0.1.0"),
			attribute.String("environment", "demo"),
		),
	)
	return r
}
您想要与 SDK 处理的所有遥测数据相关联的任何信息都可以添加到返回的资源中。这是通过向 TracerProvider 注册资源来完成的。您现在可以创造的东西！


Installing a Tracer Provider
您已对您的应用程序进行检测以生成遥测数据，并且您有一个导出器将该数据发送到控制台，但它们是如何连接的？这是使用 TracerProvider 的地方。这是
一个集中点，仪器将从中获取 Tracer，并将来自这些 Tracer 的遥测数据汇集到导出管道。

接收数据并最终将数据传输到导出器的管道称为 SpanProcessors。 TracerProvider 可以配置为具有多个 span 处理器，但对于此示例，您只需要配置一
个。使用以下内容更新 main.go 中的主要功能。
func main() {
	l := log.New(os.Stdout, "", 0)
	
	// Write telemetry data to a file.									// 将遥测数据写入文件。
	f, err := os.Create("traces.txt")
	if err != nil {
		l.Fatal(err)
	}
	defer f.Close()
	
	exp, err := newExporter(f)
	if err != nil {
		l.Fatal(err)
	}
	
	tp := trace.NewTracerProvider(
		trace.WithBatcher(exp),
		trace.WithResource(newResource()),
	)
	defer func() {
		if err := tp.Shutdown(context.Background()); err != nil {
			l.Fatal(err)
		}
	}()
	otel.SetTracerProvider(tp)
	
	/* … */
}
这里发生了很多事情。首先，您要创建一个将导出到文件的控制台导出器。然后，您将向新的 TracerProvider 注册导出器。当 BatchSpanProcessor 传递
给 trace.WithBatcher 选项时，这是通过 BatchSpanProcessor 完成的。批处理数据是一种很好的做法，有助于避免下游系统过载。最后，创建 TracerProvider
后，您将延迟一个函数来刷新和停止它，并将其注册为全局 OpenTelemetry TracerProvider。

您还记得在之前的检测部分中我们使用全局 TracerProvider 来获取 Tracer 吗？最后一步，即全局注册 TracerProvider，会将仪器的 Tracer 与此 TracerProvider
连接起来。这种使用全局 TracerProvider 的模式很方便，但并不总是合适的。 TracerProviders 可以显式传递给检测或从包含跨度的上下文中推断出来
。对于这个使用全局提供程序的简单示例是有意义的，但对于更复杂或分布式代码库，这些其他传递 TracerProviders 的方法可能更有意义。


Putting It All Together
您现在应该有一个可以生成跟踪遥测数据的工作应用程序！试一试。
	$ go run .
	What Fibonacci number would you like to know:
	42
	Fibonacci(42) = 267914296
	What Fibonacci number would you like to know:
	^C
	goodbye
应在您的工作目录中创建一个名为 traces.txt 的新文件。运行应用程序时创建的所有跟踪都应该在那里！


(Bonus) Errors
此时您有一个工作应用程序，它正在生成跟踪遥测数据。不幸的是，发现 fib 模块的核心功能有错误。
	$ go run .
	What Fibonacci number would you like to know:
	100
	Fibonacci(100) = 3736710778780434371
	# …
	
但是第 100 个斐波那契数是 354224848179261915075，而不是 3736710778780434371！此应用程序仅作为演示，但不应返回错误值。更新 Fibonacci
函数以返回错误而不是计算不正确的值。
// Fibonacci returns the n-th fibonacci number. An error is returned if the			// Fibonacci 返回第 n 个斐波那契数。如果
// fibonacci number cannot be represented as a uint64.								// 斐波那契数不能表示为 uint64。
func Fibonacci(n uint) (uint64, error) {
	if n <= 1 {
		return uint64(n), nil
	}
	
	if n > 93 {
		return 0, fmt.Errorf("unsupported fibonacci number %d: too large", n)
	}
	
	var n2, n1 uint64 = 0, 1
	for i := uint(2); i < n; i++ {
		n2, n1 = n1, n1+n2
	}
	
	return n2 + n1, nil
}

太好了，您已经修复了代码，但最好在遥测数据中包含返回给用户的错误。幸运的是，跨度可以配置为传达此信息。使用以下代码更新 app.go 中的 Write 方法。
// Write writes the n-th Fibonacci number back to the user.		// Write 将第 n 个斐波那契数写回给用户。
func (a *App) Write(ctx context.Context, n uint) {
	var span trace.Span
	ctx, span = otel.Tracer(name).Start(ctx, "Write")
	defer span.End()
	
	f, err := func(ctx context.Context) (uint64, error) {
		_, span := otel.Tracer(name).Start(ctx, "Fibonacci")
		defer span.End()
		f, err := Fibonacci(n)
		if err != nil {
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())
		}
		return f, err
	}(ctx)
	/* … */
}
通过此更改，从 Fibonacci 函数返回的任何错误都会将该跨度标记为错误并记录一个描述错误的事件。

这是一个很好的开始，但它不是应用程序返回的唯一错误。如果用户请求非无符号整数值，应用程序将失败。使用类似的修复更新 Poll 方法以捕获遥测数据中的此错误。
// Poll asks a user for input and returns the request.					// Poll 要求用户输入并返回请求。
func (a *App) Poll(ctx context.Context) (uint, error) {
	_, span := otel.Tracer(name).Start(ctx, "Poll")
	defer span.End()
	
	a.l.Print("What Fibonacci number would you like to know: ")
	
	var n uint
	_, err := fmt.Fscanf(a.r, "%d\n", &n)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return 0, err
	}
	
	// Store n as a string to not overflow an int64.					// 将 n 存储为字符串以防止 int64 溢出。
	nStr := strconv.FormatUint(uint64(n), 10)
	span.SetAttributes(attribute.String("request.n", nStr))
	
	return n, nil
}

剩下的就是更新 app.go 文件的导入以包含 go.opentelemetry.io/otel/codes 包。
import (
	"context"
	"fmt"
	"io"
	"log"
	"strconv"
	
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

在这些修复到位并更新检测后，重新触发错误。
	$ go run .
	What Fibonacci number would you like to know:
	100
	Fibonacci(100): unsupported fibonacci number 100: too large
	What Fibonacci number would you like to know:
	^C
	goodbye
	
出色的！应用程序不再返回错误值，并且查看 traces.txt 文件中的遥测数据，您应该会看到作为事件捕获的错误。
"Events": [
	{
		"Name": "exception",
		"Attributes": [
			{
				"Key": "exception.type",
				"Value": {
				"Type": "STRING",
				"Value": "*errors.errorString"
				}
			},
			{
				"Key": "exception.message",
				"Value": {
				"Type": "STRING",
				"Value": "unsupported fibonacci number 100: too large"
				}
			}
		],
		...
	}
]


What’s Next
本指南已指导您完成向应用程序添加跟踪检测以及使用控制台导出器将遥测数据发送到文件的过程。 OpenTelemetry 中还有许多其他主题要涵盖，但此时您应
该准备好开始将 OpenTelemetry Go 添加到您的项目中。去检测你的代码！

有关检测代码以及可以使用 span 执行的操作的更多信息，请参阅手动检测文档。
您还需要配置一个适当的导出器以将您的遥测数据导出到一个或多个遥测后端。
