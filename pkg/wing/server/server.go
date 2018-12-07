// Copyright Jetstack Ltd. See LICENSE for details.

package server

import (
	"fmt"
	"io"
	"net"

	"github.com/spf13/cobra"
	utilerrors "k8s.io/apimachinery/pkg/util/errors"
	"k8s.io/apiserver/pkg/admission"
	genericapiserver "k8s.io/apiserver/pkg/server"
	genericoptions "k8s.io/apiserver/pkg/server/options"

	"github.com/jetstack/tarmak/pkg/apis/wing/v1alpha1"
	"github.com/jetstack/tarmak/pkg/wing/admission/plugin/banmachine"
	"github.com/jetstack/tarmak/pkg/wing/admission/winginitializer"
	"github.com/jetstack/tarmak/pkg/wing/apiserver"
	clientset "github.com/jetstack/tarmak/pkg/wing/client/clientset/internalversion"
	informers "github.com/jetstack/tarmak/pkg/wing/client/informers/internalversion"
)

const defaultEtcdPathPrefix = "/registry/wing.tarmak.io"

type WingServerOptions struct {
	RecommendedOptions    *genericoptions.RecommendedOptions
	SharedInformerFactory informers.SharedInformerFactory
	Admission             *genericoptions.AdmissionOptions

	StdOut io.Writer
	StdErr io.Writer
}

//var defaultAdmissionControllers = []string{banmachine.PluginName}

func NewWingServerOptions(out, errOut io.Writer) *WingServerOptions {
	o := &WingServerOptions{
		RecommendedOptions: genericoptions.NewRecommendedOptions(defaultEtcdPathPrefix, apiserver.Codecs.LegacyCodec(v1alpha1.SchemeGroupVersion)),
		Admission:          genericoptions.NewAdmissionOptions(),

		StdOut: out,
		StdErr: errOut,
	}

	return o
}

// NewCommandStartMaster provides a CLI handler for 'start master' command
func NewCommandStartWingServer(out, errOut io.Writer, stopCh <-chan struct{}) *cobra.Command {
	o := NewWingServerOptions(out, errOut)
	banmachine.Register(o.Admission.Plugins)
	//o.Admission.PluginNames = defaultAdmissionControllers

	cmd := &cobra.Command{
		Short: "Launch a wing API server",
		Long:  "Launch a wing API server",
		RunE: func(c *cobra.Command, args []string) error {
			if err := o.Complete(); err != nil {
				return err
			}
			if err := o.Validate(args); err != nil {
				return err
			}
			if err := o.RunWingServer(stopCh); err != nil {
				return err
			}
			return nil
		},
	}

	flags := cmd.Flags()
	o.RecommendedOptions.Etcd.AddFlags(flags)
	o.RecommendedOptions.SecureServing.AddFlags(flags)
	o.Admission.AddFlags(flags)

	return cmd
}

func (o *WingServerOptions) Validate(args []string) error {
	errors := []error{}
	errors = append(errors, o.RecommendedOptions.Validate()...)
	errors = append(errors, o.Admission.Validate()...)
	return utilerrors.NewAggregate(errors)
}

func (o *WingServerOptions) Complete() error {
	return nil
}

//func (o WingServerOptions) Config() (*apiserver.Config, error) {
//	// TODO have a "real" external address
//	if err := o.RecommendedOptions.SecureServing.MaybeDefaultWithSelfSignedCerts("localhost", nil, []net.IP{net.ParseIP("127.0.0.1")}); err != nil {
//		return nil, fmt.Errorf("error creating self-signed certificates: %v", err)
//	}
//
//	serverConfig := genericapiserver.NewRecommendedConfig(apiserver.Codecs)
//	if err := o.RecommendedOptions.Etcd.ApplyTo(&serverConfig.Config); err != nil {
//		return nil, err
//	}
//	if err := o.RecommendedOptions.SecureServing.ApplyTo(&serverConfig.Config); err != nil {
//		return nil, err
//	}
//
//	client, err := clientset.NewForConfig(serverConfig.LoopbackClientConfig)
//	if err != nil {
//		return nil, err
//	}
//	informerFactory := informers.NewSharedInformerFactory(client, serverConfig.LoopbackClientConfig.Timeout)
//	admissionInitializer, err := winginitializer.New(informerFactory)
//	if err != nil {
//		return nil, err
//	}
//
//	if err := o.Admission.ApplyTo(&serverConfig.Config, serverConfig.SharedInformerFactory, serverConfig.LoopbackClientConfig, apiserver.Scheme, admissionInitializer); err != nil {
//		return nil, err
//	}
//
//	config := &apiserver.Config{
//		GenericConfig: serverConfig,
//	}
//
//	return config, nil
//}

func (o *WingServerOptions) Config() (*apiserver.Config, error) {
	// TODO have a "real" external address
	if err := o.RecommendedOptions.SecureServing.MaybeDefaultWithSelfSignedCerts("localhost", nil, []net.IP{net.ParseIP("127.0.0.1")}); err != nil {
		return nil, fmt.Errorf("error creating self-signed certificates: %v", err)
	}

	o.RecommendedOptions.ExtraAdmissionInitializers = func(c *genericapiserver.RecommendedConfig) ([]admission.PluginInitializer, error) {
		client, err := clientset.NewForConfig(c.LoopbackClientConfig)
		if err != nil {
			return nil, err
		}
		informerFactory := informers.NewSharedInformerFactory(client, c.LoopbackClientConfig.Timeout)
		o.SharedInformerFactory = informerFactory
		return []admission.PluginInitializer{winginitializer.New(informerFactory)}, nil
	}

	serverConfig := genericapiserver.NewRecommendedConfig(apiserver.Codecs)
	if err := o.RecommendedOptions.ApplyTo(serverConfig, apiserver.Scheme); err != nil {
		return nil, err
	}

	config := &apiserver.Config{
		GenericConfig: serverConfig,
		ExtraConfig:   apiserver.ExtraConfig{},
	}
	return config, nil
}

func (o WingServerOptions) RunWingServer(stopCh <-chan struct{}) error {
	config, err := o.Config()
	if err != nil {
		return err
	}

	server, err := config.Complete().New()
	if err != nil {
		return err
	}

	/*server.GenericAPIServer.AddPostStartHook("start-sample-server-informers", func(context genericapiserver.PostStartHookContext) error {
		config.GenericConfig.SharedInformerFactory.Start(context.StopCh)
		return nil
	})
	*/

	return server.GenericAPIServer.PrepareRun().Run(stopCh)
}
