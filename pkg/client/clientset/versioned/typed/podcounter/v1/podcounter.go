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

package v1

import (
	v1 "github.com/bitshifta/crd-controller/pkg/apis/podcounter/v1"
	scheme "github.com/bitshifta/crd-controller/pkg/client/clientset/versioned/scheme"
	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	rest "k8s.io/client-go/rest"
)

// PodCountersGetter has a method to return a PodCounterInterface.
// A group's client should implement this interface.
type PodCountersGetter interface {
	PodCounters(namespace string) PodCounterInterface
}

// PodCounterInterface has methods to work with PodCounter resources.
type PodCounterInterface interface {
	Create(*v1.PodCounter) (*v1.PodCounter, error)
	Update(*v1.PodCounter) (*v1.PodCounter, error)
	Delete(name string, options *meta_v1.DeleteOptions) error
	DeleteCollection(options *meta_v1.DeleteOptions, listOptions meta_v1.ListOptions) error
	Get(name string, options meta_v1.GetOptions) (*v1.PodCounter, error)
	List(opts meta_v1.ListOptions) (*v1.PodCounterList, error)
	Watch(opts meta_v1.ListOptions) (watch.Interface, error)
	Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *v1.PodCounter, err error)
	PodCounterExpansion
}

// podCounters implements PodCounterInterface
type podCounters struct {
	client rest.Interface
	ns     string
}

// newPodCounters returns a PodCounters
func newPodCounters(c *KhaliltV1Client, namespace string) *podCounters {
	return &podCounters{
		client: c.RESTClient(),
		ns:     namespace,
	}
}

// Get takes name of the podCounter, and returns the corresponding podCounter object, and an error if there is any.
func (c *podCounters) Get(name string, options meta_v1.GetOptions) (result *v1.PodCounter, err error) {
	result = &v1.PodCounter{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource("podcounters").
		Name(name).
		VersionedParams(&options, scheme.ParameterCodec).
		Do().
		Into(result)
	return
}

// List takes label and field selectors, and returns the list of PodCounters that match those selectors.
func (c *podCounters) List(opts meta_v1.ListOptions) (result *v1.PodCounterList, err error) {
	result = &v1.PodCounterList{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource("podcounters").
		VersionedParams(&opts, scheme.ParameterCodec).
		Do().
		Into(result)
	return
}

// Watch returns a watch.Interface that watches the requested podCounters.
func (c *podCounters) Watch(opts meta_v1.ListOptions) (watch.Interface, error) {
	opts.Watch = true
	return c.client.Get().
		Namespace(c.ns).
		Resource("podcounters").
		VersionedParams(&opts, scheme.ParameterCodec).
		Watch()
}

// Create takes the representation of a podCounter and creates it.  Returns the server's representation of the podCounter, and an error, if there is any.
func (c *podCounters) Create(podCounter *v1.PodCounter) (result *v1.PodCounter, err error) {
	result = &v1.PodCounter{}
	err = c.client.Post().
		Namespace(c.ns).
		Resource("podcounters").
		Body(podCounter).
		Do().
		Into(result)
	return
}

// Update takes the representation of a podCounter and updates it. Returns the server's representation of the podCounter, and an error, if there is any.
func (c *podCounters) Update(podCounter *v1.PodCounter) (result *v1.PodCounter, err error) {
	result = &v1.PodCounter{}
	err = c.client.Put().
		Namespace(c.ns).
		Resource("podcounters").
		Name(podCounter.Name).
		Body(podCounter).
		Do().
		Into(result)
	return
}

// Delete takes name of the podCounter and deletes it. Returns an error if one occurs.
func (c *podCounters) Delete(name string, options *meta_v1.DeleteOptions) error {
	return c.client.Delete().
		Namespace(c.ns).
		Resource("podcounters").
		Name(name).
		Body(options).
		Do().
		Error()
}

// DeleteCollection deletes a collection of objects.
func (c *podCounters) DeleteCollection(options *meta_v1.DeleteOptions, listOptions meta_v1.ListOptions) error {
	return c.client.Delete().
		Namespace(c.ns).
		Resource("podcounters").
		VersionedParams(&listOptions, scheme.ParameterCodec).
		Body(options).
		Do().
		Error()
}

// Patch applies the patch and returns the patched podCounter.
func (c *podCounters) Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *v1.PodCounter, err error) {
	result = &v1.PodCounter{}
	err = c.client.Patch(pt).
		Namespace(c.ns).
		Resource("podcounters").
		SubResource(subresources...).
		Name(name).
		Body(data).
		Do().
		Into(result)
	return
}
