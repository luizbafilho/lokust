// +build !ignore_autogenerated

/*


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

// Code generated by controller-gen. DO NOT EDIT.

package v1beta1

import (
	runtime "k8s.io/apimachinery/pkg/runtime"
)

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *LocustTest) DeepCopyInto(out *LocustTest) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	in.Status.DeepCopyInto(&out.Status)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new LocustTest.
func (in *LocustTest) DeepCopy() *LocustTest {
	if in == nil {
		return nil
	}
	out := new(LocustTest)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *LocustTest) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *LocustTestList) DeepCopyInto(out *LocustTestList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]LocustTest, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new LocustTestList.
func (in *LocustTestList) DeepCopy() *LocustTestList {
	if in == nil {
		return nil
	}
	out := new(LocustTestList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *LocustTestList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *LocustTestResources) DeepCopyInto(out *LocustTestResources) {
	*out = *in
	in.Master.DeepCopyInto(&out.Master)
	in.Workers.DeepCopyInto(&out.Workers)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new LocustTestResources.
func (in *LocustTestResources) DeepCopy() *LocustTestResources {
	if in == nil {
		return nil
	}
	out := new(LocustTestResources)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *LocustTestSpec) DeepCopyInto(out *LocustTestSpec) {
	*out = *in
	if in.Replicas != nil {
		in, out := &in.Replicas, &out.Replicas
		*out = new(int32)
		**out = **in
	}
	in.Resources.DeepCopyInto(&out.Resources)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new LocustTestSpec.
func (in *LocustTestSpec) DeepCopy() *LocustTestSpec {
	if in == nil {
		return nil
	}
	out := new(LocustTestSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *LocustTestStatus) DeepCopyInto(out *LocustTestStatus) {
	*out = *in
	if in.Replicas != nil {
		in, out := &in.Replicas, &out.Replicas
		*out = new(int32)
		**out = **in
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new LocustTestStatus.
func (in *LocustTestStatus) DeepCopy() *LocustTestStatus {
	if in == nil {
		return nil
	}
	out := new(LocustTestStatus)
	in.DeepCopyInto(out)
	return out
}
