// Copyright Jetstack Ltd. See LICENSE for details.
package install

import (
	"github.com/jetstack/tarmak/pkg/apis/wing/v1alpha1"
	"k8s.io/apimachinery/pkg/runtime"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
)

// Install registers the API group and adds types to a scheme
func Install(scheme *runtime.Scheme) {
	utilruntime.Must(v1alpha1.AddToScheme(scheme))
}
