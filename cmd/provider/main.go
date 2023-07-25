/*
Copyright 2021 The Crossplane Authors.

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

package main

import (
	"os"
	"path/filepath"
	"time"

	xpcontroller "github.com/crossplane/crossplane-runtime/pkg/controller"
	"github.com/crossplane/crossplane-runtime/pkg/logging"
	"github.com/crossplane/crossplane-runtime/pkg/ratelimiter"
	ujcontroller "github.com/upbound/upjet/pkg/controller"
	"github.com/upbound/upjet/pkg/terraform"
	"gopkg.in/alecthomas/kingpin.v2"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"

	"github.com/crossplane-contrib/provider-jet-mongodbatlas/apis"
	"github.com/crossplane-contrib/provider-jet-mongodbatlas/config"
	"github.com/crossplane-contrib/provider-jet-mongodbatlas/internal/clients"
	"github.com/crossplane-contrib/provider-jet-mongodbatlas/internal/controller"
)

func main() {
	var (
		app              = kingpin.New(filepath.Base(os.Args[0]), "Terraform based Crossplane provider for MongoDBAtlas").DefaultEnvars()
		debug            = app.Flag("debug", "Run with debug logging.").Short('d').Bool()
		syncPeriod       = app.Flag("sync", "Controller manager sync period such as 300ms, 1.5h, or 2h45m").Short('s').Default("1h").Duration()
		leaderElection   = app.Flag("leader-election", "Use leader election for the controller manager.").Short('l').Default("false").OverrideDefaultFromEnvar("LEADER_ELECTION").Bool()
		terraformVersion = app.Flag("terraform-version", "Terraform version.").Required().Envar("TERRAFORM_VERSION").String()
		providerSource   = app.Flag("terraform-provider-source", "Terraform provider source.").Required().Envar("TERRAFORM_PROVIDER_SOURCE").String()
		providerVersion  = app.Flag("terraform-provider-version", "Terraform provider version.").Required().Envar("TERRAFORM_PROVIDER_VERSION").String()
		maxReconcileRate = app.Flag("max-reconcile-rate", "The global maximum rate per second at which resources may checked for drift from the desired state.").Default("10").Int()
	)
	kingpin.MustParse(app.Parse(os.Args[1:]))

	zl := zap.New(zap.UseDevMode(*debug))
	log := logging.NewLogrLogger(zl.WithName("provider-jet-mongodbatlas"))
	if *debug {
		// The controller-runtime runs with a no-op logger by default. It is
		// *very* verbose even at info level, so we only provide it a real
		// logger when we're running in debug mode.
		ctrl.SetLogger(zl)
	}

	log.Debug("Starting", "sync-period", syncPeriod.String())

	cfg, err := ctrl.GetConfig()
	kingpin.FatalIfError(err, "Cannot get API server rest config")

	mgr, err := ctrl.NewManager(cfg, ctrl.Options{
		LeaderElection:     *leaderElection,
		LeaderElectionID:   "crossplane-leader-election-provider-jet-mongodbatlas",
		SyncPeriod:         syncPeriod,
		MetricsBindAddress: "0.0.0.0:8081",
	})
	kingpin.FatalIfError(err, "Cannot create controller manager")

	o := ujcontroller.Options{
		Options: xpcontroller.Options{
			Logger:                  log,
			GlobalRateLimiter:       ratelimiter.NewGlobal(*maxReconcileRate),
			PollInterval:            1 * time.Minute,
			MaxConcurrentReconciles: 1,
		},
		Provider: config.GetProvider(),
		// use the following WorkspaceStoreOption to enable the shared gRPC mode
		// terraform.WithProviderRunner(terraform.NewSharedProvider(log, os.Getenv("TERRAFORM_NATIVE_PROVIDER_PATH"), terraform.WithNativeProviderArgs("-debuggable")))
		WorkspaceStore: terraform.NewWorkspaceStore(log),
		SetupFn:        clients.TerraformSetupBuilder(*terraformVersion, *providerSource, *providerVersion),
	}

	kingpin.FatalIfError(apis.AddToScheme(mgr.GetScheme()), "Cannot add MongoDBAtlas APIs to scheme")
	kingpin.FatalIfError(controller.Setup(mgr, o), "Cannot setup MongoDBAtlas controllers")
	kingpin.FatalIfError(mgr.Start(ctrl.SetupSignalHandler()), "Cannot start controller manager")
}
