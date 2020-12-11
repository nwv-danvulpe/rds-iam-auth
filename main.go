package main

import (
	"database/sql"
	"fmt"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/credentials/ec2rolecreds"
	"github.com/aws/aws-sdk-go/aws/ec2metadata"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/rds/rdsutils"
	"log"
	"os"
)

func main() {
	dbName := os.Getenv("DB_NAME")
	dbHost := os.Getenv("DB_HOST")
	dbUser := os.Getenv("DB_USER")
	dbEndpoint := fmt.Sprintf("%s:5432", dbHost)
	region := os.Getenv("AWS_REGION")
	sess := session.Must(session.NewSession())
	creds := credentials.NewChainCredentials(
		[]credentials.Provider{
			&ec2rolecreds.EC2RoleProvider{
				Client: ec2metadata.New(sess),
			},
		})
	token, err := rdsutils.BuildAuthToken(dbEndpoint, region, dbUser, creds)
	if err != nil {
		log.Fatalf("error loading creds: %v", err)
	}

	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s", dbHost, 5432, dbUser, token, dbName)
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatalf("could not connect to database: %v", err)
	}
	defer db.Close()
	if err := db.Ping(); err != nil {
		log.Fatalf("could not ping database: %v", err)
	}
}
