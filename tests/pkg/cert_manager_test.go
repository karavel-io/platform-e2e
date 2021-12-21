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

func TestCertManager(t *testing.T) {
	f := features.New("cert-manager").WithLabel("provider", "all").
		Assess("operator is running", func(ctx context.Context, t *testing.T, cfg *envconf.Config) context.Context {
			client := cfg.Client().Resources("cert-manager")

			pods := corev1.PodList{}
			err := client.List(ctx, &pods, func(opts *v1.ListOptions) {
				opts.LabelSelector = "karavel.io/component-name=cert-manager"
			})
			assert.NoError(t, err)

			assert.Equal(t, 3, len(pods.Items), "expected 3 pods for cert-manager")

			for _, pod := range pods.Items {
				assert.Equal(t, pod.Status.Phase, corev1.PodRunning, "pod should be running")
			}
			return ctx
		}).
		Assess("letsencrypt cluster issuers are healthy", func(ctx context.Context, t *testing.T, cfg *envconf.Config) context.Context {
			client := cfg.Client().Resources()

			for _, iss := range []string{"prod", "staging"} {

				issuer := unstructured.Unstructured{}
				issuer.SetGroupVersionKind(schema.GroupVersionKind{
					Group:   "cert-manager.io",
					Version: "v1",
					Kind:    "ClusterIssuer",
				})
				err := client.Get(ctx, "letsencrypt-"+iss, "", &issuer)
				assert.NoError(t, err)

				conditions, found, err := unstructured.NestedSlice(issuer.UnstructuredContent(), "status", "conditions")
				assert.NoError(t, err)
				assert.True(t, found)
				assert.True(t, len(conditions) > 0)
				lastCond := conditions[0].(map[string]interface{})
				status, found, err := unstructured.NestedString(lastCond, "status")
				assert.NoError(t, err)
				assert.True(t, found)
				typ, found, err := unstructured.NestedString(lastCond, "type")
				assert.NoError(t, err)
				assert.True(t, found)
				assert.Equal(t, status, "True", "issuer status should be 'True'")
				assert.Equal(t, typ, "Ready", "issuer type should be 'Ready'")
			}

			return ctx
		}).
		Feature()

	testenv.Test(t, f)
}
