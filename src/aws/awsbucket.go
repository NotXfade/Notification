//+build !test

package aws

import (
	"errors"
	"io"
	"log"
	"mime"
	"path/filepath"
	"time"

	"git.xenonstack.com/xs-onboarding/document-manage/config"
	"git.xenonstack.com/xs-onboarding/document-manage/database"
	"git.xenonstack.com/xs-onboarding/document-manage/src/util"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

//Upload is used to upload files in s3 bucket
func Upload(sess *session.Session, filename string, body io.Reader) (map[string]interface{}, int) {
	defer util.Panic()
	mapd := make(map[string]interface{})
	uploader := s3manager.NewUploader(sess)
	result, err := uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(config.Conf.AWS.BucketName),
		Key:    aws.String(filename),
		Body:   body,
		ACL:    aws.String("private"),
	})
	if err != nil {
		log.Print(err)
		mapd["error"] = true
		mapd["message"] = "failed to upload the file"
		return mapd, 400
	}
	mapd["error"] = false
	mapd["message"] = "file upload successful"
	mapd["link"] = result.Location
	//log.Println(result.Location)
	return mapd, 200
}

//SaveUploadInDB is used to save upload data activity in DB
func SaveUploadInDB(uploadData database.UploadLinks) {
	defer util.Panic()
	db := config.DB
	linkData := database.UploadLinks{}
	db.Where("userid=? AND link_type=?", uploadData.Userid, uploadData.LinkType).Find(&linkData)
	if linkData.ID == 0 {
		uploadData.CreatedAt = time.Now()
		db.Create(&uploadData)

	} else {
		linkData.UpdatedAt = time.Now()
		db.Save(&linkData)
	}
}

//=============================== FetchFile function to get file from aws ========================================

//FetchFile is used to get presigned url from aws
func FetchFile(sess *session.Session, filename string) (string, error) {
	defer util.Panic()
	// Create S3 service client
	svc := s3.New(sess)
	contentType := GetContentType(filename)
	//log.Println(contentType)
	req, _ := svc.GetObjectRequest(&s3.GetObjectInput{
		Bucket:                     aws.String(config.Conf.AWS.BucketName),
		Key:                        aws.String(filename),
		ResponseContentDisposition: aws.String("inline"),
		ResponseContentType:        aws.String(contentType),
	})
	//create a signed url
	link, err := req.Presign(config.Conf.Service.PresignedLinkTimeout)
	if err != nil {
		log.Println(err)
		return "", errors.New("Failed to sign request")
	}
	return link, nil
}

//GetContentType is used to get mime type based on extension
func GetContentType(filename string) string {
	defer util.Panic()
	//log.Println(filename)
	extension := filepath.Ext(filename)
	ext := "." + extension
	return mime.TypeByExtension(ext)
}

//================================= DeleteFile function to delete file from aws ======================================

//DeleteFile is used to delete file from s3
func DeleteFile(sess *session.Session, filename string) error {
	defer util.Panic()
	// Create S3 service client
	svc := s3.New(sess)
	_, err := svc.DeleteObject(&s3.DeleteObjectInput{
		Bucket: aws.String(config.Conf.AWS.BucketName),
		Key:    aws.String(filename),
	})
	if err != nil {
		return err
	}

	err = svc.WaitUntilObjectNotExists(&s3.HeadObjectInput{
		Bucket: aws.String(config.Conf.AWS.BucketName),
		Key:    aws.String(filename),
	})
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}
