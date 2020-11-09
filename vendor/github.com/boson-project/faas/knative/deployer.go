package knative

import (
	"fmt"
	"strings"
	"time"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	servinglib "knative.dev/client/pkg/serving"
	"knative.dev/client/pkg/wait"
	servingv1 "knative.dev/serving/pkg/apis/serving/v1"
	v1 "knative.dev/serving/pkg/apis/serving/v1"

	"github.com/boson-project/faas"
	"github.com/boson-project/faas/k8s"
)

type Deployer struct {
	// Namespace with which to override that set on the default configuration (such as the ~/.kube/config).
	// If left blank, deployment will commence to the configured namespace.
	Namespace string
	// Verbose logging enablement flag.
	Verbose bool
}

func NewDeployer(namespaceOverride string) (deployer *Deployer, err error) {
	deployer = &Deployer{}
	namespace, err := GetNamespace(namespaceOverride)
	if err != nil {
		return
	}
	deployer.Namespace = namespace
	return
}

func (d *Deployer) Deploy(f faas.Function) (err error) {

	// k8s does not support service names with dots. so encode it such that
	// www.my-domain,com -> www-my--domain-com
	serviceName, err := k8s.ToK8sAllowedName(f.Name)
	if err != nil {
		return
	}

	client, err := NewServingClient(d.Namespace)
	if err != nil {
		return
	}

	_, err = client.GetService(serviceName)
	if err != nil {
		if errors.IsNotFound(err) {

			// Let's create a new Service
			if d.Verbose {
				fmt.Printf("Creating Knative Service: %v\n", serviceName)
			}
			err := client.CreateService(generateNewService(serviceName, f.ImageWithDigest(), f.Runtime))
			if err != nil {
				err = fmt.Errorf("knative deployer failed to deploy the service: %v", err)
				return err
			}

			if d.Verbose {
				fmt.Println("Waiting for Knative Service to become ready")
			}
			err, _ = client.WaitForService(serviceName, DefaultWaitingTimeout, wait.NoopMessageCallback())
			if err != nil {
				err = fmt.Errorf("knative deployer failed to wait for the service to become ready: %v", err)
				return err
			}

			route, err := client.GetRoute(serviceName)
			if err != nil {
				err = fmt.Errorf("knative deployer failed to get the route: %v", err)
				return err
			}

			fmt.Println("Function deployed at URL: " + route.Status.URL.String())

		} else {
			err = fmt.Errorf("knative deployer failed to get the service: %v", err)
			return err
		}
	} else {
		// Update the existing Service
		err = client.UpdateServiceWithRetry(serviceName, updateService(f.ImageWithDigest(), f.EnvVars), 3)
		if err != nil {
			err = fmt.Errorf("knative deployer failed to update the service: %v", err)
			return err
		}

		route, err := client.GetRoute(serviceName)
		if err != nil {
			err = fmt.Errorf("knative deployer failed to get the route: %v", err)
			return err
		}

		fmt.Println("Function updated at URL: " + route.Status.URL.String())
	}

	return nil
}

func probeFor(url string) *corev1.Probe {
	return &corev1.Probe{
		Handler: corev1.Handler{
			HTTPGet: &corev1.HTTPGetAction{
				Path: url,
			},
		},
	}
}

func generateNewService(name, image, runtime string) *servingv1.Service {
	containers := []corev1.Container{
		{
			Image: image,
			Env: []corev1.EnvVar{
				{Name: "VERBOSE", Value: "true"},
			},
		},
	}

	if runtime != "quarkus" {
		containers[0].LivenessProbe = probeFor("/health/liveness")
		containers[0].ReadinessProbe = probeFor("/health/readiness")
	}

	return &v1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name: name,
			Labels: map[string]string{
				"boson.dev/function": "true",
				"boson.dev/runtime":  runtime,
			},
		},
		Spec: v1.ServiceSpec{
			ConfigurationSpec: v1.ConfigurationSpec{
				Template: v1.RevisionTemplateSpec{
					Spec: v1.RevisionSpec{
						PodSpec: corev1.PodSpec{
							Containers: containers,
						},
					},
				},
			},
		},
	}
}

func updateService(image string, envVars map[string]string) func(service *servingv1.Service) (*servingv1.Service, error) {
	return func(service *servingv1.Service) (*servingv1.Service, error) {
		err := servinglib.UpdateImage(&service.Spec.Template, image)
		if err != nil {
			return service, err
		}
		return updateEnvVars(service, envVars)
	}
}

func updateEnvVars(service *servingv1.Service, envVars map[string]string) (*servingv1.Service, error) {
	builtEnvVarName := "BUILT"
	builtEnvVarValue := time.Now().Format("20060102T150405")

	toUpdate := make(map[string]string, len(envVars)+1)
	toRemove := make([]string, 0)

	for name, value := range envVars {
		if strings.HasSuffix(name, "-") {
			toRemove = append(toRemove, strings.TrimSuffix(name, "-"))
		} else {
			toUpdate[name] = value
		}
	}

	toUpdate[builtEnvVarName] = builtEnvVarValue

	return service, servinglib.UpdateEnvVars(&service.Spec.Template, toUpdate, toRemove)
}
