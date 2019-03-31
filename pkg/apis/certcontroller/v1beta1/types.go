/*
Copyright 2017 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package v1beta1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// Cert is a specification for a Cert resource
type Cert struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   CertSpec   `json:"spec"`
	Status CertStatus `json:"status"`
}



// CertSpec is the spec for a Cert resource
type CertSpec struct {
	SecretName string `json:"secretName"`
	Domain string `json:"domain"`
	Type string `json:"type"`
	ValidityPeriod int `json:"validityPeriod"`
	Env map[string]string `json:"env"`
}

// CertStatus is the status for a GenericDaemon resource
type CertStatus struct {
	SecretCreateTime string `json:"secretCreateTime"`
	RemainingValidDays int `json:"remainingValidDays"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// CertList is a list of Cert resources
type CertList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata"`

	Items []Cert `json:"items"`
}
