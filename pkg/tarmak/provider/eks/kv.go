// Copyright Jetstack Ltd. See LICENSE for details.
package eks

import (
	"fmt"

	"github.com/jetstack/vault-unsealer/pkg/kv"
	"github.com/jetstack/vault-unsealer/pkg/kv/aws_kms"
	"github.com/jetstack/vault-unsealer/pkg/kv/aws_ssm"
)

func (e *EKS) secretsKMSKeyID() (string, error) {
	output, err := e.tarmak.Cluster().Environment().Hub().TerraformOutput()
	if err != nil {
		return "", fmt.Errorf("error getting hub terraform output: %s", err)
	}

	key := "secrets_kms_arn"

	secretIDIntf, ok := output[key]
	if !ok {
		return "", fmt.Errorf("error could not find '%s' in terraform state output", key)
	}

	var secretID string
	switch v := secretIDIntf.(type) {
	// return a list (necessary for 0.11 terraform +
	case []interface{}:
		if len(v) < 1 {
			return "", fmt.Errorf("no list elements found for '%s'", key)
		}
		elem, ok := v[0].(string)
		if !ok {
			return "", fmt.Errorf("first element for '%s' is not a string", key)
		}
		secretID = elem

	case string:
		secretID = v

	default:
		return "", fmt.Errorf("error unexpected type for '%s': %T", key, secretIDIntf)
	}

	return secretID, nil
}

func (e *EKS) vaultUnsealKeyName() (string, error) {
	key := "vault_unseal_key_name"

	output, err := e.tarmak.Cluster().Environment().Hub().TerraformOutput()
	if err != nil {
		return "", fmt.Errorf("error getting hub terraform output: %s", err)
	}

	keyNameIntf, ok := output[key]
	if !ok {
		return "", fmt.Errorf("error could not find '%s' in terraform vault output", key)
	}

	keyName, ok := keyNameIntf.(string)
	if !ok {
		return "", fmt.Errorf("error unexpected type for '%s': %T", key, keyNameIntf)
	}

	return keyName, nil

}

func (e *EKS) VaultKV() (kv.Service, error) {
	session, err := e.Session()
	if err != nil {
		return nil, err
	}

	kmsKeyID, err := e.secretsKMSKeyID()
	if err != nil {
		return nil, err
	}

	unsealKeyName, err := e.vaultUnsealKeyName()
	if err != nil {
		return nil, err
	}

	ssm, err := aws_ssm.NewWithSession(session, unsealKeyName)
	if err != nil {
		return nil, fmt.Errorf("error creating EKS SSM kv store: %s", err.Error())
	}

	kms, err := aws_kms.NewWithSession(session, ssm, kmsKeyID)
	if err != nil {
		return nil, fmt.Errorf("error creating EKS KMS ID kv store: %s", err.Error())
	}

	return kms, nil
}

func (e *EKS) VaultKVWithParams(kmsKeyID, unsealKeyName string) (kv.Service, error) {
	session, err := e.Session()
	if err != nil {
		return nil, err
	}

	ssm, err := aws_ssm.NewWithSession(session, unsealKeyName)
	if err != nil {
		return nil, fmt.Errorf("error creating EKS SSM kv store: %s", err.Error())
	}

	kms, err := aws_kms.NewWithSession(session, ssm, kmsKeyID)
	if err != nil {
		return nil, fmt.Errorf("error creating EKS KMS ID kv store: %s", err.Error())
	}

	return kms, nil
}
