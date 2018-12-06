// Copyright Jetstack Ltd. See LICENSE for details.
package v1alpha1

import (
	v1alpha1 "github.com/jetstack/tarmak/pkg/apis/wing/v1alpha1"
	scheme "github.com/jetstack/tarmak/pkg/wing/client/clientset/versioned/scheme"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	rest "k8s.io/client-go/rest"
)

// MachinesGetter has a method to return a MachineInterface.
// A group's client should implement this interface.
type MachinesGetter interface {
	Machines(namespace string) MachineInterface
}

// MachineInterface has methods to work with Machine resources.
type MachineInterface interface {
	Create(*v1alpha1.Machine) (*v1alpha1.Machine, error)
	Update(*v1alpha1.Machine) (*v1alpha1.Machine, error)
	UpdateStatus(*v1alpha1.Machine) (*v1alpha1.Machine, error)
	Delete(name string, options *v1.DeleteOptions) error
	DeleteCollection(options *v1.DeleteOptions, listOptions v1.ListOptions) error
	Get(name string, options v1.GetOptions) (*v1alpha1.Machine, error)
	List(opts v1.ListOptions) (*v1alpha1.MachineList, error)
	Watch(opts v1.ListOptions) (watch.Interface, error)
	Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *v1alpha1.Machine, err error)
	MachineExpansion
}

// machines implements MachineInterface
type machines struct {
	client rest.Interface
	ns     string
}

// newMachines returns a Machines
func newMachines(c *WingV1alpha1Client, namespace string) *machines {
	return &machines{
		client: c.RESTClient(),
		ns:     namespace,
	}
}

// Get takes name of the machine, and returns the corresponding machine object, and an error if there is any.
func (c *machines) Get(name string, options v1.GetOptions) (result *v1alpha1.Machine, err error) {
	result = &v1alpha1.Machine{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource("machines").
		Name(name).
		VersionedParams(&options, scheme.ParameterCodec).
		Do().
		Into(result)
	return
}

// List takes label and field selectors, and returns the list of Machines that match those selectors.
func (c *machines) List(opts v1.ListOptions) (result *v1alpha1.MachineList, err error) {
	result = &v1alpha1.MachineList{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource("machines").
		VersionedParams(&opts, scheme.ParameterCodec).
		Do().
		Into(result)
	return
}

// Watch returns a watch.Interface that watches the requested machines.
func (c *machines) Watch(opts v1.ListOptions) (watch.Interface, error) {
	opts.Watch = true
	return c.client.Get().
		Namespace(c.ns).
		Resource("machines").
		VersionedParams(&opts, scheme.ParameterCodec).
		Watch()
}

// Create takes the representation of a machine and creates it.  Returns the server's representation of the machine, and an error, if there is any.
func (c *machines) Create(machine *v1alpha1.Machine) (result *v1alpha1.Machine, err error) {
	result = &v1alpha1.Machine{}
	err = c.client.Post().
		Namespace(c.ns).
		Resource("machines").
		Body(machine).
		Do().
		Into(result)
	return
}

// Update takes the representation of a machine and updates it. Returns the server's representation of the machine, and an error, if there is any.
func (c *machines) Update(machine *v1alpha1.Machine) (result *v1alpha1.Machine, err error) {
	result = &v1alpha1.Machine{}
	err = c.client.Put().
		Namespace(c.ns).
		Resource("machines").
		Name(machine.Name).
		Body(machine).
		Do().
		Into(result)
	return
}

// UpdateStatus was generated because the type contains a Status member.
// Add a +genclient:noStatus comment above the type to avoid generating UpdateStatus().

func (c *machines) UpdateStatus(machine *v1alpha1.Machine) (result *v1alpha1.Machine, err error) {
	result = &v1alpha1.Machine{}
	err = c.client.Put().
		Namespace(c.ns).
		Resource("machines").
		Name(machine.Name).
		SubResource("status").
		Body(machine).
		Do().
		Into(result)
	return
}

// Delete takes name of the machine and deletes it. Returns an error if one occurs.
func (c *machines) Delete(name string, options *v1.DeleteOptions) error {
	return c.client.Delete().
		Namespace(c.ns).
		Resource("machines").
		Name(name).
		Body(options).
		Do().
		Error()
}

// DeleteCollection deletes a collection of objects.
func (c *machines) DeleteCollection(options *v1.DeleteOptions, listOptions v1.ListOptions) error {
	return c.client.Delete().
		Namespace(c.ns).
		Resource("machines").
		VersionedParams(&listOptions, scheme.ParameterCodec).
		Body(options).
		Do().
		Error()
}

// Patch applies the patch and returns the patched machine.
func (c *machines) Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *v1alpha1.Machine, err error) {
	result = &v1alpha1.Machine{}
	err = c.client.Patch(pt).
		Namespace(c.ns).
		Resource("machines").
		SubResource(subresources...).
		Name(name).
		Body(data).
		Do().
		Into(result)
	return
}
