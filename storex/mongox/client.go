package mongox

import (
	"context"
	"github.com/ivalue2333/pkg/logx"
	"github.com/ivalue2333/pkg/storex/mongox/mongoc_wrapped"
	"go.mongodb.org/mongo-driver/mongo"
	mgOptions "go.mongodb.org/mongo-driver/mongo/options"
)

type (
	Model struct {
		client         *concurrentClient
		dbName         string
		collectionName string
		opts           []Option
	}
)

func MustNewModel(url, collection string, opts ...Option) *Model {
	model, err := NewModel(url, collection, opts...)
	if err != nil {
		logx.Fatalf(context.Background(), "NewModel failed: err:%+v", err)
	}
	return model
}

func NewModel(url, collection string, opts ...Option) (*Model, error) {
	concurrentClient, err := getConcurrentClient(url)
	if err != nil {
		return nil, err
	}

	return &Model{
		client:         concurrentClient,
		collectionName: collection,
		opts:           opts,
	}, nil
}

func (b *Model) DatabaseName() string {
	return b.dbName
}

func (b *Model) CollectionName() string {
	return b.collectionName
}

func (b *Model) Client(ctx context.Context) *mongoc_wrapped.WrappedClient {
	return b.client.WrappedClient
	//c, err := b.client.takeClient()
	//if err != nil {
	//	logx.Fatalf(ctx, "takeClient failed, err:%+v", err)
	//}
	//return c
}

func (b *Model) collection(ctx context.Context) *mongo.Collection {
	return b.Client(ctx).Database(b.DatabaseName()).Collection(b.CollectionName()).Collection()
}

func (b *Model) Aggregate(
	ctx context.Context,
	pipeline, results interface{},
	opts ...*mgOptions.AggregateOptions,
) error {

	cur, err := b.collection(ctx).Aggregate(ctx, pipeline)
	if err == nil {
		err = cur.All(ctx, results)
	}
	return err
}

func (b *Model) BulkWrite(
	ctx context.Context,
	models []mongo.WriteModel,
	opts ...*mgOptions.BulkWriteOptions,
) (*mongo.BulkWriteResult, error) {

	bwres, err := b.collection(ctx).BulkWrite(ctx, models, opts...)

	return bwres, err
}

func (b *Model) Count(ctx context.Context, filter interface{}, opts ...*mgOptions.CountOptions) (int64, error) {
	count, err := b.collection(ctx).CountDocuments(ctx, filter, opts...)

	return count, err
}

func (b *Model) CountDocuments(ctx context.Context, filter interface{}, opts ...*mgOptions.CountOptions) (int64, error) {

	count, err := b.collection(ctx).CountDocuments(ctx, filter, opts...)

	return count, err
}

func (b *Model) DeleteMany(
	ctx context.Context,
	filter interface{},
	opts ...*mgOptions.DeleteOptions,
) (*mongo.DeleteResult, error) {

	dmres, err := b.collection(ctx).DeleteMany(ctx, filter, opts...)

	return dmres, err
}

func (b *Model) DeleteOne(
	ctx context.Context,
	filter interface{},
	opts ...*mgOptions.DeleteOptions,
) (*mongo.DeleteResult, error) {

	dor, err := b.collection(ctx).DeleteOne(ctx, filter, opts...)

	return dor, err
}

func (b *Model) Distinct(
	ctx context.Context,
	fieldName string,
	filter interface{},
	opts ...*mgOptions.DistinctOptions,
) ([]interface{}, error) {

	distinct, err := b.collection(ctx).Distinct(ctx, fieldName, filter, opts...)

	return distinct, err
}

func (b *Model) Drop(ctx context.Context) error {

	err := b.collection(ctx).Drop(ctx)

	return err
}

func (b *Model) EstimatedDocumentCount(
	ctx context.Context,
	opts ...*mgOptions.EstimatedDocumentCountOptions,
) (int64, error) {

	count, err := b.collection(ctx).EstimatedDocumentCount(ctx, opts...)

	return count, err
}

func (b *Model) Find(ctx context.Context, filter interface{}, results interface{}, opts ...*mgOptions.FindOptions) error {

	cur, err := b.collection(ctx).Find(ctx, filter, opts...)
	if err == nil {
		err = cur.All(ctx, results)
	}

	return err
}

func (b *Model) FindOne(
	ctx context.Context,
	filter interface{}, result interface{},
	opts ...*mgOptions.FindOneOptions,
) error {

	return b.collection(ctx).FindOne(ctx, filter, opts...).Decode(result)
}

func (b *Model) FindOneAndDelete(
	ctx context.Context,
	filter interface{}, result interface{},
	opts ...*mgOptions.FindOneAndDeleteOptions,
) error {

	return b.collection(ctx).FindOneAndDelete(ctx, filter, opts...).Decode(result)
}

func (b *Model) FindOneAndReplace(
	ctx context.Context,
	filter, replacement, result interface{},
	opts ...*mgOptions.FindOneAndReplaceOptions,
) error {

	return b.collection(ctx).FindOneAndReplace(ctx, filter, replacement, opts...).Decode(result)
}

func (b *Model) FindOneAndUpdate(
	ctx context.Context,
	filter, update, result interface{},
	opts ...*mgOptions.FindOneAndUpdateOptions,
) error {

	return b.collection(ctx).FindOneAndUpdate(ctx, filter, update, opts...).Decode(result)
}

func (b *Model) FindOneAndUpsert(
	ctx context.Context,
	filter, update, result interface{},
	opts ...*mgOptions.FindOneAndUpdateOptions,
) error {

	rd := mgOptions.After
	optUpsert := mgOptions.FindOneAndUpdate().SetUpsert(true).SetReturnDocument(rd)
	opts = append(opts, optUpsert)

	return b.collection(ctx).FindOneAndUpdate(ctx, filter, update, opts...).Decode(result)
}

func (b *Model) Indexes() mongo.IndexView { return b.c(context.Background()).Indexes() }

func (b *Model) InsertMany(
	ctx context.Context,
	documents []interface{},
	opts ...*mgOptions.InsertManyOptions,
) (*mongo.InsertManyResult, error) {

	insmres, err := b.collection(ctx).InsertMany(ctx, documents, opts...)

	return insmres, err
}

func (b *Model) InsertOne(
	ctx context.Context,
	document interface{},
	opts ...*mgOptions.InsertOneOptions,
) (*mongo.InsertOneResult, error) {

	insores, err := b.collection(ctx).InsertOne(ctx, document, opts...)

	return insores, err
}

func (b *Model) ReplaceOne(
	ctx context.Context,
	filter, replacement interface{},
	opts ...*mgOptions.ReplaceOptions,
) (*mongo.UpdateResult, error) {

	repres, err := b.collection(ctx).ReplaceOne(ctx, filter, replacement, opts...)

	return repres, err
}

func (b *Model) UpdateMany(ctx context.Context, filter, replacement interface{}, opts ...*mgOptions.UpdateOptions) (*mongo.UpdateResult, error) {

	umres, err := b.collection(ctx).UpdateMany(ctx, filter, replacement, opts...)

	return umres, err
}

func (b *Model) UpdateOne(
	ctx context.Context,
	filter, replacement interface{},
	opts ...*mgOptions.UpdateOptions,
) (*mongo.UpdateResult, error) {

	uores, err := b.collection(ctx).UpdateOne(ctx, filter, replacement, opts...)

	return uores, err
}

func (b *Model) Upsert(
	ctx context.Context,
	filter, replacement interface{},
	opts ...*mgOptions.UpdateOptions,
) (*mongo.UpdateResult, error) {

	optUpsert := mgOptions.Update().SetUpsert(true)
	opts = append(opts, optUpsert)
	uores, err := b.collection(ctx).UpdateOne(ctx, filter, replacement, opts...)

	return uores, err
}
