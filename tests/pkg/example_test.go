package pkg

import (
	"context"
	corev1 "k8s.io/api/core/v1"
	"sigs.k8s.io/e2e-framework/pkg/envconf"
	"sigs.k8s.io/e2e-framework/pkg/features"
	"testing"
)

func TestListNamespaces(t *testing.T) {
	f := features.New("pod list").
		Assess("pods from kube-system", func(ctx context.Context, t *testing.T, cfg *envconf.Config) context.Context {
			var nss corev1.NamespaceList
			err := cfg.Client().Resources().List(context.TODO(), &nss)
			if err != nil {
				t.Fatal(err)
			}
			t.Logf("found %d namespaces", len(nss.Items))
			if len(nss.Items) == 0 {
				t.Fatal("no namespaces in cluster")
			}
			return ctx
		})

	testenv.Test(t, f.Feature())
}
