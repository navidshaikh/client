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

package rest

import (
	restclient "k8s.io/client-go/rest"
)

// KnRESTClient to client-go REST client. All methods are relative to the
// namespace specified during construction
type KnRESTClient interface {
	// Namespace in which this client is operating for
	Namespace() string
}

// KnRESTClient is a combination of client-go REST client interface and namespace
type knRESTClient struct {
	client    *restclient.RESTClient
	namespace string
}

// NewKnRESTClient is to invoke Eventing Sources Client API to create object
func NewKnRESTClient(client *restclient.RESTClient, namespace string) KnRESTClient {
	return &knRESTClient{
		namespace: namespace,
		client:    client,
	}
}

// Return the client's namespace
func (c *knRESTClient) Namespace() string {
	return c.namespace
}
