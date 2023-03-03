.package OpenTelemetry_GoConcepts

import (
	"context"
	"net/http"
)

Manual Instrumentation
Instrumentation 是将可观察性代码添加到您的应用程序的过程。有两种通用类型的检测 - 自动和手动 - 您应该熟悉这两种类型，以便有效地检测您的软件。


Getting a Tracer
要创建跨度，您需要先获取或初始化跟踪器。


Initializing a new tracer
确保安装了正确的软件包：
	go get go.opentelemetry.io/otel \
	go.opentelemetry.io/otel/trace \
	go.opentelemetry.io/otel/sdk \

然后初始化导出器、资源、跟踪器提供程序，最后是跟踪器。
package app

import (
	"context"
	"fmt"
	"log"
	
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"
	"go.opentelemetry.io/otel/trace"
)

var tracer trace.Tracer

func newExporter(ctx context.Context)  /* (someExporter.Exporter, error) */ {
	// Your preferred exporter: console, jaeger, zipkin, OTLP, etc.				// 你首选的导出器：console、jaeger、zipkin、OTLP 等。
}

func newTraceProvider(exp sdktrace.SpanExporter) *sdktrace.TracerProvider {
	// Ensure default SDK resources and the required service name are set.		// 确保设置了默认的 SDK 资源和所需的服务名称。
	r, err := resource.Merge(
		resource.Default(),
		resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceName("ExampleService"),
		),
	)
	
	if err != nil {
		panic(err)
	}
	
	return sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exp),
		sdktrace.WithResource(r),
	)
}

func main() {
	ctx := context.Background()
	
	exp, err := newExporter(ctx)
	if err != nil {
		log.Fatalf("failed to initialize exporter: %v", err)
	}
	
	// Create a new tracer provider with a batch span processor and the given exporter.		// 使用批处理跨度处理器和给定的导出器创建一个新的跟踪器提供程序。
	tp := newTraceProvider(exp)
	
	// Handle shutdown properly so nothing leaks.											// 正确处理关闭，所以没有泄漏。
	defer func() { _ = tp.Shutdown(ctx) }()
	
	otel.SetTracerProvider(tp)
	
	// Finally, set the tracer that can be used for this package.							// 最后，设置可用于此包的跟踪器。
	tracer = tp.Tracer("ExampleService")
}
您现在可以访问跟踪器以手动检测您的代码。


Creating Spans
跨度由跟踪器创建。如果你没有初始化，你需要这样做。
要使用跟踪器创建跨度，您还需要 context.Context 实例的句柄。这些通常来自请求对象之类的东西，并且可能已经包含来自检测库的父跨度。
func httpHandler(w http.ResponseWriter, r *http.Request) {
	ctx, span := tracer.Start(r.Context(), "hello-span")
	defer span.End()
	
	// do some work to track with hello-span												// 做一些工作来跟踪 hello-span
}
在 Go 中，上下文包用于存储活动跨度。当您启动一个跨度时，您不仅会得到创建的跨度的句柄，还会得到包含它的修改后的上下文。
一旦跨度完成，它就是不可变的并且不能再被修改。


Get the current span
要获取当前跨度，您需要将其从上下文中拉出。您可以处理的上下文：
// This context needs contain the active span you plan to extract.								// 此上下文需要包含您计划提取的活动跨度。
ctx := context.TODO()
span := trace.SpanFromContext(ctx)

// Do something with the current span, optionally calling `span.End()` if you want it to end	// 对当前跨度做一些事情，如果你想让它结束，可以选择调用 `span.End()`
如果您想在某个时间点向当前跨度添加信息，这会很有帮助。


Create nested spans
您可以创建嵌套跨度来跟踪嵌套操作中的工作。
如果当前的 context.Context 你有句柄已经在其中包含一个跨度，创建一个新的跨度会使它成为一个嵌套的跨度。例如：
func parentFunction(ctx context.Context) {
	ctx, parentSpan := tracer.Start(ctx, "parent")
	defer parentSpan.End()
	
	// call the child function and start a nested span in there											// 调用子函数并在其中启动一个嵌套的 span
	childFunction(ctx)
	
	// do more work - when this function ends, parentSpan will complete.								// 做更多的工作 - 当这个函数结束时，parentSpan 将完成。
}

func childFunction(ctx context.Context) {
	// Create a span to track `childFunction()` - this is a nested span whose parent is `parentSpan`	// 创建一个 span 来跟踪 `childFunction()` - 这是一个嵌套的 span，它的父级是 `parentSpan`
	ctx, childSpan := tracer.Start(ctx, "child")
	defer childSpan.End()
	
	// do work here, when this function returns, childSpan will complete.								// 在这里工作，当这个函数返回时，childSpan 将完成。
}
一旦跨度完成，它就是不可变的并且不能再被修改。


Span Attributes
属性是作为元数据应用于跨度的键和值，可用于聚合、过滤和分组跟踪。可以在跨度创建时添加属性，也可以在跨度完成之前的生命周期中的任何其他时间添加属性。
// setting attributes at creation...																	// 在创建时设置属性...
ctx, span = tracer.Start(ctx, "attributesAtCreation", trace.WithAttributes(attribute.String("hello", "world")))
// ... and after creation																				// ... 创建后
span.SetAttributes(attribute.Bool("isTrue", true), attribute.String("stringAttr", "hi!"))

也可以预先计算属性键：
var myKey = attribute.Key("myCoolAttribute")
span.SetAttributes(myKey.String("a value"))


Semantic Attributes
语义属性是由 OpenTelemetry 规范定义的属性，目的是为 HTTP 方法、状态代码、用户代理等常见概念提供跨多种语言、框架和运行时的一组共享属性键。这
些属性在 go.opentelemetry.io/otel/semconv/v1.12.0 包中可用。
有关详细信息，请参阅跟踪语义约定。


Events
事件是跨度上的人类可读消息，表示在其生命周期内“发生的事情”。例如，假设一个函数需要独占访问互斥量下的资源。可以在两点创建一个事件 - 一次是在我
们尝试访问资源时，另一次是在我们获取互斥锁时。
span.AddEvent("Acquiring lock")
mutex.Lock()
span.AddEvent("Got lock, doing work...")
// do stuff		// 做东西
span.AddEvent("Unlocking")
mutex.Unlock()
事件的一个有用特征是它们的时间戳显示为从跨度开始的偏移量，使您可以轻松查看它们之间经过了多少时间。
事件也可以有自己的属性——
span.AddEvent("Cancelled wait due to external signal", trace.WithAttributes(attribute.Int("pid", 4328), attribute.String("signal", "SIGHUP")))


Set span status
可以在跨度上设置状态，通常用于指定跨度正在跟踪的操作中存在错误 - .Error。
import (
	// ...
	"go.opentelemetry.io/otel/codes"
	// ...
)

// ...

result, err := operationThatCouldFail()
if err != nil {
	span.SetStatus(codes.Error, "operationThatCouldFail failed")
}
默认情况下，所有跨度的状态都是未设置。在极少数情况下，您可能还希望将状态设置为 Ok。不过，这通常不是必需的。


Record errors
如果您有一个失败的操作并且您希望捕获它产生的错误，您可以记录该错误。
import (
	// ...
	"go.opentelemetry.io/otel/codes"
	// ...
)

// ...

result, err := operationThatCouldFail()
if err != nil {
	span.SetStatus(codes.Error, "operationThatCouldFail failed")
	span.RecordError(err)
}
强烈建议您在使用 RecordError 时也将跨度的状态设置为错误，除非您不希望将跟踪失败操作的跨度视为错误跨度。 RecordError 函数在调用时不会自动设置跨度状态。


Creating Metrics
指标 API 目前不稳定，文档待定。


Propagators and Context
跟踪可以扩展到单个进程之外。这需要上下文传播，这是一种将跟踪标识符发送到远程进程的机制。
为了通过线路传播跟踪上下文，必须使用 OpenTelemetry API 注册传播者。
import (
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
)
...
otel.SetTextMapPropagator(propagation.TraceContext{})
OpenTelemetry 还支持 B3 标头格式，以便与不支持 W3C TraceContext 标准的现有跟踪系统 (go.opentelemetry.io/contrib/propagators/b3) 兼容。
配置上下文传播后，您很可能希望使用自动检测来处理实际管理序列化上下文的幕后工作。
