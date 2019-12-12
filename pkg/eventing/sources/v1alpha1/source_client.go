// Copyright © 2019 The Knative Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package v1alpha1

import (
	apis_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	kn_errors "knative.dev/client/pkg/errors"
	"knative.dev/eventing/pkg/apis/sources/v1alpha1"
	client_v1alpha1 "knative.dev/eventing/pkg/client/clientset/versioned/typed/sources/v1alpha1"
)

// KnApiServerSourcesClient to Eventing Sources. All methods are relative to the
// namespace specified during construction
type KnApiServerSourcesClient interface {
	// Namespace in which this client is operating for
	Namespace() string

	// Get ApiServerSource
	GetApiServerSource(name string) (*v1alpha1.ApiServerSource, error)

	// Get an ApiServerSource by object
	CreateApiServerSource(apisvrsrc *v1alpha1.ApiServerSource) (*v1alpha1.ApiServerSource, error)

	// Update APIServerSource
	UpdateApiServerSource(src *v1alpha1.ApiServerSource) error

	// Delete an ApiServerSource by name
	DeleteApiServerSource(name string) error
}

// KnApiServerSourcesClient is a combination of Sources client interface and namespace
// Temporarily help to add sources dependencies
// May be changed when adding real sources features
type knApiServerSourcesClient struct {
	client    client_v1alpha1.SourcesV1alpha1Interface
	namespace string
}

// NewKnApiServerSourcesClient is to invoke Eventing Sources Client API to create object
func NewKnApiServerSourcesClient(client client_v1alpha1.SourcesV1alpha1Interface, namespace string) KnApiServerSourcesClient {
	return &knApiServerSourcesClient{
		client:    client,
		namespace: namespace,
	}
}

// GetApiServerSource finds ApiServerSource object by name if present
func (c *knApiServerSourcesClient) GetApiServerSource(name string) (*v1alpha1.ApiServerSource, error) {
	src, err := c.client.ApiServerSources(c.namespace).Get(name, apis_v1.GetOptions{})
	if err != nil {
		return nil, kn_errors.GetError(err)
	}
	return src, nil
}

//CreateApiServerSource is used to create an instance of ApiServerSource
func (c *knApiServerSourcesClient) CreateApiServerSource(apisvrsrc *v1alpha1.ApiServerSource) (*v1alpha1.ApiServerSource, error) {
	ins, err := c.client.ApiServerSources(c.namespace).Create(apisvrsrc)
	if err != nil {
		return nil, kn_errors.GetError(err)
	}
	return ins, nil
}

//DeleteApiServerSource is used to create an instance of ApiServerSource
func (c *knApiServerSourcesClient) DeleteApiServerSource(name string) error {
	err := c.client.ApiServerSources(c.namespace).Delete(name, &apis_v1.DeleteOptions{})
	return err
}

//UpdateApiServerSource
func (c *knApiServerSourcesClient) UpdateApiServerSource(source *v1alpha1.ApiServerSource) error {
	_, err := c.client.ApiServerSources(c.namespace).Update(source)
	if err != nil {
		return kn_errors.GetError(err)
	}
	return nil
}

// Return the client's namespace
func (c *knApiServerSourcesClient) Namespace() string {
	return c.namespace
}
