package main

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

var s3Client *s3.S3

func main() {

	// minio uses the host header to sign the request,
	//so we have to use the same host here and later replace it with the actual domain
	sess, err := session.NewSession(&aws.Config{
		Region:           aws.String("eu-central-1"),
		Endpoint:         aws.String("http://127.0.0.1:9000"),
		S3ForcePathStyle: aws.Bool(true),
		Credentials: credentials.NewStaticCredentials(
			"e9lCjzyC3B75NoAz",
			"Vq9jgjmI3tYXOXtesVwiXci6xIFsEh4O",
			"",
		),
	})

	if err != nil {
		panic(err)
	}

	s3Client = s3.New(sess)

	bucket := "test"
	key := "test.txt"
	content := "body"
	contentType := "text/plain"

	req, _ := s3Client.PutObjectRequest(&s3.PutObjectInput{
		Bucket:      aws.String(bucket),
		Key:         aws.String(key),
		ContentType: aws.String(contentType),
	})
	uploadUri, err := req.Presign(60 * time.Minute)

	if err != nil {
		panic(err)
	}

	uploadUri = strings.Replace(uploadUri, "http://127.0.0.1:9000", "https://minio.docker.localhost", 1)

	request, err := http.NewRequest(http.MethodPut, uploadUri, bytes.NewReader([]byte(content)))
	if err != nil {
		panic(err)
	}
	request.Header.Add("Content-Type", contentType)

	response, err := http.DefaultClient.Do(request)
	if err != nil {
		panic(err)
	}
	if response.StatusCode > 299 {

		defer response.Body.Close()
		b, err := io.ReadAll(response.Body)

		if err != nil {
			panic(err)
		}

		fmt.Println("error", string(b))
		return
	}

	req, _ = s3Client.GetObjectRequest(&s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	})

	downloadUri, err := req.Presign(60 * time.Minute)

	downloadUri = strings.Replace(downloadUri, "http://127.0.0.1:9000", "https://minio.docker.localhost", 1)

	request, err = http.NewRequest(http.MethodGet, downloadUri, nil)
	if err != nil {
		panic(err)
	}

	response, err = http.DefaultClient.Do(request)
	if err != nil {
		panic(err)
	}

	defer response.Body.Close()
	b, err := io.ReadAll(response.Body)

	if err != nil {
		panic(err)
	}

	fmt.Println("body", string(b))

}
