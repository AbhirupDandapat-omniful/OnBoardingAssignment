package store

import (
    "context"

    "github.com/aws/aws-sdk-go-v2/aws"
    awsCfg "github.com/aws/aws-sdk-go-v2/config"
    "github.com/aws/aws-sdk-go-v2/service/s3"
    "github.com/omniful/go_commons/config"
)

// NewLocalStackS3Client returns an AWS SDK v2 S3 client configured
// to talk to LocalStack (path-style addressing).
func NewLocalStackS3Client(ctx context.Context) (*s3.Client, error) {
    endpoint := config.GetString(ctx, "s3.endpoint")
    region   := config.GetString(ctx, "aws.region")

    cfg, err := awsCfg.LoadDefaultConfig(ctx,
        awsCfg.WithRegion(region),
        awsCfg.WithEndpointResolver(aws.EndpointResolverFunc(func(service, region string) (aws.Endpoint, error) {
            if service == s3.ServiceID {
                return aws.Endpoint{
                    URL:               endpoint,
                    HostnameImmutable: true,
                }, nil
            }
            return aws.Endpoint{}, &aws.EndpointNotFoundError{}
        })),
    )
    if err != nil {
        return nil, err
    }

    return s3.NewFromConfig(cfg, func(o *s3.Options) {
        o.UsePathStyle = true
    }), nil
}
