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
	"testing"

	"knative.dev/eventing/pkg/apis/sources/v1alpha1"

	"knative.dev/client/pkg/util/mock"
)

type MockKnApiServerSourcesClient struct {
	t         *testing.T
	recorder  *ApiServerSourcesRecorder
	namespace string
}

// NewMockKnApiServerSourceClient returns a new mock instance which you need to record for
func NewMockKnApiServerSourceClient(t *testing.T, ns ...string) *MockKnApiServerSourcesClient {
	namespace := "default"
	if len(ns) > 0 {
		namespace = ns[0]
	}
	return &MockKnApiServerSourcesClient{
		t:        t,
		recorder: &ApiServerSourcesRecorder{mock.NewRecorder(t, namespace)},
	}
}

// ApiServerSourcesRecorder records methods for ApiServerSourcesClient
type ApiServerSourcesRecorder struct {
	r *mock.Recorder
}

// Recorder returns method recorder
func (c *MockKnApiServerSourcesClient) Recorder() *ApiServerSourcesRecorder {
	return c.recorder
}

// Namespace of this client
func (c *MockKnApiServerSourcesClient) Namespace() string {
	return c.recorder.r.Namespace()
}

func (sr *ApiServerSourcesRecorder) GetApiServerSource(name interface{}, source *v1alpha1.ApiServerSource, err error) {
	sr.r.Add("GetApiServerSource", []interface{}{name}, []interface{}{source, err})
}

func (c *MockKnApiServerSourcesClient) GetApiServerSource(name string) (*v1alpha1.ApiServerSource, error) {
	call := c.recorder.r.VerifyCall("GetApiServerSource", name)
	return call.Result[0].(*v1alpha1.ApiServerSource), mock.ErrorOrNil(call.Result[1])
}

func (sr *ApiServerSourcesRecorder) UpdateApiServerSource(source interface{}, err error) {
	sr.r.Add("UpdateService", []interface{}{source}, []interface{}{err})
}

func (c *MockKnApiServerSourcesClient) UpdateApiServerSource(source *v1alpha1.ApiServerSource) error {
	call := c.recorder.r.VerifyCall("UpdateService", source)
	return mock.ErrorOrNil(call.Result[0])
}
