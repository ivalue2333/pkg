package mongox

import (
	"context"
	"github.com/ivalue2333/pkg/src/logx"
	"github.com/ivalue2333/pkg/src/storex/mongox/mongoc_wrapped"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type (
	Client struct {
		WrappedClient  *mongoc_wrapped.WrappedClient
		dbName         string
		collectionName string
	}
)

func MustNewClient(ctx context.Context, uri, dbName, collectionName string) *Client {
	client, err := NewClient(ctx, uri, dbName, collectionName)
	if err != nil {
		logx.Fatalf(context.Background(), "NewClient failed: err:%+v", err)
	}
	return client
}

func NewClient(ctx context.Context, uri, dbName, collectionName string) (*Client, error) {
	c, err := mongoc_wrapped.NewClient(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		return nil, err
	}

	if err = c.Connect(ctx); err != nil {
		return nil, err
	}

	if err = c.Ping(ctx, nil); err != nil { // ping before use
		return nil, err
	}

	return &Client{
		WrappedClient:  c,
		dbName:         dbName,
		collectionName: collectionName,
	}, nil
}

func (b *Client) DatabaseName() string {
	return b.dbName
}

func (b *Client) CollectionName() string {
	return b.collectionName
}

func (b *Client) Client(ctx context.Context) *mongoc_wrapped.WrappedClient {
	return b.WrappedClient
}

func (b *Client) collection(ctx context.Context) *mongo.Collection {
	return b.Client(ctx).Database(b.DatabaseName()).Collection(b.CollectionName()).Collection()
}

func (b *Client) Aggregate(
	ctx context.Context,
	pipeline, results interface{},
	opts ...*options.AggregateOptions,
) error {

	cur, err := b.collection(ctx).Aggregate(ctx, pipeline)
	if err == nil {
		err = cur.All(ctx, results)
	}
	return err
}

func (b *Client) BulkWrite(
	ctx context.Context,
	clients []mongo.WriteModel,
	opts ...*options.BulkWriteOptions,
) (*mongo.BulkWriteResult, error) {

	bwres, err := b.collection(ctx).BulkWrite(ctx, clients, opts...)

	return bwres, err
}

func (b *Client) Count(ctx context.Context, filter interface{}, opts ...*options.CountOptions) (int64, error) {
	count, err := b.collection(ctx).CountDocuments(ctx, filter, opts...)

	return count, err
}

func (b *Client) CountDocuments(ctx context.Context, filter interface{}, opts ...*options.CountOptions) (int64, error) {

	count, err := b.collection(ctx).CountDocuments(ctx, filter, opts...)

	return count, err
}

func (b *Client) DeleteMany(
	ctx context.Context,
	filter interface{},
	opts ...*options.DeleteOptions,
) (*mongo.DeleteResult, error) {

	dmres, err := b.collection(ctx).DeleteMany(ctx, filter, opts...)

	return dmres, err
}

func (b *Client) DeleteOne(
	ctx context.Context,
	filter interface{},
	opts ...*options.DeleteOptions,
) (*mongo.DeleteResult, error) {

	dor, err := b.collection(ctx).DeleteOne(ctx, filter, opts...)

	return dor, err
}

func (b *Client) Distinct(
	ctx context.Context,
	fieldName string,
	filter interface{},
	opts ...*options.DistinctOptions,
) ([]interface{}, error) {

	distinct, err := b.collection(ctx).Distinct(ctx, fieldName, filter, opts...)

	return distinct, err
}

func (b *Client) Drop(ctx context.Context) error {

	err := b.collection(ctx).Drop(ctx)

	return err
}

func (b *Client) EstimatedDocumentCount(
	ctx context.Context,
	opts ...*options.EstimatedDocumentCountOptions,
) (int64, error) {

	count, err := b.collection(ctx).EstimatedDocumentCount(ctx, opts...)

	return count, err
}

func (b *Client) Find(ctx context.Context, filter interface{}, results interface{}, opts ...*options.FindOptions) error {

	cur, err := b.collection(ctx).Find(ctx, filter, opts...)
	if err == nil {
		err = cur.All(ctx, results)
	}

	return err
}

func (b *Client) FindOne(
	ctx context.Context,
	filter interface{}, result interface{},
	opts ...*options.FindOneOptions,
) error {

	return b.collection(ctx).FindOne(ctx, filter, opts...).Decode(result)
}

func (b *Client) FindOneAndDelete(
	ctx context.Context,
	filter interface{}, result interface{},
	opts ...*options.FindOneAndDeleteOptions,
) error {

	return b.collection(ctx).FindOneAndDelete(ctx, filter, opts...).Decode(result)
}

func (b *Client) FindOneAndReplace(
	ctx context.Context,
	filter, replacement, result interface{},
	opts ...*options.FindOneAndReplaceOptions,
) error {

	return b.collection(ctx).FindOneAndReplace(ctx, filter, replacement, opts...).Decode(result)
}

func (b *Client) FindOneAndUpdate(
	ctx context.Context,
	filter, update, result interface{},
	opts ...*options.FindOneAndUpdateOptions,
) error {

	return b.collection(ctx).FindOneAndUpdate(ctx, filter, update, opts...).Decode(result)
}

func (b *Client) FindOneAndUpsert(
	ctx context.Context,
	filter, update, result interface{},
	opts ...*options.FindOneAndUpdateOptions,
) error {

	rd := options.After
	optUpsert := options.FindOneAndUpdate().SetUpsert(true).SetReturnDocument(rd)
	opts = append(opts, optUpsert)

	return b.collection(ctx).FindOneAndUpdate(ctx, filter, update, opts...).Decode(result)
}

func (b *Client) Indexes() mongo.IndexView { return b.collection(context.Background()).Indexes() }

func (b *Client) InsertMany(
	ctx context.Context,
	documents []interface{},
	opts ...*options.InsertManyOptions,
) (*mongo.InsertManyResult, error) {

	insmres, err := b.collection(ctx).InsertMany(ctx, documents, opts...)

	return insmres, err
}

func (b *Client) InsertOne(
	ctx context.Context,
	document interface{},
	opts ...*options.InsertOneOptions,
) (*mongo.InsertOneResult, error) {

	insores, err := b.collection(ctx).InsertOne(ctx, document, opts...)

	return insores, err
}

func (b *Client) ReplaceOne(
	ctx context.Context,
	filter, replacement interface{},
	opts ...*options.ReplaceOptions,
) (*mongo.UpdateResult, error) {

	repres, err := b.collection(ctx).ReplaceOne(ctx, filter, replacement, opts...)

	return repres, err
}

func (b *Client) UpdateMany(ctx context.Context, filter, replacement interface{}, opts ...*options.UpdateOptions) (*mongo.UpdateResult, error) {

	umres, err := b.collection(ctx).UpdateMany(ctx, filter, replacement, opts...)

	return umres, err
}

func (b *Client) UpdateOne(
	ctx context.Context,
	filter, replacement interface{},
	opts ...*options.UpdateOptions,
) (*mongo.UpdateResult, error) {

	uores, err := b.collection(ctx).UpdateOne(ctx, filter, replacement, opts...)

	return uores, err
}

func (b *Client) Upsert(
	ctx context.Context,
	filter, replacement interface{},
	opts ...*options.UpdateOptions,
) (*mongo.UpdateResult, error) {

	optUpsert := options.Update().SetUpsert(true)
	opts = append(opts, optUpsert)
	uores, err := b.collection(ctx).UpdateOne(ctx, filter, replacement, opts...)

	return uores, err
}
