package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

const (
	GroupName    = "config.example.com"
	GroupVersion = "v1"
)

var (
	SchemeGroupVersion = schema.GroupVersion{Group: GroupName, Version: GroupVersion}
	SchemeBuilder      = runtime.NewSchemeBuilder(addKnownTypes)
	AddToScheme        = SchemeBuilder.AddToScheme
)

func Resource(resource string) schema.GroupResource {
	return SchemeGroupVersion.WithResource(resource).GroupResource()
}

func addKnownTypes(scheme *runtime.Scheme) error {
	scheme.AddKnownTypes(SchemeGroupVersion,
		&ServiceConfiguration{},
		&ServiceConfigurationList{},
	)
	metav1.AddToGroupVersion(scheme, SchemeGroupVersion)
	return nil
}

func (in *ServiceConfiguration) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

func (in *ServiceConfiguration) DeepCopy() *ServiceConfiguration {
	if in == nil {
		return nil
	}
	out := new(ServiceConfiguration)
	in.DeepCopyInto(out)
	return out
}

func (in *ServiceConfiguration) DeepCopyInto(out *ServiceConfiguration) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
}

func (in *ServiceConfigurationList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

func (in *ServiceConfigurationList) DeepCopy() *ServiceConfigurationList {
	if in == nil {
		return nil
	}
	out := new(ServiceConfigurationList)
	in.DeepCopyInto(out)
	return out
}

func (in *ServiceConfigurationList) DeepCopyInto(out *ServiceConfigurationList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]ServiceConfiguration, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

func (in *ServiceConfigurationSpec) DeepCopyInto(out *ServiceConfigurationSpec) {
	*out = *in
	if in.Config.Captures != nil {
		in, out := &in.Config.Captures, &out.Config.Captures
		*out = make([]CaptureConfig, len(*in))
		copy(*out, *in)
	}
	if in.Config.Destinations != nil {
		in, out := &in.Config.Destinations, &out.Config.Destinations
		*out = make([]DestinationConfig, len(*in))
		copy(*out, *in)
	}
}
