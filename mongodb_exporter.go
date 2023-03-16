package mongodbexporter

import (
	"context"

	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/pdata/plog"
	"go.opentelemetry.io/collector/pdata/ptrace"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var traceMarshaller = ptrace.JSONMarshaler{}
var logsMarshaller = plog.JSONMarshaler{}

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
	bsonDoc, err := ptraceToBsonDoc(&td)
	if err != nil {
		return err
	}
	_, err = mdbe.collectionTraces.InsertOne(ctx, bsonDoc, options.InsertOne().SetBypassDocumentValidation(true))
	if err != nil {
		return err
	}
	return nil
}

func (mdbe *mongoDbExporter) ConsumeLogs(ctx context.Context, ld plog.Logs) error {
	bsonDoc, err := plogToBsonDoc(&ld)
	if err != nil {
		return err
	}
	_, err = mdbe.collectionLogs.InsertOne(ctx, bsonDoc, options.InsertOne().SetBypassDocumentValidation(true))
	if err != nil {
		return err
	}
	return nil
}

func ptraceToBsonDoc(trace *ptrace.Traces) (interface{}, error) {
	data, err := traceMarshaller.MarshalTraces(*trace)
	if err != nil {
		return nil, err
	}
	var bsonDoc interface{}
	err = bson.UnmarshalExtJSON(data, true, &bsonDoc)
	if err != nil {
		return nil, err
	}
	return bsonDoc, nil
}

func plogToBsonDoc(ld *plog.Logs) (interface{}, error) {
	data, err := logsMarshaller.MarshalLogs(*ld)
	if err != nil {
		return nil, err
	}
	var bsonDoc interface{}
	err = bson.UnmarshalExtJSON(data, true, &bsonDoc)
	if err != nil {
		return nil, err
	}
	return bsonDoc, nil
}

func (mdbe *mongoDbExporter) Start(context.Context, component.Host) error {
	return nil
}

// Shutdown stops the exporter and is invoked during shutdown.
func (mdbe *mongoDbExporter) Shutdown(ctx context.Context) error {
	return mdbe.client.Disconnect(ctx)
}
