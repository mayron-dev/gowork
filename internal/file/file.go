package file

import (
	"bytes"
	"context"
	"io"
	"time"

	s3Config "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
)

type IFileService interface {
	Upload(ctx context.Context, bucketName, key string, file io.Reader) error
	Download(ctx context.Context, bucketName, key string) ([]byte, error)
	MultipartUpload(ctx context.Context, bucketName, key string, file io.Reader, partSize int64) error
	GeneratePresignedUploadURL(ctx context.Context, bucketName, key string, duration time.Duration) (string, error)
	GeneratePresignedDownloadURL(ctx context.Context, bucketName, key string, duration time.Duration) (string, error)
}

type FileService struct {
	client *s3.Client
}

func (s *FileService) Upload(ctx context.Context, bucketName, key string, file io.Reader) error {
	_, err := s.client.PutObject(ctx, &s3.PutObjectInput{
		Bucket: &bucketName,
		Key:    &key,
		Body:   file,
	})
	return err
}

func (s *FileService) Download(ctx context.Context, bucketName, key string) ([]byte, error) {
	resp, err := s.client.GetObject(ctx, &s3.GetObjectInput{
		Bucket: &bucketName,
		Key:    &key,
	})
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	buf := new(bytes.Buffer)
	_, err = io.Copy(buf, resp.Body)
	return buf.Bytes(), err
}

func (s *FileService) MultipartUpload(ctx context.Context, bucketName, key string, file io.Reader, partSize int64) error {
	createResp, err := s.client.CreateMultipartUpload(ctx, &s3.CreateMultipartUploadInput{
		Bucket: &bucketName,
		Key:    &key,
	})
	if err != nil {
		return err
	}

	var completedParts []types.CompletedPart
	partNumber := int32(1)
	buf := make([]byte, partSize)

	for {
		n, err := file.Read(buf)
		if err != nil && err != io.EOF {
			return err
		}
		if n == 0 {
			break
		}

		partResp, err := s.client.UploadPart(ctx, &s3.UploadPartInput{
			Bucket:     &bucketName,
			Key:        &key,
			PartNumber: &partNumber,
			UploadId:   createResp.UploadId,
			Body:       bytes.NewReader(buf[:n]),
		})
		if err != nil {
			return err
		}

		completedParts = append(completedParts, types.CompletedPart{
			ETag:       partResp.ETag,
			PartNumber: &partNumber,
		})
		partNumber++
	}

	_, err = s.client.CompleteMultipartUpload(ctx, &s3.CompleteMultipartUploadInput{
		Bucket:   &bucketName,
		Key:      &key,
		UploadId: createResp.UploadId,
		MultipartUpload: &types.CompletedMultipartUpload{
			Parts: completedParts,
		},
	})
	return err
}

func (s *FileService) GeneratePresignedUploadURL(ctx context.Context, bucketName, key string, duration time.Duration) (string, error) {
	psClient := s3.NewPresignClient(s.client)
	presignedURL, err := psClient.PresignPutObject(ctx, &s3.PutObjectInput{
		Bucket: &bucketName,
		Key:    &key,
	}, s3.WithPresignExpires(duration))
	return presignedURL.URL, err
}

func (s *FileService) GeneratePresignedDownloadURL(ctx context.Context, bucketName, key string, duration time.Duration) (string, error) {
	psClient := s3.NewPresignClient(s.client)
	presignedURL, err := psClient.PresignGetObject(ctx, &s3.GetObjectInput{
		Bucket: &bucketName,
		Key:    &key,
	}, s3.WithPresignExpires(duration))
	return presignedURL.URL, err
}

func NewFileService(accessKeyId, secretAccessKey, endpoint, region string) (IFileService, error) {
	cfg, err := s3Config.LoadDefaultConfig(context.TODO(),
		s3Config.WithRegion(region),
		s3Config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(accessKeyId, secretAccessKey, "")),
		s3Config.WithBaseEndpoint(endpoint),
	)
	if err != nil {
		return nil, err
	}

	return &FileService{client: s3.NewFromConfig(cfg)}, nil
}
