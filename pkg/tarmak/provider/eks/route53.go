// Copyright Jetstack Ltd. See LICENSE for details.
package eks

import (
	"errors"
	"fmt"
	"net"
	"reflect"
	"sort"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/route53"
)

func (e *EKS) PublicZone() string {
	return e.conf.EKS.PublicZone
}

// this removes an ending . in zone and converts it to lowercase
func normalizeZone(in string) string {
	return strings.ToLower(strings.TrimRight(in, "."))
}

func (e *EKS) initPublicZone() (*route53.HostedZone, error) {
	publicZone := normalizeZone(e.conf.EKS.PublicZone)
	if publicZone == "" {
		return nil, errors.New("no public zone given in provider config")
	}
	if e.conf.EKS.PublicHostedZoneID != "" {
		return nil, errors.New("can not auto create public zone as there is HostedZoneID given in provider config")
	}

	svc, err := e.Route53()
	if err != nil {
		return nil, err
	}

	result, err := svc.CreateHostedZone(&route53.CreateHostedZoneInput{
		CallerReference: aws.String(time.Now().Format(time.RFC3339)),
		Name:            aws.String(normalizeZone(publicZone)),
		HostedZoneConfig: &route53.HostedZoneConfig{
			Comment: aws.String("public zone for tarmak"),
		},
	})
	return result.HostedZone, err
}

func (e *EKS) validatePublicZone() error {
	svc, err := e.Route53()
	if err != nil {
		return err
	}

	publicZoneName := normalizeZone(e.conf.EKS.PublicZone)

	input := &route53.ListHostedZonesByNameInput{}
	if publicZoneName != "" {
		input.DNSName = aws.String(publicZoneName)
	}

	if hostedZoneID := e.conf.EKS.PublicHostedZoneID; hostedZoneID != "" {
		input.HostedZoneId = aws.String(hostedZoneID)
	}

	var zone *route53.HostedZone

	zonesResponse, err := svc.ListHostedZonesByName(input)
	if err != nil {
		return err
	}
	var zones []*route53.HostedZone
	for pos, _ := range zonesResponse.HostedZones {
		zone := zonesResponse.HostedZones[pos]
		if normalizeZone(*zone.Name) == publicZoneName {
			zones = append(zones, zone)
		}
	}

	if len(zones) > 1 {
		msg := "more than one matching zone found, "
		if input.HostedZoneId != nil {
			msg = fmt.Sprintf("%shostedZoneID = %s ", msg, *input.HostedZoneId)
		}
		if input.DNSName != nil {
			msg = fmt.Sprintf("%sdnsName = %s ", msg, *input.DNSName)
		}
		return errors.New(msg)
	} else if len(zones) == 0 {
		zone, err = e.initPublicZone()
		if err != nil {
			return err
		}
	} else {
		zone = zones[0]
	}

	// store hostedzone id
	if split := strings.Split(*zone.Id, "/"); len(split) < 2 {
		return fmt.Errorf("Unexpected ID %s", *zone.Id)
	} else {
		e.conf.EKS.PublicHostedZoneID = split[2]
	}

	// store zone information
	e.conf.EKS.PublicZone = normalizeZone(*zone.Name)

	// validate delegation
	zoneResult, err := svc.GetHostedZone(&route53.GetHostedZoneInput{Id: zone.Id})
	if err != nil {
		return fmt.Errorf("unabled to get zone with ID '%s': %s", *zone.Id, err)
	}

	zoneNameservers := make([]string, len(zoneResult.DelegationSet.NameServers))
	for pos, _ := range zoneResult.DelegationSet.NameServers {
		zoneNameservers[pos] = *zoneResult.DelegationSet.NameServers[pos]
	}

	notice := fmt.Sprintf("make sure the domain is delegated to these nameservers %+v", zoneNameservers)

	dnsResult, err := net.LookupNS(e.conf.EKS.PublicZone)
	if err != nil {
		return fmt.Errorf("error resolving NS records for %s (%s), %s", e.conf.EKS.PublicZone, err, notice)
	}

	dnsNameservers := make([]string, len(dnsResult))
	for pos, _ := range dnsResult {
		dnsNameservers[pos] = normalizeZone(dnsResult[pos].Host)
	}

	sort.Strings(dnsNameservers)
	sort.Strings(zoneNameservers)

	if !reflect.DeepEqual(dnsNameservers, zoneNameservers) {
		return fmt.Errorf("public root dns namesevers %v and zone nameservers %v mismatch", dnsNameservers, zoneNameservers)
	}

	return nil
}
