package operator

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"worker/consts"
	"worker/dal"
	"worker/model"
)

type MongoOperator struct {
	mongoClient *dal.MongoClient
}

func NewMongoOperator(mongoClient *dal.MongoClient) (*MongoOperator, error) {
	return &MongoOperator{mongoClient: mongoClient}, nil
}

func (mo *MongoOperator) DropDataBase() error {
	return mo.mongoClient.GetDataBase().Drop(context.TODO())
}

func (mo *MongoOperator) DropCollection(collection string) error {
	return mo.mongoClient.GetDataBase().Collection(collection).Drop(context.TODO())
}

func (mo *MongoOperator) DeleteStorage(storageAddress string) error {
	filter := bson.M{"storageAddress": storageAddress}
	_, err := mo.mongoClient.GetDataBase().Collection(consts.Storages).DeleteOne(context.TODO(), filter)

	return err
}

func (mo *MongoOperator) GetStorages() ([]model.Storage, error) {
	filter := bson.D{{}}
	cursor, err := mo.mongoClient.GetDataBase().Collection(consts.Storages).Find(context.TODO(), filter)
	if err != nil {
		return nil, err
	}

	defer func(cursor *mongo.Cursor, ctx context.Context) {
		_ = cursor.Close(ctx)
	}(cursor, context.TODO())
	var storages []model.Storage
	err = cursor.All(context.TODO(), &storages)
	if err != nil {
		return nil, err
	}

	return storages, nil
}

func (mo *MongoOperator) UpdateMappingInfoDeleteOldGroupID(shardIDStart, shardIDEnd uint32) error {
	filter := bson.M{"shardIDStart": shardIDStart, "shardIDEnd": shardIDEnd}
	update := bson.M{"$set": bson.M{"oldGroupID": ""}}
	_, err := mo.mongoClient.GetDataBase().Collection(consts.MappingInfos).UpdateOne(context.TODO(), filter, update)

	return err
}

func (mo *MongoOperator) GetMappingInfos() ([]model.MappingInfo, error) {
	filter := bson.D{{}}
	cursor, err := mo.mongoClient.GetDataBase().Collection(consts.MappingInfos).Find(context.TODO(), filter)
	if err != nil {
		return nil, err
	}

	defer func(cursor *mongo.Cursor, ctx context.Context) {
		_ = cursor.Close(ctx)
	}(cursor, context.TODO())
	var mappingInfos []model.MappingInfo
	err = cursor.All(context.TODO(), &mappingInfos)
	if err != nil {
		return nil, err
	}

	return mappingInfos, nil
}

func (mo *MongoOperator) UpdateGroupAvailableCap(groupID string, increase int) error {
	filter := bson.M{"groupID": groupID}
	update := bson.M{"$inc": bson.M{"availableCap": increase}}
	_, err := mo.mongoClient.GetDataBase().Collection(consts.Groups).UpdateOne(context.TODO(), filter, update)

	return err
}

func (mo *MongoOperator) GetGroup(groupID string) (*model.Group, error) {
	filter := bson.M{"groupID": groupID}
	findRes := mo.mongoClient.GetDataBase().Collection(consts.Groups).FindOne(context.TODO(), filter)

	var group model.Group
	err := findRes.Decode(&group)
	if err != nil {
		return nil, err
	}

	return &group, nil
}

func (mo *MongoOperator) GetGroups() ([]model.Group, error) {
	filter := bson.D{{}}
	cursor, err := mo.mongoClient.GetDataBase().Collection(consts.Groups).Find(context.TODO(), filter)
	if err != nil {
		return nil, err
	}

	defer func(cursor *mongo.Cursor, ctx context.Context) {
		_ = cursor.Close(ctx)
	}(cursor, context.TODO())
	var groups []model.Group
	err = cursor.All(context.TODO(), &groups)
	if err != nil {
		return nil, err
	}

	return groups, nil
}

func (mo *MongoOperator) UpdateUserAppendBucket(userName string, bucket *model.Bucket) error {
	filter := bson.M{"userName": userName}
	update := bson.M{"$push": bson.M{"buckets": bucket}}
	_, err := mo.mongoClient.GetDataBase().Collection(consts.Users).UpdateOne(context.TODO(), filter, update)

	return err
}

func (mo *MongoOperator) UpdateUserDeleteBucket(userName string, bucketName string) error {
	filter := bson.M{"userName": userName}
	update := bson.M{"$pull": bson.M{"buckets": bson.M{"bucketName": bucketName}}}
	_, err := mo.mongoClient.GetDataBase().Collection(consts.Users).UpdateOne(context.TODO(), filter, update)

	return err
}

func (mo *MongoOperator) UpdateUserCount(userName string, bucketName string, fileIncrease int64) error {
	filter := bson.M{"userName": userName, "buckets.bucketName": bucketName}
	update := bson.M{"$inc": bson.M{"fileCount": fileIncrease, "operationCount": 1, "buckets.$[item].fileCount": fileIncrease, "buckets.$[item].operationCount": 1}}
	arrayFilter := bson.M{"item.bucketName": bucketName}
	_, err := mo.mongoClient.GetDataBase().Collection(consts.Users).UpdateOne(context.TODO(), filter, update, options.Update().SetArrayFilters(
		options.ArrayFilters{
			Filters: []interface{}{
				arrayFilter,
			},
		},
	))

	return err
}

func (mo *MongoOperator) GetUser(userName string) (*model.User, error) {
	filter := bson.M{"userName": userName}
	findRes := mo.mongoClient.GetDataBase().Collection(consts.Users).FindOne(context.TODO(), filter)

	var user model.User
	err := findRes.Decode(&user)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (mo *MongoOperator) GetUsers() ([]model.User, error) {
	filter := bson.D{{}}
	cursor, err := mo.mongoClient.GetDataBase().Collection(consts.Users).Find(context.TODO(), filter)
	if err != nil {
		return nil, err
	}

	defer func(cursor *mongo.Cursor, ctx context.Context) {
		_ = cursor.Close(ctx)
	}(cursor, context.TODO())
	var users []model.User
	err = cursor.All(context.TODO(), &users)
	if err != nil {
		return nil, err
	}

	return users, nil
}

func (mo *MongoOperator) InsertChunk(chunk *model.Chunk) error {
	_, err := mo.mongoClient.GetDataBase().Collection(consts.Chunks).InsertOne(context.TODO(), chunk)

	return err
}

func (mo *MongoOperator) DeleteChunk(chunkID string) error {
	filter := bson.M{"chunkID": chunkID}
	_, err := mo.mongoClient.GetDataBase().Collection(consts.Chunks).DeleteOne(context.TODO(), filter)

	return err
}

func (mo *MongoOperator) GetChunk(chunkID string) (*model.Chunk, error) {
	filter := bson.M{"chunkID": chunkID}
	findRes := mo.mongoClient.GetDataBase().Collection(consts.Chunks).FindOne(context.TODO(), filter)

	var chunk model.Chunk
	err := findRes.Decode(&chunk)
	if err != nil {
		return nil, err
	}

	return &chunk, nil
}

func (mo *MongoOperator) GetFileNamesByUserNameAndBucketName(userName string, bucketName string) ([]string, error) {
	filter := bson.M{"userName": userName, "bucketName": bucketName}
	cursor, err := mo.mongoClient.GetDataBase().Collection(consts.Chunks).Find(context.TODO(), filter)
	if err != nil {
		return nil, err
	}

	defer func(cursor *mongo.Cursor, ctx context.Context) {
		_ = cursor.Close(ctx)
	}(cursor, context.TODO())
	var chunks []model.Chunk
	err = cursor.All(context.TODO(), &chunks)
	if err != nil {
		return nil, err
	}

	fileNamesSet := make(map[string]bool)
	for _, chunk := range chunks {
		fileNamesSet[chunk.FileName] = true
	}
	fileNames := make([]string, 0)
	for fileName, _ := range fileNamesSet {
		fileNames = append(fileNames, fileName)
	}

	return fileNames, nil
}
