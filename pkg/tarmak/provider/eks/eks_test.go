// Copyright Jetstack Ltd. See LICENSE for details.
package eks

import (
	"reflect"
	"strings"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/golang/mock/gomock"
	"github.com/sirupsen/logrus"

	clusterv1alpha1 "github.com/jetstack/tarmak/pkg/apis/cluster/v1alpha1"
	tarmakv1alpha1 "github.com/jetstack/tarmak/pkg/apis/tarmak/v1alpha1"
	"github.com/jetstack/tarmak/pkg/tarmak/mocks"
)

type fakeEKS struct {
	*EKS
	ctrl *gomock.Controller

	fakeEC2         *mocks.MockEC2
	fakeEnvironment *mocks.MockEnvironment
	fakeCluster     *mocks.MockCluster
	fakeTarmak      *mocks.MockTarmak
}

func newFakeEKS(t *testing.T) *fakeEKS {

	f := &fakeEKS{
		ctrl: gomock.NewController(t),
		EKS: &EKS{
			conf: &tarmakv1alpha1.Provider{
				EKS: &tarmakv1alpha1.ProviderEKS{
					KeyName: "myfake_key",
				},
			},
			log: logrus.WithField("test", true),
		},
	}
	f.fakeEC2 = mocks.NewMockEC2(f.ctrl)
	f.fakeEnvironment = mocks.NewMockEnvironment(f.ctrl)
	f.fakeCluster = mocks.NewMockCluster(f.ctrl)
	f.fakeTarmak = mocks.NewMockTarmak(f.ctrl)
	f.EKS.ec2 = f.fakeEC2
	f.EKS.tarmak = f.fakeTarmak
	f.fakeTarmak.EXPECT().Cluster().AnyTimes().Return(f.fakeCluster)
	f.fakeTarmak.EXPECT().Environment().AnyTimes().Return(f.fakeEnvironment)
	f.fakeCluster.EXPECT().Environment().AnyTimes().Return(f.fakeEnvironment)

	return f
}

func TestEKS_validateAvailabilityZonesNoneGiven(t *testing.T) {
	a := newFakeEKS(t)
	defer a.ctrl.Finish()

	a.fakeCluster.EXPECT().Subnets().Return([]clusterv1alpha1.Subnet{}).MinTimes(1)
	a.fakeCluster.EXPECT().Region().Return("london-north-1").AnyTimes()

	a.fakeEC2.EXPECT().DescribeAvailabilityZones(gomock.Any()).Return(&ec2.DescribeAvailabilityZonesOutput{
		AvailabilityZones: []*ec2.AvailabilityZone{
			&ec2.AvailabilityZone{
				ZoneName:   aws.String("london-north-1a"),
				State:      aws.String("available"),
				RegionName: aws.String("london-north-1"),
			},
			&ec2.AvailabilityZone{
				ZoneName:   aws.String("london-north-1b"),
				State:      aws.String("available"),
				RegionName: aws.String("london-north-1"),
			},
			&ec2.AvailabilityZone{
				ZoneName:   aws.String("london-north-1c"),
				State:      aws.String("available"),
				RegionName: aws.String("london-north-1"),
			},
		},
	}, nil)

	err := a.validateAvailabilityZones()
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}

	if act, exp := a.AvailabilityZones(), []string{"london-north-1a"}; !reflect.DeepEqual(act, exp) {
		t.Errorf("unexpected availability zones: act=%+v exp=%+v", act, exp)
	}
}

func TestEKS_validateAvailabilityZonesCorrectGiven(t *testing.T) {
	a := newFakeEKS(t)
	defer a.ctrl.Finish()

	a.fakeCluster.EXPECT().Subnets().Return([]clusterv1alpha1.Subnet{
		clusterv1alpha1.Subnet{
			Zone: "london-north-1b",
		},
		clusterv1alpha1.Subnet{
			Zone: "london-north-1c",
		},
	}).MinTimes(1)
	a.fakeCluster.EXPECT().Region().Return("london-north-1").AnyTimes()

	a.fakeEC2.EXPECT().DescribeAvailabilityZones(gomock.Any()).Return(&ec2.DescribeAvailabilityZonesOutput{
		AvailabilityZones: []*ec2.AvailabilityZone{
			&ec2.AvailabilityZone{
				ZoneName:   aws.String("london-north-1a"),
				State:      aws.String("available"),
				RegionName: aws.String("london-north-1"),
			},
			&ec2.AvailabilityZone{
				ZoneName:   aws.String("london-north-1b"),
				State:      aws.String("available"),
				RegionName: aws.String("london-north-1"),
			},
			&ec2.AvailabilityZone{
				ZoneName:   aws.String("london-north-1c"),
				State:      aws.String("available"),
				RegionName: aws.String("london-north-1"),
			},
		},
	}, nil)

	err := a.validateAvailabilityZones()
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}

	if act, exp := a.AvailabilityZones(), []string{"london-north-1b", "london-north-1c"}; !reflect.DeepEqual(act, exp) {
		t.Errorf("unexpected availability zones: act=%+v exp=%+v", act, exp)
	}
}

func TestEKS_validateAvailabilityZonesFalseGiven(t *testing.T) {
	a := newFakeEKS(t)
	defer a.ctrl.Finish()

	a.fakeCluster.EXPECT().Subnets().Return([]clusterv1alpha1.Subnet{
		clusterv1alpha1.Subnet{
			Zone: "london-north-1a",
		},
		clusterv1alpha1.Subnet{
			Zone: "london-north-1d",
		},
		clusterv1alpha1.Subnet{
			Zone: "london-north-1e",
		},
	}).MinTimes(1)
	a.fakeCluster.EXPECT().Region().Return("london-north-1").AnyTimes()
	a.fakeEnvironment.EXPECT().Location().Return("london-north-1").AnyTimes()

	a.fakeEC2.EXPECT().DescribeAvailabilityZones(gomock.Any()).Return(&ec2.DescribeAvailabilityZonesOutput{
		AvailabilityZones: []*ec2.AvailabilityZone{
			&ec2.AvailabilityZone{
				ZoneName:   aws.String("london-north-1a"),
				State:      aws.String("available"),
				RegionName: aws.String("london-north-1"),
			},
			&ec2.AvailabilityZone{
				ZoneName:   aws.String("london-north-1b"),
				State:      aws.String("available"),
				RegionName: aws.String("london-north-1"),
			},
			&ec2.AvailabilityZone{
				ZoneName:   aws.String("london-north-1c"),
				State:      aws.String("available"),
				RegionName: aws.String("london-north-1"),
			},
		},
	}, nil)

	err := a.validateAvailabilityZones()
	if err == nil {
		t.Error("expected an error")
	} else if !strings.Contains(err.Error(), "specified invalid availability zone") {
		t.Errorf("unexpected error messge: %s", err)
	}
}

func TestEKS_verifyInstanceTypeNoneGiven(t *testing.T) {
	a := newFakeEKS(t)
	defer a.ctrl.Finish()

	svc, err := a.EC2()
	if err != nil {
		t.Errorf("unexpected err:%v", err)
	}

	responce := &ec2.DescribeReservedInstancesOfferingsOutput{
		ReservedInstancesOfferings: []*ec2.ReservedInstancesOffering{
			&ec2.ReservedInstancesOffering{
				AvailabilityZone: aws.String("test-east-1a"),
			},
			&ec2.ReservedInstancesOffering{
				AvailabilityZone: aws.String("test-east-1b"),
			},
			&ec2.ReservedInstancesOffering{
				AvailabilityZone: aws.String("test-east-1c"),
			},
		},
	}

	a.fakeEC2.EXPECT().DescribeReservedInstancesOfferings(gomock.Any()).Return(responce, nil)

	err = a.verifyInstanceType("atype", []string{}, svc)
	if err != nil {
		t.Errorf("unexpected err:%v", err)
	}
}

func TestEKS_verifyInstanceTypeZonesGiven(t *testing.T) {
	a := newFakeEKS(t)
	defer a.ctrl.Finish()

	svc, err := a.EC2()
	if err != nil {
		t.Errorf("unexpected err:%v", err)
	}

	responce := &ec2.DescribeReservedInstancesOfferingsOutput{
		ReservedInstancesOfferings: []*ec2.ReservedInstancesOffering{
			&ec2.ReservedInstancesOffering{
				AvailabilityZone: aws.String("test-east-1a"),
			},
			&ec2.ReservedInstancesOffering{
				AvailabilityZone: aws.String("test-east-1b"),
			},
			&ec2.ReservedInstancesOffering{
				AvailabilityZone: aws.String("test-east-1c"),
			},
		},
	}

	a.fakeEC2.EXPECT().DescribeReservedInstancesOfferings(gomock.Any()).Return(responce, nil)

	err = a.verifyInstanceType("atype", []string{"test-east-1a", "test-east-1b", "test-east-1c"}, svc)
	if err != nil {
		t.Errorf("unexpected err:%v", err)
	}
}

func TestEKS_verifyInstanceTypeZonesGivenWrong(t *testing.T) {
	a := newFakeEKS(t)
	defer a.ctrl.Finish()

	svc, err := a.EC2()
	if err != nil {
		t.Errorf("unexpected err:%v", err)
	}

	responce := &ec2.DescribeReservedInstancesOfferingsOutput{
		ReservedInstancesOfferings: []*ec2.ReservedInstancesOffering{
			&ec2.ReservedInstancesOffering{
				AvailabilityZone: aws.String("test-wrong-1a"),
			},
			&ec2.ReservedInstancesOffering{
				AvailabilityZone: aws.String("test-east-1b"),
			},
			&ec2.ReservedInstancesOffering{
				AvailabilityZone: aws.String("test-east-1c"),
			},
		},
	}

	a.fakeEC2.EXPECT().DescribeReservedInstancesOfferings(gomock.Any()).Return(responce, nil)

	err = a.verifyInstanceType("atype", []string{"test-east-1a", "test-east-1b", "test-east-1c"}, svc)
	if err == nil {
		t.Errorf("expected an error but got none.")
	}
}

func TestEKS_verifyInstanceTypeZonesOneOne(t *testing.T) {
	a := newFakeEKS(t)
	defer a.ctrl.Finish()

	svc, err := a.EC2()
	if err != nil {
		t.Errorf("unexpected err:%v", err)
	}

	responce := &ec2.DescribeReservedInstancesOfferingsOutput{
		ReservedInstancesOfferings: []*ec2.ReservedInstancesOffering{
			&ec2.ReservedInstancesOffering{
				AvailabilityZone: aws.String("test-east-1a"),
			},
		},
	}

	a.fakeEC2.EXPECT().DescribeReservedInstancesOfferings(gomock.Any()).Return(responce, nil)

	err = a.verifyInstanceType("atype", []string{"test-east-1a"}, svc)
	if err != nil {
		t.Errorf("unexpected err:%v", err)
	}
}

func TestEKS_verifyInstanceTypeZonesOneMany(t *testing.T) {
	a := newFakeEKS(t)
	defer a.ctrl.Finish()

	svc, err := a.EC2()
	if err != nil {
		t.Errorf("unexpected err:%v", err)
	}

	responce := &ec2.DescribeReservedInstancesOfferingsOutput{
		ReservedInstancesOfferings: []*ec2.ReservedInstancesOffering{
			&ec2.ReservedInstancesOffering{
				AvailabilityZone: aws.String("test-east-1a"),
			},
			&ec2.ReservedInstancesOffering{
				AvailabilityZone: aws.String("test-east-1b"),
			},
		},
	}

	a.fakeEC2.EXPECT().DescribeReservedInstancesOfferings(gomock.Any()).Return(responce, nil)

	err = a.verifyInstanceType("atype", []string{"test-east-1a"}, svc)
	if err != nil {
		t.Errorf("unexpected err:%v", err)
	}
}

func TestEKS_verifyInstanceTypeZonesNotAllAvailable(t *testing.T) {
	a := newFakeEKS(t)
	defer a.ctrl.Finish()

	svc, err := a.EC2()
	if err != nil {
		t.Errorf("unexpected err:%v", err)
	}

	responce := &ec2.DescribeReservedInstancesOfferingsOutput{
		ReservedInstancesOfferings: []*ec2.ReservedInstancesOffering{
			&ec2.ReservedInstancesOffering{
				AvailabilityZone: aws.String("test-east-1a"),
			},
		},
	}

	a.fakeEC2.EXPECT().DescribeReservedInstancesOfferings(gomock.Any()).Return(responce, nil)

	err = a.verifyInstanceType("atype", []string{"test-east-1a", "test-east-1b"}, svc)
	if err == nil {
		t.Errorf("expected error but got none")
	}
}
