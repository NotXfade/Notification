//+build !test

package aws

import (
	"git.xenonstack.com/xs-onboarding/document-manage/config"
	"git.xenonstack.com/xs-onboarding/document-manage/src/util"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
)

// InitSession create a session for login credentials
func InitSession() (*session.Session, error) {
	defer util.Panic()
	//sess store session value
	var sess *session.Session
	var err error
	cred := credentials.NewStaticCredentialsFromCreds(credentials.Value{
		AccessKeyID:     config.Conf.AWS.AccessId,
		SecretAccessKey: config.Conf.AWS.AccessKey,
	})

	//sess store session value
	sess, err = session.NewSession(&aws.Config{
		Region:      aws.String("ap-south-1"),
		Credentials: cred,
	})
	return sess, err
}
