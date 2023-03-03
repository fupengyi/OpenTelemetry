package OpenTelemetry_Go
Instrumenting libraries
了解如何将本机检测添加到您的库中。
OpenTelemetry 为许多库提供检测库，这通常是通过库挂钩或猴子修补库代码完成的。
使用 OpenTelemetry 的本机库检测为用户提供了更好的可观察性和开发人员体验，消除了库公开和记录挂钩的需要：
	自定义日志挂钩可以替换为常见且易于使用的 OpenTelemetry API，用户将只与 OpenTelemetry 交互
	来自库和应用程序代码的跟踪、日志、指标是相关且连贯的
	通用约定允许用户在同一技术内以及跨库和语言获得相似且一致的遥测
	遥测信号可以使用各种有据可查的 OpenTelemetry 扩展点针对各种消费场景进行微调（过滤、处理、聚合）。

Semantic Conventions
查看涵盖 Web 框架、RPC 客户端、数据库、消息传递客户端、基础设施等的可用语义约定！

如果你的图书馆是其中之一 - 遵循约定，它们是事实的主要来源，并告诉哪些信息应该包含在 span 中。约定使仪器保持一致：使用遥测的用户不必学习库细节
，可观察性供应商可以为各种技术（例如数据库或消息系统）构建经验。当库遵循约定时，许多场景可以在没有用户输入或配置的情况下开箱即用。

如果您有任何反馈或想添加新约定 - 请前来贡献！ Instrumentation Slack 或 Specification repo 是一个很好的起点！

When not to instrument
一些库是包装网络调用的瘦客户端。 OpenTelemetry 可能有一个用于底层 RPC 客户端的检测库（查看注册表）。在这种情况下，可能不需要检测包装器库。
在以下情况下不要检测：
	你的库是文档化或不言自明的 API 之上的瘦代理
	并且 OpenTelemetry 具有用于底层网络调用的工具
	并且没有您的图书馆应该遵循的约定来丰富遥测
如果您有疑问 - 不要使用仪器 - 您可以在以后看到需要时随时使用。
如果您选择不进行检测，提供一种为内部 RPC 客户端实例配置 OpenTelemetry 处理程序的方法可能仍然有用。它在不支持全自动仪器的语言中是必不可少的，但在其他语言中仍然有用。
如果您决定使用什么以及如何使用，本文档的其余部分将提供指导。

OpenTelemetry API
第一步是依赖 OpenTelemetry API 包。

OpenTelemetry 有两个主要模块——API 和 SDK。 OpenTelemetry API 是一组抽象和非操作性实现。除非您的应用程序导入 OpenTelemetry SDK，否则
您的检测不会执行任何操作，也不会影响应用程序性能。

Libraries should only use the OpenTelemetry API.
您可能会担心添加新的依赖项，这里有一些注意事项可帮助您决定如何最大程度地减少依赖项地狱：
	OpenTelemetry Trace API 在 2021 年初达到稳定，它遵循语义版本控制 2.0，我们非常重视 API 稳定性。
	获取依赖项时，请使用最早稳定的 OpenTelemetry API (1.0.*) 并避免更新它，除非您必须使用新功能。
	当您的仪器稳定时，请考虑将其作为单独的包裹运送，这样就不会给不使用它的用户带来问题。您可以将其保留在您的存储库中，或将其添加到 OpenTelemetry，这样它将与其他检测包一起提供。
	语义约定是稳定的，但会不断发展：虽然这不会导致任何功能问题，但您可能需要每隔一段时间更新一次您的工具。将它放在预览插件或 OpenTelemetry contrib repo 中可能有助于使约定保持最新，而不会破坏用户的更改。

Getting a tracer
通过 Tracer API，所有应用程序配置都对您的库隐藏。默认情况下，库应从全局 TracerProvider 获取跟踪器。
	private static final Tracer tracer = GlobalOpenTelemetry.getTracer("demo-db-client", "0.1.0-beta1");
对于库来说，拥有一个允许应用程序显式传递 TracerProvider 实例的 API 是很有用的，这样可以实现更好的依赖注入并简化测试。
获取跟踪器时，请提供您的库（或跟踪插件）名称和版本 - 它们会显示在遥测数据中并帮助用户处理和过滤遥测数据、了解它的来源以及调试/报告任何检测问题。

What to instrument
Public APIs
公共 API 是跟踪的良好候选者：为公共 API 调用创建的跨度允许用户将遥测映射到应用程序代码，了解库调用的持续时间和结果。要跟踪的调用：
	在内部进行网络调用的公共方法或花费大量时间并可能失败的本地操作（例如 IO）
	处理请求或消息的处理程序

仪器示例：
private static final Tracer tracer = GlobalOpenTelemetry.getTracer("demo-db-client", "0.1.0-beta1");

private Response selectWithTracing(Query query) {
	// check out conventions for guidance on span names and attributes
	Span span = tracer.spanBuilder(String.format("SELECT %s.%s", dbName, collectionName))
		.setSpanKind(SpanKind.CLIENT)
		.setAttribute("db.name", dbName)
		...
		.startSpan();

	// makes span active and allows correlating logs and nest spans
	try (Scope unused = span.makeCurrent()) {
		Response response = query.runWithRetries();
		if (response.isSuccessful()) {
			span.setStatus(StatusCode.OK);
		}

		if (span.isRecording()) {
			// populate response attributes for response codes and other information
		}
	} catch (Exception e) {
		span.recordException(e);
		span.setStatus(StatusCode.ERROR, e.getClass().getSimpleName());
		throw e;
	} finally {
		span.end();
	}
}
按照约定填充属性！如果没有适用的，请查看一般约定。

Nested network and other spans
网络调用通常通过相应的客户端实现使用 OpenTelemetry 自动检测进行跟踪。
...
如果 OpenTelemetry 不支持跟踪您的网络客户端，请使用您的最佳判断，这里有一些注意事项可以提供帮助：
	1.跟踪网络调用是否会提高用户的可观察性或您支持他们的能力？
	2.你的库是公共的、有文档记录的 RPC API 之上的包装器吗？如果出现问题，用户是否需要从底层服务获得支持？
		(1) 检测库并确保跟踪单个网络尝试
	3.使用 span 跟踪这些调用会非常冗长吗？还是会显着影响性能？
		(1)	使用带有详细信息或跨度事件的日志：日志可以与父级（公共 API 调用）相关联，而跨度事件应该在公共 API 跨度上设置。
		(2) 如果它们必须是跨度（携带和传播唯一的跟踪上下文），请将它们放在配置选项后面并默认禁用它们。

如果 OpenTelemetry 已经支持跟踪您的网络调用，您可能不想复制它。可能有一些例外：
	支持没有自动检测的用户（在某些环境中可能无法工作，或者用户可能担心猴子补丁）
	使用底层服务启用自定义（遗留）关联和上下文传播协议
	使用自动检测未涵盖的绝对必要的库/服务特定信息丰富 RPC 跨度
警告：正在构建避免重复的通用解决方案🚧。

Events
Traces 是您的应用程序可以发出的一种信号。事件（或日志）和跟踪相互补充，而不是重复。每当你有一些应该冗长的东西时，日志是比跟踪更好的选择。

很可能您的应用程序已经使用了日志记录或一些类似的模块。您的模块可能已经集成了 OpenTelemetry——要查找，请查看注册表。集成通常会在所有日志上标
记活动跟踪上下文，因此用户可以将它们关联起来。

如果您的语言和生态系统不支持通用日志记录，请使用跨度事件来共享其他应用程序详细信息。如果您还想添加属性，事件可能会更方便。
根据经验，使用事件或日志来获取详细数据而不是跨度。始终将事件附加到您的检测创建的跨度实例。如果可以，请避免使用活动范围，因为您无法控制它所指的内容。

Context propagation
Extracting context
如果您在接收上游调用的库或服务上工作，例如一个网络框架或一个消息传递消费者，你应该从传入的请求/消息中提取上下文。 OpenTelemetry 提供了 Propagator
API，它隐藏了特定的传播标准并从线路中读取跟踪上下文。在单个响应的情况下，线路上只有一个上下文，它成为库创建的新跨度的父级。

创建跨度后，您应该通过激活跨度将新的跟踪上下文传递给应用程序代码（回调或处理程序）；如果可能，您应该明确地这样做。
// extract the context
Context extractedContext = propagator.extract(Context.current(), httpExchange, getter);
Span span = tracer.spanBuilder("receive")
			.setSpanKind(SpanKind.SERVER)
			.setParent(extractedContext)
			.startSpan();

// make span active so any nested telemetry is correlated
try (Scope unused = span.makeCurrent()) {
	userCode();
} catch (Exception e) {
	span.recordException(e);
	span.setStatus(StatusCode.ERROR);
	throw e;
} finally {
	span.end();
}
以下是 Java 中上下文提取的完整示例，请查看您的语言的 OpenTelemetry 文档。
在消息系统的情况下，您可能会同时收到多条消息。收到的消息成为您创建的跨度上的链接。有关详细信息，请参阅消息传递约定（警告：消息传递约定正在构建中🚧）。

Injecting context
当您进行出站调用时，您通常希望将上下文传播到下游服务。在这种情况下，您应该创建一个新的 span 来跟踪传出调用并使用 Propagator API 将上下文注
入到消息中。在其他情况下，您可能想要注入上下文，例如在为异步处理创建消息时。
Span span = tracer.spanBuilder("send")
			.setSpanKind(SpanKind.CLIENT)
			.startSpan();

// make span active so any nested telemetry is correlated
// even network calls might have nested layers of spans, logs or events
try (Scope unused = span.makeCurrent()) {
	// inject the context
	propagator.inject(Context.current(), transportLayer, setter);
	send();
} catch (Exception e) {
	span.recordException(e);
	span.setStatus(StatusCode.ERROR);
	throw e;
} finally {
	span.end();
}
这是 Java 中上下文注入的完整示例。
可能有一些例外：
	1.下游服务不支持元数据或禁止未知字段
	2.下游服务没有定义关联协议。未来的某个服务版本是否可能支持兼容的上下文传播？注射它！
	3.下游服务支持自定义关联协议。
		(1) 自定义传播器的最大努力：如果兼容，请使用 OpenTelemetry 跟踪上下文。
		(2) 或在跨度上生成并标记自定义相关 ID。

In-process
	1.使您的跨度处于活动状态（又名当前）：它可以将跨度与日志和任何嵌套的自动检测相关联。
	2.如果库有上下文的概念，除了活动跨度之外，还支持可选的显式跟踪上下文传播
		(1) 将库创建的跨度（跟踪上下文）明确地放在上下文中，记录如何访问它
		(2) 允许用户在您的上下文中传递跟踪上下文
	3.在库中，显式传播跟踪上下文 - 活动跨度可能会在回调期间发生变化！
		(1) 尽快从公共 API 表面上的用户捕获活动上下文，将其用作跨度的父上下文
		(2) 传递上下文并在显式传播的实例上标记属性、异常和事件
		(3) 如果您显式启动线程、进行后台处理或其他由于您的语言中的异步上下文流限制而可能中断的事情，这是必不可少的

Metrics
指标 API 还不稳定，我们还没有定义指标约定。

Misc
Instrumentation registry
请将您的检测库添加到 OpenTelemetry 注册表，以便用户可以找到它。

Performance
当应用程序中没有 SDK 时，OpenTelemetry API 是空操作并且性能非常好。配置 OpenTelemetry SDK 时，它会消耗绑定资源。

现实生活中的应用程序，尤其是大规模应用程序，通常会配置基于头部的采样。采样出的跨度很便宜，您可以检查跨度是否正在记录，以避免在填充属性时进行额
外的分配和潜在的昂贵计算。

// some attributes are important for sampling, they should be provided at creation time
Span span = tracer.spanBuilder(String.format("SELECT %s.%s", dbName, collectionName))
		.setSpanKind(SpanKind.CLIENT)
		.setAttribute("db.name", dbName)
		...
		.startSpan();

// other attributes, especially those that are expensive to calculate
// should be added if span is recording
if (span.isRecording()) {
	span.setAttribute("db.statement", sanitize(query.statement()))
}

Error handling
OpenTelemetry API 在运行时是宽容的——不会因无效参数而失败，从不抛出和吞下异常。这样，检测问题就不会影响应用程序逻辑。测试仪器以注意 OpenTelemetry 在运行时隐藏的问题。

Testing
由于 OpenTelemetry 具有多种自动检测，因此尝试您的检测如何与其他遥测交互非常有用：传入请求、传出请求、日志等。使用典型的应用程序，在尝试检测
时启用流行的框架和库并启用所有跟踪.查看与您的图书馆相似的图书馆是如何出现的。

对于单元测试，您通常可以模拟或伪造 SpanProcessor 和 SpanExporter。
@Test
public void checkInstrumentation() {
	SpanExporter exporter = new TestExporter();
	
	Tracer tracer = OpenTelemetrySdk.builder()
				.setTracerProvider(SdkTracerProvider.builder()
				.addSpanProcessor(SimpleSpanProcessor.create(exporter)).build()).build()
				.getTracer("test");
	// run test ...
	
	validateSpans(exporter.exportedSpans);
}

class TestExporter implements SpanExporter {
	public final List<SpanData> exportedSpans = Collections.synchronizedList(new ArrayList<>());
	
	@Override
	public CompletableResultCode export(Collection<SpanData> spans) {
		exportedSpans.addAll(spans);
		return CompletableResultCode.ofSuccess();
	}
	...
}
