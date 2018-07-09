// Copyright Jetstack Ltd. See LICENSE for details.
package eks

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/route53"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/hashicorp/go-multierror"
	vault "github.com/hashicorp/vault/api"
	"github.com/sirupsen/logrus"

	clusterv1alpha1 "github.com/jetstack/tarmak/pkg/apis/cluster/v1alpha1"
	tarmakv1alpha1 "github.com/jetstack/tarmak/pkg/apis/tarmak/v1alpha1"
	"github.com/jetstack/tarmak/pkg/tarmak/interfaces"
	"github.com/jetstack/tarmak/pkg/tarmak/utils/input"
)

var _ interfaces.Provider = &EKS{}

type EKS struct {
	conf *tarmakv1alpha1.Provider

	tarmak interfaces.Tarmak

	availabilityZones *[]string

	session  *session.Session
	ec2      EC2
	s3       S3
	dynamodb DynamoDB
	route53  Route53
	log      *logrus.Entry
}

type S3 interface {
	HeadBucket(input *s3.HeadBucketInput) (*s3.HeadBucketOutput, error)
	CreateBucket(input *s3.CreateBucketInput) (*s3.CreateBucketOutput, error)
	GetBucketVersioning(input *s3.GetBucketVersioningInput) (*s3.GetBucketVersioningOutput, error)
	GetBucketLocation(input *s3.GetBucketLocationInput) (*s3.GetBucketLocationOutput, error)
	PutBucketVersioning(input *s3.PutBucketVersioningInput) (*s3.PutBucketVersioningOutput, error)
	PutObject(input *s3.PutObjectInput) (*s3.PutObjectOutput, error)
}

type EC2 interface {
	DescribeInstances(input *ec2.DescribeInstancesInput) (*ec2.DescribeInstancesOutput, error)
	ImportKeyPair(input *ec2.ImportKeyPairInput) (*ec2.ImportKeyPairOutput, error)
	DescribeKeyPairs(input *ec2.DescribeKeyPairsInput) (*ec2.DescribeKeyPairsOutput, error)
	DescribeAvailabilityZones(input *ec2.DescribeAvailabilityZonesInput) (*ec2.DescribeAvailabilityZonesOutput, error)
	DescribeRegions(input *ec2.DescribeRegionsInput) (*ec2.DescribeRegionsOutput, error)
	DescribeReservedInstancesOfferings(input *ec2.DescribeReservedInstancesOfferingsInput) (*ec2.DescribeReservedInstancesOfferingsOutput, error)
}

type DynamoDB interface {
	DescribeTable(input *dynamodb.DescribeTableInput) (*dynamodb.DescribeTableOutput, error)
	CreateTable(input *dynamodb.CreateTableInput) (*dynamodb.CreateTableOutput, error)
}

type Route53 interface {
	CreateHostedZone(input *route53.CreateHostedZoneInput) (*route53.CreateHostedZoneOutput, error)
	GetHostedZone(input *route53.GetHostedZoneInput) (*route53.GetHostedZoneOutput, error)
	ListHostedZonesByName(input *route53.ListHostedZonesByNameInput) (*route53.ListHostedZonesByNameOutput, error)
}

func NewFromConfig(tarmak interfaces.Tarmak, conf *tarmakv1alpha1.Provider) (*EKS, error) {

	e := &EKS{
		conf:   conf,
		log:    tarmak.Log().WithField("provider_name", conf.ObjectMeta.Name),
		tarmak: tarmak,
	}

	return e, nil
}

func (e *EKS) Name() string {
	return e.conf.Name
}

func (e *EKS) Cloud() string {
	return clusterv1alpha1.CloudEKS
}

// this clears all cached state from the provider
func (e *EKS) Reset() {
	e.dynamodb = nil
	e.session = nil
	e.s3 = nil
	e.ec2 = nil
	e.route53 = nil
	e.availabilityZones = nil
}

// This parameters should include non sensitive information to identify a provider
func (e *EKS) Parameters() map[string]string {
	p := map[string]string{
		"name":          e.Name(),
		"cloud":         e.Cloud(),
		"public_zone":   e.conf.EKS.PublicZone,
		"bucket_prefix": e.conf.EKS.BucketPrefix,
	}
	if e.conf.EKS.VaultPath != "" {
		p["vault_path"] = e.conf.EKS.VaultPath
	}
	if e.conf.EKS.Profile != "" {
		p["amazon_profile"] = e.conf.EKS.Profile
	}
	return p
}

func (e *EKS) String() string {
	return fmt.Sprintf("%s[%s]", e.Cloud(), e.Name())
}

func (e *EKS) ListRegions() (regions []string, err error) {
	svc, err := e.EC2()
	if err != nil {
		return regions, err
	}

	regionsOutput, err := svc.DescribeRegions(&ec2.DescribeRegionsInput{})
	if err != nil {
		return regions, err
	}

	for _, region := range regionsOutput.Regions {
		regions = append(regions, *region.RegionName)
	}

	sort.Strings(regions)

	return regions, nil

}

func (e *EKS) AskEnvironmentLocation(init interfaces.Initialize) (location string, err error) {
	regions, err := e.ListRegions()
	if err != nil {
		return "", err
	}

	regionPos, err := init.Input().AskSelection(&input.AskSelection{
		Query:   "In which region should this environment reside?",
		Choices: regions,
		Default: -1,
	})
	if err != nil {
		return "", err
	}

	return regions[regionPos], nil
}

func (e *EKS) AskInstancePoolZones(init interfaces.Initialize) (zones []string, err error) {

	zones, err = e.getAvailablityZoneByRegion()
	if err != nil {
		return []string{}, fmt.Errorf("failed to get availabilty zones: %v", err)
	}

	if len(zones) == 0 {
		return []string{}, fmt.Errorf("no availability zones found for region '%s'", e.Region())
	}

	sChoices := make([]bool, len(zones))
	sChoices[0] = true

	multiSel := &input.AskMultipleSelection{
		AskSelection: &input.AskSelection{
			Query:   "Please select availabilty zones. Enter numbers to toggle selection.",
			Choices: zones,
			Default: 1,
		},
		SelectedChoices: sChoices,
		MinSelected:     1,
		MaxSelected:     len(zones),
	}

	return init.Input().AskMultipleSelection(multiSel)
}

func (e *EKS) Region() string {
	// without environment selected, fall back to default region
	if e.tarmak.Environment() == nil {
		return "us-east-1"
	}
	return e.tarmak.Environment().Location()
}

// This return the availabililty zones that are used for a cluster
func (e *EKS) AvailabilityZones() (availabiltyZones []string) {
	if e.availabilityZones != nil {
		return *e.availabilityZones
	}

	subnets := e.tarmak.Cluster().Subnets()
	zones := make(map[string]bool)

	for _, subnet := range subnets {
		zones[subnet.Zone] = true
	}

	e.availabilityZones = &availabiltyZones

	for zone, _ := range zones {
		availabiltyZones = append(availabiltyZones, zone)
	}

	sort.Strings(availabiltyZones)

	return availabiltyZones
}

func (e *EKS) EC2() (EC2, error) {
	if e.ec2 == nil {
		sess, err := e.Session()
		if err != nil {
			return nil, fmt.Errorf("error getting EKS session: %s", err)
		}
		e.ec2 = ec2.New(sess)
	}
	return e.ec2, nil
}

func (e *EKS) S3() (S3, error) {
	if e.s3 == nil {
		sess, err := e.Session()
		if err != nil {
			return nil, fmt.Errorf("error getting EKS session: %s", err)
		}
		e.s3 = s3.New(sess)
	}
	return e.s3, nil
}

func (e *EKS) DynamoDB() (DynamoDB, error) {
	if e.dynamodb == nil {
		sess, err := e.Session()
		if err != nil {
			return nil, fmt.Errorf("error getting EKS session: %s", err)
		}
		e.dynamodb = dynamodb.New(sess)
	}
	return e.dynamodb, nil
}

func (e *EKS) Route53() (Route53, error) {
	if e.route53 == nil {
		sess, err := e.Session()
		if err != nil {
			return nil, fmt.Errorf("error getting EKS session: %s", err)
		}
		e.route53 = route53.New(sess)
	}
	return e.route53, nil
}

func (e *EKS) Variables() map[string]interface{} {
	output := map[string]interface{}{}
	output["key_name"] = e.KeyName()
	if len(e.conf.EKS.AllowedAccountIDs) > 0 {
		output["allowed_account_ids"] = e.conf.EKS.AllowedAccountIDs
	}
	output["availability_zones"] = e.AvailabilityZones()
	output["region"] = e.Region()

	output["public_zone"] = e.conf.EKS.PublicZone
	output["public_zone_id"] = e.conf.EKS.PublicHostedZoneID
	output["bucket_prefix"] = e.conf.EKS.BucketPrefix

	return output
}

// This will return necessary environment variables
func (e *EKS) Environment() ([]string, error) {
	sess, err := e.Session()
	if err != nil {
		return []string{}, fmt.Errorf("error getting session: %s", err)
	}

	creds, err := sess.Config.Credentials.Get()
	if err != nil {
		return []string{}, fmt.Errorf("error getting credentials: %s", err)
	}

	return []string{
		fmt.Sprintf("AWS_ACCESS_KEY_ID=%s", creds.AccessKeyID),
		fmt.Sprintf("AWS_SECRET_ACCESS_KEY=%s", creds.SecretAccessKey),
		fmt.Sprintf("AWS_SESSION_TOKEN=%s", creds.SessionToken),
		fmt.Sprintf("AWS_DEFAULT_REGION=%s", e.Region()),
	}, nil
}

// This reads the vault token from ~/.vault-token
func (e *EKS) readVaultToken() (string, error) {
	homeDir := e.tarmak.HomeDir()

	filePath := filepath.Join(homeDir, ".vault-token")

	vaultToken, err := ioutil.ReadFile(filePath)
	if err != nil {
		return "", err
	}

	return string(vaultToken), nil
}

func (e *EKS) Validate() error {
	var result error
	var err error

	// These checks only make sense with an environment given
	if e.tarmak.Environment() != nil {
		err = e.validateRemoteStateBucket()
		if err != nil {
			result = multierror.Append(result, err)
		}

		err = e.validateRemoteStateDynamoDB()
		if err != nil {
			result = multierror.Append(result, err)
		}

		err = e.validateAvailabilityZones()
		if err != nil {
			result = multierror.Append(result, err)
		}

		err = e.validateAWSKeyPair()
		if err != nil {
			result = multierror.Append(result, err)
		}

	}

	err = e.validatePublicZone()
	if err != nil {
		result = multierror.Append(result, err)
	}

	if result != nil {
		return result
	}
	return nil

}

func (e *EKS) Verify() (result error) {
	return result
}

func (e *EKS) getAvailablityZoneByRegion() (zones []string, err error) {
	svc, err := e.EC2()
	if err != nil {
		return []string{}, fmt.Errorf("error getting AWS EC2 session: %s", err)
	}

	ec2Zones, err := svc.DescribeAvailabilityZones(&ec2.DescribeAvailabilityZonesInput{
		Filters: []*ec2.Filter{
			&ec2.Filter{
				Name:   aws.String("state"),
				Values: []*string{aws.String("available")},
			},
		},
	})
	if err != nil {
		return []string{}, err
	}

	for _, zone := range ec2Zones.AvailabilityZones {
		zones = append(zones, *zone.ZoneName)
	}

	sort.Strings(zones)

	return zones, nil
}

func (e *EKS) validateAvailabilityZones() error {
	var result error

	zones, err := e.getAvailablityZoneByRegion()
	if err != nil {
		return err
	}

	if len(zones) == 0 {
		return fmt.Errorf(
			"no availability zone found for region '%s'",
			e.Region(),
		)
	}

	availabilityZones := e.AvailabilityZones()

	for _, zoneConfigured := range availabilityZones {
		found := false
		for _, zone := range zones {
			if zone != "" && zone == zoneConfigured {
				found = true
				break
			}
		}
		if !found {
			result = multierror.Append(result, fmt.Errorf(
				"specified invalid availability zone '%s' for region '%s'",
				zoneConfigured,
				e.Region(),
			))
		}
	}
	if result != nil {
		return result
	}

	if len(availabilityZones) == 0 {
		zone := zones[0]
		if zone == "" {
			return fmt.Errorf("error determining availabilty zone")
		}
		e.log.Debugf("no availability zones specified selecting zone: %s", zone)
		availabilityZones = []string{zone}
		e.availabilityZones = &availabilityZones
	}

	return nil
}

func (e *EKS) Session() (*session.Session, error) {

	// return cached session
	if e.session != nil {
		return e.session, nil
	}

	// use default config, if vault disabled
	if e.conf.EKS.VaultPath != "" {
		sess, err := e.vaultSession()
		if err != nil {
			return nil, err
		} else {
			e.session = sess
			return e.session, nil
		}
	}

	e.session = session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
		Profile:           e.conf.EKS.Profile,
	}))
	e.session.Config.Region = aws.String(e.Region())
	return e.session, nil
}

func (e *EKS) vaultSession() (*session.Session, error) {
	vaultClient, err := vault.NewClient(nil)
	if err != nil {
		return nil, err
	}

	// without vault token lookup vault token file
	if os.Getenv("VAULT_TOKEN") == "" {
		vaultToken, err := e.readVaultToken()
		if err != nil {
			e.log.Debug("failed to read vault token file: ", err)
		} else {
			vaultClient.SetToken(vaultToken)
		}
	}

	awsSecret, err := vaultClient.Logical().Read(e.conf.EKS.VaultPath)
	if err != nil {
		return nil, err
	}
	if awsSecret == nil || awsSecret.Data == nil {
		return nil, fmt.Errorf("vault did not return data at path '%s'", e.conf.EKS.VaultPath)
	}

	values := []string{}

	for _, key := range []string{"access_key", "secret_key", "security_token"} {
		val, ok := awsSecret.Data[key]
		if !ok {
			return nil, fmt.Errorf("vault did not return data with key '%s'", key)
		}
		valString, ok := val.(string)
		if !ok {
			return nil, fmt.Errorf("vault did not return data with a string in key '%s'", key)
		}
		values = append(values, valString)
	}

	creds := credentials.NewStaticCredentials(values[0], values[1], values[2])

	sess := session.Must(session.NewSession())
	sess.Config.Region = aws.String(e.Region())
	sess.Config.Credentials = creds

	return sess, nil
}

func (e *EKS) VerifyInstanceTypes(instancePools []interfaces.InstancePool) error {
	var result error

	svc, err := e.EC2()
	if err != nil {
		return err
	}

	for _, instance := range instancePools {
		instanceType, err := e.InstanceType(instance.Config().Size)
		if err != nil {
			return err
		}

		if err := e.verifyInstanceType(instanceType, instance.Zones(), svc); err != nil {
			result = multierror.Append(result, err)
		}
	}

	return result
}

func (e *EKS) verifyInstanceType(instanceType string, zones []string, svc EC2) error {
	var result error
	var available bool

	//Request offering, filter by given instance type
	request := &ec2.DescribeReservedInstancesOfferingsInput{
		InstanceTenancy:    aws.String("default"),
		IncludeMarketplace: aws.Bool(false),
		OfferingClass:      aws.String("standard"),
		OfferingType:       aws.String("No Upfront"),
		ProductDescription: aws.String("Linux/UNIX (EKS VPC)"),
		InstanceType:       aws.String(instanceType),
	}
	response, err := svc.DescribeReservedInstancesOfferings(request)
	if err != nil {
		return fmt.Errorf("error reaching aws to verify instance type %s: %v", instanceType, err)
	}

	//Loop through the given zones
	for _, zone := range zones {
		available = false

		//Loop through every offer given. Check the zone against the current looped zone.
		for _, offer := range response.ReservedInstancesOfferings {
			if offer.AvailabilityZone != nil && *offer.AvailabilityZone == zone {
				available = true
				break
			}
		}

		//Collect non matched zones
		if !available {
			result = multierror.Append(result, fmt.Errorf("availabilty zone %s not offered for type %s", zone, instanceType))
		}
	}

	return result
}

// This methods converts and possibly validates a generic instance type to a
// provider specifc
func (e *EKS) InstanceType(typeIn string) (typeOut string, err error) {
	if typeIn == clusterv1alpha1.InstancePoolSizeTiny {
		return "t2.nano", nil
	}
	if typeIn == clusterv1alpha1.InstancePoolSizeSmall {
		return "t2.medium", nil
	}
	if typeIn == clusterv1alpha1.InstancePoolSizeMedium {
		return "m4.large", nil
	}
	if typeIn == clusterv1alpha1.InstancePoolSizeLarge {
		return "m4.xlarge", nil
	}

	// TODO: Validate custom instance type here
	return typeIn, nil
}

// This methods converts and possibly validates a generic volume type to a
// provider specifc
func (e *EKS) VolumeType(typeIn string) (typeOut string, err error) {
	if typeIn == clusterv1alpha1.VolumeTypeHDD {
		return "st2", nil
	}
	if typeIn == clusterv1alpha1.VolumeTypeSSD {
		return "gp2", nil
	}
	// TODO: Validate custom instance type here
	return typeIn, nil
}
