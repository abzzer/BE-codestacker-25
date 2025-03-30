package database

import (
	"context"
	"log"
	"os"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

var MinioClient *minio.Client

func ConnectMinIO() {
	endpoint := os.Getenv("MINIO_ENDPOINT")
	accessKey := os.Getenv("MINIO_ACCESS_KEY")
	secretKey := os.Getenv("MINIO_SECRET_KEY")

	client, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKey, secretKey, ""),
		Secure: false,
	})

	if err != nil {
		log.Fatalf("Failed to connect to MinIO: %v", err)
	}

	MinioClient = client
	log.Println("Connection to MinIO successful")

	ctx := context.Background()
	bucketName := "evidence-bucket"
	exists, errBuck := client.BucketExists(ctx, bucketName)
	if errBuck != nil {
		log.Fatalf("Couldn't connect to check for bucket: %v", errBuck)
	}

	if !exists {
		err = client.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{})
		if err != nil {
			log.Fatalf("Failed to create the backet: %v", err)
		}
		log.Println("Successfully created the bucket with the name:", bucketName)
	}
}
