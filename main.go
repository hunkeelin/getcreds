package main

import (
	"flag"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sts"
	"os"
)

var (
	tokenCode  = flag.String("t", "", "the mfa code")
	assumeRole = flag.String("p", "", "The aws profile to assume")
	roleArn    = flag.String("r", "", "the arn of your role")
)

func main() {
	flag.Parse()
	if *tokenCode == "" {
		panic(fmt.Errorf("tokencode cannot be empty"))
	}
	svc := sts.New(session.New())
	input := &sts.GetSessionTokenInput{
		DurationSeconds: aws.Int64(129600),
		SerialNumber:    aws.String(*roleArn),
		TokenCode:       aws.String(*tokenCode),
	}

	result, err := svc.GetSessionToken(input)
	if err != nil {
		panic(err)
	}
	profiles, err := grabProfiles()
	if err != nil {
		panic(err)
	}
	tmp := profiles[*assumeRole]
	tmp.awsAccessKeyId = *result.Credentials.AccessKeyId
	tmp.awsSecretAccessKey = *result.Credentials.SecretAccessKey
	tmp.awsSessionToken = *result.Credentials.SessionToken
	profiles[*assumeRole] = tmp

	var towrite []byte
	for k, v := range profiles {
		profile := []byte("[" + k + "]\n")
		towrite = append(towrite, profile...)
		accessKey := []byte("aws_access_key_id = " + v.awsAccessKeyId + "\n")
		towrite = append(towrite, accessKey...)
		secretKey := []byte("aws_secret_access_key = " + v.awsSecretAccessKey + "\n")
		towrite = append(towrite, secretKey...)
		if v.awsSessionToken != "" {
			token := []byte("aws_session_token = " + v.awsSessionToken + "\n")
			towrite = append(towrite, token...)
		}
	}
	dirname, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}

	err = os.WriteFile(dirname+"/.aws/credentials", towrite, 0644)
	if err != nil {
		panic(err)
	}
}
