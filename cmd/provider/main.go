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
	"context"
	"fmt"
	"os"
	"time"

	"github.com/alecthomas/kong"
	changelogsv1alpha1 "github.com/crossplane/crossplane-runtime/v2/apis/changelogs/proto/v1alpha1"
	xpcontroller "github.com/crossplane/crossplane-runtime/v2/pkg/controller"
	"github.com/crossplane/crossplane-runtime/v2/pkg/errors"
	"github.com/crossplane/crossplane-runtime/v2/pkg/feature"
	"github.com/crossplane/crossplane-runtime/v2/pkg/gate"
	"github.com/crossplane/crossplane-runtime/v2/pkg/logging"
	"github.com/crossplane/crossplane-runtime/v2/pkg/ratelimiter"
	"github.com/crossplane/crossplane-runtime/v2/pkg/reconciler/customresourcesgate"
	"github.com/crossplane/crossplane-runtime/v2/pkg/reconciler/managed"
	"github.com/crossplane/crossplane-runtime/v2/pkg/statemetrics"
	tjcontroller "github.com/crossplane/upjet/v2/pkg/controller"
	"github.com/crossplane/upjet/v2/pkg/terraform"
	"github.com/go-logr/logr"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	authv1 "k8s.io/api/authorization/v1"
	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/tools/leaderelection/resourcelock"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/cache"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/metrics"
	metricsserver "sigs.k8s.io/controller-runtime/pkg/metrics/server"
	"sigs.k8s.io/controller-runtime/pkg/webhook"

	apisCluster "github.com/crossplane-contrib/provider-mongodbatlas/apis/cluster"
	apisNamespaced "github.com/crossplane-contrib/provider-mongodbatlas/apis/namespaced"
	"github.com/crossplane-contrib/provider-mongodbatlas/config"
	"github.com/crossplane-contrib/provider-mongodbatlas/internal/clients"
	controllerCluster "github.com/crossplane-contrib/provider-mongodbatlas/internal/controller/cluster"
	controllerNamespaced "github.com/crossplane-contrib/provider-mongodbatlas/internal/controller/namespaced"
	"github.com/crossplane-contrib/provider-mongodbatlas/internal/features"
	"github.com/crossplane-contrib/provider-mongodbatlas/internal/version"
)

const (
	webhookTLSCertDirEnvVar = "WEBHOOK_TLS_CERT_DIR"
	tlsServerCertDirEnvVar  = "TLS_SERVER_CERTS_DIR"
	certsDirEnvVar          = "CERTS_DIR"
	tlsServerCertDir        = "/tls/server"
)

type certsDir string

func (d certsDir) BeforeApply(certsDirSet *bool) error {
	// we record whether the command-line option "--certs-dir" was supplied
	*certsDirSet = true
	return nil
}

var cli struct {
	Debug          bool `help:"Run with debug logging." short:"d"`
	LeaderElection bool `help:"Use leader election for the controller manager." short:"l" default:"false" env:"LEADER_ELECTION"`

	SyncPeriod              time.Duration `help:"Controller manager sync period such as 300ms, 1.5h, or 2h45m" short:"s" default:"1h"`
	PollInterval            time.Duration `help:"How often individual resources will be checked for drift from the desired state" default:"1m"`
	PollStateMetricInterval time.Duration `help:"State metric recording interval" default:"5s"`

	MaxReconcileRate         int           `help:"The global maximum rate per second at which resources may checked for drift from the desired state." default:"10"`
	EnableManagementPolicies bool          `help:"Enable support for Management Policies." default:"true" env:"ENABLE_MANAGEMENT_POLICIES"`
	EnableChangeLogs         bool          `help:"Enable support for capturing change logs during reconciliation." default:"false" env:"ENABLE_CHANGE_LOGS"`
	ChangelogsSocketPath     string        `help:"Path for changelogs socket (if enabled)" default:"/var/run/changelogs/changelogs.sock" env:"CHANGELOGS_SOCKET_PATH"`
	WebhookPort              int           `help:"The port the webhook listens on" default:"9443" env:"WEBHOOK_PORT"`
	MetricsBindAddress       string        `help:"The address the metrics server listens on" default:":8081" env:"METRICS_BIND_ADDRESS"`
	BrokerConnectionTimeout  time.Duration `help:"Timeout for establishing connection to Kafka brokers" default:"30s"`

	TerraformVersion string   `required:"true" help:"Terraform version" env:"TERRAFORM_VERSION"`
	ProviderSource   string   `required:"true" help:"Terraform provider source" env:"TERRAFORM_PROVIDER_SOURCE"`
	ProviderVersion  string   `required:"true" help:"Terraform provider version" env:"TERRAFORM_PROVIDER_VERSION"`
	CertsDir         certsDir `help:"The directory that contains the server key and certificate" default:"${defaultCertsDir}" env:"${defautCertsDirEnvVar}"`
}

func main() {
	certsDirSet := false
	ctx := kong.Parse(&cli,
		kong.Description("Crossplane MongoDB Atlas Provider"),
		kong.Bind(&certsDirSet),
		kong.Vars{
			"defaultCertsDir":      tlsServerCertDir,
			"defautCertsDirEnvVar": certsDirEnvVar,
		},
	)

	zl := zap.New(zap.UseDevMode(cli.Debug))
	log := logging.NewLogrLogger(zl.WithName("provider-mongodbatlas"))
	if cli.Debug {
		// The controller-runtime runs with a no-op logger by default. It is
		// *very* verbose even at info level, so we only provide it a real
		// logger when we're running in debug mode.
		ctrl.SetLogger(zl)
	} else {
		// controller-runtime v0.23+ requires SetLogger to be called, otherwise
		// background goroutines (e.g. priority queue) will warn about missing logger.
		ctrl.SetLogger(logr.Discard())
	}

	// currently, we configure the jitter to be the 5% of the poll interval
	pollJitter := time.Duration(float64(cli.PollInterval) * 0.05)

	log.Debug("Starting", "sync-period", cli.SyncPeriod.String(), "poll-interval", cli.PollInterval.String(), "poll-jitter", pollJitter, "max-reconcile-rate", cli.MaxReconcileRate)

	cfg, err := ctrl.GetConfig()
	ctx.FatalIfErrorf(err, "Cannot get API server rest config")

	// Get the TLS certs directory from the environment variables set by
	// Crossplane if they're available.
	// In older XP versions we used WEBHOOK_TLS_CERT_DIR, in newer versions
	// we use TLS_SERVER_CERTS_DIR. If an explicit certs dir is not supplied
	// via the command-line options, then these environment variables are used
	// instead.
	if !certsDirSet {
		// backwards-compatibility concerns
		xpCertsDir := os.Getenv(certsDirEnvVar)
		if xpCertsDir == "" {
			xpCertsDir = os.Getenv(tlsServerCertDirEnvVar)
		}
		if xpCertsDir == "" {
			xpCertsDir = os.Getenv(webhookTLSCertDirEnvVar)
		}
		// we probably don't need this condition but just to be on the
		// safe side, if we are missing any kong machinery details...
		if xpCertsDir != "" {
			cli.CertsDir = certsDir(xpCertsDir)
		}
	}

	mgr, err := ctrl.NewManager(cfg, ctrl.Options{
		LeaderElection:   cli.LeaderElection,
		LeaderElectionID: "crossplane-leader-election-upjet-provider-mongodbatlas",
		Cache: cache.Options{
			SyncPeriod: &cli.SyncPeriod,
		},
		Metrics: metricsserver.Options{
			BindAddress: cli.MetricsBindAddress,
		},
		WebhookServer: webhook.NewServer(
			webhook.Options{
				CertDir: string(cli.CertsDir),
				Port:    cli.WebhookPort,
			}),
		LeaderElectionResourceLock: resourcelock.LeasesResourceLock,
		LeaseDuration:              func() *time.Duration { d := 60 * time.Second; return &d }(),
		RenewDeadline:              func() *time.Duration { d := 50 * time.Second; return &d }(),
	})

	ctx.FatalIfErrorf(err, "Cannot create controller manager")
	ctx.FatalIfErrorf(apisCluster.AddToScheme(mgr.GetScheme()), "Cannot add cluster-scoped MongoDBAtlas APIs to scheme")
	ctx.FatalIfErrorf(apisNamespaced.AddToScheme(mgr.GetScheme()), "Cannot add namespaced MongoDBAtlas APIs to scheme")
	ctx.FatalIfErrorf(apiextensionsv1.AddToScheme(mgr.GetScheme()), "Cannot add api-extensions APIs to scheme")
	ctx.FatalIfErrorf(authv1.AddToScheme(mgr.GetScheme()), "Cannot add k8s authorization APIs to scheme")

	metricRecorder := managed.NewMRMetricRecorder()
	stateMetrics := statemetrics.NewMRStateMetrics()

	metrics.Registry.MustRegister(metricRecorder)
	metrics.Registry.MustRegister(stateMetrics)

	clusterOpts := tjcontroller.Options{
		Options: xpcontroller.Options{
			Logger:                  log,
			GlobalRateLimiter:       ratelimiter.NewGlobal(cli.MaxReconcileRate),
			PollInterval:            cli.PollInterval,
			MaxConcurrentReconciles: cli.MaxReconcileRate,
			Features:                &feature.Flags{},
			MetricOptions: &xpcontroller.MetricOptions{
				PollStateMetricInterval: cli.PollStateMetricInterval,
				MRMetrics:               metricRecorder,
				MRStateMetrics:          stateMetrics,
			},
		},
		Provider:       config.GetidentifierFromProvider(),
		SetupFn:        clients.TerraformSetupBuilder(cli.TerraformVersion, cli.ProviderSource, cli.ProviderVersion),
		WorkspaceStore: terraform.NewWorkspaceStore(log),
		PollJitter:     pollJitter,
		StartWebhooks:  cli.CertsDir != "",
	}

	namespacedOpts := tjcontroller.Options{
		Options: xpcontroller.Options{
			Logger:                  log,
			GlobalRateLimiter:       ratelimiter.NewGlobal(cli.MaxReconcileRate),
			PollInterval:            cli.PollInterval,
			MaxConcurrentReconciles: cli.MaxReconcileRate,
			Features:                &feature.Flags{},
			MetricOptions: &xpcontroller.MetricOptions{
				PollStateMetricInterval: cli.PollStateMetricInterval,
				MRMetrics:               metricRecorder,
				MRStateMetrics:          stateMetrics,
			},
		},
		Provider:       config.GetProviderNamespaced(),
		SetupFn:        clients.TerraformSetupBuilder(cli.TerraformVersion, cli.ProviderSource, cli.ProviderVersion),
		WorkspaceStore: terraform.NewWorkspaceStore(log),
		PollJitter:     pollJitter,
		StartWebhooks:  cli.CertsDir != "",
	}

	if cli.EnableManagementPolicies {
		clusterOpts.Features.Enable(features.EnableBetaManagementPolicies)
		namespacedOpts.Features.Enable(features.EnableBetaManagementPolicies)
		log.Info("Beta feature enabled", "flag", features.EnableBetaManagementPolicies)
	}

	if cli.EnableChangeLogs {
		clusterOpts.Features.Enable(feature.EnableAlphaChangeLogs)
		namespacedOpts.Features.Enable(feature.EnableAlphaChangeLogs)
		log.Info("Alpha feature enabled", "flag", feature.EnableAlphaChangeLogs)

		conn, err := grpc.NewClient("unix://"+cli.ChangelogsSocketPath, grpc.WithTransportCredentials(insecure.NewCredentials()))
		ctx.FatalIfErrorf(err, "failed to create change logs client connection at %s", cli.ChangelogsSocketPath)

		clo := xpcontroller.ChangeLogOptions{
			ChangeLogger: managed.NewGRPCChangeLogger(
				changelogsv1alpha1.NewChangeLogServiceClient(conn),
				managed.WithProviderVersion(fmt.Sprintf("provider-mongodbatlas:%s", version.Version))),
		}
		clusterOpts.ChangeLogOptions = &clo
		namespacedOpts.ChangeLogOptions = &clo
	}

	canSafeStart, err := canWatchCRD(context.TODO(), mgr)
	ctx.FatalIfErrorf(err, "SafeStart precheck failed")
	if canSafeStart {
		crdGate := new(gate.Gate[schema.GroupVersionKind])
		clusterOpts.Gate = crdGate
		namespacedOpts.Gate = crdGate
		ctx.FatalIfErrorf(customresourcesgate.Setup(mgr, xpcontroller.Options{
			Logger:                  log,
			Gate:                    crdGate,
			MaxConcurrentReconciles: 1,
		}), "Cannot setup CRD gate")
		ctx.FatalIfErrorf(controllerCluster.SetupGated(mgr, clusterOpts), "Cannot setup cluster-scoped MongoDBAtlas controllers")
		ctx.FatalIfErrorf(controllerNamespaced.SetupGated(mgr, namespacedOpts), "Cannot setup namespaced MongoDBAtlas controllers")
	} else {
		log.Info("Provider has missing RBAC permissions for watching CRDs, controller SafeStart capability will be disabled")
		ctx.FatalIfErrorf(controllerCluster.Setup(mgr, clusterOpts), "Cannot setup cluster-scoped MongoDBAtlas controllers")
		ctx.FatalIfErrorf(controllerNamespaced.Setup(mgr, namespacedOpts), "Cannot setup namespaced MongoDBAtlas controllers")
	}

	ctx.FatalIfErrorf(mgr.Start(ctrl.SetupSignalHandler()), "Cannot start controller manager")
}

func canWatchCRD(ctx context.Context, mgr manager.Manager) (bool, error) {
	if err := authv1.AddToScheme(mgr.GetScheme()); err != nil {
		return false, err
	}
	verbs := []string{"get", "list", "watch"}
	for _, verb := range verbs {
		sar := &authv1.SelfSubjectAccessReview{
			Spec: authv1.SelfSubjectAccessReviewSpec{
				ResourceAttributes: &authv1.ResourceAttributes{
					Group:    "apiextensions.k8s.io",
					Resource: "customresourcedefinitions",
					Verb:     verb,
				},
			},
		}
		if err := mgr.GetClient().Create(ctx, sar); err != nil {
			return false, errors.Wrapf(err, "unable to perform RBAC check for verb %s on CustomResourceDefinitions", verbs)
		}
		if !sar.Status.Allowed {
			return false, nil
		}
	}
	return true, nil
}
