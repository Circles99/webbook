package article

import (
	"context"
	"errors"
	"github.com/bwmarrin/snowflake"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

type MongoDBDAO struct {

	// 制作库
	col *mongo.Collection
	// 线上库
	liveCol *mongo.Collection
	node    *snowflake.Node
}

func InitCollections(db *mongo.Database) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()
	index := []mongo.IndexModel{
		{
			Keys:    bson.D{bson.E{Key: "id", Value: 1}},
			Options: options.Index().SetUnique(true),
		},
		{
			Keys: bson.D{bson.E{Key: "author_id", Value: 1},
				bson.E{Key: "ctime", Value: 1},
			},
			Options: options.Index(),
		},
	}
	_, err := db.Collection("articles").Indexes().
		CreateMany(ctx, index)
	if err != nil {
		return err
	}
	_, err = db.Collection("published_articles").Indexes().
		CreateMany(ctx, index)
	return err
}

func NewMongoDBDAO(db *mongo.Database, node *snowflake.Node) ArticleDAO {
	return &MongoDBDAO{

		col:     db.Collection("articles"),
		liveCol: db.Collection("published_articles"),
		node:    node,
	}
}

func (m MongoDBDAO) GetByAuthor(ctx context.Context, userId int64, offset int, limit int) ([]Article, error) {
	//TODO implement me
	panic("implement me")
}

func (m MongoDBDAO) Insert(ctx context.Context, art Article) (int64, error) {
	now := time.Now().UnixMilli()
	art.Created = now
	art.Updated = now
	id := m.node.Generate().Int64()
	art.Id = id
	_, err := m.col.InsertOne(ctx, art)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (m MongoDBDAO) Update(ctx context.Context, art Article) error {
	filter := bson.M{"id": art.Id, "author_id": art.AuthorId}
	update := bson.M{
		"$set": bson.M{
			"title":   art.Title,
			"content": art.Content,
			"status":  art.Status,
			"updated": time.Now().UnixMilli(),
		},
	}
	res, err := m.col.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}

	if res.ModifiedCount == 0 {
		return errors.New("更新失败")
	}

	return nil
}

func (m MongoDBDAO) Sync(ctx context.Context, art Article) (int64, error) {
	var (
		id  = art.Id
		err error
	)

	if id > 0 {
		err = m.Update(ctx, art)
	} else {
		id, err = m.Insert(ctx, art)
	}
	if err != nil {
		return 0, err
	}

	art.Id = id

	err = m.Upsert(ctx, PublishArticle{art})
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (m MongoDBDAO) Upsert(ctx context.Context, art PublishArticle) error {

	now := time.Now().UnixMilli()
	art.Updated = now
	update := bson.E{"$set", art}
	upsert := bson.E{"$setOnInsert", bson.D{bson.E{"created", now}}}
	filter := bson.M{"id": art.Id, "author_id": art.AuthorId}

	_, err := m.liveCol.UpdateOne(ctx, filter, bson.D{update, upsert}, options.Update().SetUpsert(true))
	if err != nil {
		return err
	}
	return nil
}

func (m MongoDBDAO) SyncStatus(ctx context.Context, id int64, authorId int64, status uint8) error {
	filter := bson.D{bson.E{Key: "id", Value: id},
		bson.E{Key: "author_id", Value: authorId}}
	sets := bson.D{bson.E{Key: "$set",
		Value: bson.D{bson.E{Key: "status", Value: status}}}}
	res, err := m.col.UpdateOne(ctx, filter, sets)
	if err != nil {
		return err
	}
	if res.ModifiedCount != 1 {
		return errors.New("更新失败")
	}
	return nil
}
