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

// MachineDeploymentsGetter has a method to return a MachineDeploymentInterface.
// A group's client should implement this interface.
type MachineDeploymentsGetter interface {
	MachineDeployments(namespace string) MachineDeploymentInterface
}

// MachineDeploymentInterface has methods to work with MachineDeployment resources.
type MachineDeploymentInterface interface {
	Create(*wing.MachineDeployment) (*wing.MachineDeployment, error)
	Update(*wing.MachineDeployment) (*wing.MachineDeployment, error)
	Delete(name string, options *v1.DeleteOptions) error
	DeleteCollection(options *v1.DeleteOptions, listOptions v1.ListOptions) error
	Get(name string, options v1.GetOptions) (*wing.MachineDeployment, error)
	List(opts v1.ListOptions) (*wing.MachineDeploymentList, error)
	Watch(opts v1.ListOptions) (watch.Interface, error)
	Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *wing.MachineDeployment, err error)
	MachineDeploymentExpansion
}

// machinedeployments implements MachineDeploymentInterface
type machinedeployments struct {
	client rest.Interface
	ns     string
}

// newMachineDeployments returns a MachineDeployments
func newMachineDeployments(c *WingClient, namespace string) *machinedeployments {
	return &machinedeployments{
		client: c.RESTClient(),
		ns:     namespace,
	}
}

// Get takes name of the machinedeployment, and returns the corresponding machinedeployment object, and an error if there is any.
func (c *machinedeployments) Get(name string, options v1.GetOptions) (result *wing.MachineDeployment, err error) {
	result = &wing.MachineDeployment{}
	err = c.client.Get().
		Resource("machinedeployments").
		Name(name).
		VersionedParams(&options, scheme.ParameterCodec).
		Do().
		Into(result)
	return
}

// List takes label and field selectors, and returns the list of MachineDeployments that match those selectors.
func (c *machinedeployments) List(opts v1.ListOptions) (result *wing.MachineDeploymentList, err error) {
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	result = &wing.MachineDeploymentList{}
	err = c.client.Get().
		Resource("machinedeployments").
		VersionedParams(&opts, scheme.ParameterCodec).
		Timeout(timeout).
		Do().
		Into(result)
	return
}

// Watch returns a watch.Interface that watches the requested machinedeployments.
func (c *machinedeployments) Watch(opts v1.ListOptions) (watch.Interface, error) {
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	opts.Watch = true
	return c.client.Get().
		Resource("machinedeployments").
		VersionedParams(&opts, scheme.ParameterCodec).
		Timeout(timeout).
		Watch()
}

// Create takes the representation of a machinedeployment and creates it.  Returns the server's representation of the machinedeployment, and an error, if there is any.
func (c *machinedeployments) Create(machinedeployment *wing.MachineDeployment) (result *wing.MachineDeployment, err error) {
	result = &wing.MachineDeployment{}
	err = c.client.Post().
		Resource("machinedeployments").
		Body(machinedeployment).
		Do().
		Into(result)
	return
}

// Update takes the representation of a machinedeployment and updates it. Returns the server's representation of the machinedeployment, and an error, if there is any.
func (c *machinedeployments) Update(machinedeployment *wing.MachineDeployment) (result *wing.MachineDeployment, err error) {
	result = &wing.MachineDeployment{}
	err = c.client.Put().
		Resource("machinedeployments").
		Name(machinedeployment.Name).
		Body(machinedeployment).
		Do().
		Into(result)
	return
}

// Delete takes name of the machinedeployment and deletes it. Returns an error if one occurs.
func (c *machinedeployments) Delete(name string, options *v1.DeleteOptions) error {
	return c.client.Delete().
		Resource("machinedeployments").
		Name(name).
		Body(options).
		Do().
		Error()
}

// DeleteCollection deletes a collection of objects.
func (c *machinedeployments) DeleteCollection(options *v1.DeleteOptions, listOptions v1.ListOptions) error {
	var timeout time.Duration
	if listOptions.TimeoutSeconds != nil {
		timeout = time.Duration(*listOptions.TimeoutSeconds) * time.Second
	}
	return c.client.Delete().
		Resource("machinedeployments").
		VersionedParams(&listOptions, scheme.ParameterCodec).
		Timeout(timeout).
		Body(options).
		Do().
		Error()
}

// Patch applies the patch and returns the patched machinedeployment.
func (c *machinedeployments) Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *wing.MachineDeployment, err error) {
	result = &wing.MachineDeployment{}
	err = c.client.Patch(pt).
		Resource("machinedeployments").
		SubResource(subresources...).
		Name(name).
		Body(data).
		Do().
		Into(result)
	return
}
