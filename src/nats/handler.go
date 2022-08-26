//+build !test

package nats

import (
	"encoding/json"
	"log"

	"git.xenonstack.com/xs-onboarding/document-manage/src/accounts"
	"git.xenonstack.com/xs-onboarding/document-manage/src/aws"
	"git.xenonstack.com/xs-onboarding/document-manage/src/util"
)

//Getlink is used to return back
func Getlink(data []byte) []byte {
	defer util.Panic()
	link := string(data)
	val := accounts.GetLink(link)
	var value string
	if val == true {
		value = "true"
	} else {
		value = "false"
	}
	payload := []byte(value)
	return payload
}

func GetPresignedUrl(data []byte) []byte {
	defer util.Panic()
	mapd := make(map[string]interface{})
	filename := string(data)
	sess, err := aws.InitSession()
	if err != nil {
		mapd["error"] = true
		mapd["message"] = err
		mapd["link"] = ""
		payload, err := json.Marshal(mapd)
		if err != nil {
			log.Println(err)
		}
		return payload
	}
	link, err := aws.FetchFile(sess, filename)
	if err != nil {
		mapd["error"] = true
		mapd["message"] = err
		mapd["link"] = ""
		payload, err := json.Marshal(mapd)
		if err != nil {
			log.Println(err)
		}
		return payload
	}
	mapd["error"] = false
	mapd["message"] = "Operation Successful"
	mapd["link"] = link
	payload, err := json.Marshal(mapd)
	if err != nil {
		log.Println(err)
	}
	return payload
}

func DeleteFile(data []byte) []byte {
	filename := string(data)
	sess, err := aws.InitSession()
	var isDeleted bool
	if err != nil {
		log.Println(err)
		isDeleted = false
	}
	err = aws.DeleteFile(sess, filename)
	if err != nil {
		log.Println(err)
		isDeleted = false
	} else {
		isDeleted = true
	}
	payload, err := json.Marshal(isDeleted)
	if err != nil {
		log.Println(err)
	}
	return payload
}
