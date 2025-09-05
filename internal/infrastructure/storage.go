package infrastructure

import (
	"FitByte/configs"
	"FitByte/pkg/log"
	"context"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

func InitMinioStorage(appConfig configs.Config) *minio.Client {
	endpoint := appConfig.Minio.Endpoint
	accessKeyID := appConfig.Minio.AccessKeyID
	secretAccessKey := appConfig.Minio.SecretAccessKey
	useSSL := appConfig.Minio.UseSSL

	// Initialize minio client object.
	minioClient, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
		Secure: useSSL,
	})
	if err != nil {
		log.Logger.Fatal().Err(err).Msg("minio init failed")
	}

	//Validate Bucket Exists
	exists, err := minioClient.BucketExists(context.Background(), appConfig.Minio.Bucket)
	if err != nil {
		log.Logger.Fatal().Err(err).Msg("minio bucket exists check failed")
	}

	if !exists {
		err = minioClient.MakeBucket(context.Background(), appConfig.Minio.Bucket, minio.MakeBucketOptions{})
		if err != nil {
			log.Logger.Fatal().Err(err).Msg("minio bucket creation failed")
		} else {
			log.Logger.Info().Str("bucket", appConfig.Minio.Bucket).Msg("minio bucket created")
		}
	}
	log.Logger.Info().Msg("minio init success")
	return minioClient
}
