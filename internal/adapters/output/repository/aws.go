package repository

import (
	"context"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
)

func DefaultiLocalAWSClientConfig() aws.Config {

	host := "http://localhost:4566" // default value pointing for local stack
	region := "us-east-1"

	if v, ok := os.LookupEnv("AWS_ENDPOINT"); ok {
		host = v
	}
	if v, ok := os.LookupEnv("AWS_REGION"); ok {
		region = v
	}

	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion(region),
		config.WithCredentialsProvider(aws.AnonymousCredentials{}),
		config.WithEndpointResolverWithOptions(aws.EndpointResolverWithOptionsFunc(
			func(service, region string, opts ...any) (aws.Endpoint, error) {
				return aws.Endpoint{
					URL: host,
				}, nil
			},
		)),
	)
	if err != nil {
		panic(err)
	}

	return cfg
}
