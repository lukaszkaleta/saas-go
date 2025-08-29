package _examples_test

import (
	"context"
	"os"
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/lukaszkaleta/saas-go/database/pg"
	"github.com/lukaszkaleta/saas-go/filestore"
	aws_s3 "github.com/lukaszkaleta/saas-go/filestore/integration/aws-s3"
	"github.com/lukaszkaleta/saas-go/universal"
	"github.com/lukaszkaleta/saas-go/user"
	pguser "github.com/lukaszkaleta/saas-go/user/pg"
)

func TestMakingAwsAvatar(t *testing.T) {
	avatarFile, err := os.Open("cat.png")

	if err != nil {
		panic(err)
	}
	users := pguser.NewPgUsers(pg.NewPg())
	newUser, err := users.Add(&universal.PersonModel{Phone: "123"})
	if err != nil {
		panic(err)
	}
	fileSystemName := string(user.UserAvatarFs)
	fileSystem, err := newUser.FileSystem(fileSystemName)
	if err != nil {
		panic(err)
	}
	s3client := &s3.Client{}
	records := aws_s3.AmazonS3RecordsFromClient(s3client, fileSystemName, *fileSystem.Records())

	avatarModel, err := filestore.SingleDescriptionFromFile(context.Background(), records)(*avatarFile)
	if err != nil {
		panic(err)
	}

	err = newUser.Settings().Avatar().Update(avatarModel)
	if err != nil {
		panic(err)
	}
}
