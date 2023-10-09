package aliyun

import (
	"fmt"
	"io"
	"io/fs"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	// "strconv"
	"strings"
	"sync"
	// "time"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/sts"
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
)

type OssClient struct {
	// *oss.Client
	config OssConfig
	bucket *oss.Bucket
	sts    *sts.Client
}

type StsResult struct {
	AccessKeyId     string `json:"accessKeyId"`
	AccessKeySecret string `json:"accessKeySecret"`
	SecurityToken   string `json:"securityToken"`
	Expiration      string `json:"expiration"`
	Region          string `json:"region"`
	Bucket          string `json:"bucket"`
}

func NewOssClient(fp, field string) (client *OssClient, err error) {
	var config OssConfig
	if config, err = NewOssConfig(fp, field); err != nil {
		return nil, err
	}

	return ossClientFromConfig(config)
}

func ossClientFromConfig(config OssConfig) (client *OssClient, err error) {
	var cli *oss.Client

	client = &OssClient{config: config}
	cli, err = oss.New(
		regionUrl(config.Region), config.AccessKeyId, config.AccessKeySecret,
	)
	if err != nil {
		return nil, err
	}

	if client.bucket, err = cli.Bucket(client.config.Bucket); err != nil {
		return nil, err
	}

	if config.RoleArn != "" {
		client.sts, err = sts.NewClientWithAccessKey(
			config.Region, config.AccessKeyId, config.AccessKeySecret,
		)
		if err != nil {
			return nil, err
		}
	}

	return
}

func (client *OssClient) BucketUrl() string {
	return client.config.bucketUrl
}

func (client *OssClient) GetSite() string {
	return client.config.Site
}

// OSS, oss.ForbidOverWrite(forbidWrite bool) oss.Option
func (client *OssClient) UploadLocal(source, target string, opts ...oss.Option) (
	link string, err error) {

	// urlpath = strings.Trim(fmt.Sprintf("%s/%s", strings.Trim(config.Path, "/"), subpath), "/")
	// additional slash in middle will cause error: The specified object is not valid.
	// subpath with slash tail will create a new directory
	if target, err = ValidUrlFilepath(target); err != nil {
		return "", err
	}

	if err = client.bucket.PutObjectFromFile(target, source, opts...); err != nil {
		return "", err
	}

	return client.config.Url(target), nil
}

func (client *OssClient) UploadFromUrl(source string, target string, opts ...oss.Option) (
	link string, err error) {
	var (
		httpRes *http.Response
		ossReq  *oss.PutObjectRequest
		// ossRes  *oss.Response
	)

	if target, err = ValidUrlFilepath(target); err != nil {
		return "", err
	}

	//	if httpRes, err = http.Head(source); err != nil {
	//		return "", 0, err
	//	}
	//	fsize, _ = strconv.ParseInt(httpRes.Header.Get("Content-Length"), 10, 64)

	if httpRes, err = http.Get(source); err != nil {
		return "", err
	}
	// ?? no Content-Length here
	// fsize, _ = strconv.ParseInt(httpRes.Header.Get("Content-Length"), 10, 64)
	if httpRes.StatusCode != http.StatusOK {
		return "", fmt.Errorf("http status: %d, %s", httpRes.StatusCode, httpRes.Status)
	}
	defer httpRes.Body.Close()

	ossReq = &oss.PutObjectRequest{
		ObjectKey: target,
		Reader:    httpRes.Body,
	}

	if _, err = client.bucket.DoPutObject(ossReq, opts); err != nil {
		return "", err
	}
	// fmt.Printf("%#v\n", ossRes)

	return client.config.Url(target), nil
}

func (client *OssClient) PutObject(reader io.Reader, target string, opts ...oss.Option) (
	link string, err error) {
	var (
		ossReq *oss.PutObjectRequest
		// ossRes  *oss.Response
	)

	if target, err = ValidUrlFilepath(target); err != nil {
		return "", err
	}

	ossReq = &oss.PutObjectRequest{
		ObjectKey: target,
		Reader:    reader,
	}

	if _, err = client.bucket.DoPutObject(ossReq, opts); err != nil {
		return "", err
	}
	// fmt.Printf("%#v\n", ossRes)
	return client.config.Url(target), nil
}

func (client *OssClient) Upload(source string, target string) (link string, err error) {
	if strings.HasPrefix(source, "http://") || strings.HasPrefix(source, "https://") {
		return client.UploadFromUrl(source, target, oss.ForbidOverWrite(false))
	}

	var info fs.FileInfo

	if info, err = os.Stat(source); err != nil {
		return "", err
	}
	//	fsize = info.Size()

	if info.IsDir() {
		return client.UploadDir(source, target)
	}
	return client.UploadLocal(source, target, oss.ForbidOverWrite(false))
}

func (client *OssClient) UploadDir(source string, target string) (link string, err error) {
	conc, once := 2*uint(runtime.NumCPU()), new(sync.Once)
	n, done := uint(0), make(chan bool, conc)

	// err = filepath.Walk(...)
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
			_, e2 := client.UploadLocal(fp, t, oss.ForbidOverWrite(false))
			if e2 != nil {
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

	// link is a dir path, try to access html like dir + "index/html"
	return client.config.Url(strings.Trim(target, "/")), nil
}

func (client *OssClient) Copy(source, target string, opts ...oss.Option) (
	code string, err error) {
	var serviceError oss.ServiceError

	// result oss.CopyObjectResult
	if _, err = client.bucket.CopyObject(source, target, opts...); err != nil {
		serviceError, _ = err.(oss.ServiceError)
		return serviceError.Code, err
	}

	return "", nil
}

func (client *OssClient) Delete(source string, opts ...oss.Option) (err error) {
	// result oss.CopyObjectResult
	if err = client.bucket.DeleteObject(source, opts...); err != nil {
		return err
	}

	return nil
}

func (client *OssClient) Get(source, target string, opts ...oss.Option) (code string, err error) {
	source = strings.Trim(source, "/")
	if target == "" {
		target = filepath.Base(source)
	}

	err = client.bucket.GetObjectToFile(source, target, opts...)
	if err != nil {
		serviceError, _ := err.(oss.ServiceError)
		code = serviceError.Code
	}

	return code, err
}

// STS
func (client *OssClient) GetSts(sessionName string) (result *StsResult, err error) {
	var (
		request  *sts.AssumeRoleRequest
		response *sts.AssumeRoleResponse
	)

	if client.config.RoleArn == "" {
		return nil, fmt.Errorf("role_arn is unset")
	}

	request = sts.CreateAssumeRoleRequest()
	request.Scheme = "https"
	request.RoleArn = client.config.RoleArn
	request.RoleSessionName = sessionName
	if client.config.DurationSeconds > 0 {
		request.DurationSeconds = requests.NewInteger(client.config.DurationSeconds)
	}

	if response, err = client.sts.AssumeRole(request); err != nil {
		return nil, err
	}
	// fmt.Printf("~~ %v\n", response)

	result = &StsResult{
		AccessKeyId:     response.Credentials.AccessKeyId,
		AccessKeySecret: response.Credentials.AccessKeySecret,
		SecurityToken:   response.Credentials.SecurityToken,
		Expiration:      response.Credentials.Expiration,
		Region:          client.config.Region,
		Bucket:          client.config.Bucket,
	}

	return result, nil
}

// test Sts client
func (client *OssClient) StsUploadLocal(sessionName string, source, target string,
	options ...oss.Option) (link string, err error) {

	var (
		result *StsResult
		bucket *oss.Bucket
		cli    *oss.Client
	)

	target = strings.Trim(target, "/")

	if result, err = client.GetSts(sessionName); err != nil {
		return "", err
	}

	cli, err = oss.New(
		regionUrl(result.Region), result.AccessKeyId, result.AccessKeySecret,
		oss.SecurityToken(result.SecurityToken),
	)
	if err != nil {
		return "", err
	}

	if bucket, err = cli.Bucket(result.Bucket); err != nil {
		return "", err
	}

	if target, err = ValidUrlFilepath(target); err != nil {
		return "", err
	}
	if err = bucket.PutObjectFromFile(target, source, options...); err != nil {
		return "", err
	}

	return client.config.bucketUrl + "/" + target, nil
}
