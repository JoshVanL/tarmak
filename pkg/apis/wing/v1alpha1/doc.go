// Copyright Jetstack Ltd. See LICENSE for details.

// Api versions allow the api contract for a resource to be changed while keeping
// backward compatibility by support multiple concurrent versions
// of the same resource

// +k8s:openapi-gen=true
// +k8s:deepcopy-gen=package,register
// +k8s:conversion-gen=github.com/jetstack/wing-api/pkg/apis/wing
// +k8s:defaulter-gen=TypeMeta
// +groupName=wing.tarmak.io
package v1alpha1
