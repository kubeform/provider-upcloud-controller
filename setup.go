/*
Copyright AppsCode Inc. and Contributors

Licensed under the AppsCode Community License 1.0.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    https://github.com/appscode/licenses/raw/1.0.0/AppsCode-Community-1.0.0.md

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

// Code generated by Kubeform. DO NOT EDIT.

package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"sync"
	"time"

	upcloud "github.com/UpCloudLtd/terraform-provider-upcloud/upcloud"
	"github.com/gobuffalo/flect"
	auditlib "go.bytebuilders.dev/audit/lib"
	arv1 "k8s.io/api/admissionregistration/v1"
	"k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset"
	informers "k8s.io/apiextensions-apiserver/pkg/client/informers/externalversions"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	admissionregistrationv1 "k8s.io/client-go/kubernetes/typed/admissionregistration/v1"
	"k8s.io/client-go/tools/cache"
	"k8s.io/klog/v2"
	firewallv1alpha1 "kubeform.dev/provider-upcloud-api/apis/firewall/v1alpha1"
	floatingv1alpha1 "kubeform.dev/provider-upcloud-api/apis/floating/v1alpha1"
	managedv1alpha1 "kubeform.dev/provider-upcloud-api/apis/managed/v1alpha1"
	networkv1alpha1 "kubeform.dev/provider-upcloud-api/apis/network/v1alpha1"
	objectv1alpha1 "kubeform.dev/provider-upcloud-api/apis/object/v1alpha1"
	routerv1alpha1 "kubeform.dev/provider-upcloud-api/apis/router/v1alpha1"
	serverv1alpha1 "kubeform.dev/provider-upcloud-api/apis/server/v1alpha1"
	storagev1alpha1 "kubeform.dev/provider-upcloud-api/apis/storage/v1alpha1"
	tagv1alpha1 "kubeform.dev/provider-upcloud-api/apis/tag/v1alpha1"
	controllersfirewall "kubeform.dev/provider-upcloud-controller/controllers/firewall"
	controllersfloating "kubeform.dev/provider-upcloud-controller/controllers/floating"
	controllersmanaged "kubeform.dev/provider-upcloud-controller/controllers/managed"
	controllersnetwork "kubeform.dev/provider-upcloud-controller/controllers/network"
	controllersobject "kubeform.dev/provider-upcloud-controller/controllers/object"
	controllersrouter "kubeform.dev/provider-upcloud-controller/controllers/router"
	controllersserver "kubeform.dev/provider-upcloud-controller/controllers/server"
	controllersstorage "kubeform.dev/provider-upcloud-controller/controllers/storage"
	controllerstag "kubeform.dev/provider-upcloud-controller/controllers/tag"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/manager"
)

var _provider = upcloud.Provider()

var runningControllers = struct {
	sync.RWMutex
	mp map[schema.GroupVersionKind]bool
}{mp: make(map[schema.GroupVersionKind]bool)}

func watchCRD(ctx context.Context, crdClient *clientset.Clientset, vwcClient *admissionregistrationv1.AdmissionregistrationV1Client, stopCh <-chan struct{}, mgr manager.Manager, auditor *auditlib.EventPublisher, restrictToNamespace string) error {
	informerFactory := informers.NewSharedInformerFactory(crdClient, time.Second*30)
	i := informerFactory.Apiextensions().V1().CustomResourceDefinitions().Informer()
	l := informerFactory.Apiextensions().V1().CustomResourceDefinitions().Lister()

	i.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			var key string
			key, err := cache.MetaNamespaceKeyFunc(obj)
			if err != nil {
				klog.Error(err)
				return
			}

			_, name, err := cache.SplitMetaNamespaceKey(key)
			if err != nil {
				klog.Error(err)
				return
			}

			crd, err := l.Get(name)
			if err != nil {
				klog.Error(err)
				return
			}
			if strings.Contains(crd.Spec.Group, "upcloud.kubeform.com") {
				gvk := schema.GroupVersionKind{
					Group:   crd.Spec.Group,
					Version: crd.Spec.Versions[0].Name,
					Kind:    crd.Spec.Names.Kind,
				}

				// check whether this gvk came before, if no then start the controller
				runningControllers.RLock()
				_, ok := runningControllers.mp[gvk]
				runningControllers.RUnlock()

				if !ok {
					runningControllers.Lock()
					runningControllers.mp[gvk] = true
					runningControllers.Unlock()

					if enableValidatingWebhook {
						// add dynamic ValidatingWebhookConfiguration

						// create empty VWC if the group has come for the first time
						err := createEmptyVWC(vwcClient, gvk)
						if err != nil {
							klog.Error(err)
							return
						}

						// update
						err = updateVWC(vwcClient, gvk)
						if err != nil {
							klog.Error(err)
							return
						}

						err = SetupWebhook(mgr, gvk)
						if err != nil {
							setupLog.Error(err, "unable to enable webhook")
							os.Exit(1)
						}
					}

					err = SetupManager(ctx, mgr, gvk, auditor, restrictToNamespace)
					if err != nil {
						setupLog.Error(err, "unable to start manager")
						os.Exit(1)
					}
				}
			}
		},
	})

	informerFactory.Start(stopCh)

	return nil
}

func createEmptyVWC(vwcClient *admissionregistrationv1.AdmissionregistrationV1Client, gvk schema.GroupVersionKind) error {
	vwcName := strings.ReplaceAll(strings.ToLower(gvk.Group), ".", "-")
	_, err := vwcClient.ValidatingWebhookConfigurations().Get(context.TODO(), vwcName, metav1.GetOptions{})
	if err == nil || !(errors.IsNotFound(err)) {
		return err
	}

	emptyVWC := &arv1.ValidatingWebhookConfiguration{
		TypeMeta: metav1.TypeMeta{
			Kind:       "ValidatingWebhookConfiguration",
			APIVersion: "admissionregistration.k8s.io/v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: strings.ReplaceAll(strings.ToLower(gvk.Group), ".", "-"),
			Labels: map[string]string{
				"app.kubernetes.io/instance": "upcloud.kubeform.com",
				"app.kubernetes.io/part-of":  "kubeform.com",
			},
		},
	}
	_, err = vwcClient.ValidatingWebhookConfigurations().Create(context.TODO(), emptyVWC, metav1.CreateOptions{})
	if err != nil {
		return err
	}

	return nil
}

func updateVWC(vwcClient *admissionregistrationv1.AdmissionregistrationV1Client, gvk schema.GroupVersionKind) error {
	vwcName := strings.ReplaceAll(strings.ToLower(gvk.Group), ".", "-")
	vwc, err := vwcClient.ValidatingWebhookConfigurations().Get(context.TODO(), vwcName, metav1.GetOptions{})
	if err != nil {
		return err
	}

	path := "/validate-" + strings.ReplaceAll(strings.ToLower(gvk.Group), ".", "-") + "-v1alpha1-" + strings.ToLower(gvk.Kind)
	fail := arv1.Fail
	sideEffects := arv1.SideEffectClassNone
	admissionReviewVersions := []string{"v1beta1"}

	rules := []arv1.RuleWithOperations{
		{
			Operations: []arv1.OperationType{
				arv1.Delete,
				arv1.Update,
			},
			Rule: arv1.Rule{
				APIGroups:   []string{strings.ToLower(gvk.Group)},
				APIVersions: []string{gvk.Version},
				Resources:   []string{strings.ToLower(flect.Pluralize(gvk.Kind))},
			},
		},
	}

	data, err := ioutil.ReadFile("/tmp/k8s-webhook-server/serving-certs/ca.crt")
	if err != nil {
		return err
	}

	name := strings.ToLower(gvk.Kind) + "." + gvk.Group
	for _, webhook := range vwc.Webhooks {
		if webhook.Name == name {
			return nil
		}
	}

	newWebhook := arv1.ValidatingWebhook{
		Name: name,
		ClientConfig: arv1.WebhookClientConfig{
			Service: &arv1.ServiceReference{
				Namespace: webhookNamespace,
				Name:      webhookName,
				Path:      &path,
			},
			CABundle: data,
		},
		Rules:                   rules,
		FailurePolicy:           &fail,
		SideEffects:             &sideEffects,
		AdmissionReviewVersions: admissionReviewVersions,
	}

	vwc.Webhooks = append(vwc.Webhooks, newWebhook)

	_, err = vwcClient.ValidatingWebhookConfigurations().Update(context.TODO(), vwc, metav1.UpdateOptions{})
	if err != nil {
		return err
	}

	return nil
}

func SetupManager(ctx context.Context, mgr manager.Manager, gvk schema.GroupVersionKind, auditor *auditlib.EventPublisher, restrictToNamespace string) error {
	switch gvk {
	case schema.GroupVersionKind{
		Group:   "firewall.upcloud.kubeform.com",
		Version: "v1alpha1",
		Kind:    "Rules",
	}:
		if err := (&controllersfirewall.RulesReconciler{
			Client:   mgr.GetClient(),
			Log:      ctrl.Log.WithName("controllers").WithName("Rules"),
			Scheme:   mgr.GetScheme(),
			Gvk:      gvk,
			Provider: _provider,
			Resource: _provider.ResourcesMap["upcloud_firewall_rules"],
			TypeName: "upcloud_firewall_rules",
		}).SetupWithManager(ctx, mgr, auditor, restrictToNamespace); err != nil {
			setupLog.Error(err, "unable to create controller", "controller", "Rules")
			return err
		}
	case schema.GroupVersionKind{
		Group:   "floating.upcloud.kubeform.com",
		Version: "v1alpha1",
		Kind:    "IpAddress",
	}:
		if err := (&controllersfloating.IpAddressReconciler{
			Client:   mgr.GetClient(),
			Log:      ctrl.Log.WithName("controllers").WithName("IpAddress"),
			Scheme:   mgr.GetScheme(),
			Gvk:      gvk,
			Provider: _provider,
			Resource: _provider.ResourcesMap["upcloud_floating_ip_address"],
			TypeName: "upcloud_floating_ip_address",
		}).SetupWithManager(ctx, mgr, auditor, restrictToNamespace); err != nil {
			setupLog.Error(err, "unable to create controller", "controller", "IpAddress")
			return err
		}
	case schema.GroupVersionKind{
		Group:   "managed.upcloud.kubeform.com",
		Version: "v1alpha1",
		Kind:    "DatabaseLogicalDatabase",
	}:
		if err := (&controllersmanaged.DatabaseLogicalDatabaseReconciler{
			Client:   mgr.GetClient(),
			Log:      ctrl.Log.WithName("controllers").WithName("DatabaseLogicalDatabase"),
			Scheme:   mgr.GetScheme(),
			Gvk:      gvk,
			Provider: _provider,
			Resource: _provider.ResourcesMap["upcloud_managed_database_logical_database"],
			TypeName: "upcloud_managed_database_logical_database",
		}).SetupWithManager(ctx, mgr, auditor, restrictToNamespace); err != nil {
			setupLog.Error(err, "unable to create controller", "controller", "DatabaseLogicalDatabase")
			return err
		}
	case schema.GroupVersionKind{
		Group:   "managed.upcloud.kubeform.com",
		Version: "v1alpha1",
		Kind:    "DatabaseMysql",
	}:
		if err := (&controllersmanaged.DatabaseMysqlReconciler{
			Client:   mgr.GetClient(),
			Log:      ctrl.Log.WithName("controllers").WithName("DatabaseMysql"),
			Scheme:   mgr.GetScheme(),
			Gvk:      gvk,
			Provider: _provider,
			Resource: _provider.ResourcesMap["upcloud_managed_database_mysql"],
			TypeName: "upcloud_managed_database_mysql",
		}).SetupWithManager(ctx, mgr, auditor, restrictToNamespace); err != nil {
			setupLog.Error(err, "unable to create controller", "controller", "DatabaseMysql")
			return err
		}
	case schema.GroupVersionKind{
		Group:   "managed.upcloud.kubeform.com",
		Version: "v1alpha1",
		Kind:    "DatabasePostgresql",
	}:
		if err := (&controllersmanaged.DatabasePostgresqlReconciler{
			Client:   mgr.GetClient(),
			Log:      ctrl.Log.WithName("controllers").WithName("DatabasePostgresql"),
			Scheme:   mgr.GetScheme(),
			Gvk:      gvk,
			Provider: _provider,
			Resource: _provider.ResourcesMap["upcloud_managed_database_postgresql"],
			TypeName: "upcloud_managed_database_postgresql",
		}).SetupWithManager(ctx, mgr, auditor, restrictToNamespace); err != nil {
			setupLog.Error(err, "unable to create controller", "controller", "DatabasePostgresql")
			return err
		}
	case schema.GroupVersionKind{
		Group:   "managed.upcloud.kubeform.com",
		Version: "v1alpha1",
		Kind:    "DatabaseUser",
	}:
		if err := (&controllersmanaged.DatabaseUserReconciler{
			Client:   mgr.GetClient(),
			Log:      ctrl.Log.WithName("controllers").WithName("DatabaseUser"),
			Scheme:   mgr.GetScheme(),
			Gvk:      gvk,
			Provider: _provider,
			Resource: _provider.ResourcesMap["upcloud_managed_database_user"],
			TypeName: "upcloud_managed_database_user",
		}).SetupWithManager(ctx, mgr, auditor, restrictToNamespace); err != nil {
			setupLog.Error(err, "unable to create controller", "controller", "DatabaseUser")
			return err
		}
	case schema.GroupVersionKind{
		Group:   "network.upcloud.kubeform.com",
		Version: "v1alpha1",
		Kind:    "Network",
	}:
		if err := (&controllersnetwork.NetworkReconciler{
			Client:   mgr.GetClient(),
			Log:      ctrl.Log.WithName("controllers").WithName("Network"),
			Scheme:   mgr.GetScheme(),
			Gvk:      gvk,
			Provider: _provider,
			Resource: _provider.ResourcesMap["upcloud_network"],
			TypeName: "upcloud_network",
		}).SetupWithManager(ctx, mgr, auditor, restrictToNamespace); err != nil {
			setupLog.Error(err, "unable to create controller", "controller", "Network")
			return err
		}
	case schema.GroupVersionKind{
		Group:   "object.upcloud.kubeform.com",
		Version: "v1alpha1",
		Kind:    "Storage",
	}:
		if err := (&controllersobject.StorageReconciler{
			Client:   mgr.GetClient(),
			Log:      ctrl.Log.WithName("controllers").WithName("Storage"),
			Scheme:   mgr.GetScheme(),
			Gvk:      gvk,
			Provider: _provider,
			Resource: _provider.ResourcesMap["upcloud_object_storage"],
			TypeName: "upcloud_object_storage",
		}).SetupWithManager(ctx, mgr, auditor, restrictToNamespace); err != nil {
			setupLog.Error(err, "unable to create controller", "controller", "Storage")
			return err
		}
	case schema.GroupVersionKind{
		Group:   "router.upcloud.kubeform.com",
		Version: "v1alpha1",
		Kind:    "Router",
	}:
		if err := (&controllersrouter.RouterReconciler{
			Client:   mgr.GetClient(),
			Log:      ctrl.Log.WithName("controllers").WithName("Router"),
			Scheme:   mgr.GetScheme(),
			Gvk:      gvk,
			Provider: _provider,
			Resource: _provider.ResourcesMap["upcloud_router"],
			TypeName: "upcloud_router",
		}).SetupWithManager(ctx, mgr, auditor, restrictToNamespace); err != nil {
			setupLog.Error(err, "unable to create controller", "controller", "Router")
			return err
		}
	case schema.GroupVersionKind{
		Group:   "server.upcloud.kubeform.com",
		Version: "v1alpha1",
		Kind:    "Server",
	}:
		if err := (&controllersserver.ServerReconciler{
			Client:   mgr.GetClient(),
			Log:      ctrl.Log.WithName("controllers").WithName("Server"),
			Scheme:   mgr.GetScheme(),
			Gvk:      gvk,
			Provider: _provider,
			Resource: _provider.ResourcesMap["upcloud_server"],
			TypeName: "upcloud_server",
		}).SetupWithManager(ctx, mgr, auditor, restrictToNamespace); err != nil {
			setupLog.Error(err, "unable to create controller", "controller", "Server")
			return err
		}
	case schema.GroupVersionKind{
		Group:   "storage.upcloud.kubeform.com",
		Version: "v1alpha1",
		Kind:    "Storage",
	}:
		if err := (&controllersstorage.StorageReconciler{
			Client:   mgr.GetClient(),
			Log:      ctrl.Log.WithName("controllers").WithName("Storage"),
			Scheme:   mgr.GetScheme(),
			Gvk:      gvk,
			Provider: _provider,
			Resource: _provider.ResourcesMap["upcloud_storage"],
			TypeName: "upcloud_storage",
		}).SetupWithManager(ctx, mgr, auditor, restrictToNamespace); err != nil {
			setupLog.Error(err, "unable to create controller", "controller", "Storage")
			return err
		}
	case schema.GroupVersionKind{
		Group:   "tag.upcloud.kubeform.com",
		Version: "v1alpha1",
		Kind:    "Tag",
	}:
		if err := (&controllerstag.TagReconciler{
			Client:   mgr.GetClient(),
			Log:      ctrl.Log.WithName("controllers").WithName("Tag"),
			Scheme:   mgr.GetScheme(),
			Gvk:      gvk,
			Provider: _provider,
			Resource: _provider.ResourcesMap["upcloud_tag"],
			TypeName: "upcloud_tag",
		}).SetupWithManager(ctx, mgr, auditor, restrictToNamespace); err != nil {
			setupLog.Error(err, "unable to create controller", "controller", "Tag")
			return err
		}

	default:
		return fmt.Errorf("Invalid CRD")
	}

	return nil
}

func SetupWebhook(mgr manager.Manager, gvk schema.GroupVersionKind) error {
	switch gvk {
	case schema.GroupVersionKind{
		Group:   "firewall.upcloud.kubeform.com",
		Version: "v1alpha1",
		Kind:    "Rules",
	}:
		if err := (&firewallv1alpha1.Rules{}).SetupWebhookWithManager(mgr); err != nil {
			setupLog.Error(err, "unable to create webhook", "webhook", "Rules")
			return err
		}
	case schema.GroupVersionKind{
		Group:   "floating.upcloud.kubeform.com",
		Version: "v1alpha1",
		Kind:    "IpAddress",
	}:
		if err := (&floatingv1alpha1.IpAddress{}).SetupWebhookWithManager(mgr); err != nil {
			setupLog.Error(err, "unable to create webhook", "webhook", "IpAddress")
			return err
		}
	case schema.GroupVersionKind{
		Group:   "managed.upcloud.kubeform.com",
		Version: "v1alpha1",
		Kind:    "DatabaseLogicalDatabase",
	}:
		if err := (&managedv1alpha1.DatabaseLogicalDatabase{}).SetupWebhookWithManager(mgr); err != nil {
			setupLog.Error(err, "unable to create webhook", "webhook", "DatabaseLogicalDatabase")
			return err
		}
	case schema.GroupVersionKind{
		Group:   "managed.upcloud.kubeform.com",
		Version: "v1alpha1",
		Kind:    "DatabaseMysql",
	}:
		if err := (&managedv1alpha1.DatabaseMysql{}).SetupWebhookWithManager(mgr); err != nil {
			setupLog.Error(err, "unable to create webhook", "webhook", "DatabaseMysql")
			return err
		}
	case schema.GroupVersionKind{
		Group:   "managed.upcloud.kubeform.com",
		Version: "v1alpha1",
		Kind:    "DatabasePostgresql",
	}:
		if err := (&managedv1alpha1.DatabasePostgresql{}).SetupWebhookWithManager(mgr); err != nil {
			setupLog.Error(err, "unable to create webhook", "webhook", "DatabasePostgresql")
			return err
		}
	case schema.GroupVersionKind{
		Group:   "managed.upcloud.kubeform.com",
		Version: "v1alpha1",
		Kind:    "DatabaseUser",
	}:
		if err := (&managedv1alpha1.DatabaseUser{}).SetupWebhookWithManager(mgr); err != nil {
			setupLog.Error(err, "unable to create webhook", "webhook", "DatabaseUser")
			return err
		}
	case schema.GroupVersionKind{
		Group:   "network.upcloud.kubeform.com",
		Version: "v1alpha1",
		Kind:    "Network",
	}:
		if err := (&networkv1alpha1.Network{}).SetupWebhookWithManager(mgr); err != nil {
			setupLog.Error(err, "unable to create webhook", "webhook", "Network")
			return err
		}
	case schema.GroupVersionKind{
		Group:   "object.upcloud.kubeform.com",
		Version: "v1alpha1",
		Kind:    "Storage",
	}:
		if err := (&objectv1alpha1.Storage{}).SetupWebhookWithManager(mgr); err != nil {
			setupLog.Error(err, "unable to create webhook", "webhook", "Storage")
			return err
		}
	case schema.GroupVersionKind{
		Group:   "router.upcloud.kubeform.com",
		Version: "v1alpha1",
		Kind:    "Router",
	}:
		if err := (&routerv1alpha1.Router{}).SetupWebhookWithManager(mgr); err != nil {
			setupLog.Error(err, "unable to create webhook", "webhook", "Router")
			return err
		}
	case schema.GroupVersionKind{
		Group:   "server.upcloud.kubeform.com",
		Version: "v1alpha1",
		Kind:    "Server",
	}:
		if err := (&serverv1alpha1.Server{}).SetupWebhookWithManager(mgr); err != nil {
			setupLog.Error(err, "unable to create webhook", "webhook", "Server")
			return err
		}
	case schema.GroupVersionKind{
		Group:   "storage.upcloud.kubeform.com",
		Version: "v1alpha1",
		Kind:    "Storage",
	}:
		if err := (&storagev1alpha1.Storage{}).SetupWebhookWithManager(mgr); err != nil {
			setupLog.Error(err, "unable to create webhook", "webhook", "Storage")
			return err
		}
	case schema.GroupVersionKind{
		Group:   "tag.upcloud.kubeform.com",
		Version: "v1alpha1",
		Kind:    "Tag",
	}:
		if err := (&tagv1alpha1.Tag{}).SetupWebhookWithManager(mgr); err != nil {
			setupLog.Error(err, "unable to create webhook", "webhook", "Tag")
			return err
		}

	default:
		return fmt.Errorf("Invalid Webhook")
	}

	return nil
}
