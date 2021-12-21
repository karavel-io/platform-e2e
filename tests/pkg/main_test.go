package tests

import (
	"context"
	"fmt"
	corev1 "k8s.io/api/core/v1"
	"os"
	"testing"

	helper "github.com/vmware-tanzu/sonobuoy-plugins/plugin-helper"
	"sigs.k8s.io/e2e-framework/pkg/env"
	"sigs.k8s.io/e2e-framework/pkg/envconf"
)

const (
	ProgressReporterCtxKey = "SONOBUOY_PROGRESS_REPORTER"
	TestNsCtxKey           = "E2E_TEST_NS"
)

var testenv env.Environment

func TestMain(m *testing.M) {
	cfg, err := envconf.NewFromFlags()
	if err != nil {
		panic(err)
	}

	testenv = env.NewWithConfig(cfg)

	// Specifying a run ID so that multiple runs wouldn't collide. Allow a prefix to be set via env var
	// so that a plugin configuration (yaml file) can easily set that without code changes.
	nsPrefix := os.Getenv("NS_PREFIX")
	if nsPrefix == "" {
		nsPrefix = "karavel-e2e"
	}
	runID := envconf.RandomName(nsPrefix, 4)

	// Create updateReporter; will also place into context during Setup for use in features.
	updateReporter := helper.NewProgressReporter(0)

	testenv.Setup(func(ctx context.Context, cfg *envconf.Config) (context.Context, error) {
		// Try and create the client; doing it before all the tests allows the tests to assume
		// it can be created without error and they can just use cfg.Client().
		_, err := cfg.NewClient()
		if err != nil {
			return ctx, err
		}

		ctx = context.WithValue(ctx, ProgressReporterCtxKey, updateReporter)
		return createNSForTest(ctx, cfg, runID)
	})

	testenv.BeforeEachTest(func(ctx context.Context, cfg *envconf.Config, t *testing.T) (context.Context, error) {
		updateReporter.StartTest(t.Name())
		return ctx, nil
	})
	testenv.AfterEachTest(func(ctx context.Context, cfg *envconf.Config, t *testing.T) (context.Context, error) {
		updateReporter.StopTest(t.Name(), t.Failed(), t.Skipped(), nil)
		return ctx, nil
	})

	testenv.Finish(func(ctx context.Context, cfg *envconf.Config) (context.Context, error) {
		return deleteNSForTest(ctx, cfg)
	})

	os.Exit(testenv.Run(m))
}

// CreateNSForTest creates a random namespace with the runID as a prefix. It is stored in the context
// so that the deleteNSForTest routine can look it up and delete it.
func createNSForTest(ctx context.Context, cfg *envconf.Config, runID string) (context.Context, error) {
	ns := envconf.RandomName(runID, 10)
	ctx = context.WithValue(ctx, TestNsCtxKey, ns)

	nsObj := corev1.Namespace{}
	nsObj.Name = ns
	return ctx, cfg.Client().Resources().Create(ctx, &nsObj)
}

// DeleteNSForTest looks up the namespace corresponding to the given test and deletes it.
func deleteNSForTest(ctx context.Context, cfg *envconf.Config) (context.Context, error) {
	ns := fmt.Sprint(ctx.Value(TestNsCtxKey))

	nsObj := corev1.Namespace{}
	nsObj.Name = ns
	return ctx, cfg.Client().Resources().Delete(ctx, &nsObj)
}
