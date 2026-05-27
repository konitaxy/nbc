package upload

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"strconv"
	"strings"
	"time"

	"gitlab.com/ucard/config"
	"gitlab.com/ucard/global"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"go.uber.org/zap"
)

type AwsS3 struct{}

//@author: [WqyJh](https://github.com/WqyJh)
//@object: *AwsS3
//@function: UploadFile
//@description: Upload file to Aws S3 using aws-sdk-go. See https://docs.aws.amazon.com/sdk-for-go/v1/developer-guide/s3-example-basic-bucket-operations.html#s3-examples-bucket-ops-upload-file-to-bucket
//@param: file *multipart.FileHeader
//@return: string, string, error

func (*AwsS3) UploadFile(file *multipart.FileHeader) (string, string, error) {
	session := newSession()
	uploader := s3manager.NewUploader(session)
	uploader.UploadWithIterator(aws.BackgroundContext(), &s3manager.UploadObjectsIterator{
		Objects: []s3manager.BatchUploadObject{
			{
				Object: &s3manager.UploadInput{},
			},
		},
	})
	fileKey := fmt.Sprintf("%d%s", time.Now().Unix(), file.Filename)
	filename := global.GVA_CONFIG.AwsS3.PathPrefix + "/" + fileKey
	f, openError := file.Open()
	if openError != nil {
		global.GVA_LOG.Error("function file.Open() Filed", zap.Any("err", openError.Error()))
		return "", "", errors.New("function file.Open() Filed, err:" + openError.Error())
	}
	defer f.Close() // 创建文件 defer 关闭

	_, err := uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(global.GVA_CONFIG.AwsS3.Bucket),
		Key:    aws.String(filename),
		Body:   f,
	})
	if err != nil {
		global.GVA_LOG.Error("function uploader.Upload() Filed", zap.Any("err", err.Error()))
		return "", "", err
	}

	return global.GVA_CONFIG.AwsS3.BaseURL + "/" + filename, fileKey, nil
}

// 获取s3的临时链接,有效时间3分钟
func (*AwsS3) GetObjectTempRequest(key string) (string, error) {
	session := newSession()
	s33 := s3.New(session)

	req, _ := s33.GetObjectRequest(&s3.GetObjectInput{
		Bucket: aws.String(global.GVA_CONFIG.AwsS3.Bucket),
		Key:    aws.String(key),
	})

	return req.Presign(3 * time.Minute)
}

func (*AwsS3) UploadFileWitName(file *multipart.FileHeader, name, contentType string) (string, string, error) {
	session := newSession()
	uploader := s3manager.NewUploader(session)

	fileKey := name
	filename := global.GVA_CONFIG.AwsS3.PathPrefix + "/" + fileKey
	f, openError := file.Open()
	if openError != nil {
		global.GVA_LOG.Error("function file.Open() Filed", zap.Any("err", openError.Error()))
		return "", "", errors.New("function file.Open() Filed, err:" + openError.Error())
	}
	defer f.Close() // 创建文件 defer 关闭

	_, err := uploader.Upload(&s3manager.UploadInput{
		Bucket:      aws.String(global.GVA_CONFIG.AwsS3.Bucket),
		Key:         aws.String(filename),
		Body:        f,
		ContentType: aws.String(contentType),
	})
	if err != nil {
		global.GVA_LOG.Error("function uploader.Upload() Filed", zap.Any("err", err.Error()))
		return "", "", err
	}

	return "https://" + global.GVA_CONFIG.AwsS3.Bucket + "." + global.GVA_CONFIG.AwsS3.BaseURL + "/" + filename, fileKey, nil
}

func (*AwsS3) UploadFileWitNameWithoutPrefix(file *multipart.FileHeader, name, contentType string) (string, string, error) {
	session := newSession()
	uploader := s3manager.NewUploader(session)

	fileKey := name
	filename := fileKey
	f, openError := file.Open()
	if openError != nil {
		global.GVA_LOG.Error("function file.Open() Filed", zap.Any("err", openError.Error()))
		return "", "", errors.New("function file.Open() Filed, err:" + openError.Error())
	}
	defer f.Close() // 创建文件 defer 关闭

	_, err := uploader.Upload(&s3manager.UploadInput{
		Bucket:      aws.String(global.GVA_CONFIG.AwsS3.Bucket),
		Key:         aws.String(filename),
		Body:        f,
		ContentType: aws.String(contentType),
	})
	if err != nil {
		global.GVA_LOG.Error("function uploader.Upload() Filed", zap.Any("err", err.Error()))
		return "", "", err
	}

	return "https://" + global.GVA_CONFIG.AwsS3.Bucket + "." + global.GVA_CONFIG.AwsS3.BaseURL + "/" + filename, fileKey, nil
}

func (*AwsS3) CopyObject(source, target string) error {

	session := newSession()
	svc := s3.New(session)
	_, err := svc.CopyObject(&s3.CopyObjectInput{
		Bucket:     aws.String(global.GVA_CONFIG.AwsS3.Bucket),
		Key:        aws.String(target),
		CopySource: aws.String(source),
	})
	return err

}

func (*AwsS3) BatchUploadObject(name string, bs []bytes.Buffer) error {

	session := newSession()
	uploader := s3manager.NewUploader(session)
	filename := name
	var objects []s3manager.BatchUploadObject
	for i, b := range bs {
		key := strings.ReplaceAll(filename, ".", "_"+strconv.Itoa(i)+".")
		objects = append(objects, s3manager.BatchUploadObject{

			Object: &s3manager.UploadInput{
				Body:        &b,
				Bucket:      aws.String(global.GVA_CONFIG.AwsS3.Bucket),
				Key:         aws.String(key),
				ContentType: aws.String("image/png"),
			},
		})
	}
	return uploader.UploadWithIterator(context.Background(), &s3manager.UploadObjectsIterator{
		Objects: objects,
	})
}

func (*AwsS3) UploadBytesWitName(buff []byte, name, contentType string) (string, string, error) {
	session := newSession()
	uploader := s3manager.NewUploader(session)
	filename := global.GVA_CONFIG.AwsS3.PathPrefix + "/" + name
	_, err := uploader.Upload(&s3manager.UploadInput{
		Bucket:      aws.String(global.GVA_CONFIG.AwsS3.Bucket),
		Key:         aws.String(filename),
		Body:        bytes.NewReader(buff),
		ContentType: aws.String(contentType),
	})
	if err != nil {
		global.GVA_LOG.Error("function uploader.Upload() Filed", zap.Any("err", err.Error()))
		return "", "", err
	}
	return "https://" + global.GVA_CONFIG.AwsS3.Bucket + "." + global.GVA_CONFIG.AwsS3.BaseURL + "/" + filename, name, nil
}
func (*AwsS3) UploadTmpFile(buff []byte, name, contentType string) (string, string, error) {
	session := newSession()
	uploader := s3manager.NewUploader(session)
	filename := "tmp/" + name

	_, err := uploader.Upload(&s3manager.UploadInput{
		Bucket:      aws.String(global.GVA_CONFIG.AwsS3.Bucket),
		Key:         aws.String(filename),
		Body:        bytes.NewReader(buff),
		ContentType: aws.String(contentType),
	})
	if err != nil {
		global.GVA_LOG.Error("function uploader.Upload() Filed", zap.Any("err", err.Error()))
		return "", "", err
	}
	return "https://" + global.GVA_CONFIG.AwsS3.Bucket + "." + global.GVA_CONFIG.AwsS3.BaseURL + "/" + filename, name, nil
}

func (*AwsS3) UploadFile2(path string, name string) (string, string, error) {
	session := newSession()
	uploader := s3manager.NewUploader(session)

	fileKey := fmt.Sprintf("%d%s", time.Now().Unix(), name)
	filename := global.GVA_CONFIG.AwsS3.PathPrefix + "/" + fileKey
	reader, _ := os.OpenFile(path, os.O_RDONLY, 0666)
	defer reader.Close()
	_, err := uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(global.GVA_CONFIG.AwsS3.Bucket),
		Key:    aws.String(filename),
		Body:   reader,
	})
	if err != nil {
		// global.GVA_LOG.Error("function uploader.Upload() Filed", zap.Any("err", err.Error()))
		return "", "", err
	}

	return global.GVA_CONFIG.AwsS3.BaseURL + "/" + filename, fileKey, nil
}

//@author: [WqyJh](https://github.com/WqyJh)
//@object: *AwsS3
//@function: DeleteFile
//@description: Delete file from Aws S3 using aws-sdk-go. See https://docs.aws.amazon.com/sdk-for-go/v1/developer-guide/s3-example-basic-bucket-operations.html#s3-examples-bucket-ops-delete-bucket-item
//@param: file *multipart.FileHeader
//@return: string, string, error

func (*AwsS3) DeleteFile(key string) error {
	session := newSession()
	svc := s3.New(session)
	filename := key
	if strings.HasPrefix(filename, global.GVA_CONFIG.AwsS3.PathPrefix) {
		filename = strings.Join([]string{global.GVA_CONFIG.AwsS3.PathPrefix, filename}, "/")
	}
	bucket := global.GVA_CONFIG.AwsS3.Bucket

	_, err := svc.DeleteObject(&s3.DeleteObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(filename),
	})
	if err != nil {
		global.GVA_LOG.Error("function svc.DeleteObject() Filed", zap.Any("err", err.Error()))
		return errors.New("function svc.DeleteObject() Filed, err:" + err.Error())
	}
	_ = svc.WaitUntilObjectNotExists(&s3.HeadObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(filename),
	})

	return nil
}

func (*AwsS3) GetFile(key string) ([]byte, error) {
	session := newSession()
	svc := s3.New(session)

	bucket := global.GVA_CONFIG.AwsS3.Bucket
	output, err := svc.GetObject(&s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	})

	if err != nil {
		global.GVA_LOG.Error("function svc.GetFile() Filed", zap.Any("err", err.Error()))
		return nil, errors.New("function svc.GetFile() Filed, err:" + err.Error())
	}
	return io.ReadAll(output.Body)
}

func (*AwsS3) GetFileWithCfg(cfg config.AwsS3, key string) ([]byte, error) {
	session := newSessionWithCfg(cfg)
	svc := s3.New(session)

	output, err := svc.GetObject(&s3.GetObjectInput{
		Bucket: aws.String(cfg.Bucket),
		Key:    aws.String(key),
	})

	if err != nil {
		global.GVA_LOG.Error("function svc.GetObject() Filed", zap.Any("err", err.Error()))
		return nil, errors.New("function svc.GetObject() Filed, err:" + err.Error())
	}
	return io.ReadAll(output.Body)
}

func (*AwsS3) ListFiles() ([]string, error) {
	session := newSession()
	svc := s3.New(session)
	bucket := global.GVA_CONFIG.AwsS3.Bucket
	output, err := svc.ListObjectsV2(&s3.ListObjectsV2Input{
		Bucket: aws.String(bucket),
	})
	if err != nil {
		global.GVA_LOG.Error("function svc.DeleteObject() Filed", zap.Any("err", err.Error()))
		return nil, errors.New("function svc.DeleteObject() Filed, err:" + err.Error())
	}
	var keys []string
	for _, key := range output.Contents {

		keys = append(keys, *key.Key)
	}
	for output.NextContinuationToken != nil {
		output, err = svc.ListObjectsV2(&s3.ListObjectsV2Input{
			Bucket: aws.String(bucket),
		})
		if err != nil {
			global.GVA_LOG.Error("function svc.DeleteObject() Filed", zap.Any("err", err.Error()))
			return nil, errors.New("function svc.DeleteObject() Filed, err:" + err.Error())
		}
		for _, key := range output.Contents {
			keys = append(keys, *key.Key)
		}
	}
	return keys, err
}

func (*AwsS3) GetPutS3SignUrl(cfg config.AwsS3, fileName string) (string, error) {
	session := newSessionWithCfg(cfg)

	uploader := s3.New(session)
	key := fmt.Sprintf("%s/%s", cfg.PathPrefix, fileName)
	req, _ := uploader.PutObjectRequest(&s3.PutObjectInput{
		Bucket: aws.String(cfg.Bucket),
		Key:    aws.String(key),
	})
	return req.Presign(time.Minute * 3)
}

// newSession Create S3 session
func newSession() *session.Session {
	sess, _ := session.NewSession(&aws.Config{
		Region:           aws.String(global.GVA_CONFIG.AwsS3.Region),
		Endpoint:         aws.String(global.GVA_CONFIG.AwsS3.Endpoint), //minio在这里设置地址,可以兼容
		S3ForcePathStyle: aws.Bool(global.GVA_CONFIG.AwsS3.S3ForcePathStyle),
		DisableSSL:       aws.Bool(global.GVA_CONFIG.AwsS3.DisableSSL),
		Credentials: credentials.NewStaticCredentials(
			global.GVA_CONFIG.AwsS3.SecretID,
			global.GVA_CONFIG.AwsS3.SecretKey,
			"",
		),
	})
	return sess
}

func newSessionWithCfg(cfg config.AwsS3) *session.Session {
	sess, _ := session.NewSession(&aws.Config{
		Region:           aws.String(cfg.Region),
		Endpoint:         aws.String(cfg.Endpoint), //minio在这里设置地址,可以兼容
		S3ForcePathStyle: aws.Bool(cfg.S3ForcePathStyle),
		DisableSSL:       aws.Bool(cfg.DisableSSL),
		Credentials: credentials.NewStaticCredentials(
			cfg.SecretID,
			cfg.SecretKey,
			"",
		),
	})
	return sess
}
