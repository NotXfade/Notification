package accounts

import (
	"log"
	"os"
	"testing"

	"git.xenonstack.com/xs-onboarding/document-manage/config"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	_ "github.com/mattn/go-sqlite3"
)

func init() {

	os.Remove(os.Getenv("HOME") + "/documents-testing.db")
	db, err := gorm.Open("sqlite3", os.Getenv("HOME")+"/documents-testing.db")
	if err != nil {
		log.Println(err)
		log.Println("Exit")
		os.Exit(1)
	}
	config.DB = db

}
func TestGetLink(t *testing.T) {
	isExist := GetLink("test@example.com/someexample.png")
	if isExist == true {
		t.Error(isExist)
	}
}

func TestSaveuploadedpolicy(t *testing.T) {
	_, code := Saveuploadedpolicy(1, "traning", "test@example.com/someexample.png", "L1")
	if code == 200 {
		t.Error(code)
	}
}

func TestGetPolicyForLevel(t *testing.T) {
	_, err := GetPolicyForLevel("L1")
	if err == nil {
		t.Error("Test case for get policy for level failed")
	}
}
