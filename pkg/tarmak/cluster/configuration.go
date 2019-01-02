// Copyright Jetstack Ltd. See LICENSE for details.
package cluster

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"time"

	wingv1alpha1 "github.com/jetstack/tarmak/pkg/apis/wing/v1alpha1"
)

const (
	retries = 100
)

// This upload the puppet.tar.gz to the cluster, warning there is some duplication as terraform is also uploading this puppet.tar.gz
func (c *Cluster) UploadConfiguration() error {

	buffer := new(bytes.Buffer)

	// get puppet config
	err := c.Environment().Tarmak().Puppet().TarGz(buffer)
	if err != nil {
		return err
	}

	// build reader from config
	reader := bytes.NewReader(buffer.Bytes())
	hasher := md5.New()
	hasher.Write(buffer.Bytes())

	return c.Environment().Provider().UploadConfiguration(
		c,
		reader,
		hex.EncodeToString(hasher.Sum(nil)),
	)
}

// This enforces a reapply of the puppet.tar.gz on every instance in the cluster
func (c *Cluster) ReapplyConfiguration() error {
	c.log.Infof("making sure all instances apply the latest manifest")

	// connect to wing
	client, err := c.wingMachineClient()
	if err != nil {
		return fmt.Errorf("failed to connect to wing API on bastion: %s", err)
	}

	// list instances
	instances, err := c.listMachines()
	if err != nil {
		return fmt.Errorf("failed to list instances: %s", err)
	}

	for pos, _ := range instances {
		instance := instances[pos]
		if instance.Spec == nil {
			instance.Spec = &wingv1alpha1.MachineSpec{}
		}
		instance.Spec.Converge = &wingv1alpha1.MachineSpecManifest{}

		if _, err := client.Update(instance); err != nil {
			c.log.Warnf("error updating instance %s in wing API: %s", instance.Name, err)
		}
	}

	// TODO: solve this on the API server side
	time.Sleep(time.Second * 5)

	return nil
}

// This waits until all instances have congverged successfully
func (c *Cluster) WaitForConvergance() error {
	c.log.Debugf("making sure all instances have converged using puppet")

	retries := retries
	for {
		instances, err := c.listMachines()
		if err != nil {
			return fmt.Errorf("failed to list instances: %s", err)
		}

		instanceByState := make(map[wingv1alpha1.MachineManifestState][]*wingv1alpha1.Machine)

		for pos, _ := range instances {
			instance := instances[pos]

			// index by instance convergance state
			if instance.Status == nil || instance.Status.Converge == nil || instance.Status.Converge.State == "" {
				continue
			}

			state := instance.Status.Converge.State
			if _, ok := instanceByState[state]; !ok {
				instanceByState[state] = []*wingv1alpha1.Machine{}
			}

			instanceByState[state] = append(
				instanceByState[state],
				instance,
			)
		}

		err = c.checkAllMachinesConverged(instanceByState)
		if err == nil {
			c.log.Info("all instances converged")
			return nil
		} else {
			c.log.Debug(err)
		}

		retries--
		if retries == 0 {
			break
		}
		time.Sleep(time.Second * 5)

	}

	return fmt.Errorf("instances failed to converge in time")
}
