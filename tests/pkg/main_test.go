package tests

import (
	"context"
	"os"
	"testing"

	helper "github.com/vmware-tanzu/sonobuoy-plugins/plugin-helper"
	"sigs.k8s.io/e2e-framework/pkg/env"
	"sigs.k8s.io/e2e-framework/pkg/envconf"
)

const (
	ProgressReporterCtxKey = "SONOBUOY_PROGRESS_REPORTER"
)

var testenv env.Environment

func TestMain(m *testing.M) {
	cfg, err := envconf.NewFromFlags()
	if err != nil {
		panic(err)
	}

	testenv = env.NewWithConfig(cfg)

	// Create updateReporter; will also place into context during Setup for use in features.
	updateReporter := helper.NewProgressReporter(0)

	testenv.Setup(func(ctx context.Context, config *envconf.Config) (context.Context, error) {
		// Try and create the client; doing it before all the tests allows the tests to assume
		// it can be created without error and they can just use config.Client().
		_, err := config.NewClient()
		return context.WithValue(ctx, ProgressReporterCtxKey, updateReporter), err
	})

	testenv.BeforeEachTest(func(ctx context.Context, cfg *envconf.Config, t *testing.T) (context.Context, error) {
		updateReporter.StartTest(t.Name())
		return ctx, nil
	})
	testenv.AfterEachTest(func(ctx context.Context, cfg *envconf.Config, t *testing.T) (context.Context, error) {
		updateReporter.StopTest(t.Name(), t.Failed(), t.Skipped(), nil)
		return ctx, nil
	})

	os.Exit(testenv.Run(m))
}
