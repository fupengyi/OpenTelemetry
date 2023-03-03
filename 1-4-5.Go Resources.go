package OpenTelemetry_Go
Resources
Resources 是一种特殊类型的属性，适用于流程生成的所有跨度。这些应该用于表示关于非临时进程的底层元数据——例如，进程的主机名或其实例 ID。
Resources 应该在初始化时分配给跟踪器提供者，并且创建起来很像属性：
	resources := resource.NewWithAttributes(
		semconv.SchemaURL,
		semconv.ServiceNameKey.String("myService"),
		semconv.ServiceVersionKey.String("1.0.0"),
		semconv.ServiceInstanceIDKey.String("abcdef12345"),
	)
	
	provider := sdktrace.NewTracerProvider(
		...
		sdktrace.WithResource(resources),
	)
请注意使用 semconv 包为资源属性提供常规名称。这有助于确保使用这些语义约定生成的遥测数据的消费者可以轻松发现相关属性并理解它们的含义。

Resources 也可以通过 resource.Detector 实现自动检测。这些 Detectors 可能会发现有关当前正在运行的进程、它正在运行的操作系统、托管该操作
系统实例的云提供商或任何数量的其他资源属性的信息。
	resources := resource.New(context.Background(),
		resource.WithFromEnv(), // pull attributes from OTEL_RESOURCE_ATTRIBUTES and OTEL_SERVICE_NAME environment variables	// 从 OTEL_RESOURCE_ATTRIBUTES 和 OTEL_SERVICE_NAME 环境变量中提取属性
		resource.WithProcess(), // This option configures a set of Detectors that discover process information					// 该选项配置了一组发现进程信息的检测器
		resource.WithOS(), // This option configures a set of Detectors that discover OS information							// 此选项配置一组发现操作系统信息的检测器
		resource.WithContainer(), // This option configures a set of Detectors that discover container information				// 该选项配置一组发现容器信息的检测器
		resource.WithHost(), // This option configures a set of Detectors that discover host information						// 该选项配置一组发现主机信息的检测器
		resource.WithDetectors(thirdparty.Detector{}), // Bring your own external Detector implementation						// 自带外部检测器实现
		resource.WithAttributes(attribute.String("foo", "bar")), // Or specify resource attributes directly						// 或者直接指定资源属性
)
