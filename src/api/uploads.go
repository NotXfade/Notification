package api

import (
	"bytes"
	"fmt"
	"log"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"git.xenonstack.com/xs-onboarding/document-manage/database"
	"git.xenonstack.com/xs-onboarding/document-manage/src/accounts"
	"git.xenonstack.com/xs-onboarding/document-manage/src/aws"
	"git.xenonstack.com/xs-onboarding/document-manage/src/methods"
	"git.xenonstack.com/xs-onboarding/document-manage/src/util"
	jwt "github.com/appleboy/gin-jwt"
	"github.com/gin-gonic/gin"
)

//UploadDocuments is used to upload document on aws
func UploadDocuments(c *gin.Context) {
	defer util.Panic()
	claims := jwt.ExtractClaims(c)
	//	log.Println(claims)
	uid := int(claims["id"].(float64))
	email := claims["email"].(string)

	form, err := c.MultipartForm()
	if err != nil {
		log.Println(err)
		c.JSON(400, gin.H{
			"error":   true,
			"message": "Some error occured while uploading files, Please try again later",
		})
	}
	name, exist := form.Value["type"]

	if exist == false {
		c.JSON(400, gin.H{
			"error":   true,
			"message": "Some error occured while uploading files, Please try again later",
		})
	}
	if len(name) == 0 {
		c.JSON(400, gin.H{
			"error":   true,
			"message": "Please enter valid document type",
		})
	}
	files, exist := form.File["files"]
	if exist == false {
		c.JSON(400, gin.H{
			"error":   true,
			"message": "Some error occured while uploading files, Please try again later",
		})
	}
	if len(files) == 0 {
		c.JSON(400, gin.H{
			"error":   true,
			"message": "No files found to upload",
		})
	}
	sess, err := aws.InitSession()
	if err != nil {
		log.Println(err)
		c.JSON(500, gin.H{
			"error":   true,
			"message": "unable to upload document",
		})
	}
	userid := strconv.Itoa(uid)
	link := make([]string, 0)
	unsuccessfull := ""
	previewLinks := make([]string, 0)
	i := 0
	fileFormats := []string{".doc", ".docx", ".html", ".htm", ".odt", ".pdf", ".xls", ".xlsx", ".ods", ".ppt", ".pptx", ".txt", ".jpg", ".png", ".BMP", ".svg", ".jpeg", ".psd", ".bmp"}
	for index, fileheader := range files {
		file, err := fileheader.Open()
		if err != nil {
			log.Println(err)
			unsuccessfull = fileheader.Filename + "," + unsuccessfull
			continue
		}
		extension := filepath.Ext(fileheader.Filename)
		if !methods.Contains(fileFormats, extension) {
			log.Println("reached here")
			unsuccessfull = fileheader.Filename + "," + unsuccessfull
			continue
		}
		file_name := strings.Split(fileheader.Filename, ".")
		var filename string
		if name[0] == "modules" {
			filename = name[0] + "/" + file_name[0] + "_" + strconv.FormatInt(time.Now().Unix(), 10) + fmt.Sprint(index) + extension
		} else {
			filename = email + "/" + userid + "_" + name[0] + strconv.FormatInt(time.Now().Unix(), 10) + fmt.Sprint(index) + extension
		}

		////
		//log.Println(filename)
		size := fileheader.Size
		//CheckFileSize for files
		err = methods.CheckFileSize(int(size))
		if err != nil {
			log.Println(err)
			unsuccessfull = fileheader.Filename + "," + unsuccessfull
			continue
		}
		buffer := make([]byte, size)
		file.Read(buffer)
		body := bytes.NewReader(buffer)
		mapd, code := aws.Upload(sess, filename, body)
		//log.Println(mapd)
		if code == 200 {
			var uploadData database.UploadLinks
			uploadData.Userid = uid
			uploadData.LinkType = fileheader.Filename
			uploadData.Link = mapd["link"].(string)
			uploadData.UpdatedAt = time.Now()
			aws.SaveUploadInDB(uploadData)
			link = append(link, filename)
		} else {
			unsuccessfull = fileheader.Filename + "," + unsuccessfull
		}
		viewLink, _ := aws.FetchFile(sess, filename)
		previewLinks = append(previewLinks, viewLink)
		i++
	}
	mapd := make(map[string]interface{})
	mapd["error"] = false
	mapd["message"] = "operation successful"
	mapd["uploadlink"] = link
	mapd["preview_links"] = previewLinks
	mapd["unsuccessful_uploads"] = unsuccessfull
	c.JSON(200, mapd)
}

//UploadDocuments : It is used to store documents is aws s3
/* func UploadDocuments(c *gin.Context) {
	defer util.Panic()
	claims := jwt.ExtractClaims(c)
	//log.Println(claims)
	uid := int(claims["id"].(float64))
	//log.Println(uid)
	fileheader, err := c.FormFile("uploadfile")
	file, err := fileheader.Open()
	if err != nil {
		c.JSON(400, "could'nt read the uploaded document ")
	}
	userid := strconv.Itoa(uid)
	filename := userid + "_" + fileheader.Filename
	//log.Println(filename)
	size := fileheader.Size
	buffer := make([]byte, size)
	file.Read(buffer)
	body := bytes.NewReader(buffer)
	sess, err := aws.InitSession()
	if err != nil {
		c.JSON(500, "Internal Server Error : could'nt upload the document ")
	}
	mapd, code := aws.Upload(sess, filename, body)
	if code == 200 {
		var uploadData database.UploadLinks
		uploadData.Userid = uid
		uploadData.LinkType = fileheader.Filename
		uploadData.Link = mapd["link"].(string)
		uploadData.UpdatedAt = time.Now()
		aws.SaveUploadInDB(uploadData)
	}
	mapd["link"] = filename
	c.JSON(code, mapd)
} */

//UploadPolicy is used to upload policy details
func UploadPolicy(c *gin.Context) {
	defer util.Panic()
	claims := jwt.ExtractClaims(c)
	uid := int(claims["id"].(float64))
	level := claims["level"].(string)
	fileheader, err := c.FormFile("uploadfile")
	file, err := fileheader.Open()
	if err != nil {
		c.JSON(400, "could'nt read the uploaded document ")
	}
	filename := fileheader.Filename
	size := fileheader.Size
	buffer := make([]byte, size)
	file.Read(buffer)
	body := bytes.NewReader(buffer)
	sess, err := aws.InitSession()
	if err != nil {
		c.JSON(500, "Internal Server Error : could'nt upload the document ")
	}
	mapd, code := aws.Upload(sess, filename, body)
	if code != 200 {
		c.JSON(code, mapd)
	}
	link := mapd["link"].(string)
	mapd, code = accounts.Saveuploadedpolicy(uid, filename, link, level)
	mapd["link"] = filename
	c.JSON(code, mapd)
}

//==================================== GETFILE FROM S3 =================================================

//GetFile is used to fetch file for both policy and uploads file
func GetFile(c *gin.Context) {
	defer util.Panic()
	filename := c.Param("filename")
	foldername := c.Param("foldername")
	sess, err := aws.InitSession()
	if err != nil {
		c.JSON(500, "Internal Server Error : could'nt upload the document ")
	}
	key := foldername + "/" + filename
	link, err := aws.FetchFile(sess, key)
	if err != nil {
		c.JSON(400, gin.H{
			"error":   true,
			"message": err,
		})
	}
	c.JSON(200, gin.H{
		"error":   false,
		"message": "operation successful",
		"link":    link,
	})
}

//=========================================== GET POLICY ==================================================

//GetPolicy is used to get policy
func GetPolicy(c *gin.Context) {
	defer util.Panic()
	claims := jwt.ExtractClaims(c)
	level := claims["level"].(string)
	policynames, err := accounts.GetPolicyForLevel(level)
	//log.Println(policynames, level)
	if err != nil {
		c.JSON(400, gin.H{
			"error":   true,
			"message": "please pass correct level",
		})
		return
	}
	sess, err := aws.InitSession()
	if err != nil {
		c.JSON(500, "Internal Server Error : could'nt initialise session ")
		return
	}
	links := make([]string, 0)
	for i := 0; i < len(policynames); i++ {
		link, err := aws.FetchFile(sess, policynames[i])
		if err != nil {
			c.JSON(400, gin.H{
				"error":   true,
				"message": "error ocurred while fetching file",
			})
			return
		}
		links = append(links, link)
	}
	c.JSON(200, gin.H{
		"error":   false,
		"message": "operation successful",
		"links":   links,
	})
}

//DeleteDocument is used to delete document using key name
func DeleteDocument(c *gin.Context) {
	//handler panic and Alerts
	defer util.Panic()
	filename := c.Param("filename")
	foldername := c.Param("foldername")
	key := foldername + "/" + filename
	sess, err := aws.InitSession()
	if err != nil {
		c.JSON(500, "Internal Server Error : could'nt initialise session ")
		return
	}
	err = aws.DeleteFile(sess, key)
	if err != nil {
		c.JSON(400, gin.H{
			"error":   true,
			"message": "could not delete file" + err.Error(),
		})
		return
	}
	c.JSON(200, gin.H{
		"error":   false,
		"message": "file deleted successfuly",
	})
}
