/*
Copyright 2022 The SODA Authors.

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

	v1beta1 "github.com/soda-cdm/kahu/apis/kahu/v1beta1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	labels "k8s.io/apimachinery/pkg/labels"
	schema "k8s.io/apimachinery/pkg/runtime/schema"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	testing "k8s.io/client-go/testing"
)

// FakeBackups implements BackupInterface
type FakeBackups struct {
	Fake *FakeKahuV1beta1
}

var backupsResource = schema.GroupVersionResource{Group: "kahu.io", Version: "v1beta1", Resource: "backups"}

var backupsKind = schema.GroupVersionKind{Group: "kahu.io", Version: "v1beta1", Kind: "Backup"}

// Get takes name of the backup, and returns the corresponding backup object, and an error if there is any.
func (c *FakeBackups) Get(ctx context.Context, name string, options v1.GetOptions) (result *v1beta1.Backup, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootGetAction(backupsResource, name), &v1beta1.Backup{})
	if obj == nil {
		return nil, err
	}
	return obj.(*v1beta1.Backup), err
}

// List takes label and field selectors, and returns the list of Backups that match those selectors.
func (c *FakeBackups) List(ctx context.Context, opts v1.ListOptions) (result *v1beta1.BackupList, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootListAction(backupsResource, backupsKind, opts), &v1beta1.BackupList{})
	if obj == nil {
		return nil, err
	}

	label, _, _ := testing.ExtractFromListOptions(opts)
	if label == nil {
		label = labels.Everything()
	}
	list := &v1beta1.BackupList{ListMeta: obj.(*v1beta1.BackupList).ListMeta}
	for _, item := range obj.(*v1beta1.BackupList).Items {
		if label.Matches(labels.Set(item.Labels)) {
			list.Items = append(list.Items, item)
		}
	}
	return list, err
}

// Watch returns a watch.Interface that watches the requested backups.
func (c *FakeBackups) Watch(ctx context.Context, opts v1.ListOptions) (watch.Interface, error) {
	return c.Fake.
		InvokesWatch(testing.NewRootWatchAction(backupsResource, opts))
}

// Create takes the representation of a backup and creates it.  Returns the server's representation of the backup, and an error, if there is any.
func (c *FakeBackups) Create(ctx context.Context, backup *v1beta1.Backup, opts v1.CreateOptions) (result *v1beta1.Backup, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootCreateAction(backupsResource, backup), &v1beta1.Backup{})
	if obj == nil {
		return nil, err
	}
	return obj.(*v1beta1.Backup), err
}

// Update takes the representation of a backup and updates it. Returns the server's representation of the backup, and an error, if there is any.
func (c *FakeBackups) Update(ctx context.Context, backup *v1beta1.Backup, opts v1.UpdateOptions) (result *v1beta1.Backup, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootUpdateAction(backupsResource, backup), &v1beta1.Backup{})
	if obj == nil {
		return nil, err
	}
	return obj.(*v1beta1.Backup), err
}

// UpdateStatus was generated because the type contains a Status member.
// Add a +genclient:noStatus comment above the type to avoid generating UpdateStatus().
func (c *FakeBackups) UpdateStatus(ctx context.Context, backup *v1beta1.Backup, opts v1.UpdateOptions) (*v1beta1.Backup, error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootUpdateSubresourceAction(backupsResource, "status", backup), &v1beta1.Backup{})
	if obj == nil {
		return nil, err
	}
	return obj.(*v1beta1.Backup), err
}

// Delete takes name of the backup and deletes it. Returns an error if one occurs.
func (c *FakeBackups) Delete(ctx context.Context, name string, opts v1.DeleteOptions) error {
	_, err := c.Fake.
		Invokes(testing.NewRootDeleteAction(backupsResource, name), &v1beta1.Backup{})
	return err
}

// DeleteCollection deletes a collection of objects.
func (c *FakeBackups) DeleteCollection(ctx context.Context, opts v1.DeleteOptions, listOpts v1.ListOptions) error {
	action := testing.NewRootDeleteCollectionAction(backupsResource, listOpts)

	_, err := c.Fake.Invokes(action, &v1beta1.BackupList{})
	return err
}

// Patch applies the patch and returns the patched backup.
func (c *FakeBackups) Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts v1.PatchOptions, subresources ...string) (result *v1beta1.Backup, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootPatchSubresourceAction(backupsResource, name, pt, data, subresources...), &v1beta1.Backup{})
	if obj == nil {
		return nil, err
	}
	return obj.(*v1beta1.Backup), err
}
