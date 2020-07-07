package peering_request_operator

import (
	discoveryv1 "github.com/liqoTech/liqo/api/discovery/v1"
	"github.com/liqoTech/liqo/pkg/clusterID"
	"github.com/liqoTech/liqo/pkg/crdClient/v1alpha1"
	"k8s.io/apimachinery/pkg/runtime"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	"k8s.io/klog"
	"os"
	"path/filepath"
	ctrl "sigs.k8s.io/controller-runtime"
)

var (
	scheme = runtime.NewScheme()
)

func init() {
	_ = clientgoscheme.AddToScheme(scheme)
	_ = discoveryv1.AddToScheme(scheme)
	// +kubebuilder:scaffold:scheme
}

func StartOperator(namespace string, configMapName string, broadcasterImage string, broadcasterServiceAccount string) {
	mgr, err := ctrl.NewManager(ctrl.GetConfigOrDie(), ctrl.Options{
		Scheme:           scheme,
		Port:             9443,
		LeaderElection:   false,
		LeaderElectionID: "b3156c4e.liqo.io",
	})
	if err != nil {
		klog.Error(err, "unable to start manager")
		os.Exit(1)
	}

	config, err := v1alpha1.NewKubeconfig(filepath.Join(os.Getenv("HOME"), ".kube", "config"), &discoveryv1.GroupVersion)
	if err != nil {
		klog.Error(err, "unable to get kube config")
		os.Exit(1)
	}
	crdClient, err := v1alpha1.NewFromConfig(config)
	if err != nil {
		klog.Error(err, "unable to create crd client")
		os.Exit(1)
	}

	clusterId, err := clusterID.NewClusterID()
	if err != nil {
		klog.Error(err, "unable to get clusterID")
		os.Exit(1)
	}

	if err = (GetPRReconciler(
		mgr.GetScheme(),
		crdClient,
		namespace,
		clusterId,
		configMapName,
		broadcasterImage,
		broadcasterServiceAccount,
	)).SetupWithManager(mgr); err != nil {
		klog.Error(err, "unable to create controller")
		os.Exit(1)
	}
	// +kubebuilder:scaffold:builder

	if err := mgr.Start(ctrl.SetupSignalHandler()); err != nil {
		klog.Error(err, "problem running manager")
		os.Exit(1)
	}
}

func GetPRReconciler(scheme *runtime.Scheme, crdClient *v1alpha1.CRDClient, namespace string, clusterId *clusterID.ClusterID, configMapName string, broadcasterImage string, broadcasterServiceAccount string) *PeeringRequestReconciler {
	return &PeeringRequestReconciler{
		Scheme:                    scheme,
		crdClient:                 crdClient,
		Namespace:                 namespace,
		clusterId:                 clusterId,
		configMapName:             configMapName,
		broadcasterImage:          broadcasterImage,
		broadcasterServiceAccount: broadcasterServiceAccount,
	}
}