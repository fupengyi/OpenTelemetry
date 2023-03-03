package OpenTelemetry_Go

import (
	"context"
	"time"
)

exptrace "go.opentelemetry.io/otel/sdk/export/trace"

Types ¶
type SpanExporter interface {														// SpanExporter 处理 SpanSnapshot 结构到外部接收器的传递。这是跟踪导出管道中的最后一个组件。
	// ExportSpans exports a batch of SpanSnapshots.								// ExportSpans 导出一批 SpanSnapshots。
	//
	// This function is called synchronously, so there is no concurrency			// 这个函数是同步调用的，所以没有并发
	// safety requirement. However, due to the synchronous calling pattern,			// 安全要求。但是，由于同步调用模式，
	// it is critical that all timeouts and cancellations contained in the			// 至关重要的是，所有超时和取消都包含在
	// passed context must be honored.												// 必须遵守传递的上下文。
	//
	// Any retry logic must be contained in this function. The SDK that				// 任何重试逻辑都必须包含在此函数中。开发工具包
	// calls this function will not implement any retry logic. All errors			// 调用此函数不会实现任何重试逻辑。所有错误
	// returned by this function are considered unrecoverable and will be			// 此函数返回的结果被认为是不可恢复的，将被
	// reported to a configured error Handler.										// 报告给配置的错误处理程序。
	ExportSpans(ctx context.Context, ss []*SpanSnapshot) error
	// Shutdown notifies the exporter of a pending halt to operations. The			// Shutdown 通知出口商操作暂停。这
	// exporter is expected to preform any cleanup or synchronization it			// 出口商应该执行任何清理或同步它
	// requires while honoring all timeouts and cancellations contained in			// 要求同时遵守包含在中的所有超时和取消
	// the passed context.															// 传递的上下文。
	Shutdown(ctx context.Context) error
}



// SpanSnapshot 是一个 span 的快照，它包含了 span 收集的所有信息。它的主要目的是导出完成的跨度。尽管可以访问和修改 SpanSnapshot 字段，
// 但 SpanSnapshot 应该被视为不可变的。对创建 SpanSnapshot 的跨度的更改不会反映在 SpanSnapshot 中。
type SpanSnapshot struct {
	SpanContext  trace.SpanContext
	ParentSpanID trace.SpanID
	SpanKind     trace.SpanKind
	Name         string
	StartTime    time.Time
	// The wall clock time of EndTime will be adjusted to always be offset			// EndTime 的挂钟时间将调整为始终偏移
	// from StartTime by the duration of the span.									// 从 StartTime 到跨度的持续时间。
	EndTime         time.Time
	Attributes      []attribute.KeyValue
	MessageEvents   []trace.Event
	Links           []trace.Link
	StatusCode      codes.Code
	StatusMessage   string
	HasRemoteParent bool
	
	// DroppedAttributeCount contains dropped attributes for the span itself, events and links.	// DroppedAttributeCount 包含跨度本身、事件和链接的已删除属性。
	DroppedAttributeCount    int
	DroppedMessageEventCount int
	DroppedLinkCount         int
	
	// ChildSpanCount holds the number of child span created for this span.			// ChildSpanCount 保存为此跨度创建的子跨度数。
	ChildSpanCount int
	
	// Resource contains attributes representing an entity that produced this span.	// 资源包含表示生成此跨度的实体的属性。
	Resource *resource.Resource
	
	// InstrumentationLibrary defines the instrumentation library used to			// InstrumentationLibrary 定义用于
	// provide instrumentation.														// 提供仪器。
	InstrumentationLibrary instrumentation.Library
}
