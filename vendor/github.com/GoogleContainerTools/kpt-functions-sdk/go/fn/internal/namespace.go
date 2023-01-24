// Copyright 2022 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package internal

import (
	"sigs.k8s.io/kustomize/kyaml/yaml"
)

// PrecomputedIsNamespaceScoped copies the sigs.k8s.io/kustomize/kyaml/openapi precomputedIsNamespaceScoped
var PrecomputedIsNamespaceScoped = map[yaml.TypeMeta]bool{
	{APIVersion: "admissionregistration.k8s.io/v1", Kind: "MutatingWebhookConfiguration"}:        false,
	{APIVersion: "admissionregistration.k8s.io/v1", Kind: "ValidatingWebhookConfiguration"}:      false,
	{APIVersion: "admissionregistration.k8s.io/v1beta1", Kind: "MutatingWebhookConfiguration"}:   false,
	{APIVersion: "admissionregistration.k8s.io/v1beta1", Kind: "ValidatingWebhookConfiguration"}: false,
	{APIVersion: "apiextensions.k8s.io/v1", Kind: "CustomResourceDefinition"}:                    false,
	{APIVersion: "apiextensions.k8s.io/v1beta1", Kind: "CustomResourceDefinition"}:               false,
	{APIVersion: "apiregistration.k8s.io/v1", Kind: "APIService"}:                                false,
	{APIVersion: "apiregistration.k8s.io/v1beta1", Kind: "APIService"}:                           false,
	{APIVersion: "apps/v1", Kind: "ControllerRevision"}:                                          true,
	{APIVersion: "apps/v1", Kind: "DaemonSet"}:                                                   true,
	{APIVersion: "apps/v1", Kind: "Deployment"}:                                                  true,
	{APIVersion: "apps/v1", Kind: "ReplicaSet"}:                                                  true,
	{APIVersion: "apps/v1", Kind: "StatefulSet"}:                                                 true,
	{APIVersion: "autoscaling/v1", Kind: "HorizontalPodAutoscaler"}:                              true,
	{APIVersion: "autoscaling/v1", Kind: "Scale"}:                                                true,
	{APIVersion: "autoscaling/v2beta1", Kind: "HorizontalPodAutoscaler"}:                         true,
	{APIVersion: "autoscaling/v2beta2", Kind: "HorizontalPodAutoscaler"}:                         true,
	{APIVersion: "batch/v1", Kind: "CronJob"}:                                                    true,
	{APIVersion: "batch/v1", Kind: "Job"}:                                                        true,
	{APIVersion: "batch/v1beta1", Kind: "CronJob"}:                                               true,
	{APIVersion: "certificates.k8s.io/v1", Kind: "CertificateSigningRequest"}:                    false,
	{APIVersion: "certificates.k8s.io/v1beta1", Kind: "CertificateSigningRequest"}:               false,
	{APIVersion: "coordination.k8s.io/v1", Kind: "Lease"}:                                        true,
	{APIVersion: "coordination.k8s.io/v1beta1", Kind: "Lease"}:                                   true,
	{APIVersion: "discovery.k8s.io/v1", Kind: "EndpointSlice"}:                                   true,
	{APIVersion: "discovery.k8s.io/v1beta1", Kind: "EndpointSlice"}:                              true,
	{APIVersion: "events.k8s.io/v1", Kind: "Event"}:                                              true,
	{APIVersion: "events.k8s.io/v1beta1", Kind: "Event"}:                                         true,
	{APIVersion: "extensions/v1beta1", Kind: "Ingress"}:                                          true,
	{APIVersion: "flowcontrol.apiserver.k8s.io/v1beta1", Kind: "FlowSchema"}:                     false,
	{APIVersion: "flowcontrol.apiserver.k8s.io/v1beta1", Kind: "PriorityLevelConfiguration"}:     false,
	{APIVersion: "networking.k8s.io/v1", Kind: "Ingress"}:                                        true,
	{APIVersion: "networking.k8s.io/v1", Kind: "IngressClass"}:                                   false,
	{APIVersion: "networking.k8s.io/v1", Kind: "NetworkPolicy"}:                                  true,
	{APIVersion: "networking.k8s.io/v1beta1", Kind: "Ingress"}:                                   true,
	{APIVersion: "networking.k8s.io/v1beta1", Kind: "IngressClass"}:                              false,
	{APIVersion: "node.k8s.io/v1", Kind: "RuntimeClass"}:                                         false,
	{APIVersion: "node.k8s.io/v1beta1", Kind: "RuntimeClass"}:                                    false,
	{APIVersion: "policy/v1", Kind: "PodDisruptionBudget"}:                                       true,
	{APIVersion: "policy/v1beta1", Kind: "PodDisruptionBudget"}:                                  true,
	{APIVersion: "policy/v1beta1", Kind: "PodSecurityPolicy"}:                                    false,
	{APIVersion: "rbac.authorization.k8s.io/v1", Kind: "ClusterRole"}:                            false,
	{APIVersion: "rbac.authorization.k8s.io/v1", Kind: "ClusterRoleBinding"}:                     false,
	{APIVersion: "rbac.authorization.k8s.io/v1", Kind: "Role"}:                                   true,
	{APIVersion: "rbac.authorization.k8s.io/v1", Kind: "RoleBinding"}:                            true,
	{APIVersion: "rbac.authorization.k8s.io/v1beta1", Kind: "ClusterRole"}:                       false,
	{APIVersion: "rbac.authorization.k8s.io/v1beta1", Kind: "ClusterRoleBinding"}:                false,
	{APIVersion: "rbac.authorization.k8s.io/v1beta1", Kind: "Role"}:                              true,
	{APIVersion: "rbac.authorization.k8s.io/v1beta1", Kind: "RoleBinding"}:                       true,
	{APIVersion: "scheduling.k8s.io/v1", Kind: "PriorityClass"}:                                  false,
	{APIVersion: "scheduling.k8s.io/v1beta1", Kind: "PriorityClass"}:                             false,
	{APIVersion: "storage.k8s.io/v1", Kind: "CSIDriver"}:                                         false,
	{APIVersion: "storage.k8s.io/v1", Kind: "CSINode"}:                                           false,
	{APIVersion: "storage.k8s.io/v1", Kind: "StorageClass"}:                                      false,
	{APIVersion: "storage.k8s.io/v1", Kind: "VolumeAttachment"}:                                  false,
	{APIVersion: "storage.k8s.io/v1beta1", Kind: "CSIDriver"}:                                    false,
	{APIVersion: "storage.k8s.io/v1beta1", Kind: "CSINode"}:                                      false,
	{APIVersion: "storage.k8s.io/v1beta1", Kind: "CSIStorageCapacity"}:                           true,
	{APIVersion: "storage.k8s.io/v1beta1", Kind: "StorageClass"}:                                 false,
	{APIVersion: "storage.k8s.io/v1beta1", Kind: "VolumeAttachment"}:                             false,
	{APIVersion: "v1", Kind: "ComponentStatus"}:                                                  false,
	{APIVersion: "v1", Kind: "ConfigMap"}:                                                        true,
	{APIVersion: "v1", Kind: "Endpoints"}:                                                        true,
	{APIVersion: "v1", Kind: "Event"}:                                                            true,
	{APIVersion: "v1", Kind: "LimitRange"}:                                                       true,
	{APIVersion: "v1", Kind: "Namespace"}:                                                        false,
	{APIVersion: "v1", Kind: "Node"}:                                                             false,
	{APIVersion: "v1", Kind: "NodeProxyOptions"}:                                                 false,
	{APIVersion: "v1", Kind: "PersistentVolume"}:                                                 false,
	{APIVersion: "v1", Kind: "PersistentVolumeClaim"}:                                            true,
	{APIVersion: "v1", Kind: "Pod"}:                                                              true,
	{APIVersion: "v1", Kind: "PodAttachOptions"}:                                                 true,
	{APIVersion: "v1", Kind: "PodExecOptions"}:                                                   true,
	{APIVersion: "v1", Kind: "PodPortForwardOptions"}:                                            true,
	{APIVersion: "v1", Kind: "PodProxyOptions"}:                                                  true,
	{APIVersion: "v1", Kind: "PodTemplate"}:                                                      true,
	{APIVersion: "v1", Kind: "ReplicationController"}:                                            true,
	{APIVersion: "v1", Kind: "ResourceQuota"}:                                                    true,
	{APIVersion: "v1", Kind: "Secret"}:                                                           true,
	{APIVersion: "v1", Kind: "Service"}:                                                          true,
	{APIVersion: "v1", Kind: "ServiceAccount"}:                                                   true,
	{APIVersion: "v1", Kind: "ServiceProxyOptions"}:                                              true,
}
