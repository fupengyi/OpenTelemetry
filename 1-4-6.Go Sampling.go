package OpenTelemetry_Go
Sampling
采样是一个限制系统生成的跟踪数量的过程。您应该使用的确切采样器取决于您的特定需求，但通常您应该在跟踪开始时做出决定，并允许采样决定传播到其他服务。
配置时需要在跟踪器提供程序上设置采样器，如下所示：
	provider := sdktrace.NewTracerProvider(
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
	)
AlwaysSample 和 NeverSample 是不言自明的。 Always 意味着每条迹线都将被采样，反之亦然。当您刚开始或在开发环境中时，您几乎总是想使用 AlwaysSample。
其他采样器包括：
	TraceIDRatioBased，它将根据提供给采样器的分数对一部分轨迹进行采样。因此，如果将其设置为 .5，将对一半的迹线进行采样。
	ParentBased，它根据传入的抽样决策表现不同。通常，这将对具有已采样父级的跨度进行采样，并且不会对未采样其父级的跨度进行采样。
在生产环境中，您应该考虑将 TraceIDRatioBased 采样器与 ParentBased 采样器一起使用。
