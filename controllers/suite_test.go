/*
Copyright 2021.

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

package controllers

import (
	//"os"
	"path/filepath"
	ctrl "sigs.k8s.io/controller-runtime"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/envtest"
	"sigs.k8s.io/controller-runtime/pkg/envtest/printer"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"

	batchv1 "github.com/chanshik/cronjob/api/v1"
	//+kubebuilder:scaffold:imports
)

// These tests use Ginkgo (BDD-style Go testing framework). Refer to
// http://onsi.github.io/ginkgo/ to learn more about Ginkgo.

var cfg *rest.Config
var k8sClient client.Client
var testEnv *envtest.Environment

func TestAPIs(t *testing.T) {
	RegisterFailHandler(Fail)

	RunSpecsWithDefaultAndCustomReporters(t,
		"Controller Suite",
		[]Reporter{printer.NewlineReporter{}})
}

var _ = BeforeSuite(func() {
	//os.Setenv("KUBEBUILDER_ATTACH_CONTROL_PLANE_OUTPUT", "true")
	//os.Setenv("KUBEBUILDER_CONTROLPLANE_START_TIMEOUT", "20s")
	//os.Setenv("TEST_ASSET_KUBE_APISERVER", "/opt/kubebuilder/testbin/kube-apiserver")
	//os.Setenv("TEST_ASSET_ETCD", "/opt/kubebuilder/testbin/etcd")
	//os.Setenv("TEST_ASSET_KUBECTL", "/opt/kubebuilder/testbin/kubectl")

	logf.SetLogger(zap.New(zap.WriteTo(GinkgoWriter), zap.UseDevMode(true)))

	//flags := []string{
	//	"-v=1",
	//	//"--cert-dir=/tmp/k8s/",
	//	//"--service-account-issuer=api",
	//	//"--service-account-key-file=/tmp/k8s/sa.pub",
	//	//"--service-account-signing-key-file=/tmp/k8s/sa.key",
	//	//"--bind-address=0.0.0.0",
	//	//"--advertise-address=192.168.50.193",
	//	//"--feature-gates=ServiceAccountIssuerDiscovery=false",
	//	//"--authorization-mode=AlwaysAllow",
	//	//"--insecure-port=0",
	//	//"--insecure-bind-address=127.0.0.1",
	//}
	////for _, serverFlag := range envtest.DefaultKubeAPIServerFlags {
	////	if strings.HasPrefix(serverFlag, "--insecure") {
	////		continue
	////	}
	////	if strings.HasPrefix(serverFlag, "--advertise-address") {
	////		continue
	////	}
	////
	////	flags = append(flags, serverFlag)
	////}
	//flags = append(flags, envtest.DefaultKubeAPIServerFlags...)

	By("bootstrapping test environment")
	testEnv = &envtest.Environment{
		CRDDirectoryPaths:     []string{filepath.Join("..", "config", "crd", "bases")},
		ErrorIfCRDPathMissing: true,
		//KubeAPIServerFlags:    flags,
	}

	cfg, err := testEnv.Start()
	Expect(err).NotTo(HaveOccurred())
	Expect(cfg).NotTo(BeNil())

	err = batchv1.AddToScheme(scheme.Scheme)
	Expect(err).NotTo(HaveOccurred())

	//+kubebuilder:scaffold:scheme

	k8sClient, err = client.New(cfg, client.Options{Scheme: scheme.Scheme})
	Expect(err).NotTo(HaveOccurred())
	Expect(k8sClient).NotTo(BeNil())

	k8sManager, err := ctrl.NewManager(cfg, ctrl.Options{
		Scheme: scheme.Scheme,
	})
	Expect(err).ToNot(HaveOccurred())

	err = (&CronJobReconciler{
		Client: k8sManager.GetClient(),
		Scheme: k8sManager.GetScheme(),
		Log:    ctrl.Log.WithName("controllers").WithName("CronJob"),
	}).SetupWithManager(k8sManager)
	Expect(err).ToNot(HaveOccurred())

	go func() {
		err = k8sManager.Start(ctrl.SetupSignalHandler())
		Expect(err).ToNot(HaveOccurred())
	}()

}, 60)

var _ = AfterSuite(func() {
	By("tearing down the test environment")
	err := testEnv.Stop()
	Expect(err).NotTo(HaveOccurred())
})
