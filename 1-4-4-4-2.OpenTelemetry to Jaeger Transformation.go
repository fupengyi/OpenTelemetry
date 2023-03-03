package OpenTelemetry_Go

import "encoding/binary"

OpenTelemetry to Jaeger Transformation
Status: Stable
本文档定义了 OpenTelemetry 和 Jaeger Spans 之间的转换。此处指定的通用转换规则也适用。如果特定的通用转换规则与本文档中的规则相矛盾，则必须使用本文档中的规则。
Jaeger 接受以下格式的跨度：
OpenTelemetry Protocol (OTLP)，在开放遥测协议中定义
Thrift Batch，在 jaeger-idl/.../jaeger.thrift 中定义，通过 UDP 或 HTTP 接受
Protobuf Batch，在 jaeger-idl/.../model.proto 中定义，通过 gRPC 接受
See also:
Jaeger API
OpenTelemetry Collector Contrib 存储库中此翻译的参考实现

Summary
下表总结了 OpenTelemetry 和 Jaeger 之间的主要转换。
OpenTelemetry					Jaeger Thrift				Jaeger Proto			Notes
Span.TraceId					Span.traceIdLow/High		Span.trace_id			See IDs
Span.ParentId					Span.parentSpanId			as SpanReference		See Parent ID
Span.SpanId						Span.spanId					Span.span_id
Span.TraceState					TBD							TBD
Span.Name						Span.operationName			Span.operation_name
Span.Kind						Span.tags["span.kind"]		same					See SpanKind for values mapping
Span.StartTime					Span.startTime				Span.start_time			See Unit of time
Span.EndTime					Span.duration				same					Calculated as EndTime - StartTime. See also Unit of time
Span.Attributes					Span.tags					same					See Attributes for data types for the mapping.
Span.DroppedAttributesCount		Add to Span.tags			same					See Dropped Attributes Count for tag name to use.
Span.Events						Span.logs					same					See Events for the mapping format.
Span.DroppedEventsCount			Add to Span.tags			same					See Dropped Events Count for tag name to use.
Span.Links						Span.references				same					See Links
Span.DroppedLinksCount			Add to Span.tags			same					See Dropped Links Count for tag name to use.
Span.Status						Add to Span.tags			same					See Status for tag names to use.

Mappings
本节讨论 OpenTelemetry 和 Jaeger 之间转换的细节。

Resource
OpenTelemetry 资源必须映射到 Jaeger 的 Span.Process 标签。单个进程可以存在多个资源，exporters 需要相应地处理这种情况。
至关重要的是，Jaeger 后端依赖于 Span.Process.ServiceName 来识别产生跨度的服务。该字段必须从 service resource 的 service.name 属性中
填充。如果 Span 的资源中不包含 service.name，则该字段必须从默认 Resource 中填充。

IDs
Jaeger 中的跟踪和跨度 ID 是随机的字节序列。但是，Thrift 模型使用 i64 类型表示 ID，或者在 128 位宽跟踪 ID 的情况下表示为两个 i64 字段 traceIdLow
和 traceIdHigh。必须使用 Big Endian 字节顺序将字节转换为无符号整数或从无符号整数转换，例如[0x10, 0x00, 0x00, 0x00] == 268435456。无
符号整数必须通过将现有的 64 位值重新解释为带符号的 i64 来转换为 i64。例如（在 Go 中）：
var (
	id       []byte = []byte{0xFF, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}
	unsigned uint64 = binary.BigEndian.Uint64(id)
	signed   int64  = int64(unsigned)
)
fmt.Println("unsigned:", unsigned)
fmt.Println("  signed:", signed)
// Output:
// unsigned: 18374686479671623680
//   signed: -72057594037927936

Parent ID
Jaeger Thrift 格式允许在顶级 Span 字段中捕获父 ID。 Jaeger Proto 格式不支持 parent ID 字段；相反，父代必须记录为 CHILD_OF 类型的 SpanReference，例如：
SpanReference(
ref_type=opentracing.CHILD_OF,
trace_id=span.context.trace_id,
span_id=parent_id,
)
这个跨度引用必须是引用列表中的第一个。

SpanKind
OpenTelemetry SpanKind 字段必须在 Jaeger span 中编码为 span.kind 标签，但 SpanKind.INTERNAL 除外，它不应被转换为标签。
OpenTelemetry			Jaeger
SpanKind.CLIENT			"client"
SpanKind.SERVER			"server"
SpanKind.CONSUMER		"consumer"
SpanKind.PRODUCER		"producer"
SpanKind.INTERNAL		do not add span.kind tag

Unit of time
在 Jaeger Thrift 格式中，时间戳和持续时间必须以微秒表示（从时间戳的纪元开始）。如果 OpenTelemetry 中的原始值以纳秒表示，则必须四舍五入或截断为微秒。
在 Jaeger Proto 格式中，时间戳和持续时间必须使用 google.protobuf.Timestamp 和 google.protobuf.Duration 类型以纳秒精度表示。

Status
Status 被记录为 Span 标签。有关要使用的标签名称，请参阅状态。

Error flag
当 Span Status 设置为 ERROR 时，必须添加一个带有布尔值 true 的错误 span 标签。添加的错误标签可以覆盖任何以前的值。

Attributes
OpenTelemetry Span 属性必须作为标签报告给 Jaeger。
原始类型必须由相应类型的 Jaeger 标签表示。
数组值必须像语义约定中提到的 JSON 列表一样序列化为字符串。

Links
必须使用 FOLLOWS_FROM 引用类型在 Jaeger 中将 OpenTelemetry 链接转换为 SpanReference。 Link 的属性不能在 Jaeger 中显式表示。出口商可以另外将链接转换为跨度日志：
使用 Span 开始时间作为 Log 的时间戳
设置日志标签 event=link
从各自的 SpanContext 字段中设置日志标签 trace_id 和 span_id
将 Link 的属性存储为日志标签
从链接生成的跨度引用必须添加到从父 ID 生成的跨度引用之后，如果有的话。

Events
事件必须转换为 Jaeger 日志。 OpenTelemetry Event 的 time_unix_nano 和 attributes 字段直接映射到 Jaeger Log 的 timestamp 和 fields
字段。 Jaeger Log 没有直接等效于 OpenTelemetry Event 的名称字段，但 OpenTracing 语义约定在此处指定了一些特殊的属性名称。 OpenTelemetry Event
的名称字段应添加到 Jaeger Log 的字段映射中，如下所示：

OpenTelemetry Event Field			Jaeger Attribute
name								event
如果 OpenTelemetry Event 包含带有键 event 的属性，它应该优先于 Event 的名称字段。
