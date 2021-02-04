package s3

import (
	"bytes"
	"encoding/base64"

	"github.com/aws/aws-sdk-go/aws"

	//"github.com/aws/aws-sdk-go/aws/awsutil"
	//"github.com/aws/aws-sdk-go/aws/credentials"
	core "cloudrack-lambda-core/core"
	"encoding/json"
	"errors"
	_ "image/jpeg"
	"io/ioutil"
	"log"
	"math"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

const S3_HOST = "https://s3.amazonaws.com"
const LIST_PART_MAX_PAGE = 1000

type S3Config struct {
	Region    string
	Bucket    string
	AccessKey string
	Secret    string
	Token     string
	Path      string
}

type S3MultipartConfig struct {
	S3Config          S3Config
	MultipartUploadId string
	S3Key             string
	ChunkSize         int64
	NThreads          int64
}

func Init(bucket string, path string, region string) S3Config {
	return S3Config{
		Bucket: bucket,
		Region: region,
		Path:   path}
}

func (s3c S3Config) Url() string {
	return S3_HOST + "/" + s3c.Bucket + "/" + s3c.Path
}

func (s3c S3Config) Key(key string) string {
	if key == "" {
		panic("Must provide a key to this function")
	}
	if s3c.Path == "" {
		return key
	}
	return s3c.Path + "/" + key
}

func (mpc S3MultipartConfig) Key() string {
	return mpc.S3Config.Key(mpc.S3Key)
}

func (s3c S3Config) UploadJpegToS3(path string, id string, rawData string) (string, error) {
	cfg := aws.NewConfig().WithRegion(s3c.Region)
	svc := s3.New(session.New(), cfg)
	//decoding
	reader := base64.NewDecoder(base64.StdEncoding, strings.NewReader(rawData))
	byteArray, _ := ioutil.ReadAll(reader)
	fileType := http.DetectContentType(byteArray)
	fileBytes := bytes.NewReader(byteArray)
	params := &s3.PutObjectInput{
		Bucket:        aws.String(s3c.Bucket),
		Key:           aws.String(s3c.Path + "/" + path + "/" + id),
		Body:          fileBytes,
		ContentLength: aws.Int64(int64(len(byteArray))),
		ContentType:   aws.String(fileType)}
	//ACL:           aws.String("public-read")
	log.Printf("UPLOADINg TO S3 %+v", params)
	result, err := svc.PutObject(params)
	log.Printf("RESULT OF UPLOAD TO S3 %+v", result)
	log.Printf("RESULT OF UPLOAD TO S3 (ERROR) %+v", err)
	return s3c.Url() + "/" + path + "/" + id, err
}

func (s3c S3Config) SaveJson(path string, id string, data []byte) error {
	return s3c.Save(path, id, data, "application/json")
}
func (s3c S3Config) SaveXml(path string, id string, data []byte) error {
	return s3c.Save(path, id, data, "application/xml")
}

func (s3c S3Config) Save(path string, id string, data []byte, contentType string) error {
	cfg := aws.NewConfig().WithRegion(s3c.Region) //.WithCredentials(creds)
	svc := s3.New(session.New(), cfg)
	dataBytes := bytes.NewReader(data)
	key := s3c.Path + "/" + path + "/" + id
	if path == "" {
		key = s3c.Path + "/" + id
	}
	params := &s3.PutObjectInput{
		Bucket:      aws.String(s3c.Bucket),
		Key:         aws.String(key),
		Body:        dataBytes,
		ContentType: aws.String(contentType)}
	log.Printf("Uploading data to S3 %+v", params)
	result, err := svc.PutObject(params)
	if err == nil {
		log.Printf("Result of upload to S3: %+v", result)
	} else {
		log.Printf("Error while uploading to S3: %+v", err)
	}
	return err
}

//Function to get JSON encode object and unmarshall it
func (s3c S3Config) Get(key string, object interface{}) error {
	sess, _ := session.NewSession(&aws.Config{
		Region: aws.String(s3c.Region)},
	)
	downloader := s3manager.NewDownloader(sess)
	buf := aws.NewWriteAtBuffer([]byte{})
	log.Printf("DOWNLOADING FROM S3 bucket %s at key %s", s3c.Bucket, s3c.Path+"/"+key)
	numBytes, err := downloader.Download(buf,
		&s3.GetObjectInput{
			Bucket: aws.String(s3c.Bucket),
			Key:    aws.String(s3c.Path + "/" + key),
		})
	log.Printf("Downloaded %v Bytes", numBytes)
	log.Printf("[S3] Object Downloaded %s", buf.Bytes())
	json.Unmarshal(buf.Bytes(), object)
	if err != nil {
		log.Printf("[S3] ERROR while downloading object %+v", err)
	}
	return err
}

//Function to get JSON encode object and unmarshall it
func (s3c S3Config) GetTextObj(key string) (string, error) {
	sess, _ := session.NewSession(&aws.Config{
		Region: aws.String(s3c.Region)},
	)
	downloader := s3manager.NewDownloader(sess)
	buf := aws.NewWriteAtBuffer([]byte{})
	log.Printf("DOWNLOADING FROM S3 bucket %s at key %s", s3c.Bucket, s3c.Path+"/"+key)
	numBytes, err := downloader.Download(buf,
		&s3.GetObjectInput{
			Bucket: aws.String(s3c.Bucket),
			Key:    aws.String(s3c.Path + "/" + key),
		})
	log.Printf("Downloaded %v Bytes", numBytes)
	//log.Printf("[S3] Object Downloaded %s", buf.Bytes())
	if err != nil {
		log.Printf("[S3] ERROR while downloading object %+v", err)
	}
	return string(buf.Bytes()), err
}

func (s3c S3Config) GetManyAsText(keys []string) ([]string, error) {
	log.Printf("[S3][GetManyAsText] GetManyAsText from S3 bucket %s", s3c.Bucket)
	var wg sync.WaitGroup
	responses := make([]string, len(keys), len(keys))
	wg.Add(len(keys))
	var lastErr error
	for ind, _ := range keys {
		go func(ind int, responses *[]string, lastErr *error) {
			log.Printf("[S3][GetManyAsText]] getting object: %s\n", keys[ind])
			defer wg.Done()
			res, err := s3c.GetTextObj(keys[ind])
			if err == nil {
				(*responses)[ind] = res
			} else {
				*lastErr = err
			}
		}(ind, &responses, &lastErr)
	}
	wg.Wait()
	return responses, lastErr
}

func (s3c S3Config) Search(prefix string) ([]string, error) {
	sess, _ := session.NewSession(&aws.Config{
		Region: aws.String(s3c.Region)},
	)
	svc := s3.New(sess)
	input := &s3.ListObjectsV2Input{
		Bucket: aws.String("examplebucket"),
		Prefix: aws.String(prefix),
	}
	result, err := svc.ListObjectsV2(input)
	res := []string{}
	for _, obj := range result.Contents {
		res = append(res, *obj.Key)
	}
	if err != nil {
		log.Printf("[S3] ERROR while downloading object %+v", err)
	}
	return res, err
}

func (s3c S3Config) UploadFile(key string, filePath string) error {
	sess, _ := session.NewSession(&aws.Config{
		Region: aws.String(s3c.Region)},
	)
	// Create an uploader with the session and custom options
	uploader := s3manager.NewUploader(sess, func(u *s3manager.Uploader) {
		u.PartSize = 5 * 1024 * 1024 // The minimum/default allowed part size is 5MB
		u.Concurrency = 2            // default is 5
	})

	//open the file
	f, err := os.Open(filePath)
	if err != nil {
		log.Printf("failed to open file %q, %v", filePath, err)
		return err
	}

	// Upload the file to S3.
	_, err2 := uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(s3c.Bucket),
		Key:    aws.String(key),
		Body:   f,
	})

	if err2 != nil {
		log.Printf("File upload faile with error %+v", err)
		return err2
	}

	return nil
}

//TODO: compute the chunksize from the filezise to avoid the 10000 parts limit
//TODO: get default Nthread from system
func (s3c S3Config) StartMPUpload(key string) (S3MultipartConfig, error) {
	sess, _ := session.NewSession(&aws.Config{
		Region: aws.String(s3c.Region)},
	)
	svc := s3.New(sess)
	mpConfig := S3MultipartConfig{S3Config: s3c, S3Key: key, ChunkSize: 5, NThreads: 2}

	params := &s3.CreateMultipartUploadInput{
		Bucket: aws.String(s3c.Bucket),   // Required
		Key:    aws.String(s3c.Key(key)), // Required
	}
	resp, err := svc.CreateMultipartUpload(params)

	if err != nil {
		log.Printf("Failed to create new Multipart Upload: %+v", err)
		return mpConfig, err
	}
	mpConfig.MultipartUploadId = *resp.UploadId
	return mpConfig, nil
}

func (s3c S3Config) ListMpUploads() ([]S3MultipartConfig, error) {
	sess, _ := session.NewSession(&aws.Config{
		Region: aws.String(s3c.Region)},
	)
	svc := s3.New(sess)

	uploads := []S3MultipartConfig{}

	params := &s3.ListMultipartUploadsInput{
		Bucket: aws.String(s3c.Bucket), // Required
	}
	res, err := svc.ListMultipartUploads(params)

	if err != nil {
		// Print the error, cast err to awserr.Error to get the Code and
		// Message from an error.
		log.Printf("Failed to list Multipart Uploads: %+v", err)
		return uploads, err
	}

	for _, upload := range res.Uploads {
		uploads = append(uploads, S3MultipartConfig{S3Config: s3c,
			MultipartUploadId: *upload.UploadId,
			S3Key:             *upload.Key,
		})
	}

	return uploads, nil
}

func (mpc S3MultipartConfig) Abort() error {
	sess, _ := session.NewSession(&aws.Config{
		Region: aws.String(mpc.S3Config.Region)},
	)
	svc := s3.New(sess)

	params := &s3.AbortMultipartUploadInput{
		Bucket:   aws.String(mpc.S3Config.Bucket),   // Required
		Key:      aws.String(mpc.Key()),             // Required
		UploadId: aws.String(mpc.MultipartUploadId), // Required
	}
	_, err := svc.AbortMultipartUpload(params)

	if err != nil {
		log.Printf("Failed to abort Multipart Upload: %+v", err)
		return err
	}
	return nil
}

func (mpc S3MultipartConfig) Parts() ([]*s3.CompletedPart, error) {
	sess, _ := session.NewSession(&aws.Config{
		Region: aws.String(mpc.S3Config.Region)},
	)
	svc := s3.New(sess)

	log.Printf("[Parts] Listing parts for Multipart Upload %+v", mpc.MultipartUploadId)

	parts := []*s3.CompletedPart{}

	params := &s3.ListPartsInput{
		Bucket:   aws.String(mpc.S3Config.Bucket),   // Required
		Key:      aws.String(mpc.Key()),             // Required
		UploadId: aws.String(mpc.MultipartUploadId), // Required
	}
	resp, err := svc.ListParts(params)
	page := 0
	for *resp.IsTruncated && err == nil {
		for _, part := range resp.Parts {
			parts = append(parts, &s3.CompletedPart{ETag: part.ETag, PartNumber: part.PartNumber})
		}
		page++
		log.Printf("[Parts] List is paginated. Getting page %v", page)
		params.PartNumberMarker = resp.NextPartNumberMarker
		resp, err = svc.ListParts(params)

		if page > LIST_PART_MAX_PAGE {
			log.Printf("[Parts] Max number of pages (%v) exeeded. exiting to avoid infinite loop", LIST_PART_MAX_PAGE)
			return parts, errors.New("Max number of pages exeeded. exiting to avoid infinite loop")
		}
	}
	if err != nil {
		// Print the error, cast err to awserr.Error to get the Code and
		// Message from an error.
		log.Printf("Failed to list Uploaded Parts: %+v", err)
		return parts, err
	}
	for _, part := range resp.Parts {
		parts = append(parts, &s3.CompletedPart{ETag: part.ETag, PartNumber: part.PartNumber})
	}

	log.Printf("[Parts] Found %v parts", len(parts))
	return parts, nil
}

func (mpc S3MultipartConfig) Done() error {
	sess, _ := session.NewSession(&aws.Config{
		Region: aws.String(mpc.S3Config.Region)},
	)
	svc := s3.New(sess)

	log.Printf("[Done] Completing Multipart Upload %+v", mpc.MultipartUploadId)

	uploadedParts, err := mpc.Parts()

	log.Printf("[Done] Found %+v uploaded parts", len(uploadedParts))

	if err != nil {
		log.Printf("Failed to complete Multipart Upload since the listPart operation failed with error: %+v", err)
		return err
	}

	if len(uploadedParts) == 0 {
		log.Printf("Failed to complete Multipart upload: no part uploaded. please abort")
		return errors.New("No part uploaded")
	}

	params := &s3.CompleteMultipartUploadInput{
		Bucket:   aws.String(mpc.S3Config.Bucket),   // Required
		Key:      aws.String(mpc.Key()),             // Required
		UploadId: aws.String(mpc.MultipartUploadId), // Required
		MultipartUpload: &s3.CompletedMultipartUpload{
			Parts: uploadedParts,
		},
	}
	_, err2 := svc.CompleteMultipartUpload(params)

	if err2 != nil {
		log.Printf("Failed to complete Multipart Upload: %+v", err2)
		return err2
	}
	return nil
}

func (mpc S3MultipartConfig) Send(partNumber int64, payload []byte) error {
	sess, _ := session.NewSession(&aws.Config{
		Region: aws.String(mpc.S3Config.Region)},
	)
	svc := s3.New(sess)

	params := &s3.UploadPartInput{
		Bucket:     aws.String(mpc.S3Config.Bucket),   // Required
		Key:        aws.String(mpc.Key()),             // Required
		UploadId:   aws.String(mpc.MultipartUploadId), // Required
		PartNumber: aws.Int64(partNumber),             // Required
		Body:       bytes.NewReader(payload),
	}
	_, err := svc.UploadPart(params)

	if err != nil {
		log.Printf("Failed to send part number %v, error: %+v", partNumber, err)
		return err
	}

	return nil
}

func (mpc S3MultipartConfig) Upload(fileName string) {
	log.Printf("[Upload] Starting Upload of file %s", fileName)
	//measuring time
	start := time.Now()
	//0-Opening File
	file, err := os.Open(fileName)

	if err != nil {
		log.Printf("[Upload] error file opening the file: %+v", err)
		os.Exit(1)
	}
	defer file.Close()
	//0-Opening File
	fileInfo, _ := file.Stat()
	var fileSize int64 = fileInfo.Size()
	log.Printf("[Upload] Starting Upload of file %s of size %v", fileName, fileSize)
	fileChunk := mpc.ChunkSize * (1 << 20) // 1 MB, change this to your requirement
	// calculate total number of parts the file will be chunked into
	totalPartsNum := int64(math.Ceil(float64(fileSize) / float64(fileChunk)))
	chunked := core.ChunkArray(totalPartsNum, int64(math.Ceil(float64(totalPartsNum)/float64(mpc.NThreads))))
	log.Printf("[Upload] Splitting to %d pieces: %+v.\n", totalPartsNum, chunked)
	var wg sync.WaitGroup
	log.Printf("[Upload] Adding %+v thead to waitgroup.\n", len(chunked))
	wg.Add(len(chunked))
	for _, parts := range chunked {
		//async
		go func(parts []interface{}) {
			for _, part := range parts {
				partSize := int(math.Min(float64(fileChunk), float64(fileSize-int64(part.(int64)*fileChunk))))
				log.Printf("[Upload] Uploading part %v of size %v bytes", part, partSize)
				partBuffer := make([]byte, partSize)
				file.Read(partBuffer)
				//part number must be greater than 0
				mpc.Send(part.(int64)+1, partBuffer)
				log.Printf("[Upload] Uploading part %v Completed", part)
			}
			wg.Done()
		}(parts)

	}
	wg.Wait()
	elapsed := time.Since(start)
	log.Printf("[Upload] Upload of file %s (size %v) completed in %s", fileName, fileSize, elapsed)

}
