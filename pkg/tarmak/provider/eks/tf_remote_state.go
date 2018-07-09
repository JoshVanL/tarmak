// Copyright Jetstack Ltd. See LICENSE for details.
package eks

import (
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/s3"
)

func (e *EKS) RemoteStateName() string {
	return fmt.Sprintf(
		"%s%s-terraform-state",
		e.conf.EKS.BucketPrefix,
		e.Region(),
	)
}

const DynamoDBKey = "LockID"

// TODO: remove me, deprecated
func (e *EKS) RemoteStateBucketName() string {
	return e.RemoteStateName()
}

func (e *EKS) RemoteState(namespace string, clusterName string, stackName string) string {
	return fmt.Sprintf(`terraform {
  backend "s3" {
    bucket = "%s"
    key = "%s"
    region = "%s"
    dynamodb_table ="%s"
  }
}`,
		e.RemoteStateName(),
		fmt.Sprintf("%s/%s/%s.tfstate", namespace, clusterName, stackName),
		e.Region(),
		e.RemoteStateName(),
	)
}

func (e *EKS) RemoteStateBucketAvailable() (bool, error) {
	svc, err := e.S3()
	if err != nil {
		return false, err
	}

	_, err = svc.HeadBucket(&s3.HeadBucketInput{
		Bucket: aws.String(e.RemoteStateName()),
	})
	if err == nil {
		return true, nil
	} else if strings.HasPrefix(err.Error(), "NotFound:") {
		return false, nil
	}

	return false, fmt.Errorf("error while checking if remote state is available: %s", err)
}

func (e *EKS) RemoteStateAvailable(bucketName string) (bool, error) {
	sess, err := e.Session()
	if err != nil {
		return false, fmt.Errorf("error getting session: %s", err)
	}

	svc := s3.New(sess)
	_, err = svc.HeadBucket(&s3.HeadBucketInput{
		Bucket: &bucketName,
	})
	if err == nil {
		return true, nil
	} else if strings.HasPrefix(err.Error(), "NotFound:") {
		return false, nil
	} else {
		return false, fmt.Errorf("error while checking if remote state is available: %s", err)
	}
}
func (e *EKS) initRemoteStateBucket() error {
	svc, err := e.S3()
	if err != nil {
		return err
	}

	createBucketInput := &s3.CreateBucketInput{
		Bucket: aws.String(e.RemoteStateName()),
	}

	if e.Region() != "us-east-1" {
		createBucketInput.CreateBucketConfiguration = &s3.CreateBucketConfiguration{
			LocationConstraint: aws.String(e.Region()),
		}
	}

	_, err = svc.CreateBucket(createBucketInput)
	if err != nil {
		return err
	}

	_, err = svc.PutBucketVersioning(&s3.PutBucketVersioningInput{
		Bucket: aws.String(e.RemoteStateName()),
		VersioningConfiguration: &s3.VersioningConfiguration{
			Status: aws.String("Enabled"),
		},
	})
	return err
}

func (e *EKS) validateRemoteStateBucket() error {
	svc, err := e.S3()
	if err != nil {
		return err
	}

	_, err = svc.HeadBucket(&s3.HeadBucketInput{
		Bucket: aws.String(e.RemoteStateName()),
	})
	if err != nil {
		if awsErr, ok := err.(awserr.Error); ok {
			if awsErr.Code() == "NotFound" {
				return e.initRemoteStateBucket()
			}
		}
		return fmt.Errorf("error looking for terraform state bucket: %s", err)
	}

	location, err := svc.GetBucketLocation(&s3.GetBucketLocationInput{
		Bucket: aws.String(e.RemoteStateName()),
	})
	if err != nil {
		return err
	}

	var bucketRegion string
	if location.LocationConstraint == nil {
		bucketRegion = "us-east-1"
	} else {
		bucketRegion = *location.LocationConstraint
	}

	if myRegion := e.Region(); bucketRegion != myRegion {
		return fmt.Errorf("bucket region is wrong, actual: %s expected: %s", bucketRegion, myRegion)
	}

	versioning, err := svc.GetBucketVersioning(&s3.GetBucketVersioningInput{
		Bucket: aws.String(e.RemoteStateName()),
	})
	if err != nil {
		return err
	}
	if *versioning.Status != "Enabled" {
		e.log.Warnf("state bucket %s has versioning disabled", e.RemoteStateName())
	}

	return nil

}

func (e *EKS) initRemoteStateDynamoDB() error {
	svc, err := e.DynamoDB()
	if err != nil {
		return err
	}

	_, err = svc.CreateTable(&dynamodb.CreateTableInput{
		TableName: aws.String(e.RemoteStateName()),
		AttributeDefinitions: []*dynamodb.AttributeDefinition{
			&dynamodb.AttributeDefinition{
				AttributeName: aws.String(DynamoDBKey),
				AttributeType: aws.String("S"),
			},
		},
		KeySchema: []*dynamodb.KeySchemaElement{
			&dynamodb.KeySchemaElement{
				AttributeName: aws.String(DynamoDBKey),
				KeyType:       aws.String(dynamodb.KeyTypeHash),
			},
		},
		ProvisionedThroughput: &dynamodb.ProvisionedThroughput{
			ReadCapacityUnits:  aws.Int64(1),
			WriteCapacityUnits: aws.Int64(1),
		},
	})
	return err
}

func (e *EKS) validateRemoteStateDynamoDB() error {
	svc, err := e.DynamoDB()
	if err != nil {
		return err
	}

	describeOut, err := svc.DescribeTable(&dynamodb.DescribeTableInput{
		TableName: aws.String(e.RemoteStateName()),
	})
	if err != nil {
		if awsErr, ok := err.(awserr.Error); ok {
			if awsErr.Code() == "ResourceNotFoundException" {
				return e.initRemoteStateDynamoDB()
			}
		}
		return fmt.Errorf("error looking for terraform state dynamodb: %s", err)
	}

	attributeFound := false
	for _, params := range describeOut.Table.AttributeDefinitions {
		if *params.AttributeName == DynamoDBKey {
			attributeFound = true
		}
	}
	if !attributeFound {
		return fmt.Errorf("the DynamoDB table '%s' doesn't contain a parameter named '%s'", e.RemoteStateName(), DynamoDBKey)
	}

	return nil
}
