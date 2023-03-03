package OpenTelemetry_Go
Overview
本文档概述了 OpenTelemetry 项目并定义了重要的基本术语。
可以在词汇表中找到其他术语定义。


OpenTelemetry Client Architecture
在最高架构级别，OpenTelemetry 客户端被组织成信号。每个信号都提供一种特殊形式的可观察性。例如，跟踪、指标和行李是三个独立的信号。信号共享一个
公共子系统——上下文传播——但它们彼此独立运行。

每个信号都为软件提供了一种描述自身的机制。代码库（例如 Web 框架或数据库客户端）依赖于各种信号来描述自身。然后可以将 OpenTelemetry 检测代码
混合到该代码库中的其他代码中。这使得 OpenTelemetry 成为一个横切关注点——一个混合到许多其他软件中以提供价值的软件。就其本质而言，横切关注点违
反了核心设计原则——关注点分离。因此，OpenTelemetry 客户端设计需要格外小心和注意，以避免为依赖于这些横切 API 的代码库带来问题。

OpenTelemetry 客户端旨在将每个信号中必须作为横切关注点导入的部分与可以独立管理的部分分开。 OpenTelemetry 客户端也被设计成一个可扩展的框架
。为了实现这些目标，每个信号都包含四种类型的包：API、SDK、语义约定和 Contrib。


API
API 包由用于检测的横切公共接口组成。导入第三方库和应用程序代码的 OpenTelemetry 客户端的任何部分都被视为 API 的一部分。


SDK
SDK 是 OpenTelemetry 项目提供的 API 的实现。在应用程序中，SDK 由应用程序所有者安装和管理。请注意，SDK 包括额外的公共接口，这些接口不被视
为 API 包的一部分，因为它们不是横切关注点。这些公共接口被定义为构造函数和插件接口。应用程序所有者使用 SDK 构造函数；插件作者使用 SDK 插件接
口。仪器作者不得直接引用任何类型的任何 SDK 包，只能引用 API。


Semantic Conventions
语义约定定义了描述应用程序使用的常见概念、协议和操作的键和值。
	Resource Conventions
	Span Conventions
	Metrics Conventions
collector 和客户端库都应该将语义约定键和枚举值自动生成为常量（或语言惯用的等价物）。在语义约定稳定之前，不应将生成的值分发到稳定的包中。 YAML 文件
必须用作生成的真实来源。每种语言实现都应该为代码生成器提供特定于语言的支持。


Contrib Packages
OpenTelemetry 项目与流行的 OSS 项目保持集成，这些项目已被确定为对观察现代 Web 服务很重要。示例 API 集成包括用于 Web 框架、数据库客户端和
消息队列的检测。示例 SDK 集成包括用于将遥测导出到流行的分析工具和遥测存储系统的插件。

OpenTelemetry 规范需要一些插件，例如 OTLP Exporters 和 TraceContext Propagators。这些必需的插件包含在 SDK 中。
可选且独立于 SDK 的插件和检测包称为 Contrib 包。 API Contrib 是指完全依赖于 API 的包； SDK Contrib 是指也依赖于 SDK 的包。
术语 Contrib 特指由 OpenTelemetry 项目维护的插件和工具的集合；它不涉及在别处托管的第三方插件。


Versioning and Stability
OpenTelemetry 重视稳定性和向后兼容性。有关详细信息，请参阅版本控制和稳定性指南。


Tracing Signal
分布式跟踪是一组事件，由单个逻辑操作触发，跨应用程序的各个组件进行整合。分布式跟踪包含跨进程、网络和安全边界的事件。当有人按下按钮开始网站上的
操作时，可能会启动分布式跟踪 - 在此示例中，跟踪将表示下游服务之间进行的调用，这些服务处理由按下此按钮启动的请求链。


Traces
OpenTelemetry 中的跟踪由其跨度隐式定义。特别地，Trace 可以被认为是 Spans 的有向无环图 (DAG)，其中 Spans 之间的边被定义为父/子关系。
例如，下面是由 6 个 Spans 组成的示例 Trace：
		单个 Trace 中 Span 之间的因果关系
		[Span A]  ←←←(the root span)
           |
	+------+------+
	|             |
[Span B]      [Span C] ←←←(Span C is a `child` of Span A)
	|             |
[Span D]      +---+-------+
	          |           |
			[Span E]    [Span F]

有时使用时间轴更容易可视化跟踪，如下图所示：
单个 Trace 中 Span 之间的时间关系
––|–––––––|–––––––|–––––––|–––––––|–––––––|–––––––|–––––––|–> time

[Span A···················································]
  [Span B··········································]
     [Span D······································]
   [Span C····················································]
        [Span E·······]        [Span F··]


Spans
跨度表示事务中的操作。每个 Span 封装了以下状态：
	操作名称
	开始和结束时间戳
	属性：键值对列表。
	一组零个或多个事件，每个事件本身就是一个元组（时间戳、名称、属性）。名称必须是字符串。
	父母的 Span 标识符。
	链接到零个或多个因果相关的 Spans（通过这些相关 Spans 的 SpanContext）。
	引用 Span 所需的 SpanContext 信息。见下文。
	
	
SpanContext
表示在 Trace 中标识 Span 的所有信息，并且必须传播到子 Span 并跨越进程边界。 SpanContext 包含跟踪标识符和从父 Span 传播到子 Span 的选项。
	TraceId 是跟踪的标识符。它由 16 个随机生成的字节组成，具有几乎足够的概率，是全球唯一的。 TraceId 用于将特定跟踪的所有跨度分组到所有进程中。
	SpanId 是跨度的标识符。通过将其制成 8 个随机生成的字节，它具有几乎足够的概率是全局唯一的。当传递给子跨度时，此标识符将成为子跨度的父跨度 ID。
	TraceFlags 表示跟踪的选项。它表示为 1 个字节（位图）。
		采样位 - 表示轨迹是否被采样的位（掩码 0x1）。
	Tracestate 在键值对列表中携带跟踪系统特定的上下文。 Tracestate 允许不同的供应商传播附加信息并与他们的旧 Id 格式进行互操作。有关详细信息，请参阅此。


Links between spans
一个跨度可以链接到零个或多个因果相关的其他跨度（由 SpanContext 定义）。链接可以指向单个 Trace 内或跨不同 Trace 的 Spans。链接可用于表示批
处理操作，其中一个 Span 由多个启动 Span 启动，每个 Span 表示批处理中正在处理的单个传入项。

使用链接的另一个示例是声明原始跟踪和后续跟踪之间的关系。当 Trace 进入服务的受信任边界并且服务策略需要生成新的 Trace 而不是信任传入的 Trace
上下文时，可以使用此方法。新链接的 Trace 也可能代表一个长时间运行的异步数据处理操作，它是由许多快速传入请求之一发起的。

当使用 scatter/gather（也称为 fork/join）模式时，root 操作会启动多个下游处理操作，并且所有这些操作都聚合回一个 Span 中。最后一个 Span
链接到它聚合的许多操作。它们都是来自同一个 Trace 的 Span。类似于 Span 的 Parent 字段。但是，建议不要在这种情况下设置 Span 的 parent，因
为从语义上讲，parent 字段表示单个父场景，在许多情况下，父 Span 完全包含子 Span。在分散/聚集和批处理场景中不是这种情况。


Metric Signal
OpenTelemetry 允许使用预定义的聚合和一组属性来记录原始测量或指标。

使用 OpenTelemetry API 记录原始测量值允许最终用户决定应将哪种聚合算法应用于该指标以及定义属性（维度）。它将在 gRPC 等客户端库中用于记录原
始测量值“server_latency”或“received_bytes”。因此，最终用户将决定应从这些原始测量中收集哪种类型的聚合值。它可能是简单的平均或复杂的直方图
计算。

使用 OpenTelemetry API 通过预定义聚合记录指标同样重要。它允许收集诸如 cpu 和内存使用率之类的值，或诸如“队列长度”之类的简单指标。


Recording raw measurements
用于记录原始测量的主要类是 Measure 和 Measurement。可以使用 OpenTelemetry API 记录测量列表以及附加上下文。因此，用户可以定义聚合这些测量
并使用一起传递的上下文来定义结果指标的其他属性。


Measure
度量描述了库记录的单个值的类型。它定义了公开度量的库和将这些单独的度量聚合为度量的应用程序之间的契约。度量由名称、描述和值的单位标识。


Measurement
Measurement 描述了要为 Measure 收集的单个值。测量是 API 表面中的一个空接口。该接口在 SDK 中定义。


Recording metrics with predefined aggregation
所有类型的预聚合指标的基类称为指标。它定义基本的指标属性，如名称和属性。从 Metric 继承的类定义了它们的聚合类型以及单个测量或点的结构。 API
定义了以下类型的预聚合指标：
	用于报告瞬时测量的计数器指标。计数器值可以上升或保持不变，但永远不会下降。计数器值不能为负数。有两种类型的计数器度量值 - double 和 long。
	Gauge metric 报告数值的瞬时测量。仪表既可以上升也可以下降。仪表值可以是负数。有两种类型的仪表度量值 - double 和 long。
API 允许构建所选类型的 Metric 。 SDK 定义了查询导出 Metric 当前值的方式。
每种类型的 Metric 都有它的 API 来记录要聚合的值。 API 同时支持 - 设置 Metric 值的推送和拉取模型。


Metrics data model and SDK
Metrics 数据模型在此处指定，基于 metrics.proto。该数据模型定义了三种语义：API 使用的事件模型、SDK 和 OTLP 使用的飞行中数据模型，以及表示
导出器应如何解释飞行中模型的 TimeSeries 模型。

不同的 exporters 具有不同的功能（例如支持哪些数据类型）和不同的约束（例如属性键中允许使用哪些字符）。指标旨在成为可能的超集，而不是随处支持的
最低公分母。所有出口商都通过 OpenTelemetry SDK 中定义的指标生产者接口使用指标数据模型中的数据。

因此，指标对数据施加了最小的限制（例如，键中允许使用哪些字符），处理指标的代码应避免对指标数据进行验证和清理。相反，将数据传递到后端，依靠后端
执行验证，并从后端传回任何错误。

有关详细信息，请参阅指标数据模型规范。


Log Signal
Data model
日志数据模型定义了 OpenTelemetry 如何理解日志和事件。


Baggage Signal
除了跟踪传播之外，OpenTelemetry 还提供了一种用于传播名称/值对的简单机制，称为 Baggage。 Baggage 用于索引一项服务中的可观察性事件，该服务
具有同一事务中先前服务提供的属性。这有助于在这些事件之间建立因果关系。

虽然 Baggage 可用于对其他横切关注点进行原型设计，但该机制主要旨在为 OpenTelemetry 可观察性系统传达价值。
这些值可以从 Baggage 中使用，并用作指标的附加属性，或日志和跟踪的附加上下文。一些例子：
	Web 服务可以受益于包含有关发送请求的服务的上下文
	SaaS 提供商可以包含有关负责该请求的 API 用户或令牌的上下文
	确定特定浏览器版本与图像处理服务中的故障相关联
为了与 OpenTracing 向后兼容，Baggage 在使用 OpenTracing 桥时作为 Baggage 传播。具有不同标准的新关注点应该考虑创建一个新的横切关注点来涵
盖他们的用例；它们可能受益于 W3C 编码格式，但使用新的 HTTP 标头在整个分布式跟踪中传送数据。


Resources
资源捕获有关为其记录遥测的实体的信息。例如，Kubernetes 容器公开的指标可以链接到指定集群、命名空间、pod 和容器名称的资源。
资源可以捕获实体标识的整个层次结构。它可能描述了云中的主机和特定的容器或进程中运行的应用程序。
请注意，某些进程标识信息可以通过 OpenTelemetry SDK 自动与遥测相关联。


Context Propagation
所有 OpenTelemetry 横切关注点（例如跟踪和指标）共享一个底层上下文机制，用于在分布式事务的整个生命周期中存储状态和访问数据。
See the Context


Propagators
OpenTelemetry 使用 Propagators 序列化和反序列化横切关注值，例如 Spans（通常只有 SpanContext 部分）和 Baggage。不同的传播器类型定义了
特定传输施加的限制并绑定到数据类型。
Propagators API 目前定义了一种 Propagator 类型：
	TextMapPropagator 将值作为文本注入载体并从载体中提取值。
	
	
Collector
OpenTelemetry 收集器是一组组件，可以从 OpenTelemetry 或其他监控/跟踪库（Jaeger、Prometheus 等）检测的进程中收集跟踪、指标和最终其他遥
测数据（例如日志），进行聚合和智能采样，以及将跟踪和指标导出到一个或多个监控/跟踪后端。收集器将允许丰富和转换收集的遥测数据（例如添加额外的属
性或擦除个人信息）。

OpenTelemetry 收集器有两种主要的操作模式：代理（与应用程序一起在本地运行的守护进程）和收集器（独立运行的服务）。
在 OpenTelemetry 服务长期愿景中阅读更多内容。


Instrumentation Libraries
See Instrumentation Library
该项目的灵感是让每个库和应用程序直接调用 OpenTelemetry API，从而开箱即用。然而，许多库不会有这样的集成，因此需要一个单独的库来注入这样的调
用，使用包装接口、订阅特定于库的回调或将现有遥测数据转换为 OpenTelemetry 模型等机制。

为另一个库启用 OpenTelemetry 可观察性的库称为 Instrumentation 库。
检测库的命名应遵循检测库的任何命名约定（例如，Web 框架的“中间件”）。
如果没有确定的名称，建议在包前加上“opentelemetry-instrumentation”，然后是检测库名称本身。例子包括：
	opentelemetry-instrumentation-flask (Python)
	@opentelemetry/instrumentation-grpc (Javascript)



Glossary
本文档定义了本规范中使用的一些术语。
概述文档中记录了一些其他基本术语。


User Roles
Application Owner				// 应用程序或服务的维护者，负责配置和管理 OpenTelemetry SDK 的生命周期。
Library Author					// 许多应用程序所依赖的共享库的维护者，并且是 OpenTelemetry 检测的目标。
Instrumentation Author			// 针对 OpenTelemetry API 编写的 OpenTelemetry 工具的维护者。这可能是在应用程序代码、共享库或检测库中编写的检测。
Plugin Author					// OpenTelemetry SDK 插件的维护者，针对 OpenTelemetry SDK 插件接口编写。


Common
Signals
OpenTelemetry 围绕信号或遥测类别构建。指标、日志、跟踪和行李是信号的示例。每个信号代表一组连贯的、独立的功能。每个信号都遵循一个单独的生命周
期，定义其当前的稳定性级别。


Packages
在本规范中，术语包描述了一组表示单个依赖项的代码，这些代码可以独立于其他包导入到程序中。这个概念可能映射到某些语言中的不同术语，例如“模块”。请
注意，在某些语言中，术语“包”指的是不同的概念。


ABI Compatibility
ABI（应用程序二进制接口）是一种接口，它定义了机器代码级别的软件组件之间的交互，例如应用程序可执行文件和共享对象库的已编译二进制文件之间的交互
。 ABI 兼容性意味着库的新编译版本可以正确链接到目标可执行文件，而无需重新编译该可执行文件。
ABI 兼容性对于某些语言很重要，尤其是那些提供机器代码形式的语言。对于其他语言，ABI 兼容性可能不是相关要求。


In-band and Out-of-band Data
在电信中，in-band signaling 是在用于语音或视频等数据的同一频带或信道内发送控制信息。这与通过不同信道甚至通过单独网络（维基百科）发送的
out-of-band signaling 形成对比。
在 OpenTelemetry 中，我们将 in-band data 称为作为业务消息的一部分在分布式系统的组件之间传递的数据，例如，当 trace 或 baggages 以 HTTP
headers 的形式包含在 HTTP 请求中时。这些数据通常不包含遥测数据，但用于关联和连接各个组件产生的遥测数据。遥测本身被称为 out-of-band data：
它通过专用消息从应用程序传输，通常由后台例程异步传输，而不是从业务逻辑的关键路径传输。导出到遥测后端的指标、日志和跟踪是 out-of-band data
的示例。


Manual Instrumentation
针对 OpenTelemetry API（例如 Tracing API、Metrics API 或其他 API）进行编码，以从最终用户代码或共享框架（例如 MongoDB、Redis 等）收集遥测数据。


Automatic Instrumentation
指不需要最终用户修改应用程序源代码的遥测收集方法。方法因编程语言而异，示例包括代码操作（在编译期间或运行时）、猴子修补或运行 eBPF 程序。
同义词：Auto-instrumentation


Telemetry SDK
表示实现 OpenTelemetry API 的库。
See Library Guidelines and Library resource semantic conventions.


Constructors
构造函数是应用程序所有者用来初始化和配置 OpenTelemetry SDK 和贡献包的公共代码。构造函数的示例包括配置对象、环境变量和构建器。


SDK Plugins
插件是扩展 OpenTelemetry SDK 的库。插件接口的示例是 SpanProcessor、Exporter 和 Sampler 接口。


Exporter Library
导出器是实现导出器接口并向消费者发出遥测数据的 SDK 插件。


Instrumented Library
表示为其收集遥测信号（跟踪、指标、日志）的库。
对 OpenTelemetry API 的调用可以由 Instrumented Library 本身完成，也可以由另一个 Instrumentation Library 完成。
示例：org.mongodb.client。


Instrumentation Library
表示为给定的 Instrumented Library 提供检测的库。如果 Instrumented Library 和 Instrumentation Library 具有内置的 OpenTelemetry 工具，则它可能是同一个库。
有关更详细的定义和命名指南，请参阅概述。
示例：io.opentelemetry.contrib.mongodb。
同义词：Instrumenting Library.


Instrumentation Scope
应用程序代码的逻辑单元，可与发出的遥测数据相关联。决定什么表示合理的检测范围通常是开发人员的选择。最常见的方法是使用仪器库作为范围，但是其他范
围也很常见，例如可以选择模块、包或类作为检测范围。

如果代码单元有版本，则检测范围由 (name,version) 对定义，否则省略版本，仅使用名称。名称或（名称、版本）对唯一标识发出遥测的代码的逻辑单元。确
保唯一性的典型方法是使用发出代码的完全限定名称（例如完全限定库名称或完全限定类名称）。

仪器示波器用于获取 Tracer 或 Meter。

检测范围可能具有零个或多个附加属性，这些属性提供有关范围的附加信息。例如，对于指定仪器库的范围，可以记录一个附加属性以表示存储库源代码的存储库
URL 的 URL。由于作用域是一个构建时概念，因此作用域的属性不能在运行时更改。


Tracer Name / Meter Name
这是指在创建新的 Tracer 或 Meter 时指定的名称和（可选）版本参数（请参阅获取跟踪器/获取 Meter）。名称/版本对标识 Instrumentation Scope
，例如 Instrumentation Library 或在其范围内发出遥测的另一个应用程序单元。


Execution Unit
顺序代码执行的最小单元的总称，用于不同的多任务处理概念。例如线程、协程或纤程。


Logs
Log Record
一个事件的记录。通常，记录包括指示事件发生时间的时间戳以及描述发生的事情、发生的地点等的其他数据。
Synonyms: Log Entry.


Log
有时用于指代日志记录的集合。可能有歧义，因为人们有时也使用 Log 来指代单个日志记录，因此应谨慎使用该术语，并且在可能产生歧义的上下文中应使用其
他限定词（例如日志记录）。


Embedded Log
事件列表中嵌入 Span 对象的日志记录。


Standalone Log
未嵌入 Span 内并记录在其他地方的日志记录。


Log Attributes
日志记录中包含的键/值对。


Structured Logs
以具有明确定义的结构的格式记录的日志，允许区分日志记录的不同元素（例如时间戳、属性等）。例如，Syslog 协议 (RFC 5424) 定义了一种结构化数据格式。


Flat File Logs
记录在文本文件中的日志，通常每条日志记录一行（尽管多行记录也是可能的）。以更结构化的格式（例如 JSON 文件）写入文本文件的日志是否被视为平面文
件日志，尚无通用的行业协议。如果这种区别很重要，建议特别指出。
