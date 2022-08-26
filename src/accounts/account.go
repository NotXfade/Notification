package accounts

import (
	"time"

	"git.xenonstack.com/xs-onboarding/document-manage/config"
	"git.xenonstack.com/xs-onboarding/document-manage/database"
	"git.xenonstack.com/xs-onboarding/document-manage/src/util"
)

//GetLink : It is used to get links related to a user id
func GetLink(link string) bool {
	defer util.Panic()
	db := config.DB
	links := []database.UploadLinks{}
	db.Where("link=?", link).Find(&links)
	if len(links) == 0 {
		return false
	}
	return true
}

//Saveuploadedpolicy : this function is use to  create entry to save uploaded document
func Saveuploadedpolicy(userid int, policyname, policylink, level string) (map[string]interface{}, int) {
	defer util.Panic()
	mapd := make(map[string]interface{})
	var link []database.PolicyLinks
	db := config.DB
	//check if this link is already present and being updated
	db.Where("policy_link=?", policylink).Find(&link)
	if len(link) != 0 {
		//if yes then update the database
		link[0].UpdatedAt = time.Now()
		link[0].PolicyName = policyname
		link[0].PolicyLink = policylink
		link[0].AddedBy = userid
		link[0].Level = level
		err := db.Save(&link[0]).Error
		if err != nil {
			mapd["error"] = true
			mapd["message"] = "database Error | Please enter correct Data"
			return mapd, 400
		}
		mapd["error"] = false
		mapd["message"] = "file uploaded successfuly"
		return mapd, 200
	}
	//else create a new entry for uploaded new policy
	policyLink := database.PolicyLinks{
		AddedBy:    userid,
		PolicyName: policyname,
		PolicyLink: policylink,
		Level:      level,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}
	err := db.Create(&policyLink).Error
	if err != nil {
		mapd["error"] = true
		mapd["message"] = "database Error | Please enter correct Data"
		return mapd, 400
	}
	mapd["error"] = false
	mapd["message"] = "file uploaded successfuly"
	return mapd, 200
}

//GetPolicyForLevel is used to return policyfilenames
func GetPolicyForLevel(level string) ([]string, error) {
	defer util.Panic()
	db := config.DB
	policies := []database.PolicyLinks{}
	filenames := make([]string, 0)
	db = db.Debug()
	err := db.Where("level=?", level).Find(&policies).Error
	if err != nil {
		return filenames, err
	}
	//log.Println(policies)
	for i := 0; i < len(policies); i++ {
		filenames = append(filenames, policies[i].PolicyName)
	}
	return filenames, nil
}
