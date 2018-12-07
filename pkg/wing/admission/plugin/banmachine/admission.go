// Copyright Jetstack Ltd. See LICENSE for details.
package banmachine

import (
	"fmt"
	"io"

	//wing "github.com/jetstack/tarmak/pkg/apis/wing/v1alpha1"
	"github.com/jetstack/tarmak/pkg/wing/admission/winginitializer"
	informers "github.com/jetstack/tarmak/pkg/wing/client/informers/internalversion"
	listers "github.com/jetstack/tarmak/pkg/wing/client/listers/wing/internalversion"
	//"k8s.io/apimachinery/pkg/api/errors"
	//"k8s.io/apimachinery/pkg/api/meta"
	//"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apiserver/pkg/admission"
)

// Register registers a plugin
func Register(plugins *admission.Plugins) {
	plugins.Register("BanMachine", func(config io.Reader) (admission.Interface, error) {
		return New()
	})
}

type DisallowMachine struct {
	*admission.Handler
	lister listers.MachineLister
}

var _ = winginitializer.WantsInternalWingInformerFactory(&DisallowMachine{})

// Admit ensures that the object in-flight is of kind Machine.
// In addition checks that the Name is not on the banned list.
// The list is stored in Machines API objects.
func (d *DisallowMachine) Admit(a admission.Attributes) error {
	// we are only interested in flunders
	//if a.GetKind().GroupKind() != wing.Kind("Machine") {
	//	return nil
	//}

	//if !d.WaitForReady() {
	//	return admission.NewForbidden(a, fmt.Errorf("not yet ready to handle request"))
	//}

	//metaAccessor, err := meta.Accessor(a.GetObject())
	//if err != nil {
	//	return err
	//}
	//flunderName := metaAccessor.GetName()

	//machines, err := d.lister.List(labels.Everything())
	//if err != nil {
	//	return err
	//}

	//for _, machine := range machines {
	//	for _, disallowedMachine := range machine.DisallowedMachines {
	//		if flunderName == disallowedMachine {
	//			return errors.NewForbidden(
	//				a.GetResource().GroupResource(),
	//				a.GetName(),
	//				fmt.Errorf("this name may not be used, please change the resource name"),
	//			)
	//		}
	//	}
	//}
	return nil
}

// SetInternalWingInformerFactory gets Lister from SharedInformerFactory.
// The lister knows how to lists Machines.
func (d *DisallowMachine) SetInternalWingInformerFactory(f informers.SharedInformerFactory) {
	d.lister = f.Wing().InternalVersion().Machines().Lister()
	d.SetReadyFunc(f.Wing().InternalVersion().Machines().Informer().HasSynced)
}

// ValidaValidateInitializationte checks whether the plugin was correctly initialized.
func (d *DisallowMachine) ValidateInitialization() error {
	if d.lister == nil {
		return fmt.Errorf("missing machine lister")
	}
	return nil
}

// New creates a new ban flunder admission plugin
func New() (*DisallowMachine, error) {
	return &DisallowMachine{
		Handler: admission.NewHandler(admission.Create),
	}, nil
}
