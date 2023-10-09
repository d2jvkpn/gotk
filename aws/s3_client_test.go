package aws

import (
	"fmt"
	"testing"

	. "github.com/stretchr/testify/require"
)

func TestS3ClientUpload(t *testing.T) {
	var (
		link string
		err  error
	)

	//
	fmt.Println("~~~ BucketUrl:", testS3Client.BucketUrl())

	//
	link, err = testS3Client.UploadLocal(testContext, "examples/hello.txt", "/test/hello.txt")
	NoError(t, err)
	fmt.Println(">>> UploadLocal:", link)

	link, err = testS3Client.UploadLocal(testContext, "examples/hello.txt", "/not_allowed/hello.txt")
	NoError(t, err)

	// target: /test/b/ -> create directory only
	link, err = testS3Client.UploadLocal(testContext, "examples/hello.txt", "///")
	Error(t, err)

	link, err = testS3Client.UploadDir(testContext, "examples/", "test/aaaa_dir")
	NoError(t, err)
	fmt.Println(">>> UploadDir:", link)
}

func TestS3ClientSts(t *testing.T) {
	var (
		link   string
		err    error
		result *StsResult
	)

	result, err = testS3Client.GetSts(testContext, "xxxxxxxx")
	NoError(t, err)
	fmt.Println(">>> StsResult:", result)

	link, err = testS3Client.StsUploadLocal(
		testContext, "test", "examples/hello.txt", "/test/hello.txt",
	)
	NoError(t, err)
	fmt.Println(">>> StsUploadLocal:", link)

	link, err = testS3Client.StsUploadLocal(
		testContext, "test", "examples/hello.txt", "/not_allowed/hello.txt",
	)
	Error(t, err)
}

func TestS3ClientMisc(t *testing.T) {
	var (
		code string
		err  error
	)

	// Copy
	_, err = testS3Client.Copy(testContext, "/test/hello.txt", "/test/hello_2022.txt")
	NoError(t, err)

	code, err = testS3Client.Copy(testContext, "test/NotExists", "test/NotExists2")
	Error(t, err)
	fmt.Println("~~~ Copy Error Code:", code)
	Equal(t, CODE_NoSuchKey, code)

	code, err = testS3Client.Copy(testContext, "/test/NotExits", "/test/hello_2022.txt")
	fmt.Println("~~~ Copy Error Code:", code)
	Error(t, err)

	// Delete
	err = testS3Client.Delete(testContext, "/test/hello_2022.txt")
	NoError(t, err)

	err = testS3Client.Delete(testContext, "/test/NotExits") // ??
	NoError(t, err)

	// Get
	code, err = testS3Client.Get(testContext, "/test/hello.txt", "/tmp/hello1.txt")
	NoError(t, err)
	Equal(t, "", code)

	code, err = testS3Client.Get(testContext, "/test/NotExits", "/test/hello_2022.txt")
	Error(t, err)
	Equal(t, CODE_NoSuchKey, code)
}
