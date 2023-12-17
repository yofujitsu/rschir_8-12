package storage

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/gridfs"
	"go.mongodb.org/mongo-driver/mongo/options"
	"io/ioutil"
	"os"
)

func ConnectToMongoDB() (*mongo.Client, error) {
	clientOptions := options.Client().ApplyURI(os.Getenv("MONGO_URL"))
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		return nil, err
	}
	return client, nil
}

func ReadGridFSFile(fs *gridfs.Bucket, fileID interface{}) (string, error) {
	downloadStream, err := fs.OpenDownloadStream(fileID)
	if err != nil {
		return "", err
	}
	defer func(downloadStream *gridfs.DownloadStream) {
		err := downloadStream.Close()
		if err != nil {

		}
	}(downloadStream)

	data, err := ioutil.ReadAll(downloadStream)
	if err != nil {
		return "", err
	}

	return string(data), nil
}
