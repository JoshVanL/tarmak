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

// This enforces a reapply of the puppet.tar.gz on every machine in the cluster
func (c *Cluster) ReapplyConfiguration() error {
	c.log.Infof("making sure all machines apply the latest manifest")

	if err := c.deleteUnusedMachines(); err != nil {
		return err
	}

	if err := c.updateMachineDeployments(); err != nil {
		return err
	}

	client, err := c.wingMachineClient()
	if err != nil {
		return err
	}

	// here we need to start the mechanism to trigger a re-converge
	machines, err := c.listMachines()
	if err != nil {
		return fmt.Errorf("failed to list machines: %s", err)
	}

	for pos, _ := range machines {
		machine := machines[pos]
		if machine.Spec == nil {
			machine.Spec = &wingv1alpha1.MachineSpec{}
		}
		machine.Spec.Converge = &wingv1alpha1.MachineSpecManifest{}

		if _, err := client.Update(machine); err != nil {
			c.log.Warnf("error updating machine %s in wing API: %s", machine.Name, err)
		}
	}

	// TODO: solve this on the API server side
	time.Sleep(time.Second * 5)

	return nil
}

// This waits until all machines have congverged successfully
func (c *Cluster) WaitForConvergance() error {
	c.log.Debugf("making sure all machine have converged using puppet")

	retries := retries
	for {
		deployments, err := c.listMachineDeployments()
		if err != nil {
			return fmt.Errorf("failed to list machines: %s", err)
		}

		var converged []*wingv1alpha1.MachineDeployment
		var converging []*wingv1alpha1.MachineDeployment
		for pos, _ := range deployments {
			deployment := deployments[pos]

			if deployment.Status == nil {
				converging = append(converging, deployment)
				continue
			}

			if deployment.Status.ReadyReplicas >= deployment.Status.Replicas &&
				deployment.Status.ReadyReplicas >= *deployment.Spec.MinReplicas {
				converged = append(converged, deployment)
				continue
			}

			converging = append(converging, deployment)
		}

		var convergedStr string
		for _, d := range converged {
			convergedStr = fmt.Sprintf("%s %s", convergedStr, d.Name)
		}
		convergedStr = fmt.Sprintf("converged deployments [%s]", convergedStr)

		if len(converging) == 0 {
			c.log.Info("all deployments converged")
			c.log.Info(convergedStr)
			return nil
		}

		c.log.Debug("--------")
		if len(converged) > 0 {
			c.log.Debug(convergedStr)
		}

		for _, d := range converging {
			var readyReplicas int32
			if d.Status != nil {
				readyReplicas = d.Status.Replicas
			}
			c.log.Debugf("converging %s [%v/%v]", d.Name, readyReplicas, d.Status.Replicas)
		}

		retries--
		if retries == 0 {
			break
		}

		tok := time.Tick(time.Second * 5)

		select {
		case <-c.ctx.Done():
			return c.ctx.Err()
		case <-tok:
		}
	}

	return fmt.Errorf("machines failed to converge in time")
}
