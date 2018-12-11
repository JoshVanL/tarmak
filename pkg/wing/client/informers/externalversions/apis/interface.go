// Copyright Jetstack Ltd. See LICENSE for details.

// This file was automatically generated by informer-gen

package wing

import (
	wing "github.com/jetstack/tarmak/pkg/wing/client/informers/externalversions/apis/wing"
	internalinterfaces "github.com/jetstack/tarmak/pkg/wing/client/informers/externalversions/internalinterfaces"
)

// Interface provides access to each of this group's versions.
type Interface interface {
	// Wing provides access to shared informers for resources in Wing.
	Wing() wing.Interface
}

type group struct {
	factory          internalinterfaces.SharedInformerFactory
	namespace        string
	tweakListOptions internalinterfaces.TweakListOptionsFunc
}

// New returns a new Interface.
func New(f internalinterfaces.SharedInformerFactory, namespace string, tweakListOptions internalinterfaces.TweakListOptionsFunc) Interface {
	return &group{factory: f, namespace: namespace, tweakListOptions: tweakListOptions}
}

// Wing returns a new wing.Interface.
func (g *group) Wing() wing.Interface {
	return wing.New(g.factory, g.namespace, g.tweakListOptions)
}