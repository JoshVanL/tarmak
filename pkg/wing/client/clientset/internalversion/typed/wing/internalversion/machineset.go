// Copyright Jetstack Ltd. See LICENSE for details.
package internalversion

import (
	"time"

	wing "github.com/jetstack/tarmak/pkg/apis/wing/v1alpha1"
	scheme "github.com/jetstack/tarmak/pkg/wing/client/clientset/internalversion/scheme"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	rest "k8s.io/client-go/rest"
)

// MachineSetsGetter has a method to return a MachineSetInterface.
// A group's client should implement this interface.
type MachineSetsGetter interface {
	MachineSets(namespace string) MachineSetInterface
}

// MachineSetInterface has methods to work with MachineSet resources.
type MachineSetInterface interface {
	Create(*wing.MachineSet) (*wing.MachineSet, error)
	Update(*wing.MachineSet) (*wing.MachineSet, error)
	Delete(name string, options *v1.DeleteOptions) error
	DeleteCollection(options *v1.DeleteOptions, listOptions v1.ListOptions) error
	Get(name string, options v1.GetOptions) (*wing.MachineSet, error)
	List(opts v1.ListOptions) (*wing.MachineSetList, error)
	Watch(opts v1.ListOptions) (watch.Interface, error)
	Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *wing.MachineSet, err error)
	MachineSetExpansion
}

// machinesets implements MachineSetInterface
type machinesets struct {
	client rest.Interface
	ns     string
}

// newMachineSets returns a MachineSets
func newMachineSets(c *WingClient, namespace string) *machinesets {
	return &machinesets{
		client: c.RESTClient(),
		ns:     namespace,
	}
}

// Get takes name of the machineset, and returns the corresponding machineset object, and an error if there is any.
func (c *machinesets) Get(name string, options v1.GetOptions) (result *wing.MachineSet, err error) {
	result = &wing.MachineSet{}
	err = c.client.Get().
		Resource("machinesets").
		Name(name).
		VersionedParams(&options, scheme.ParameterCodec).
		Do().
		Into(result)
	return
}

// List takes label and field selectors, and returns the list of MachineSets that match those selectors.
func (c *machinesets) List(opts v1.ListOptions) (result *wing.MachineSetList, err error) {
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	result = &wing.MachineSetList{}
	err = c.client.Get().
		Resource("machinesets").
		VersionedParams(&opts, scheme.ParameterCodec).
		Timeout(timeout).
		Do().
		Into(result)
	return
}

// Watch returns a watch.Interface that watches the requested machinesets.
func (c *machinesets) Watch(opts v1.ListOptions) (watch.Interface, error) {
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	opts.Watch = true
	return c.client.Get().
		Resource("machinesets").
		VersionedParams(&opts, scheme.ParameterCodec).
		Timeout(timeout).
		Watch()
}

// Create takes the representation of a machineset and creates it.  Returns the server's representation of the machineset, and an error, if there is any.
func (c *machinesets) Create(machineset *wing.MachineSet) (result *wing.MachineSet, err error) {
	result = &wing.MachineSet{}
	err = c.client.Post().
		Resource("machinesets").
		Body(machineset).
		Do().
		Into(result)
	return
}

// Update takes the representation of a machineset and updates it. Returns the server's representation of the machineset, and an error, if there is any.
func (c *machinesets) Update(machineset *wing.MachineSet) (result *wing.MachineSet, err error) {
	result = &wing.MachineSet{}
	err = c.client.Put().
		Resource("machinesets").
		Name(machineset.Name).
		Body(machineset).
		Do().
		Into(result)
	return
}

// Delete takes name of the machineset and deletes it. Returns an error if one occurs.
func (c *machinesets) Delete(name string, options *v1.DeleteOptions) error {
	return c.client.Delete().
		Resource("machinesets").
		Name(name).
		Body(options).
		Do().
		Error()
}

// DeleteCollection deletes a collection of objects.
func (c *machinesets) DeleteCollection(options *v1.DeleteOptions, listOptions v1.ListOptions) error {
	var timeout time.Duration
	if listOptions.TimeoutSeconds != nil {
		timeout = time.Duration(*listOptions.TimeoutSeconds) * time.Second
	}
	return c.client.Delete().
		Resource("machinesets").
		VersionedParams(&listOptions, scheme.ParameterCodec).
		Timeout(timeout).
		Body(options).
		Do().
		Error()
}

// Patch applies the patch and returns the patched machineset.
func (c *machinesets) Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *wing.MachineSet, err error) {
	result = &wing.MachineSet{}
	err = c.client.Patch(pt).
		Resource("machinesets").
		SubResource(subresources...).
		Name(name).
		Body(data).
		Do().
		Into(result)
	return
}
