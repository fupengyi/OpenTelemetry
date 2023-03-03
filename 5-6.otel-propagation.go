package OpenTelemetry_Go

import (
	"context"
	"net/http"
)

go.opentelemetry.io/otel/propagation

Overview ¶
包 propagation 包含 OpenTelemetry 上下文传播器。
OpenTelemetry 传播器用于从应用程序交换的消息中提取和注入上下文数据。此包支持的传播器是 W3C Trace Context 编码（https://www.w3.org/TR/trace-context/）和 W3C Baggage（https://www.w3.org/TR/baggage/） .


Types ¶
type Baggage struct{}					// Baggage 是一个支持 W3C Baggage 格式的传播者。
										// 这 propagates 与 trace 关联的用户定义的 baggage。完整的规范定义在 https://www.w3.org/TR/baggage/。
1.func (b Baggage) Extract(parent context.Context, carrier TextMapCarrier) context.Context		// Extract 返回一份 parent 的副本，其中添加了承运人的 baggage。
2.func (b Baggage) Fields() []string															// Fields 返回使用 Inject 设置值的键。
3.func (b Baggage) Inject(ctx context.Context, carrier TextMapCarrier)							// Inject 将来自 ctx 的 baggage 键值设置到 carrier 中。



type HeaderCarrier http.Header							// HeaderCarrier 适配 http.Header 以满足 TextMapCarrier 接口。
1.func (hc HeaderCarrier) Get(key string) string		// Get 返回与传递的键关联的值。
2.func (hc HeaderCarrier) Keys() []string				// Keys 列出存储在该载体中的密钥。
3.func (hc HeaderCarrier) Set(key string, value string)	// Set 存储键值对。



type MapCarrier map[string]string						// MapCarrier 是一种 TextMapCarrier，它使用保存在内存中的映射作为传播的键值对的存储介质。
1.func (c MapCarrier) Get(key string) string			// Get 返回与传递的键关联的值。
2.func (c MapCarrier) Keys() []string					// Keys 列出存储在该载体中的密钥。
3.func (c MapCarrier) Set(key, value string)			// Set 存储键值对。



type TextMapCarrier interface {									// TextMapCarrier 是 TextMapPropagator 使用的存储介质。
	// Get returns the value associated with the passed key.	// Get 返回与传递的键关联的值。
	Get(key string) string
	
	// Set stores the key-value pair.							// Set 存储键值对。
	Set(key string, value string)
	
	// Keys lists the keys stored in this carrier.				// Keys 列出存储在这个载体中的密钥。
	Keys() []string
}



type TextMapPropagator interface {												// TextMapPropagator 将横切关注点传播为跨进程边界带内传输的载体中的键值文本对。
	// Inject set cross-cutting concerns from the Context into the carrier.		// 将 Context 中的集合横切关注点注入到载体中。
	Inject(ctx context.Context, carrier TextMapCarrier)
	
	// Extract reads cross-cutting concerns from the carrier into a Context.	// Extract 从载体中读取横切关注点到一个 Context 中。
	Extract(ctx context.Context, carrier TextMapCarrier) context.Context
	
	// Fields returns the keys whose values are set with Inject.				// Fields 返回其值使用 Inject 设置的键。
	Fields() []string
}
1.func NewCompositeTextMapPropagator(p ...TextMapPropagator) TextMapPropagator	// NewCompositeTextMapPropagator 从传入的 TextMapPropagator 组中返回一个统一的 TextMapPropagator。这允许以统一的方式传播不同的横切关注点。
																				// 返回的 TextMapPropagator 将按照提供 TextMapPropagator 的顺序注入和提取横切关注点。此外，Fields 方法将返回使用 Inject 方法设置的键的去重切片。




type TraceContext struct{}				// TraceContext 是一个支持 W3C Trace Context 格式的传播器 (https://www.w3.org/TR/trace-context/)
										// 此传播器将传播 traceparent 和 tracestate 标头以保证跟踪不被破坏。该传播器的用户可以选择是否要通过修改 traceparent 标头和包含其专有信息的 tracestate 标头的相关部分来参与跟踪。
1.func (tc TraceContext) Extract(ctx context.Context, carrier TextMapCarrier) context.Context		// Extract 从载体中读取 tracecontext 到返回的 Context 中。
																									// 返回的 Context 将是 ctx 的副本，并包含提取的 tracecontext 作为远程 SpanContext。如果提取的 tracecontext 无效，则直接返回传递的 ctx。
2.func (tc TraceContext) Fields() []string															// Fields 返回使用 Inject 设置值的键。
3.func (tc TraceContext) Inject(ctx context.Context, carrier TextMapCarrier)						// 将 Context 中的 set tracecontext 注入到载体中。

Example ¶
package main

import (
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
)

func main() {
	tc := propagation.TraceContext{}
	// Register the TraceContext propagator globally.		// 全局注册 TraceContext 传播器。
	otel.SetTextMapPropagator(tc)
}
