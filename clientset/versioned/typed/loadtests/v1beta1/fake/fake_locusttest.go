/*
Copyright The Kubernetes Authors.

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

// Code generated by client-gen. DO NOT EDIT.

package fake

import (
	v1beta1 "github.com/luizbafilho/lokust/apis/loadtests/v1beta1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	labels "k8s.io/apimachinery/pkg/labels"
	schema "k8s.io/apimachinery/pkg/runtime/schema"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	testing "k8s.io/client-go/testing"
)

// FakeLocustTests implements LocustTestInterface
type FakeLocustTests struct {
	Fake *FakeLoadtestsV1beta1
	ns   string
}

var locusttestsResource = schema.GroupVersionResource{Group: "loadtests", Version: "v1beta1", Resource: "locusttests"}

var locusttestsKind = schema.GroupVersionKind{Group: "loadtests", Version: "v1beta1", Kind: "LocustTest"}

// Get takes name of the locustTest, and returns the corresponding locustTest object, and an error if there is any.
func (c *FakeLocustTests) Get(name string, options v1.GetOptions) (result *v1beta1.LocustTest, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewGetAction(locusttestsResource, c.ns, name), &v1beta1.LocustTest{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1beta1.LocustTest), err
}

// List takes label and field selectors, and returns the list of LocustTests that match those selectors.
func (c *FakeLocustTests) List(opts v1.ListOptions) (result *v1beta1.LocustTestList, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewListAction(locusttestsResource, locusttestsKind, c.ns, opts), &v1beta1.LocustTestList{})

	if obj == nil {
		return nil, err
	}

	label, _, _ := testing.ExtractFromListOptions(opts)
	if label == nil {
		label = labels.Everything()
	}
	list := &v1beta1.LocustTestList{ListMeta: obj.(*v1beta1.LocustTestList).ListMeta}
	for _, item := range obj.(*v1beta1.LocustTestList).Items {
		if label.Matches(labels.Set(item.Labels)) {
			list.Items = append(list.Items, item)
		}
	}
	return list, err
}

// Watch returns a watch.Interface that watches the requested locustTests.
func (c *FakeLocustTests) Watch(opts v1.ListOptions) (watch.Interface, error) {
	return c.Fake.
		InvokesWatch(testing.NewWatchAction(locusttestsResource, c.ns, opts))

}

// Create takes the representation of a locustTest and creates it.  Returns the server's representation of the locustTest, and an error, if there is any.
func (c *FakeLocustTests) Create(locustTest *v1beta1.LocustTest) (result *v1beta1.LocustTest, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewCreateAction(locusttestsResource, c.ns, locustTest), &v1beta1.LocustTest{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1beta1.LocustTest), err
}

// Update takes the representation of a locustTest and updates it. Returns the server's representation of the locustTest, and an error, if there is any.
func (c *FakeLocustTests) Update(locustTest *v1beta1.LocustTest) (result *v1beta1.LocustTest, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewUpdateAction(locusttestsResource, c.ns, locustTest), &v1beta1.LocustTest{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1beta1.LocustTest), err
}

// UpdateStatus was generated because the type contains a Status member.
// Add a +genclient:noStatus comment above the type to avoid generating UpdateStatus().
func (c *FakeLocustTests) UpdateStatus(locustTest *v1beta1.LocustTest) (*v1beta1.LocustTest, error) {
	obj, err := c.Fake.
		Invokes(testing.NewUpdateSubresourceAction(locusttestsResource, "status", c.ns, locustTest), &v1beta1.LocustTest{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1beta1.LocustTest), err
}

// Delete takes name of the locustTest and deletes it. Returns an error if one occurs.
func (c *FakeLocustTests) Delete(name string, options *v1.DeleteOptions) error {
	_, err := c.Fake.
		Invokes(testing.NewDeleteAction(locusttestsResource, c.ns, name), &v1beta1.LocustTest{})

	return err
}

// DeleteCollection deletes a collection of objects.
func (c *FakeLocustTests) DeleteCollection(options *v1.DeleteOptions, listOptions v1.ListOptions) error {
	action := testing.NewDeleteCollectionAction(locusttestsResource, c.ns, listOptions)

	_, err := c.Fake.Invokes(action, &v1beta1.LocustTestList{})
	return err
}

// Patch applies the patch and returns the patched locustTest.
func (c *FakeLocustTests) Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *v1beta1.LocustTest, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewPatchSubresourceAction(locusttestsResource, c.ns, name, pt, data, subresources...), &v1beta1.LocustTest{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1beta1.LocustTest), err
}
