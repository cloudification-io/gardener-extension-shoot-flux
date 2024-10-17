// SPDX-FileCopyrightText: 2024 SAP SE or an SAP affiliate company and Gardener contributors
//
// SPDX-License-Identifier: Apache-2.0

package v1

import (
	"strings"

	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +genclient
// +genclient:nonNamespaced
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// ControllerDeployment contains information about how this controller is deployed.
type ControllerDeployment struct {
	metav1.TypeMeta `json:",inline"`
	// Standard object metadata.
	metav1.ObjectMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`
	// Helm configures that an extension controller is deployed using helm.
	// +optional
	Helm *HelmControllerDeployment `json:"helm,omitempty" protobuf:"bytes,2,opt,name=helm"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// ControllerDeploymentList is a collection of ControllerDeployments.
type ControllerDeploymentList struct {
	metav1.TypeMeta `json:",inline"`
	// Standard list object metadata.
	// +optional
	metav1.ListMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`
	// Items is the list of ControllerDeployments.
	Items []ControllerDeployment `json:"items" protobuf:"bytes,2,rep,name=items"`
}

// HelmControllerDeployment configures how an extension controller is deployed using helm.
type HelmControllerDeployment struct {
	// RawChart is the base64-encoded, gzip'ed, tar'ed extension controller chart.
	// +optional
	RawChart []byte `json:"rawChart,omitempty" protobuf:"bytes,1,opt,name=rawChart"`
	// Values are the chart values.
	// +optional
	Values *apiextensionsv1.JSON `json:"values,omitempty" protobuf:"bytes,2,opt,name=values"`
	// OCIRepository defines where to pull the chart.
	// +optional
	OCIRepository *OCIRepository `json:"ociRepository,omitempty" protobuf:"bytes,3,opt,name=ociRepository"`
}

// OCIRepository configures where to pull an OCI Artifact, that could contain for example a Helm Chart.
type OCIRepository struct {
	// Ref is the full artifact Ref and takes precedence over all other fields.
	// +optional
	Ref *string `json:"ref,omitempty" protobuf:"bytes,1,name=ref"`
	// Repository is a reference to an OCI artifact repository.
	// +optional
	Repository *string `json:"repository,omitempty" protobuf:"bytes,2,name=repository"`
	// Tag is the image tag to pull.
	// +optional
	Tag *string `json:"tag,omitempty" protobuf:"bytes,3,opt,name=tag"`
	// Digest of the image to pull, takes precedence over tag.
	// The value should be in the format 'sha256:<HASH>'.
	// +optional
	Digest *string `json:"digest,omitempty" protobuf:"bytes,4,opt,name=digest"`
}

// GetURL returns the fully-qualified OCIRepository URL of the artifact.
func (r *OCIRepository) GetURL() string {
	var ref string

	switch {
	case r.Ref != nil:
		ref = *r.Ref
	case r.Digest != nil:
		// when digest is set we ignore the tag
		ref = *r.Repository + "@" + *r.Digest
	case r.Tag != nil:
		ref = *r.Repository + ":" + *r.Tag
	}
	return strings.TrimPrefix(ref, "oci://")
}
