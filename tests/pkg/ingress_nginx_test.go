package tests

import (
	"context"
	"fmt"
	"github.com/stretchr/testify/assert"
	corev1 "k8s.io/api/core/v1"
	netv1 "k8s.io/api/networking/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"sigs.k8s.io/e2e-framework/pkg/envconf"
	"sigs.k8s.io/e2e-framework/pkg/features"
	"testing"
	"time"
)

//goland:noinspection GoUnreachableCode
func TestIngressNginx(t *testing.T) {
	f := features.New("ingress-nginx").WithLabel("provider", "all").
		Assess("controller is running", func(ctx context.Context, t *testing.T, cfg *envconf.Config) context.Context {
			client := cfg.Client().Resources("ingress-nginx")

			pods := corev1.PodList{}
			if err := client.List(ctx, &pods, func(opts *v1.ListOptions) {
				opts.LabelSelector = "karavel.io/component-name=ingress-nginx,app.kubernetes.io/component=controller"
			}); err != nil {
				t.Fatal(err)
			}

			assert.GreaterOrEqual(t, len(pods.Items), 3, "expected at least 3 replicas for the nginx controller")

			for _, pod := range pods.Items {
				assert.Equal(t, pod.Status.Phase, corev1.PodRunning, "pod should be running")
			}
			return ctx
		}).
		Assess("default backend is running", func(ctx context.Context, t *testing.T, cfg *envconf.Config) context.Context {
			client := cfg.Client().Resources("ingress-nginx")

			pods := corev1.PodList{}
			err := client.List(ctx, &pods, func(opts *v1.ListOptions) {
				opts.LabelSelector = "karavel.io/component-name=ingress-nginx,app.kubernetes.io/component=default-backend"
			})
			assert.NoError(t, err)

			assert.GreaterOrEqual(t, len(pods.Items), 1, "default backend is not running")

			for _, pod := range pods.Items {
				assert.Equal(t, pod.Status.Phase, corev1.PodRunning, "pod should be running")
			}
			return ctx
		}).
		Assess("load balancer service is ready", func(ctx context.Context, t *testing.T, cfg *envconf.Config) context.Context {
			client := cfg.Client().Resources("ingress-nginx")

			svc := corev1.Service{}
			err := client.Get(ctx, "ingress-nginx-controller", "ingress-nginx", &svc)
			assert.NoError(t, err)

			assert.GreaterOrEqual(t, len(svc.Status.LoadBalancer.Ingress), 1)

			for _, ing := range svc.Status.LoadBalancer.Ingress {
				endpoint := ing.IP
				if endpoint == "" {
					endpoint = ing.Hostname
				}

				assert.NotEmpty(t, endpoint, "load balancer should have an endpoint")
			}

			return ctx
		}).
		Assess("ingress class exists", func(ctx context.Context, t *testing.T, cfg *envconf.Config) context.Context {
			client := cfg.Client().Resources()

			cls := netv1.IngressClass{}
			err := client.Get(ctx, "nginx", "", &cls)
			assert.NoError(t, err)

			assert.Equal(t, "k8s.io/ingress-nginx", cls.Spec.Controller)

			return ctx
		}).
		Assess("can process ingresses", func(ctx context.Context, t *testing.T, cfg *envconf.Config) context.Context {
			ns := fmt.Sprint(ctx.Value(TestNsCtxKey))
			client := cfg.Client().Resources(ns)

			meta := v1.ObjectMeta{
				Name:      "karavel-e2e-test",
				Namespace: ns,
			}
			svc := corev1.Service{
				ObjectMeta: meta,
				Spec: corev1.ServiceSpec{
					ExternalName: "example.com",
					Ports: []corev1.ServicePort{
						{
							Name:       "https",
							Port:       443,
							TargetPort: intstr.FromInt(443),
							Protocol:   corev1.ProtocolTCP,
						},
					},
				},
			}
			cls := "nginx"
			prefix := netv1.PathTypePrefix
			ing := netv1.Ingress{
				ObjectMeta: meta,
				Spec: netv1.IngressSpec{
					IngressClassName: &cls,
					Rules: []netv1.IngressRule{
						{
							Host: "test.example.com",
							IngressRuleValue: netv1.IngressRuleValue{
								HTTP: &netv1.HTTPIngressRuleValue{
									Paths: []netv1.HTTPIngressPath{
										{
											Path:     "/",
											PathType: &prefix,
											Backend: netv1.IngressBackend{
												Service: &netv1.IngressServiceBackend{
													Name: meta.Name,
													Port: netv1.ServiceBackendPort{
														Name: "https",
													},
												},
											},
										},
									},
								},
							},
						},
					},
				},
			}
			err := client.Create(ctx, &svc)
			assert.NoError(t, err)
			err = client.Create(ctx, &ing)
			assert.NoError(t, err)

			ticker := time.NewTicker(10 * time.Second)
			waitCtx, cancel := context.WithTimeout(ctx, 60*time.Second)
			t.Log("waiting for test Ingress to come online, this can take up to a minute")
			done := make(chan error)
			go func(ctx context.Context, log func(format string, args ...interface{})) {
				for {
					select {
					case <-waitCtx.Done():
						return
					case <-ticker.C:
						err = client.Get(ctx, meta.Name, meta.Namespace, &ing)
						if err != nil {
							log("error %s", err)
							done <- err
							return
						}
						if len(ing.Status.LoadBalancer.Ingress) > 0 {
							log("Ingress is ready")
							cancel()
							done <- nil
							return
						}
						log("Ingress not ready, waiting...")
					}
				}
			}(waitCtx, t.Logf)

			err = <-done
			assert.NoError(t, err)

			assert.Greater(t, len(ing.Status.LoadBalancer.Ingress), 0, "Ingress has no load balancer endpoint")

			for _, lb := range svc.Status.LoadBalancer.Ingress {
				endpoint := lb.IP
				if endpoint == "" {
					endpoint = lb.Hostname
				}

				assert.NotEmpty(t, endpoint, "load balancer should have an endpoint")
			}
			_ = client.Delete(ctx, &svc)
			_ = client.Delete(ctx, &ing)
			return ctx
		}).
		Feature()

	testenv.Test(t, f)
}
