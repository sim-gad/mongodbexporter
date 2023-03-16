package mongodbexporter

import (
	"context"

	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/pdata/plog"
	"go.opentelemetry.io/collector/pdata/ptrace"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// mongoDbExporter is the implementation of mongoDb exporter that writes telemetry data to a mongoDb collection
type mongoDbExporter struct {
	client           *mongo.Client
	collectionLogs   *mongo.Collection
	collectionTraces *mongo.Collection
}

func newMongoDbExporter(ctx context.Context, config *Config) (*mongoDbExporter, error) {
	if client, err := mongo.Connect(ctx, options.Client().ApplyURI(config.ConnectionURI)); err == nil {
		collectionLogs := client.Database(config.Database).Collection(config.CollectionLogs)
		collectionTraces := client.Database(config.Database).Collection(config.CollectionTraces)
		return &mongoDbExporter{
			client:           client,
			collectionLogs:   collectionLogs,
			collectionTraces: collectionTraces,
		}, nil
	} else {
		return nil, err
	}

}

func (mdbe *mongoDbExporter) ConsumeTraces(ctx context.Context, td ptrace.Traces) error {
	_, err := mdbe.collectionTraces.InsertOne(ctx, newTraceDoc(td), options.InsertOne().SetBypassDocumentValidation(true))
	if err != nil {
		return err
	}
	return nil
}

func (mdbe *mongoDbExporter) ConsumeLogs(ctx context.Context, ld plog.Logs) error {
	_, err := mdbe.collectionLogs.InsertOne(ctx, ld)
	if err != nil {
		return err
	}
	return nil
}

func (mdbe *mongoDbExporter) Start(context.Context, component.Host) error {
	return nil
}

// Shutdown stops the exporter and is invoked during shutdown.
func (mdbe *mongoDbExporter) Shutdown(ctx context.Context) error {
	return mdbe.client.Disconnect(ctx)
}
