package aliyun

import (
	"fmt"
	"testing"

	. "github.com/stretchr/testify/require"
)

func TestOssClientUploadLocal(t *testing.T) {
	var (
		link string
		err  error
	)

	//
	fmt.Println("~~~ BucketUrl:", testOssClient.BucketUrl())

	//
	link, err = testOssClient.UploadLocal("examples/hello.txt", "test/oss/hello1.txt")
	NoError(t, err)
	fmt.Println(">>> OssUploadLocal:", link)

	//
	link, err = testOssClient.UploadFromUrl(
		"http://127.0.0.1:8000/hello.txt",
		"test/oss/hello2.txt",
	)
	NoError(t, err)
	fmt.Println(">>> UploadFromUrl:", link)

	//
	link, err = testOssClient.UploadDir("examples", "test/oss/dir-target")
	NoError(t, err)
	fmt.Printf(">>> UploadDir: %s\n", link)
}

func TestOssMisc(t *testing.T) {
	var (
		code string
		err  error
	)

	// Copy
	code, err = testOssClient.Copy("test/oss/hello1.txt", "test/oss/hello1.copy.txt")
	if code != "" {
		fmt.Println("!!! Copy Error Code:", code)
	}
	NoError(t, err)

	code, err = testOssClient.Copy("test/oss/NotExists.txt", "test/oss/NotExists2.txt")
	Error(t, err)
	Equal(t, CODE_NoSuchKey, code)

	// Delete
	err = testOssClient.Delete("test/oss/hello1.copy.txt")
	NoError(t, err)

	// Get
	code, err = testOssClient.Get("test/oss/hello1.txt", "/tmp/hello1.txt")
	NoError(t, err)
	Equal(t, "", code)

	code, err = testOssClient.Get("test/oss/NotExists", "")
	Error(t, err)
	Equal(t, CODE_NoSuchKey, code)
}

func TestSts(t *testing.T) {
	var (
		link   string
		err    error
		result *StsResult
	)

	result, err = testOssClient.GetSts("xxxxxxxx")
	NoError(t, err)
	fmt.Println(">>> GetSts:", result)

	link, err = testOssClient.StsUploadLocal(
		"xxxxxxxx", "examples/hello.txt", "test/sts/hello1.txt",
	)
	NoError(t, err)
	fmt.Println(">>> StsUploadLocal:", link)

	link, err = testOssClient.StsUploadLocal(
		"xxxxxxxx", "examples/hello.txt", "aaaa/sts/hello1.txt",
	)
	Error(t, err)
}
