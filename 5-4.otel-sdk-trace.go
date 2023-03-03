package OpenTelemetry_Go

import (
	"context"
	"fmt"
	"time"
)

"go.opentelemetry.io/otel/sdk/trace"

Overview
包 trace 包含对 OpenTelemetry 分布式跟踪的支持。
以下假设基本熟悉 OpenTelemetry 概念。请参阅 https://opentelemetry.io。


Constants ¶
const (														// BatchSpanProcessorOptions 的默认值。
	DefaultMaxQueueSize       = 2048						// 默认最大队列大小
	DefaultScheduleDelay      = 5000						// 默认计划延迟
	DefaultExportTimeout      = 30000						// 默认导出超时
	DefaultMaxExportBatchSize = 512							// 默认最大导出批量大小
)

const (
	// DefaultAttributeValueLengthLimit is the default maximum allowed attribute value length, unlimited.
	// DefaultAttributeValueLengthLimit 是默认的最大允许属性值长度，无限制。
	DefaultAttributeValueLengthLimit = -1
	
	// DefaultAttributeCountLimit is the default maximum number of attributes a span can have.
	// DefaultAttributeCountLimit 是一个 span 可以拥有的默认最大属性数。
	DefaultAttributeCountLimit = 128
	
	// DefaultEventCountLimit is the default maximum number of events a span can have.
	// DefaultEventCountLimit 是一个 span 可以拥有的默认最大事件数。
	DefaultEventCountLimit = 128
	
	// DefaultLinkCountLimit is the default maximum number of links a span can have.
	// DefaultLinkCountLimit 是一个 span 可以拥有的默认最大链接数。
	DefaultLinkCountLimit = 128
	
	// DefaultAttributePerEventCountLimit is the default maximum number of attributes a span event can have.
	// DefaultAttributePerEventCountLimit 是 span 事件可以拥有的默认最大属性数。
	DefaultAttributePerEventCountLimit = 128
	
	// DefaultAttributePerLinkCountLimit is the default maximum number of attributes a span link can have.
	// DefaultAttributePerLinkCountLimit 是跨链接可以拥有的默认最大属性数。
	DefaultAttributePerLinkCountLimit = 128
)



Types ¶
type BatchSpanProcessorOption func(o *BatchSpanProcessorOptions)					// BatchSpanProcessorOption 配置 BatchSpanProcessor。
1.func WithBatchTimeout(delay time.Duration) BatchSpanProcessorOption				// WithBatchTimeout 返回一个 BatchSpanProcessorOption，它配置 BatchSpanProcessor 在导出任何保留的跨度（无论队列是否已满）之前允许的最大延迟。
2.func WithBlocking() BatchSpanProcessorOption										// WithBlocking 返回一个 BatchSpanProcessorOption，它将 BatchSpanProcessor 配置为等待入队操作成功，而不是在队列已满时丢弃数据。
3.func WithExportTimeout(timeout time.Duration) BatchSpanProcessorOption			// WithExportTimeout 返回一个 BatchSpanProcessorOption，它配置 BatchSpanProcessor 在放弃导出之前等待导出器导出的时间量。
4.func WithMaxExportBatchSize(size int) BatchSpanProcessorOption					// WithMaxExportBatchSize 返回 BatchSpanProcessorOption，它配置 BatchSpanProcessor 允许的最大导出批量大小。
5.func WithMaxQueueSize(size int) BatchSpanProcessorOption							// WithMaxQueueSize 返回一个 BatchSpanProcessorOption，它配置 BatchSpanProcessor 允许的最大队列大小。



type BatchSpanProcessorOptions struct {															// BatchSpanProcessorOptions 是 BatchSpanProcessor 的配置设置。
	// MaxQueueSize is the maximum queue size to buffer spans for delayed processing. If the	// MaxQueueSize 是延迟处理缓冲跨度的最大队列大小。如果
	// queue gets full it drops the spans. Use BlockOnQueueFull to change this behavior.		// 队列满了它会丢弃跨度。使用 BlockOnQueueFull 更改此行为。
	// The default value of MaxQueueSize is 2048.												// MaxQueueSize 的默认值为 2048。
	MaxQueueSize int
	
	// BatchTimeout is the maximum duration for constructing a batch. Processor					// BatchTimeout 是构建批处理的最大持续时间。处理器
	// forcefully sends available spans when timeout is reached.								// 达到超时时强制发送可用跨度。
	// The default value of BatchTimeout is 5000 msec.											// BatchTimeout 的默认值为 5000 毫秒。
	BatchTimeout time.Duration
	
	// ExportTimeout specifies the maximum duration for exporting spans. If the timeout			// ExportTimeout 指定导出跨度的最大持续时间。如果超时
	// is reached, the export will be cancelled.												// 到达，导出将被取消。
	// The default value of ExportTimeout is 30000 msec.										// ExportTimeout 的默认值为 30000 毫秒。
	ExportTimeout time.Duration
	
	// MaxExportBatchSize is the maximum number of spans to process in a single batch.			// MaxExportBatchSize 是单个批处理中要处理的最大跨度数。
	// If there are more than one batch worth of spans then it processes multiple batches		// 如果有超过一批的跨度，那么它会处理多个批次
	// of spans one batch after the other without any delay.									// of spans 一批又一批，没有任何延迟。
	// The default value of MaxExportBatchSize is 512.											// MaxExportBatchSize 的默认值为 512。
	MaxExportBatchSize int
	
	// BlockOnQueueFull blocks onEnd() and onStart() method if the queue is full				// 如果队列已满，BlockOnQueueFull 会阻塞 onEnd() 和 onStart() 方法
	// AND if BlockOnQueueFull is set to true.													// 并且如果 BlockOnQueueFull 设置为 true。
	// Blocking option should be used carefully as it can severely affect the performance of an	// 应谨慎使用阻塞选项，因为它会严重影响程序的性能
	// application.																				// 应用。
	BlockOnQueueFull bool
}



type Event struct {															// Event 是在 Span 的生命周期中发生的事情。
	// Name is the name of this event										// name是这个事件的名字
	Name string
	
	// Attributes describe the aspects of the event.						// 属性描述了事件的各个方面。
	Attributes []attribute.KeyValue
	
	// DroppedAttributeCount is the number of attributes that were not		// DroppedAttributeCount 是由于达到配置限制而未记录的属性数。
	// recorded due to configured limits being reached.
	DroppedAttributeCount int
	
	// Time at which this event was recorded.								// 记录此事件的时间。
	Time time.Time
}



type IDGenerator interface {												// IDGenerator 允许为 TraceID 和 SpanID 自定义生成器。
	
	// NewIDs returns a new trace and span ID.								// NewIDs 返回一个新的跟踪和跨度 ID。
	NewIDs(ctx context.Context) (trace.TraceID, trace.SpanID)
	
	// NewSpanID returns a ID for a new span in the trace with traceID.		// NewSpanID 使用 traceID 返回跟踪中新跨度的 ID。
	NewSpanID(ctx context.Context, traceID trace.TraceID) trace.SpanID
}



type Link struct {															// Link是两个Span之间的关系。关系可以在同一个 Trace 内，也可以跨越不同的 Trace。
	// SpanContext of the linked Span.										// 链接 Span 的 SpanContext。
	SpanContext trace.SpanContext
	
	// Attributes describe the aspects of the link.							// 属性描述链接的各个方面。
	Attributes []attribute.KeyValue
	
	// DroppedAttributeCount is the number of attributes that were not		// DroppedAttributeCount 是由于达到配置限制而未记录的属性数。
	// recorded due to configured limits being reached.
	DroppedAttributeCount int
}



type ParentBasedSamplerOption interface {									// ParentBasedSamplerOption 为特定的采样案例配置采样器。
	// contains filtered or unexported methods
}
1.func WithLocalParentNotSampled(s Sampler) ParentBasedSamplerOption		// WithLocalParentNotSampled 为未采样的本地父级设置采样器。
2.func WithLocalParentSampled(s Sampler) ParentBasedSamplerOption			// WithLocalParentSampled 为采样的本地父级设置采样器。
3.func WithRemoteParentNotSampled(s Sampler) ParentBasedSamplerOption		// WithRemoteParentNotSampled 为未采样的远程父级设置采样器。
4.func WithRemoteParentSampled(s Sampler) ParentBasedSamplerOption			// WithRemoteParentSampled 为采样的远程父级设置采样器。



// ReadOnlySpan 允许从 trace.Span 底层的数据结构中读取信息。它用于需要从跨度中读取信息但不需要或不允许更改跨度的地方。
// 警告：方法可能会在次要版本中添加到此接口。
type ReadOnlySpan interface {
	// Name returns the name of the span.											// Name 返回跨度的名称。
	Name() string
	
	// SpanContext returns the unique SpanContext that identifies the span.			// SpanContext 返回标识跨度的唯一 SpanContext。
	SpanContext() trace.SpanContext
	
	// Parent returns the unique SpanContext that identifies the parent of the		// Parent 返回标识父级的唯一 SpanContext
	// span if one exists. If the span has no parent the returned SpanContext		// span 如果存在的话。如果跨度没有父级，则返回 SpanContext
	// will be invalid.																// 将无效。
	Parent() trace.SpanContext
	
	// SpanKind returns the role the span plays in a Trace.							// SpanKind 返回 span 在 Trace 中扮演的角色。
	SpanKind() trace.SpanKind
	
	// StartTime returns the time the span started recording.						// StartTime 返回 span 开始记录的时间。
	StartTime() time.Time

	// EndTime returns the time the span stopped recording. It will be zero if		// EndTime 返回 span 停止记录的时间。它将为零，如果
	// the span has not ended.														// 跨度还没有结束。
	EndTime() time.Time
	
	// Attributes returns the defining attributes of the span.										// Attributes 返回 span 的定义属性。
	// The order of the returned attributes is not guaranteed to be stable across invocations.		// 不保证返回属性的顺序在调用中是稳定的。
	Attributes() []attribute.KeyValue
	
	// Links returns all the links the span has to other spans.						// Links 返回 span 到其他 span 的所有链接。
	Links() []Link
	
	// Events returns all the events that occurred within in the spans				// Events 返回 span 中发生的所有事件
	// lifetime.																	// 寿命。
	Events() []Event
	
	// Status returns the spans status.												// Status 返回跨度状态。
	Status() Status
	
	// InstrumentationScope returns information about the instrumentation			// InstrumentationScope 返回有关创建跨度的检测范围的信息。
	// scope that created the span.
	InstrumentationScope() instrumentation.Scope
	
	// InstrumentationLibrary returns information about the instrumentation			// InstrumentationLibrary 返回有关创建跨度的检测库的信息。
	// library that created the span.
	// Deprecated: please use InstrumentationScope instead.							// 已弃用：请改用 InstrumentationScope。
	InstrumentationLibrary() instrumentation.Library
	
	// Resource returns information about the entity that produced the span.		// Resource 返回有关生成跨度的实体的信息。
	Resource() *resource.Resource
	
	// DroppedAttributes returns the number of attributes dropped by the span		// DroppedAttributes 返回 span 丢弃的属性数
	// due to limits being reached.													// 由于达到限制。
	DroppedAttributes() int
	
	// DroppedLinks returns the number of links dropped by the span due to			// DroppedLinks 返回由于以下原因而被跨度丢弃的链接数
	// limits being reached.														// 达到限制。
	DroppedLinks() int
	
	// DroppedEvents returns the number of events dropped by the span due to		// DroppedEvents 返回由于以下原因而被跨度丢弃的事件数
	// limits being reached.														// 达到限制。
	DroppedEvents() int
	
	// ChildSpanCount returns the count of spans that consider the span a			// ChildSpanCount 返回考虑跨度 a 的跨度数
	// direct parent.																// 直接父级。
	ChildSpanCount() int
	
	// contains filtered or unexported methods
}



// ReadWriteSpan 公开了与 trace.Span 相同的方法，此外还允许从底层数据结构读取信息。此接口公开了 trace.Span（这是一个“只写”跨度）和 ReadOnlySpan
// 方法的联合。应分别在 trace.Span 或 ReadOnlySpan 下添加写入或读取跨度信息的新方法。
// 警告：方法可能会在次要版本中添加到此接口。
type ReadWriteSpan interface {
	trace.Span
	ReadOnlySpan
}



type Sampler interface {															// Sampler 决定是否应该对轨迹进行采样和导出。
	// ShouldSample returns a SamplingResult based on a decision made from the
	// passed parameters.
	ShouldSample(parameters SamplingParameters) SamplingResult						// ShouldSample 根据传递的参数做出的决定返回一个 SamplingResult。
	
	// Description returns information describing the Sampler.						// Description 返回描述采样器的信息。
	Description() string
}
1.func AlwaysSample() Sampler						// AlwaysSample 返回一个对每条轨迹进行采样的采样器。在具有大量流量的生产应用程序中使用此采样器时要小心：将为每个请求启动并导出新的跟踪。
2.func NeverSample() Sampler						// NeverSample 返回一个不对任何痕迹进行采样的采样器。
3.func ParentBased(root Sampler, samplers ...ParentBasedSamplerOption) Sampler
// ParentBased 返回一个复合采样器，它根据跨度的父级表现不同。如果跨度没有父级，则使用根（采样器）进行采样决策。如果 span 有父级，则根据父级
// 是否远程以及是否对其进行采样，将应用以下采样器之一：
//	remoteParentSampled(Sampler) (default: AlwaysOn)
//	remoteParentNotSampled(Sampler) (default: AlwaysOff)
//	localParentSampled(Sampler) (default: AlwaysOn)
//	localParentNotSampled(Sampler) (default: AlwaysOff)
4.func TraceIDRatioBased(fraction float64) Sampler	// TraceIDRatioBased 对给定部分的轨迹进行采样。 >= 1 的分数将始终采样。 < 0 的分数被视为零。为了尊重父跟踪的“SampledFlag”，“TraceIDRatioBased”采样器应该用作“Parent”采样器的委托。



type SamplingDecision uint8															// SamplingDecision 指示跨度是否被删除、记录和/或采样。
const (																				// 有效的抽样决定。

	// Drop will not record the span and all attributes/events will be dropped.		// Drop 不会记录跨度，所有属性/事件都会被丢弃。
	Drop SamplingDecision = iota
	
	// Record indicates the span's `IsRecording() == true`, but `Sampled` flag		// Record 表示 span 的 `IsRecording() == true`，但是 `Sampled` 标志
	// *must not* be set.															// *不得*设置。
	RecordOnly
	
	// RecordAndSample has span's `IsRecording() == true` and `Sampled` flag		// RecordAndSample 有 span 的 `IsRecording() == true` 和 `Sampled` 标志
	// *must* be set.																// *必须*设置。
	RecordAndSample
)


type SamplingParameters struct {				// SamplingParameters 包含传递给采样器的值。
	ParentContext context.Context
	TraceID       trace.TraceID
	Name          string
	Kind          trace.SpanKind
	Attributes    []attribute.KeyValue
	Links         []trace.Link
}


type SamplingResult struct {				// SamplingResult 传送一个 SamplingDecision、一组属性和一个 Tracestate。
	Decision   SamplingDecision
	Attributes []attribute.KeyValue
	Tracestate trace.TraceState
}



type SpanExporter interface {													// SpanExporter 处理跨度到外部接收器的传递。这是跟踪导出管道中的最后一个组件。
	// ExportSpans exports a batch of spans.									// ExportSpans 导出一批 span。
	//
	// This function is called synchronously, so there is no concurrency		// 这个函数是同步调用的，所以没有并发
	// safety requirement. However, due to the synchronous calling pattern,		// 安全要求。但是，由于同步调用模式，
	// it is critical that all timeouts and cancellations contained in the		// 至关重要的是，所有超时和取消都包含在
	// passed context must be honored.											// 必须遵守传递的上下文。
	//
	// Any retry logic must be contained in this function. The SDK that			// 任何重试逻辑都必须包含在此函数中。开发工具包
	// calls this function will not implement any retry logic. All errors		// 调用此函数不会实现任何重试逻辑。所有错误
	// returned by this function are considered unrecoverable and will be		// 此函数返回的结果被认为是不可恢复的，将被
	// reported to a configured error Handler.									// 报告给配置的错误处理程序。
	ExportSpans(ctx context.Context, spans []ReadOnlySpan) error
	
	// Shutdown notifies the exporter of a pending halt to operations. The		// Shutdown 通知出口商操作暂停。这
	// exporter is expected to preform any cleanup or synchronization it		// 出口商应该执行任何清理或同步它
	// requires while honoring all timeouts and cancellations contained in		// 需要同时遵守包含在中的所有超时和取消
	// the passed context.														// 传递的上下文。
	Shutdown(ctx context.Context) error
}



type SpanLimits struct {															// SpanLimits 表示跨度的限制。
	// AttributeValueLengthLimit is the maximum allowed attribute value length.		// AttributeValueLengthLimit 是允许的最大属性值长度。
	//
	// This limit only applies to string and string slice attribute values.			// 此限制仅适用于字符串和字符串切片属性值。
	// Any string longer than this value will be truncated to this length.			// 任何长于该值的字符串都将被截断为该长度。
	//
	// Setting this to a negative value means no limit is applied.					// 将此设置为负值意味着不应用限制。
	AttributeValueLengthLimit int
	
	// AttributeCountLimit is the maximum allowed span attribute count. Any			// AttributeCountLimit 是允许的最大跨度属性计数。任何
	// attribute added to a span once this limit is reached will be dropped.		// 一旦达到此限制，添加到跨度的属性将被删除。
	//
	// Setting this to zero means no attributes will be recorded.					// 将其设置为零意味着不会记录任何属性。
	//
	// Setting this to a negative value means no limit is applied.					// 将此设置为负值意味着不应用限制。
	AttributeCountLimit int
	
	// EventCountLimit is the maximum allowed span event count. Any event			// EventCountLimit 是允许的最大跨度事件计数。任何事件
	// added to a span once this limit is reached means it will be added but		// 一旦达到这个限制就添加到跨度意味着它将被添加但是
	// the oldest event will be dropped.											// 最早的事件将被丢弃。
	//
	// Setting this to zero means no events we be recorded.							// 将其设置为零意味着我们不会记录任何事件。
	//
	// Setting this to a negative value means no limit is applied.					// 将此设置为负值意味着不应用限制。
	EventCountLimit int
	
	// LinkCountLimit is the maximum allowed span link count. Any link added		// LinkCountLimit 是允许的最大跨度链接数。添加的任何链接
	// to a span once this limit is reached means it will be added but the			// 一旦达到这个限制就意味着它会被添加到一个跨度，但是
	// oldest link will be dropped.													// 最旧的链接将被丢弃。
	//
	// Setting this to zero means no links we be recorded.							// 将此设置为零意味着我们不会记录任何链接。
	//
	// Setting this to a negative value means no limit is applied.					// 将此设置为负值意味着不应用限制。
	LinkCountLimit int
	
	// AttributePerEventCountLimit is the maximum number of attributes allowed		// AttributePerEventCountLimit 是允许的最大属性数
	// per span event. Any attribute added after this limit reached will be			// 每个跨度事件。达到此限制后添加的任何属性都将
	// dropped.																		// 掉了。
	//
	// Setting this to zero means no attributes will be recorded for events.		// 将此设置为零意味着不会为事件记录任何属性。
	//
	// Setting this to a negative value means no limit is applied.					// 将此设置为负值意味着不应用限制。
	AttributePerEventCountLimit int
	
	// AttributePerLinkCountLimit is the maximum number of attributes allowed		// AttributePerLinkCountLimit 是允许的最大属性数
	// per span link. Any attribute added after this limit reached will be			// 每个跨度链接。达到此限制后添加的任何属性都将
	// dropped.																		// 掉了。
	//
	// Setting this to zero means no attributes will be recorded for links.			// 将此设置为零意味着不会为链接记录任何属性。
	//
	// Setting this to a negative value means no limit is applied.					// 将此设置为负值意味着不应用限制。
	AttributePerLinkCountLimit int
}
1.func NewSpanLimits() SpanLimits		// NewSpanLimits 返回一个 SpanLimits，所有限制都设置为其相应环境变量所持有的值，如果未设置则为默认值。
	• AttributeValueLengthLimit: OTEL_SPAN_ATTRIBUTE_VALUE_LENGTH_LIMIT (default: unlimited)
	• AttributeCountLimit: OTEL_SPAN_ATTRIBUTE_COUNT_LIMIT (default: 128)
	• EventCountLimit: OTEL_SPAN_EVENT_COUNT_LIMIT (default: 128)
	• AttributePerEventCountLimit: OTEL_EVENT_ATTRIBUTE_COUNT_LIMIT (default: 128)
	• LinkCountLimit: OTEL_SPAN_LINK_COUNT_LIMIT (default: 128)
	• AttributePerLinkCountLimit: OTEL_LINK_ATTRIBUTE_COUNT_LIMIT (default: 128)



type SpanProcessor interface {													// SpanProcessor 是跟踪信号中跨度的处理管道。 SpanProcessors 注册到 TracerProvider 并在 Span 生命周期的开始和结束时被调用，并按照它们注册的顺序被调用。
	// OnStart is called when a span is started. It is called synchronously		// 跨度启动时调用 OnStart。它被同步调用
	// and should not block.													// 并且不应该阻塞。
	OnStart(parent context.Context, s ReadWriteSpan)
	
	// OnEnd is called when span is finished. It is called synchronously and	// 跨度完成时调用 OnEnd。它被同步调用并且
	// hence not block.															// 因此不会阻塞。
	OnEnd(s ReadOnlySpan)
	
	// Shutdown is called when the SDK shuts down. Any cleanup or release of	// SDK 关闭时调用 Shutdown。任何清理或释放
	// resources held by the processor should be done in this call.				// 处理器持有的资源应该在这个调用中完成。
	//
	// Calls to OnStart, OnEnd, or ForceFlush after this has been called		// 调用 OnStart、OnEnd 或 ForceFlush 后调用
	// should be ignored.														// 应该被忽略。
	//
	// All timeouts and cancellations contained in ctx must be honored, this	// 必须遵守 ctx 中包含的所有超时和取消，这
	// should not block indefinitely.											// 不应该无限期地阻塞。
	Shutdown(ctx context.Context) error
	
	// ForceFlush exports all ended spans to the configured Exporter that have not yet		// ForceFlush 将所有结束的 span 导出到尚未配置的 Exporter
	// been exported.  It should only be called when absolutely necessary, such as when		// 已导出。它只应在绝对必要时调用，例如
	// using a FaaS provider that may suspend the process after an invocation, but before	// 使用 FaaS 提供者可能会在调用之后暂停进程，但之前
	// the Processor can export the completed spans.										// 处理器可以导出完成的跨度。
	ForceFlush(ctx context.Context) error
}
1.func NewBatchSpanProcessor(exporter SpanExporter, options ...BatchSpanProcessorOption) SpanProcessor
// NewBatchSpanProcessor 创建一个新的 SpanProcessor，它将使用提供的选项将完成的跨度批次发送到导出器。
// 如果导出器为零，跨度处理器将不执行任何操作。

2.func NewSimpleSpanProcessor(exporter SpanExporter) SpanProcessor
// NewSimpleSpanProcessor 返回一个新的 SpanProcessor，它将立即将完成的 span 同步发送到导出器。
// 不建议将此 SpanProcessor 用于生产用途。此 SpanProcessor 的同步特性使其非常适合测试、调试或显示其他功能的示例，但它会很慢并且具有很高的
// 计算资源使用开销。建议将 BatchSpanProcessor 用于生产用途。

Example (Annotated) ¶
package main

import (
	"context"
	"fmt"
	
	"go.opentelemetry.io/otel/attribute"
)

/*
Sometimes information about a runtime environment can change dynamically or be
delayed from startup. Instead of continuously recreating and distributing a
TracerProvider with an immutable Resource or delaying the startup of your
application on a slow-loading piece of information, annotate the created spans
dynamically using a SpanProcessor.
有时有关运行时环境的信息可以动态变化或被 从启动延迟。而不是不断地重新创建和分发一个 具有不可变资源或延迟您的启动的 TracerProvider 在缓慢加载
的信息上应用，注释创建的跨度 动态使用 SpanProcessor。
*/

var (
	// owner represents the owner of the application. In this example it is		// owner 表示应用程序的所有者。在这个例子中是
	// stored as a simple string, but in real-world use this may be the			// 存储为一个简单的字符串，但在实际使用中这可能是
	// response to an asynchronous request.										// 响应一个异步请求。
	owner    = "unknown"
	ownerKey = attribute.Key("owner")
)

// Annotator is a SpanProcessor that adds attributes to all started spans.		// Annotator 是一个 SpanProcessor，它向所有启动的 span 添加属性。
type Annotator struct {
	// AttrsFunc is called when a span is started. The attributes it returns	// 跨度启动时调用 AttrsFunc。它返回的属性
	// are set on the Span being started.										// 设置在正在启动的 Span 上。
	AttrsFunc func() []attribute.KeyValue
}

func (a Annotator) OnStart(_ context.Context, s ReadWriteSpan) { s.SetAttributes(a.AttrsFunc()...) }
func (a Annotator) Shutdown(context.Context) error             { return nil }
func (a Annotator) ForceFlush(context.Context) error           { return nil }
func (a Annotator) OnEnd(s ReadOnlySpan) {
	attr := s.Attributes()[0]
	fmt.Printf("%s: %s\n", attr.Key, attr.Value.AsString())
}

func main() {
	a := Annotator{
		AttrsFunc: func() []attribute.KeyValue {
			return []attribute.KeyValue{ownerKey.String(owner)}
		},
	}
	tracer := NewTracerProvider(WithSpanProcessor(a)).Tracer("annotated")
	
	// Simulate the situation where we want to annotate spans with an owner,		// 模拟我们想要用所有者注释 span 的情况，
	// but at startup we do not now this information. Instead of waiting for		// 但是在启动时我们现在没有这个信息。而不是等待
	// the owner to be known before starting and blocking here, start doing			// 在开始和阻塞之前要知道的所有者，开始做
	// work and update when the information becomes available.						// 在信息可用时工作并更新。
	ctx := context.Background()
	_, s0 := tracer.Start(ctx, "span0")
	
	// Simulate an asynchronous call to determine the owner succeeding. We now		// 模拟异步调用以确定所有者成功。我们现在
	// know that the owner of this application has been determined to be			// 知道此应用程序的所有者已确定为
	// Alice. Make sure all subsequent spans are annotated appropriately.			// 爱丽丝。确保对所有后续跨度进行适当注释。
	owner = "alice"
	
	_, s1 := tracer.Start(ctx, "span1")
	s0.End()
	s1.End()
}
Output:

owner: unknown
owner: alice




Example (Filtered) ¶
package main

import (
	"context"
	"time"
)

// DurationFilter is a SpanProcessor that filters spans that have lifetimes		// DurationFilter 是一个 SpanProcessor，用于过滤具有生命周期的跨度
// outside of a defined range.													// 在定义的范围之外。
type DurationFilter struct {
	// Next is the next SpanProcessor in the chain.								// 接下来是链中的下一个 SpanProcessor。
	Next SpanProcessor
	
	// Min is the duration under which spans are dropped.						// Min 是删除跨度的持续时间。
	Min time.Duration
	// Max is the duration over which spans are dropped.						// Max 是删除 span 的持续时间。
	Max time.Duration
}

func (f DurationFilter) OnStart(parent context.Context, s ReadWriteSpan) {
	f.Next.OnStart(parent, s)
}
func (f DurationFilter) Shutdown(ctx context.Context) error   { return f.Next.Shutdown(ctx) }
func (f DurationFilter) ForceFlush(ctx context.Context) error { return f.Next.ForceFlush(ctx) }
func (f DurationFilter) OnEnd(s ReadOnlySpan) {
	if f.Min > 0 && s.EndTime().Sub(s.StartTime()) < f.Min {
		// Drop short lived spans.												// 删除短暂的跨度。
		return
	}
	if f.Max > 0 && s.EndTime().Sub(s.StartTime()) > f.Max {
		// Drop long lived spans.												// 删除长寿命跨度。
		return
	}
	f.Next.OnEnd(s)
}

// InstrumentationBlacklist is a SpanProcessor that drops all spans from		// InstrumentationBlacklist 是一个 SpanProcessor，它从
// certain instrumentation.														// 某些仪器。
type InstrumentationBlacklist struct {
	// Next is the next SpanProcessor in the chain.								// 接下来是链中的下一个 SpanProcessor。
	Next SpanProcessor
	
	// Blacklist is the set of instrumentation names for which spans will be	// 黑名单是 span 将用于的仪器名称集
	// dropped.																	// 掉了。
	Blacklist map[string]bool
}

func (f InstrumentationBlacklist) OnStart(parent context.Context, s ReadWriteSpan) {
	f.Next.OnStart(parent, s)
}
func (f InstrumentationBlacklist) Shutdown(ctx context.Context) error { return f.Next.Shutdown(ctx) }
func (f InstrumentationBlacklist) ForceFlush(ctx context.Context) error {
	return f.Next.ForceFlush(ctx)
}
func (f InstrumentationBlacklist) OnEnd(s ReadOnlySpan) {
	if f.Blacklist != nil && f.Blacklist[s.InstrumentationScope().Name] {
		// Drop spans from this instrumentation									// 从此仪器中删除跨度
		return
	}
	f.Next.OnEnd(s)
}

type noopExporter struct{}

func (noopExporter) ExportSpans(context.Context, []ReadOnlySpan) error { return nil }
func (noopExporter) Shutdown(context.Context) error                    { return nil }

func main() {
	exportSP := NewSimpleSpanProcessor(noopExporter{})
	
	// Build a SpanProcessor chain to filter out all spans from the pernicious	// 构建一个 SpanProcessor 链来过滤掉所有有害的 span
	// "naughty-instrumentation" dependency and only allow spans shorter than	// "naughty-instrumentation" 依赖，只允许短于
	// an minute and longer than a second to be exported with the exportSP.		// 用 exportSP 导出一分多于一秒。
	filter := DurationFilter{
		Next: InstrumentationBlacklist{
			Next: exportSP,
			Blacklist: map[string]bool{
				"naughty-instrumentation": true,
			},
		},
		Min: time.Second,
		Max: time.Minute,
	}
	
	_ = NewTracerProvider(WithSpanProcessor(filter))
	// ...
}



type Status struct {															// Status 是 Span 的分类状态。
	// Code is an identifier of a Spans state classification.					// Code 是一个 Spans 状态分类的标识符。
	Code codes.Code
	// Description is a user hint about why that status was set. It is only		// Description 是关于为什么设置该状态的用户提示。它只是
	// applicable when Code is Error.											// 当 Code 为 Error 时适用。
	Description string
}


type TracerProvider struct {						// TracerProvider 是一个 OpenTelemetry TracerProvider。它为仪器提供跟踪器，因此它可以跟踪通过系统的操作流程。
	// contains filtered or unexported fields
}
1.func NewTracerProvider(opts ...TracerProviderOption) *TracerProvider		// NewTracerProvider 返回一个新的和配置的 TracerProvider。
																			// 默认情况下，返回的 TracerProvider 配置为：
																			// a ParentBased(AlwaysSample) Sampler
																			// a random number IDGenerator
																			// the resource.Default() Resource
																			// the default SpanLimits.
																			// 传递的选项用于覆盖这些默认值并适当配置返回的 TracerProvider。
2.func (p *TracerProvider) ForceFlush(ctx context.Context) error			// ForceFlush 立即为所有已注册的跨度处理器导出所有尚未导出的跨度。
3.func (p *TracerProvider) RegisterSpanProcessor(sp SpanProcessor)			// RegisterSpanProcessor 将给定的 SpanProcessor 添加到 SpanProcessor 列表中。
4.func (p *TracerProvider) Shutdown(ctx context.Context) error				// Shutdown 关闭 TracerProvider。所有注册的 span 处理器都按照注册的顺序关闭，并释放所有持有的计算资源。
5.func (p *TracerProvider) Tracer(name string, opts ...trace.TracerOption) trace.Tracer		// Tracer 返回具有给定名称和选项的 Tracer。如果给定名称和选项的 Tracer 不存在，则创建它，否则返回现有的 Tracer。
																							// 如果名称为空，则使用 DefaultTracerName。
																							// 此方法可以安全地并发调用。
6.func (p *TracerProvider) UnregisterSpanProcessor(sp SpanProcessor)		// UnregisterSpanProcessor 从 SpanProcessor 列表中删除给定的 SpanProcessor。



type TracerProviderOption interface {					// TracerProviderOption 配置一个 TracerProvider。
	// contains filtered or unexported methods
}
1.func WithBatcher(e SpanExporter, opts ...BatchSpanProcessorOption) TracerProviderOption	// WithBatcher 使用配置有传递的选项的 BatchSpanProcessor 向 TracerProvider 注册导出器。
2.func WithIDGenerator(g IDGenerator) TracerProviderOption			// WithIDGenerator 返回一个 TracerProviderOption，它将 IDGenerator g 配置为 TracerProvider 的 IDGenerator。 TracerProvider 创建的跟踪器使用配置的 IDGenerator 来生成新的跨度和跟踪 ID。
																	// 如果不使用此选项，TracerProvider 将默认使用随机数 IDGenerator。
3.func WithRawSpanLimits(limits SpanLimits) TracerProviderOption	// WithRawSpanLimits 返回一个 TracerProviderOption，它将 TracerProvider 配置为使用这些限制。这些限制限制了 Tracer 从 TracerProvider 创建的任何 Span。
																	// 这些限制将按原样使用。 零值或负值不会像 WithSpanLimits 那样更改为默认值。 将限制设置为零将有效地禁用它限制的相关资源，设置为负值将意味着资源是无限的。 因此，这意味着零值 SpanLimits 将禁用所有跨度资源。 因此，应使用 NewSpanLimits 构建限制并相应更新。
																	// 如果未提供此或 WithSpanLimits，则 TracerProvider 将使用环境变量定义的限制，或者如果未设置则使用默认值。有关此关系的信息，请参阅 NewSpanLimits 文档。
4.func WithResource(r *resource.Resource) TracerProviderOption		// WithResource 返回一个 TracerProviderOption，它将资源 r 配置为 TracerProvider 的资源。配置的资源由 TracerProvider 创建的所有跟踪器引用。它代表生成遥测的实体。
																	//  如果不使用此选项，TracerProvider 将默认使用 resource.Default() 资源。
5.func WithSampler(s Sampler) TracerProviderOption					// WithSampler 返回一个 TracerProviderOption，它将 Sampler 配置为 TracerProvider 的 Sampler。 TracerProvider 创建的跟踪器使用配置的采样器来为他们创建的跨度做出采样决策。
																	// 此选项会覆盖通过 OTEL_TRACES_SAMPLER 和 OTEL_TRACES_SAMPLER_ARG 环境变量配置的采样器。 如果未使用此选项并且未通过环境变量配置采样器或环境包含无效/不受支持的配置，则 TracerProvider 将默认使用 ParentBased(AlwaysSample) 采样器。
6.func WithSpanProcessor(sp SpanProcessor) TracerProviderOption		// WithSpanProcessor 将 SpanProcessor 注册到 TracerProvider。
7.func WithSyncer(e SpanExporter) TracerProviderOption				// WithSyncer 使用 SimpleSpanProcessor 向 TracerProvider 注册导出器。
																	// 不建议将其用于生产。将包装导出器的 SimpleSpanProcessor 的同步特性使其有利于测试、调试或显示其他功能的示例，但它会很慢并且具有很高的计算资源使用开销。建议将 WithBatcher 选项用于生产用途。
