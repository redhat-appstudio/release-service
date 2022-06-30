/*
Copyright 2022.

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

package release

import (
	"context"
	"reflect"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	ctrl "sigs.k8s.io/controller-runtime"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	appstudiov1alpha1 "github.com/redhat-appstudio/release-service/api/v1alpha1"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	clientsetscheme "k8s.io/client-go/kubernetes/scheme"
)

var _ = Describe("Release Controller", func() {
	var (
		release    *appstudiov1alpha1.Release
		reconciler *Reconciler
		scheme     runtime.Scheme
		req        ctrl.Request

		//strategy *v1alpha1.ReleaseStrategy
	)

	BeforeEach(func() {

		release = &appstudiov1alpha1.Release{
			TypeMeta: metav1.TypeMeta{
				APIVersion: testApiVersion,
				Kind:       "Release",
			},
			ObjectMeta: metav1.ObjectMeta{
				GenerateName: "myrelease-",
				Namespace:    testNamespace,
			},
			Spec: appstudiov1alpha1.ReleaseSpec{
				ApplicationSnapshot: "testsnapshot",
				ReleaseLink:         "testreleaselink",
			},
		}

		ctx := context.Background()
		Expect(cacheClient.Create(ctx, release)).Should(Succeed())

		req = ctrl.Request{
			NamespacedName: types.NamespacedName{
				Namespace: release.Namespace,
				Name:      release.Name,
			},
		}

		webhookInstallOptions := &testEnv.WebhookInstallOptions

		k8sManager, _ = ctrl.NewManager(cfg, ctrl.Options{
			Scheme:             clientsetscheme.Scheme,
			Host:               webhookInstallOptions.LocalServingHost,
			Port:               webhookInstallOptions.LocalServingPort,
			CertDir:            webhookInstallOptions.LocalServingCertDir,
			MetricsBindAddress: "0", // this disables metrics
			LeaderElection:     false,
		})
		err := (&appstudiov1alpha1.Release{}).SetupWebhookWithManager(k8sManager)
		Expect(err).NotTo(HaveOccurred())

		err = (&appstudiov1alpha1.ReleaseLink{}).SetupWebhookWithManager(k8sManager)
		Expect(err).NotTo(HaveOccurred())

		reconciler = NewReleaseReconciler(cacheClient, &logf.Log, &scheme)
	})

	AfterEach(func() {
		ctx := context.Background()
		Expect(cacheClient.Delete(ctx, release)).Should(Succeed())
	})

	It("can create and return a new Reconciler object", func() {
		Expect(reflect.TypeOf(reconciler)).To(Equal(reflect.TypeOf(&Reconciler{})))
	})

	It("can ReconcileHandler receive an adapter and return the result for the handling operation", func() {
		adapter := NewAdapter(release, ctrl.Log, cacheClient, ctx)
		result, err := reconciler.ReconcileHandler(adapter)
		Expect(reflect.TypeOf(result)).To(Equal(reflect.TypeOf(reconcile.Result{})))
		Expect(err).To(BeNil())
	})

	It("can Reconcile function prepare the adapter and return the result of the reconcile handling operation", func() {
		result, err := reconciler.Reconcile(ctx, req)
		Expect(reflect.TypeOf(result)).To(Equal(reflect.TypeOf(reconcile.Result{})))
		Expect(err).To(BeNil())
	})

	It("can setup the cache by adding a new index field to search for ReleaseLinks", func() {
		err := setupCache(k8sManager)
		Expect(err).ToNot(HaveOccurred())
	})

	It("can setup a new controller manager with the given reconciler", func() {
		err := setupControllerWithManager(k8sManager, reconciler)
		Expect(err).NotTo(HaveOccurred())
	})

	It("can setup a new Controller manager and start it", func() {
		// the ctrl.Complete() ignores the object returned by ctrl.Build()
		// and returns `nil` in case of success. It returns an error otherwise.
		err := SetupController(k8sManager, &ctrl.Log)
		Expect(err).To(BeNil())
		go func() {
			defer GinkgoRecover()
			// can it start the controller?
			err = k8sManager.Start(ctx)
			Expect(err).NotTo(HaveOccurred())
		}()
	})
})
