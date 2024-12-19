package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type CaptureConfig struct {
	Interface string `json:"interface"`
	Filter    string `json:"filter"`
	IMSI      string `json:"imsi"`
	SnapLen   int    `json:"snaplen"`
}

type DestinationConfig struct {
	PathPrefix string `json:"path_prefix"`
	Rotate     bool   `json:"rotate"`
	MaxPackets int    `json:"max_packets"`
	Compress   bool   `json:"compress"`
}

type VPPCaptureConfig struct {
	Deployment   string              `json:"deployment"`
	Captures     []CaptureConfig     `json:"captures"`
	BatchSize    int                 `json:"batchSize"`
	Destinations []DestinationConfig `json:"destinations"`
}

type ServiceConfiguration struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec ServiceConfigurationSpec `json:"spec"`
}

type ServiceConfigurationSpec struct {
	Organization string           `json:"organization"`
	Product      string           `json:"product"`
	Version      string           `json:"version"`
	ServiceName  string           `json:"serviceName"`
	ConfigName   string           `json:"configName"`
	Config       VPPCaptureConfig `json:"config"`
}

type ServiceConfigurationList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []ServiceConfiguration `json:"items"`
}
