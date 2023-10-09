package aws

import (
	"context"
	"fmt"
	"io"
	"io/fs"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/aws/aws-sdk-go/service/sts"
)

// TODO: s3manager.BatchUploadObject, s3manager.UploadObjectsIterator

type S3Client struct {
	sess *session.Session
	*s3manager.Uploader
	s3     *s3.S3
	config S3Config
	bucket *string
}

type StsResult struct {
	AccessKeyId     string `json:"accessKeyId"`
	SecretAccessKey string `json:"secretAccessKey"`
	SessionToken    string `json:"sessionToken"`
	Expiration      string `json:"expiration"`
	Region          string `json:"region"`
	Bucket          string `json:"bucket"`
}

func NewS3Client(fp, field string) (client *S3Client, err error) {
	var config S3Config
	if config, err = NewS3Config(fp, field); err != nil {
		return nil, err
	}

	client = &S3Client{config: config}

	client.sess, err = session.NewSession(
		&aws.Config{
			Region: aws.String(client.config.Region),
			Credentials: credentials.NewStaticCredentials(
				client.config.AccessKeyId,
				client.config.SecretAccessKey,
				"", // a token will be created when the session it's used.
			),
		},
	)
	if err != nil {
		return nil, err
	}

	client.Uploader = s3manager.NewUploader(client.sess)
	client.s3 = s3.New(client.sess)
	client.bucket = aws.String(client.config.Bucket)
	return client, nil
}

func (client *S3Client) BucketUrl() string {
	return client.config.bucketUrl
}

func (client *S3Client) Upload(ctx context.Context, source, target string) (
	link string, err error) {
	if strings.HasPrefix(source, "http://") || strings.HasPrefix(source, "https://") {
		return client.UploadFromUrl(ctx, source, target)
	}

	var info fs.FileInfo

	if info, err = os.Stat(source); err != nil {
		return "", err
	}
	//	fsize = info.Size()

	if info.IsDir() {
		return client.UploadDir(ctx, source, target)
	}
	return client.UploadLocal(ctx, source, target)
}

func (client *S3Client) UploadFromUrl(ctx context.Context, source, target string) (
	link string, err error) {

	var (
		httpRes *http.Response
		input   *s3manager.UploadInput
		outoput *s3manager.UploadOutput
	)

	if httpRes, err = http.Get(source); err != nil {
		return "", err
	}
	// ?? no Content-Length here
	// fsize, _ = strconv.ParseInt(httpRes.Header.Get("Content-Length"), 10, 64)
	if httpRes.StatusCode != http.StatusOK {
		return "", fmt.Errorf("http status: %d, %s", httpRes.StatusCode, httpRes.Status)
	}
	defer httpRes.Body.Close()

	input = &s3manager.UploadInput{
		Bucket: client.bucket,
		Key:    aws.String(target),
		Body:   httpRes.Body,
	}

	if strings.HasSuffix(target, ".html") {
		input.ContentType = aws.String("text/html")
	}

	outoput, err = client.UploadWithContext(ctx, input)
	if err != nil {
		return "", err
	}
	return outoput.Location, nil
}

func (client *S3Client) UploadLocal(ctx context.Context, source, target string) (
	link string, err error) {
	var (
		file    *os.File
		input   *s3manager.UploadInput
		outoput *s3manager.UploadOutput
	)

	if file, err = os.Open(source); err != nil {
		return "", err
	}
	defer file.Close()

	input = &s3manager.UploadInput{
		Bucket: client.bucket,
		Key:    aws.String(target),
		Body:   file,
	}

	if strings.HasSuffix(target, ".html") {
		input.ContentType = aws.String("text/html")
	}

	outoput, err = client.UploadWithContext(ctx, input)
	if err != nil {
		return "", err
	}
	return outoput.Location, nil
}

func (client *S3Client) UploadDir(ctx context.Context, source, target string) (
	link string, err error) {

	source = strings.TrimRight(source, "/")
	conc, once := uint(runtime.NumCPU()), new(sync.Once)
	n, done := uint(0), make(chan bool, conc)

	filepath.Walk(source, func(fp string, info fs.FileInfo, e1 error) error {
		if e1 != nil {
			return e1
		}
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		// fmt.Println(fp, info.Size())

		if n++; n > conc {
			<-done
			n--
		}

		go func() {
			t := filepath.Join(target, strings.Replace(fp, source, "", 1))
			// fmt.Println("~~~", time.Now().Format(time.RFC3339), fp)
			if _, e2 := client.UploadLocal(ctx, fp, t); e2 != nil {
				once.Do(func() {
					err = e2
				})
			}
			done <- true
		}()

		return nil
	})

	for ; n > 0; n-- {
		<-done
	}

	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%s/%s/", client.config.bucketUrl, strings.Trim(target, "/")), nil
}

func (client *S3Client) GetSts(ctx context.Context, sessionName string) (
	result *StsResult, err error) {
	var (
		svc    *sts.STS
		output *sts.AssumeRoleOutput
		input  *sts.AssumeRoleInput
	)

	if client.config.AssumeArn == "" {
		return nil, fmt.Errorf("assume_arn is unset")
	}

	svc = sts.New(client.sess)
	input = &sts.AssumeRoleInput{
		RoleArn:         &client.config.AssumeArn,
		RoleSessionName: &sessionName,
	}
	if client.config.DurationSeconds > 0 {
		input.DurationSeconds = &client.config.DurationSeconds
	}

	if output, err = svc.AssumeRoleWithContext(ctx, input); err != nil {
		return nil, err
	}

	result = &StsResult{
		AccessKeyId:     *output.Credentials.AccessKeyId,
		SecretAccessKey: *output.Credentials.SecretAccessKey,
		SessionToken:    *output.Credentials.SessionToken,
		Expiration:      output.Credentials.Expiration.Format(time.RFC3339),
		Region:          client.config.Region,
		Bucket:          client.config.Bucket,
	}

	return result, nil
}

func (client *S3Client) StsUploadLocal(ctx context.Context, sessionName, source, target string) (
	link string, err error) {
	var (
		file     *os.File
		sess     *session.Session
		result   *StsResult
		uploader *s3manager.Uploader
		output   *s3manager.UploadOutput
	)

	if result, err = client.GetSts(ctx, sessionName); err != nil {
		return "", err
	}

	sess, err = session.NewSession(&aws.Config{
		Region: aws.String(client.config.Region),
		Credentials: credentials.NewStaticCredentials(
			result.AccessKeyId,
			result.SecretAccessKey,
			result.SessionToken,
		),
	})
	if err != nil {
		return "", nil
	}
	uploader = s3manager.NewUploader(sess)

	if file, err = os.Open(source); err != nil {
		return "", err
	}
	defer file.Close()

	output, err = uploader.Upload(&s3manager.UploadInput{
		Bucket: client.bucket,
		Key:    aws.String(target),
		Body:   file,
	})
	if err != nil {
		return "", err
	}
	return output.Location, nil
}

func (client *S3Client) Copy(ctx context.Context, source, target string, opts ...request.Option) (
	code string, err error) {
	copySource := aws.String(url.PathEscape(
		client.config.Bucket + "/" + strings.Trim(source, "/"),
	))

	var key *string

	key = aws.String(target)
	_, err = client.s3.CopyObjectWithContext( // *s3.CopyObjectOutput
		ctx,
		&s3.CopyObjectInput{
			Bucket:     client.bucket,
			CopySource: copySource,
			Key:        key,
		},
		opts...,
	)

	if err != nil {
		err2, _ := err.(awserr.RequestFailure)
		code = err2.Code() // !! code is AccessDenied if key not exists
		return code, err
	}

	// err = client.s3.WaitUntilObjectExists(&s3.HeadObjectInput{Bucket: client.bucket, Key: key})
	return "", nil
}

func (client *S3Client) Delete(ctx context.Context, source string, opts ...request.Option) (
	err error) {
	key := aws.String(source)

	_, err = client.s3.DeleteObjectWithContext(
		ctx,
		&s3.DeleteObjectInput{
			Bucket: client.bucket,
			Key:    key,
		},
		opts...,
	)

	return
}

func (client *S3Client) Get(ctx context.Context, source, target string, opts ...request.Option) (
	code string, err error) {
	key := aws.String(source)

	var (
		output *s3.GetObjectOutput
		file   *os.File
	)

	output, err = client.s3.GetObjectWithContext(
		ctx,
		&s3.GetObjectInput{
			Bucket: client.bucket,
			Key:    key,
		},
		opts...,
	)
	if err != nil {
		// fmt.Printf("??? %[1]T, %[1]v\n", err)
		err2, _ := err.(awserr.RequestFailure)
		code = err2.Code()
		return code, err
	}
	defer output.Body.Close()

	if target == "" {
		target = filepath.Base(source)
	}

	if file, err = os.Create(target); err != nil {
		return "", err
	}
	defer file.Close() // _ = file.Close()

	_, err = io.Copy(file, output.Body)
	return "", err
}
