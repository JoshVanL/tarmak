// Copyright Jetstack Ltd. See LICENSE for details.
package internalversion

import (
	"github.com/jetstack/tarmak/pkg/wing/client/clientset/internalversion/scheme"
	rest "k8s.io/client-go/rest"
)

type WingInterface interface {
	RESTClient() rest.Interface
	MachinesGetter
	MachineSetsGetter
	MachineDeploymentsGetter
}

// WingClient is used to interact with features provided by the wing.k8s.io group.
type WingClient struct {
	restClient rest.Interface
}

func (c *WingClient) Machines(namespace string) MachineInterface {
	return newMachines(c, namespace)
}

func (c *WingClient) MachineSets(namespace string) MachineSetInterface {
	return newMachineSets(c, namespace)
}

func (c *WingClient) MachineDeployments(namespace string) MachineDeploymentInterface {
	return newMachineDeployments(c, namespace)
}

// NewForConfig creates a new WingClient for the given config.
func NewForConfig(c *rest.Config) (*WingClient, error) {
	config := *c
	if err := setConfigDefaults(&config); err != nil {
		return nil, err
	}
	client, err := rest.RESTClientFor(&config)
	if err != nil {
		return nil, err
	}
	return &WingClient{client}, nil
}

// NewForConfigOrDie creates a new WingClient for the given config and
// panics if there is an error in the config.
func NewForConfigOrDie(c *rest.Config) *WingClient {
	client, err := NewForConfig(c)
	if err != nil {
		panic(err)
	}
	return client
}

// New creates a new WingClient for the given RESTClient.
func New(c rest.Interface) *WingClient {
	return &WingClient{c}
}

func setConfigDefaults(config *rest.Config) error {
	config.APIPath = "/apis"
	if config.UserAgent == "" {
		config.UserAgent = rest.DefaultKubernetesUserAgent()
	}
	if config.GroupVersion == nil || config.GroupVersion.Group != scheme.Scheme.PrioritizedVersionsForGroup("wing.k8s.io")[0].Group {
		gv := scheme.Scheme.PrioritizedVersionsForGroup("wing.k8s.io")[0]
		config.GroupVersion = &gv
	}
	config.NegotiatedSerializer = scheme.Codecs

	if config.QPS == 0 {
		config.QPS = 5
	}
	if config.Burst == 0 {
		config.Burst = 10
	}

	return nil
}

// RESTClient returns a RESTClient that is used to communicate
// with API server by this client implementation.
func (c *WingClient) RESTClient() rest.Interface {
	if c == nil {
		return nil
	}
	return c.restClient
}
