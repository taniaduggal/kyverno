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
	"context"

	v1beta1 "github.com/kyverno/kyverno/api/kyverno/v1beta1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	labels "k8s.io/apimachinery/pkg/labels"
	schema "k8s.io/apimachinery/pkg/runtime/schema"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	testing "k8s.io/client-go/testing"
)

// FakeUpdateRequests implements UpdateRequestInterface
type FakeUpdateRequests struct {
	Fake *FakeKyvernoV1beta1
	ns   string
}

var updaterequestsResource = schema.GroupVersionResource{Group: "kyverno.io", Version: "v1beta1", Resource: "updaterequests"}

var updaterequestsKind = schema.GroupVersionKind{Group: "kyverno.io", Version: "v1beta1", Kind: "UpdateRequest"}

// Get takes name of the updateRequest, and returns the corresponding updateRequest object, and an error if there is any.
func (c *FakeUpdateRequests) Get(ctx context.Context, name string, options v1.GetOptions) (result *v1beta1.UpdateRequest, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewGetAction(updaterequestsResource, c.ns, name), &v1beta1.UpdateRequest{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1beta1.UpdateRequest), err
}

// List takes label and field selectors, and returns the list of UpdateRequests that match those selectors.
func (c *FakeUpdateRequests) List(ctx context.Context, opts v1.ListOptions) (result *v1beta1.UpdateRequestList, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewListAction(updaterequestsResource, updaterequestsKind, c.ns, opts), &v1beta1.UpdateRequestList{})

	if obj == nil {
		return nil, err
	}

	label, _, _ := testing.ExtractFromListOptions(opts)
	if label == nil {
		label = labels.Everything()
	}
	list := &v1beta1.UpdateRequestList{ListMeta: obj.(*v1beta1.UpdateRequestList).ListMeta}
	for _, item := range obj.(*v1beta1.UpdateRequestList).Items {
		if label.Matches(labels.Set(item.Labels)) {
			list.Items = append(list.Items, item)
		}
	}
	return list, err
}

// Watch returns a watch.Interface that watches the requested updateRequests.
func (c *FakeUpdateRequests) Watch(ctx context.Context, opts v1.ListOptions) (watch.Interface, error) {
	return c.Fake.
		InvokesWatch(testing.NewWatchAction(updaterequestsResource, c.ns, opts))

}

// Create takes the representation of a updateRequest and creates it.  Returns the server's representation of the updateRequest, and an error, if there is any.
func (c *FakeUpdateRequests) Create(ctx context.Context, updateRequest *v1beta1.UpdateRequest, opts v1.CreateOptions) (result *v1beta1.UpdateRequest, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewCreateAction(updaterequestsResource, c.ns, updateRequest), &v1beta1.UpdateRequest{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1beta1.UpdateRequest), err
}

// Update takes the representation of a updateRequest and updates it. Returns the server's representation of the updateRequest, and an error, if there is any.
func (c *FakeUpdateRequests) Update(ctx context.Context, updateRequest *v1beta1.UpdateRequest, opts v1.UpdateOptions) (result *v1beta1.UpdateRequest, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewUpdateAction(updaterequestsResource, c.ns, updateRequest), &v1beta1.UpdateRequest{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1beta1.UpdateRequest), err
}

// UpdateStatus was generated because the type contains a Status member.
// Add a +genclient:noStatus comment above the type to avoid generating UpdateStatus().
func (c *FakeUpdateRequests) UpdateStatus(ctx context.Context, updateRequest *v1beta1.UpdateRequest, opts v1.UpdateOptions) (*v1beta1.UpdateRequest, error) {
	obj, err := c.Fake.
		Invokes(testing.NewUpdateSubresourceAction(updaterequestsResource, "status", c.ns, updateRequest), &v1beta1.UpdateRequest{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1beta1.UpdateRequest), err
}

// Delete takes name of the updateRequest and deletes it. Returns an error if one occurs.
func (c *FakeUpdateRequests) Delete(ctx context.Context, name string, opts v1.DeleteOptions) error {
	_, err := c.Fake.
		Invokes(testing.NewDeleteActionWithOptions(updaterequestsResource, c.ns, name, opts), &v1beta1.UpdateRequest{})

	return err
}

// DeleteCollection deletes a collection of objects.
func (c *FakeUpdateRequests) DeleteCollection(ctx context.Context, opts v1.DeleteOptions, listOpts v1.ListOptions) error {
	action := testing.NewDeleteCollectionAction(updaterequestsResource, c.ns, listOpts)

	_, err := c.Fake.Invokes(action, &v1beta1.UpdateRequestList{})
	return err
}

// Patch applies the patch and returns the patched updateRequest.
func (c *FakeUpdateRequests) Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts v1.PatchOptions, subresources ...string) (result *v1beta1.UpdateRequest, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewPatchSubresourceAction(updaterequestsResource, c.ns, name, pt, data, subresources...), &v1beta1.UpdateRequest{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1beta1.UpdateRequest), err
}
