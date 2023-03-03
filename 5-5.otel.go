package OpenTelemetry_Go

import (
	"log"
	"os"
)

"go.opentelemetry.io/otel"

Functions ¶
1.func GetTextMapPropagator() propagation.TextMapPropagator			// GetTextMapPropagator 返回全局 TextMapPropagator。如果未设置，则返回 No-Op TextMapPropagator。
2.func GetTracerProvider() trace.TracerProvider						// GetTracerProvider 返回已注册的全局跟踪提供程序。如果没有注册，则返回一个 NoopTracerProvider 实例。
																	// 使用跟踪提供程序创建命名跟踪器。例如。
																	// tracer := otel.GetTracerProvider().Tracer("example.com/foo")
																	// or
																	// tracer := otel.Tracer("example.com/foo")
3.func Handle(err error)											// Handle 是 ErrorHandler().Handle(err) 的便利函数。
4.func SetErrorHandler(h ErrorHandler)								// SetErrorHandler 将全局 ErrorHandler 设置为 h。
																	// 第一次调用之前从 GetErrorHandler 返回的所有 ErrorHandler 都会将错误发送到 h 而不是默认的日志记录 ErrorHandler。
																	// 后续调用将设置全局 ErrorHandler，但不会将错误委托给 h。
5.func SetLogger(logger logr.Logger)								// SetLogger 配置内部使用的记录器以打开遥测。
Example:
package main

import (
	 "log"
	 "os"

	 "github.com/go-logr/stdr"

	 "go.opentelemetry.io/otel"
)

func main() {
	 logger := stdr.New(log.New(os.Stdout, "", log.LstdFlags|log.Lshortfile))
	 otel.SetLogger(logger)
}

6.func SetTextMapPropagator(propagator propagation.TextMapPropagator)	// SetTextMapPropagator 将传播器设置为全局 TextMapPropagator。
7.func SetTracerProvider(tp trace.TracerProvider)						// SetTracerProvider 将 `tp` 注册为global trace provider。
8.func Tracer(name string, opts ...trace.TracerOption) trace.Tracer		// Tracer 创建一个实现 Tracer 接口的命名跟踪器。如果名称为空字符串，则提供程序使用默认名称。
																		// 这是 GetTracerProvider().Tracer(name, opts...) 的缩写
9.func Version() string													// Version 是正在使用的 OpenTelemetry 的当前发行版本。



Types ¶
type ErrorHandler interface {
	// Handle handles any error deemed irremediable by an OpenTelemetry component.	// Handle 处理 OpenTelemetry 组件认为无法修复的任何错误。
	Handle(error)
}
1.func GetErrorHandler() ErrorHandler		// GetErrorHandler 返回全局 ErrorHandler 实例。
											// 返回的默认 ErrorHandler 实例会将所有错误记录到 STDERR，直到使用 SetErrorHandler 设置覆盖 ErrorHandler。
											// 在此之前返回的所有 ErrorHandler 将自动将错误转发到设置的实例而不是 logging。
											// 在第一次调用之后对 SetErrorHandler 的后续调用不会将错误转发给先前返回的实例的新 ErrorHandler。

										
											
type ErrorHandlerFunc func(error)			// ErrorHandlerFunc 是一个方便的适配器，允许将函数用作 ErrorHandler。
1.func (f ErrorHandlerFunc) Handle(err error)// Handle 通过调用 ErrorHandlerFunc 本身来处理无法补救的错误。
