package tests

import (
	"context"
	"github.com/stretchr/testify/assert"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"sigs.k8s.io/e2e-framework/pkg/envconf"
	"sigs.k8s.io/e2e-framework/pkg/features"
	"testing"
)

func TestExternalSecrets(t *testing.T) {
	f := features.New("external-secrets").WithLabel("provider", "all").
		Assess("operator is running", func(ctx context.Context, t *testing.T, cfg *envconf.Config) context.Context {
			client := cfg.Client().Resources("external-secrets")

			pods := corev1.PodList{}
			if err := client.List(ctx, &pods, func(opts *v1.ListOptions) {
				opts.LabelSelector = "karavel.io/component-name=external-secrets"
			}); err != nil {
				t.Fatal(err)
			}

			assert.Equal(t, 1, len(pods.Items), "expected only one pod for the External Secrets operator")

			pod := pods.Items[0]
			assert.Equal(t, pod.Status.Phase, corev1.PodRunning, "pod should be running")
			return ctx
		}).
		Assess("default cluster secret store is healthy", func(ctx context.Context, t *testing.T, cfg *envconf.Config) context.Context {
			client := cfg.Client().Resources()

			store := unstructured.Unstructured{}
			store.SetGroupVersionKind(schema.GroupVersionKind{
				Group:   "external-secrets.io",
				Version: "v1alpha1",
				Kind:    "ClusterSecretStore",
			})
			if err := client.Get(ctx, "default", "", &store); err != nil {
				t.Fatal(err)
			}

			ts := store.GetCreationTimestamp()
			assert.False(t, (&ts).IsZero(), "default secret store should exist")

			return ctx
		}).
		Feature()

	testenv.Test(t, f)
}
