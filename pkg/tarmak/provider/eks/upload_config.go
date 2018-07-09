// Copyright Jetstack Ltd. See LICENSE for details.
package eks

import (
	"fmt"
	"io"
	"path/filepath"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"

	"github.com/jetstack/tarmak/pkg/tarmak/interfaces"
)

// This uploads the main configuration to the S3 bucket
func (e *EKS) UploadConfiguration(cluster interfaces.Cluster, configReader io.ReadSeeker) error {
	svc, err := e.S3()
	if err != nil {
		return err
	}

	bucketName := fmt.Sprintf(
		"%s%s-%s-secrets",
		e.conf.EKS.BucketPrefix,
		cluster.Environment().Name(),
		e.Region(),
	)

	bucketPath := filepath.Join(cluster.ClusterName(), "puppet.tar.gz")

	_, err = svc.PutObject(&s3.PutObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(bucketPath),
		Body:   configReader,
	})
	return err
}
