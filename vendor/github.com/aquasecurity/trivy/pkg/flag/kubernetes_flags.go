package flag

var (
	ClusterContextFlag = Flag{
		Name:       "context",
		ConfigName: "kubernetes.context",
		Value:      "",
		Usage:      "specify a context to scan",
	}
	K8sNamespaceFlag = Flag{
		Name:       "namespace",
		ConfigName: "kubernetes.namespace",
		Shorthand:  "n",
		Value:      "",
		Usage:      "specify a namespace to scan",
	}
	KubeConfigFlag = Flag{
		Name:       "kubeconfig",
		ConfigName: "kubernetes.kubeconfig",
		Value:      "",
		Usage:      "specify the kubeconfig file path to use",
	}
)

type K8sFlagGroup struct {
	ClusterContext *Flag
	Namespace      *Flag
	KubeConfig     *Flag
}

type K8sOptions struct {
	ClusterContext string
	Namespace      string
	KubeConfig     string
}

func NewK8sFlagGroup() *K8sFlagGroup {
	return &K8sFlagGroup{
		ClusterContext: &ClusterContextFlag,
		Namespace:      &K8sNamespaceFlag,
		KubeConfig:     &KubeConfigFlag,
	}
}

func (f *K8sFlagGroup) Name() string {
	return "Kubernetes"
}

func (f *K8sFlagGroup) Flags() []*Flag {
	return []*Flag{f.ClusterContext, f.Namespace, f.KubeConfig}
}

func (f *K8sFlagGroup) ToOptions() K8sOptions {
	return K8sOptions{
		ClusterContext: getString(f.ClusterContext),
		Namespace:      getString(f.Namespace),
		KubeConfig:     getString(f.KubeConfig),
	}
}
