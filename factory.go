package mongodbexporter

import (
	"context"

	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/exporter"
	"go.opentelemetry.io/collector/exporter/exporterhelper"
)

const (
	typeStr                 = "mongodb"
	stability               = component.StabilityLevelAlpha
	collectionTracesDefault = "Requests"
	collectionLogsDefault   = "Logs"
	databaseDefault         = "OtelDB"
)

// NewFactory creates a factory for OTLP exporter.
func NewFactory() exporter.Factory {
	return exporter.NewFactory(
		typeStr,
		createDefaultConfig,
		exporter.WithTraces(createTracesExporter, stability),
		exporter.WithLogs(createLogsExporter, stability))
}

func createDefaultConfig() component.Config {
	return &Config{
		CollectionLogs:   collectionLogsDefault,
		CollectionTraces: collectionTracesDefault,
		Database:         databaseDefault,
	}
}

func createTracesExporter(
	ctx context.Context,
	set exporter.CreateSettings,
	cfg component.Config,
) (exporter.Traces, error) {
	conf := cfg.(*Config)
	mongoDbExporter, err := newMongoDbExporter(ctx, conf)
	if err != nil {
		return nil, err
	}
	return exporterhelper.NewTracesExporter(
		ctx,
		set,
		cfg,
		mongoDbExporter.ConsumeTraces,
		exporterhelper.WithStart(mongoDbExporter.Start),
		exporterhelper.WithShutdown(mongoDbExporter.Shutdown),
	)
}

func createLogsExporter(
	ctx context.Context,
	set exporter.CreateSettings,
	cfg component.Config,
) (exporter.Logs, error) {
	conf := cfg.(*Config)
	mongoDbExporter, err := newMongoDbExporter(ctx, conf)
	if err != nil {
		return nil, err
	}
	return exporterhelper.NewLogsExporter(
		ctx,
		set,
		cfg,
		mongoDbExporter.ConsumeLogs,
		exporterhelper.WithStart(mongoDbExporter.Start),
		exporterhelper.WithShutdown(mongoDbExporter.Shutdown),
	)
}
