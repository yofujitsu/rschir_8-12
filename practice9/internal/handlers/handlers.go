package handlers

import (
	"context"
	"encoding/json"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/gridfs"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"practice9/internal/storage"
)

// FileListHandler Обработчик для получения списка файлов
func FileListHandler(w http.ResponseWriter, r *http.Request) {
	client, err := storage.ConnectToMongoDB()
	if err != nil {
		http.Error(w, "Failed to connect to MongoDB", http.StatusInternalServerError)
		return
	}
	defer func(client *mongo.Client, ctx context.Context) {
		err := client.Disconnect(ctx)
		if err != nil {

		}
	}(client, context.Background())

	db := client.Database("mydb")

	filesCollection := db.Collection("fs.files")

	filter := bson.D{}

	cursor, err := filesCollection.Find(context.Background(), filter)
	if err != nil {
		log.Fatal(err)
	}
	defer func(cursor *mongo.Cursor, ctx context.Context) {
		err := cursor.Close(ctx)
		if err != nil {

		}
	}(cursor, context.Background())

	var filesInfo []map[string]interface{}

	for cursor.Next(context.Background()) {
		var fileInfo bson.M
		if err := cursor.Decode(&fileInfo); err != nil {
			log.Fatal(err)
		}

		fileInfoMap := map[string]interface{}{
			"id":       fileInfo["_id"],
			"filename": fileInfo["filename"],
			"length":   fileInfo["length"],
		}
		filesInfo = append(filesInfo, fileInfoMap)
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	jsonData, err := json.MarshalIndent(filesInfo, "", "  ")
	if err != nil {
		http.Error(w, "Failed to marshal JSON", http.StatusInternalServerError)
		return
	}
	_, err = w.Write(jsonData)
	if err != nil {
		return
	}
}

// FileInfoHandler Обработчик для получения информации о файле по id
func FileInfoHandler(w http.ResponseWriter, r *http.Request) {
	client, err := storage.ConnectToMongoDB()
	if err != nil {
		http.Error(w, "Failed to connect to MongoDB", http.StatusInternalServerError)
		return
	}
	defer func(client *mongo.Client, ctx context.Context) {
		err := client.Disconnect(ctx)
		if err != nil {

		}
	}(client, context.Background())

	vars := mux.Vars(r)
	fileID, err := primitive.ObjectIDFromHex(vars["id"])
	if err != nil {
		http.Error(w, "Failed to connect to MongoDB", http.StatusInternalServerError)
		return
	}

	db := client.Database("mydb")

	filesCollection := db.Collection("fs.files")

	filter := bson.M{"_id": fileID}

	var fileInfo bson.M
	err = filesCollection.FindOne(context.Background(), filter).Decode(&fileInfo)
	if err != nil {
		http.Error(w, "File not found", http.StatusNotFound)
		return
	}

	fileInfoMap := map[string]interface{}{
		"id":       fileInfo["_id"],
		"filename": fileInfo["filename"],
		"length":   fileInfo["length"],
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	jsonData, err := json.MarshalIndent(fileInfoMap, "", "  ")
	if err != nil {
		http.Error(w, "Failed to marshal JSON", http.StatusInternalServerError)
		return
	}
	_, err = w.Write(jsonData)
	if err != nil {
		return
	}
}

// UploadFileHandler Обработчик для загрузки файла
func UploadFileHandler(w http.ResponseWriter, r *http.Request) {
	file, handler, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Failed to get file from request", http.StatusBadRequest)
		return
	}
	defer func(file multipart.File) {
		err := file.Close()
		if err != nil {

		}
	}(file)

	filename := handler.Filename

	client, err := storage.ConnectToMongoDB()
	if err != nil {
		println(1)
		http.Error(w, "Failed to connect to MongoDB", http.StatusInternalServerError)
		return
	}
	defer func(client *mongo.Client, ctx context.Context) {
		err := client.Disconnect(ctx)
		if err != nil {

		}
	}(client, context.Background())

	db := client.Database("mydb")

	fs, err := gridfs.NewBucket(db)
	if err != nil {
		http.Error(w, "Failed to create GridFS bucket", http.StatusInternalServerError)
		return
	}

	uploadStream, err := fs.OpenUploadStream(filename)
	if err != nil {
		http.Error(w, "Failed to open upload stream", http.StatusInternalServerError)
		return
	}
	defer func(uploadStream *gridfs.UploadStream) {
		err := uploadStream.Close()
		if err != nil {

		}
	}(uploadStream)

	_, err = io.Copy(uploadStream, file)
	if err != nil {
		http.Error(w, "Failed to copy file to upload stream", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	_, err = w.Write([]byte("File uploaded successfully"))
	if err != nil {
		return
	}
}

// FileTextHandler Обработчик для получения текста из файла по id
func FileTextHandler(w http.ResponseWriter, r *http.Request) {
	client, err := storage.ConnectToMongoDB()
	if err != nil {
		http.Error(w, "Failed to connect to MongoDB", http.StatusInternalServerError)
		return
	}
	defer func(client *mongo.Client, ctx context.Context) {
		err := client.Disconnect(ctx)
		if err != nil {

		}
	}(client, context.Background())

	vars := mux.Vars(r)
	fileID, err := primitive.ObjectIDFromHex(vars["id"])
	if err != nil {
		http.Error(w, "Failed to connect to MongoDB", http.StatusInternalServerError)
		return
	}

	db := client.Database("mydb")

	filesCollection := db.Collection("fs.files")

	filter := bson.M{"_id": fileID}

	var fileInfo bson.M
	err = filesCollection.FindOne(context.Background(), filter).Decode(&fileInfo)
	if err != nil {
		http.Error(w, "File not found", http.StatusNotFound)
		return
	}

	fs, err := gridfs.NewBucket(
		db,
	)
	content, err := storage.ReadGridFSFile(fs, fileInfo["_id"])
	if err != nil {
		log.Fatal(err)
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	jsonData, err := json.MarshalIndent(content, "", "  ")
	if err != nil {
		http.Error(w, "Failed to marshal JSON", http.StatusInternalServerError)
		return
	}
	_, err = w.Write(jsonData)
	if err != nil {
		return
	}
}

// DeleteFileHandler Обработчик для удаления файла по id
func DeleteFileHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	fileID, err := primitive.ObjectIDFromHex(vars["id"])
	if err != nil {
		http.Error(w, "Invalid file ID format", http.StatusBadRequest)
		return
	}

	client, err := storage.ConnectToMongoDB()
	if err != nil {
		http.Error(w, "Failed to connect to MongoDB", http.StatusInternalServerError)
		return
	}
	defer func(client *mongo.Client, ctx context.Context) {
		err := client.Disconnect(ctx)
		if err != nil {

		}
	}(client, context.Background())

	db := client.Database("mydb")
	fs, err := gridfs.NewBucket(db)
	if err != nil {
		http.Error(w, "Failed to create GridFS bucket", http.StatusInternalServerError)
		return
	}

	err = fs.Delete(fileID)
	if err != nil {
		http.Error(w, "Failed to delete file", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	_, err = w.Write([]byte("File deleted successfully"))
	if err != nil {
		return
	}
}

// UpdateFileHandler Обработчик для обновления файла по id
func UpdateFileHandler(w http.ResponseWriter, r *http.Request) {

	file, handler, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Failed to get file from request", http.StatusBadRequest)
		return
	}
	defer func(file multipart.File) {
		err := file.Close()
		if err != nil {

		}
	}(file)

	client, err := storage.ConnectToMongoDB()
	if err != nil {
		http.Error(w, "Failed to connect to MongoDB", http.StatusInternalServerError)
		return
	}
	defer func(client *mongo.Client, ctx context.Context) {
		err := client.Disconnect(ctx)
		if err != nil {

		}
	}(client, context.Background())

	vars := mux.Vars(r)
	fileID, err := primitive.ObjectIDFromHex(vars["id"])
	if err != nil {
		http.Error(w, "Invalid file ID format", http.StatusBadRequest)
		return
	}

	db := client.Database("mydb")

	fs, err := gridfs.NewBucket(db)
	if err != nil {
		http.Error(w, "Failed to create GridFS bucket", http.StatusInternalServerError)
		return
	}
	err = fs.Delete(fileID)
	if err != nil {
		http.Error(w, "Failed to delete the old file", http.StatusInternalServerError)
		return
	}

	newFileID := fileID
	newFilename := handler.Filename

	uploadStream, err := fs.OpenUploadStreamWithID(newFileID, newFilename)
	if err != nil {
		http.Error(w, "Failed to open upload stream", http.StatusInternalServerError)
		return
	}
	defer func(uploadStream *gridfs.UploadStream) {
		err := uploadStream.Close()
		if err != nil {

		}
	}(uploadStream)

	_, err = io.Copy(uploadStream, file)
	if err != nil {
		http.Error(w, "Failed to copy file to upload stream", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	_, err = w.Write([]byte("File updated successfully"))
	if err != nil {
		return
	}
}
