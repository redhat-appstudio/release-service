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
	"go/build"
	"path/filepath"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"k8s.io/client-go/rest"
	ctrl "sigs.k8s.io/controller-runtime"

	"sigs.k8s.io/controller-runtime/pkg/client"

	"sigs.k8s.io/controller-runtime/pkg/envtest"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"

	appstudioshared "github.com/redhat-appstudio/managed-gitops/appstudio-shared/apis/appstudio.redhat.com/v1alpha1"
	appstudiov1alpha1 "github.com/redhat-appstudio/release-service/api/v1alpha1"
	tektonv1beta1 "github.com/tektoncd/pipeline/pkg/apis/pipeline/v1beta1"
	clientsetscheme "k8s.io/client-go/kubernetes/scheme"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
)

const (
	testApiVersion            = "appstudio.redhat.com/v1alpha1"
	testNamespace             = "default"
	testPersistentVolumeClaim = "test-volume"
	testServiceAccountName    = "test-account"
)

var (
	cfg         *rest.Config
	k8sManager  ctrl.Manager
	cacheClient client.Client
	testEnv     *envtest.Environment
	ctx         context.Context
	cancel      context.CancelFunc
)

func TestControllersRelease(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Release Controller Test Suite")
}

var _ = BeforeSuite(func() {
	logf.SetLogger(zap.New(zap.WriteTo(GinkgoWriter), zap.UseDevMode(true)))
	ctx, cancel = context.WithCancel(context.TODO())

	// adding required CRDs, including tekton for PipelineRun Kind
	testEnv = &envtest.Environment{
		CRDDirectoryPaths: []string{
			filepath.Join("..", "..", "config", "crd", "bases"),
			filepath.Join(
				build.Default.GOPATH,
				"pkg", "mod", "github.com", "tektoncd",
				"pipeline@v0.32.2", "config",
			),
			filepath.Join(
				build.Default.GOPATH,
				"pkg", "mod", "github.com", "redhat-appstudio", "managed-gitops",
				"appstudio-shared@v0.0.0-20220603115212-1fb4d804a8c2", "config", "crd", "bases",
			),
		},
		ErrorIfCRDPathMissing: true,
	}

	var err error
	// cfg is defined in this file globally.
	cfg, err = testEnv.Start()
	Expect(err).NotTo(HaveOccurred())
	Expect(cfg).NotTo(BeNil())

	err = appstudiov1alpha1.AddToScheme(clientsetscheme.Scheme)
	Expect(err).NotTo(HaveOccurred())

	err = tektonv1beta1.AddToScheme(clientsetscheme.Scheme)
	Expect(err).NotTo(HaveOccurred())

	err = appstudioshared.AddToScheme(clientsetscheme.Scheme)
	Expect(err).NotTo(HaveOccurred())

	//webhookInstallOptions := &testEnv.WebhookInstallOptions

	k8sManager, _ = ctrl.NewManager(cfg, ctrl.Options{
		Scheme:             clientsetscheme.Scheme,
		MetricsBindAddress: "0", // this disables metrics
		LeaderElection:     false,
	})

	cacheClient = k8sManager.GetClient()
	go func() {
		defer GinkgoRecover()

		err := setupCache(k8sManager)
		Expect(err).ToNot(HaveOccurred())

		err = k8sManager.Start(ctx)
		Expect(err).ToNot(HaveOccurred())
	}()

}, 60)

var _ = AfterSuite(func() {
	cancel()
	By("tearing down the test environment")
	err := testEnv.Stop()
	Expect(err).NotTo(HaveOccurred())
})
