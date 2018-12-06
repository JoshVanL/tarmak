// Copyright Jetstack Ltd. See LICENSE for details.
package v1alpha1

type MachineManifestState string
type MachineDeploymentStrategyType string

const (
	MachineManifestStateConverging = MachineManifestState("converging")
	MachineManifestStateConverged  = MachineManifestState("converged")
	MachineManifestStateError      = MachineManifestState("error")
)

const (
	// Replace the old MachineSet by new one using rolling update
	// i.e gradually scale down the old MachineSet and scale up the new one.
	RollingUpdateMachineDeploymentStrategyType MachineDeploymentStrategyType = "RollingUpdate"
)
