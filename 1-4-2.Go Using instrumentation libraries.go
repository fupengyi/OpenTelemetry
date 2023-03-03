package OpenTelemetry_Go

import (
	"context"
	"fmt"
	"net/http"
	"time"
)

Using instrumentation libraries
Go 不像今天的其他语言那样支持真正的自动检测。相反，您需要依赖为特定检测库生成遥测数据的检测库。例如，一旦您在代码中配置了 net/http 的检测库
，它就会自动创建跟踪入站和出站请求的跨度。


Setup
每个检测库都是一个包。一般来说，这意味着你需要去获取合适的包：
	go get go.opentelemetry.io/contrib/instrumentation/{import-path}/otel{package-name}
然后根据库需要激活的内容在您的代码中配置它。


Example with net/http
例如，您可以通过以下方式为 net/http 的入站 HTTP 请求设置自动检测：
首先，获取 net/http 检测库：
	go get go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp

接下来，使用该库在您的代码中包装 HTTP 处理程序：
package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"
	
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
)

// Package-level tracer.														// 包级跟踪器。
// This should be configured in your code setup instead of here.				// 这应该在您的代码设置中而不是在这里配置。
var tracer = otel.Tracer("github.com/full/path/to/mypkg")

// sleepy mocks work that your application does.								// sleepy mocks 你的应用程序所做的工作。
func sleepy(ctx context.Context) {
	_, span := tracer.Start(ctx, "sleep")
	defer span.End()
	
	sleepTime := 1 * time.Second
	time.Sleep(sleepTime)
	span.SetAttributes(attribute.Int("sleep.duration", int(sleepTime)))
}

// httpHandler is an HTTP handler function that is going to be instrumented.	// httpHandler 是将要检测的 HTTP 处理程序函数。
func httpHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, World! I am instrumented automatically!")
	ctx := r.Context()
	sleepy(ctx)
}

func main() {
	// Wrap your httpHandler function.											// 包装你的 httpHandler 函数。
	handler := http.HandlerFunc(httpHandler)
	wrappedHandler := otelhttp.NewHandler(handler, "hello-instrumented")
	http.Handle("/hello-instrumented", wrappedHandler)
	
	// And start the HTTP serve.												// 并启动 HTTP 服务。
	log.Fatal(http.ListenAndServe(":3030", nil))
}
假设您配置了 Tracer 和导出器，此代码将：
	在端口 3030 上启动 HTTP 服务器
	为每个到 /hello-instrumented 的入站 HTTP 请求自动生成一个跨度
	创建自动生成的子跨度，跟踪在 sleepy 中完成的工作
将您在应用程序中编写的手动检测与从库中生成的检测相连接对于在您的应用和服务中获得良好的可观察性至关重要。


Available packages
可以在 OpenTelemetry 注册表中找到可用仪器库的完整列表。


Next steps
检测库可以执行诸如为入站和出站 HTTP 请求生成遥测数据之类的操作，但它们不会检测您的实际应用程序。
要获得更丰富的遥测数据，请使用手动检测通过运行应用程序中的检测来丰富检测库中的遥测数据。
