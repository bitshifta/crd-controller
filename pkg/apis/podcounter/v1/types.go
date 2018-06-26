package v1

import (
	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +genclient
// +genclient:noStatus
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// PodCounter describes a PodCounter resource
type PodCounter struct {
	// TypeMeta is the metadata for the resource, like kind and apiversion
	meta_v1.TypeMeta `json:",inline"`
	// ObjectMeta contains the metadata for the particular object, including
	// things like...
	//  - name
	//  - namespace
	//  - self link
	//  - labels
	//  - ... etc ...
	meta_v1.ObjectMeta `json:"metadata,omitempty"`

	// Spec is the custom resource spec
	Spec PodCounterSpec `json:"spec"`
}

// PodCounterSpec is the spec for a PodCounter resource
type PodCounterSpec struct {
	// Message and SomeValue are example custom spec fields
	//
	// this is where you would put your custom resource data
	Current    *int32 `json:"current"`
	Historical *int32 `json:"historical"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// PodCounterList is a list of PodCounter resources
type PodCounterList struct {
	meta_v1.TypeMeta `json:",inline"`
	meta_v1.ListMeta `json:"metadata"`

	Items []PodCounter `json:"items"`
}
