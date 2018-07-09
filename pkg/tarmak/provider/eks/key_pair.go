// Copyright Jetstack Ltd. See LICENSE for details.
package eks

import (
	"crypto/md5"
	"crypto/rsa"
	"crypto/x509"
	"encoding/hex"
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/ec2"
	"golang.org/x/crypto/ssh"
)

func (e *EKS) KeyName() string {
	if e.conf.EKS.KeyName == "" {
		return fmt.Sprintf("tarmak_%s", e.tarmak.Cluster().Environment().Name())
	}
	return e.conf.EKS.KeyName
}

func fingerprintAWSStyle(signer interface{}) (string, error) {
	switch v := signer.(type) {
	case *rsa.PrivateKey:
		pubKeyBytes, err := x509.MarshalPKIXPublicKey(v.Public())
		if err != nil {
			return "", err
		}
		md5sum := md5.Sum(pubKeyBytes)
		hexarray := make([]string, len(md5sum))
		for i, c := range md5sum {
			hexarray[i] = hex.EncodeToString([]byte{c})
		}
		return strings.Join(hexarray, ":"), nil
	default:
		return "", fmt.Errorf("unsupported key type %t", v)
	}
}

func (e *EKS) validateAWSKeyPair() error {
	svc, err := e.EC2()
	if err != nil {
		return err
	}

	keypairs, err := svc.DescribeKeyPairs(&ec2.DescribeKeyPairsInput{
		KeyNames: []*string{aws.String(e.KeyName())},
	})
	if err != nil {
		if aerr, ok := err.(awserr.Error); !ok || aerr.Code() != "InvalidKeyPair.NotFound" {
			return err
		}
	}

	var awsKeyPair *ec2.KeyPairInfo
	if len(keypairs.KeyPairs) == 0 {
		signer, err := ssh.NewSignerFromKey(e.tarmak.Cluster().Environment().SSHPrivateKey())
		if err != nil {
			return fmt.Errorf("unable to generate public key from private key: %s", err)
		}
		_, err = svc.ImportKeyPair(&ec2.ImportKeyPairInput{
			KeyName:           aws.String(e.KeyName()),
			PublicKeyMaterial: []byte(ssh.MarshalAuthorizedKey(signer.PublicKey())),
		})
		if err != nil {
			return err
		}
		return nil
	} else if len(keypairs.KeyPairs) != 1 {
		return fmt.Errorf("unexpected number of keypairs found: %d", len(keypairs.KeyPairs))
	} else {
		awsKeyPair = keypairs.KeyPairs[0]
	}

	if err != nil {
		return fmt.Errorf("failed to parse private key: %s", err)
	}

	// warn if cannot generate fingerprint, fail if fingerprints are not matching
	fingerprintExpected, err := fingerprintAWSStyle(e.tarmak.Cluster().Environment().SSHPrivateKey())
	if err != nil {
		e.log.Warn("failed to generate local fingerprint: ", err)
	} else if act, exp := *awsKeyPair.KeyFingerprint, fingerprintExpected; act != exp {
		return fmt.Errorf("AWS key pair does not match the local key pair, aws_fingerprint=%s local_fingerprint=%s", act, exp)
	}

	return nil
}
