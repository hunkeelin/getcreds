package main

import (
	"bytes"
	"io/ioutil"
	"os/user"
)

func grabProfiles() (map[string]awsCredentials, error) {
	credMap := make(map[string]awsCredentials)
	user, err := user.Current()
	if err != nil {
		return credMap, err
	}
	awsDir := user.HomeDir + "/.aws/credentials"
	credentialFile, err := ioutil.ReadFile(awsDir)
	if err != nil {
		return credMap, err
	}
	splitByLines := bytes.Split(credentialFile, []byte("\n"))
	var cp string
	for _, line := range splitByLines {
		switch {
		case bytes.HasPrefix(line, []byte("[")):
			cp = string(bytes.Trim(bytes.Trim(line, "["), "]"))
		case bytes.HasPrefix(line, []byte("aws_access_key_id")):
			tmp := credMap[cp]
			tmp.awsAccessKeyId = string(bytes.Replace(line, []byte("aws_access_key_id = "), []byte(""), -1))
			credMap[cp] = tmp
		case bytes.HasPrefix(line, []byte("aws_secret_access_key")):
			tmp := credMap[cp]
			tmp.awsSecretAccessKey = string(bytes.Replace(line, []byte("aws_secret_access_key = "), []byte(""), -1))
			credMap[cp] = tmp
		case bytes.HasPrefix(line, []byte("aws_session_token")):
			tmp := credMap[cp]
			tmp.awsSessionToken = string(bytes.Replace(line, []byte("aws_session_token = "), []byte(""), -1))
			credMap[cp] = tmp
		default:
			continue
		}
	}
	return credMap, nil
}

type awsCredentials struct {
	awsAccessKeyId     string
	awsSecretAccessKey string
	awsSessionToken    string
}
