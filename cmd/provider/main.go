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
	"github.com/go-logr/logr"
	"github.com/mongodb/terraform-provider-mongodbatlas/xpshim"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	authv1 "k8s.io/api/authorization/v1"
	corev1 "k8s.io/api/core/v1"
	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/tools/leaderelection/resourcelock"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/cache"
	"sigs.k8s.io/controller-runtime/pkg/client"
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
	CertsDir                 certsDir      `default:"${defaultCertsDir}"                  env:"${defautCertsDirEnvVar}"                                                                              help:"The directory that contains the server key and certificate"`
	ChangelogsSocketPath     string        `default:"/var/run/changelogs/changelogs.sock" env:"CHANGELOGS_SOCKET_PATH"                                                                               help:"Path for changelogs socket (if enabled)"`
	Debug                    bool          `help:"Run with debug logging."                short:"d"`
	EnableChangeLogs         bool          `default:"false"                               env:"ENABLE_CHANGE_LOGS"                                                                                   help:"Enable support for capturing change logs during reconciliation."`
	EnableManagementPolicies bool          `default:"true"                                env:"ENABLE_MANAGEMENT_POLICIES"                                                                           help:"Enable support for Management Policies."`
	EnableSecretCache        bool          `default:"true"                                env:"ENABLE_SECRET_CACHE"                                                                                  help:"Enable caching of Secrets via an informer. Disabling it routes Secret reads through live API calls, trading memory for API server QPS."`
	LeaderElection           bool          `default:"false"                               env:"LEADER_ELECTION"                                                                                      help:"Use leader election for the controller manager."                                                                                        short:"l"`
	MaxReconcileRate         int           `default:"10"                                  help:"The global maximum rate per second at which resources may checked for drift from the desired state."`
	MetricsBindAddress       string        `default:":8081"                               env:"METRICS_BIND_ADDRESS"                                                                                 help:"The address the metrics server listens on"`
	PollInterval             time.Duration `default:"10m"                                 help:"How often individual resources will be checked for drift from the desired state"`
	PollStateMetricInterval  time.Duration `default:"5s"                                  help:"State metric recording interval"`
	PprofBindAddress         string        `default:""                                    env:"PPROF_BIND_ADDRESS"                                                                                   help:"The address the pprof profiling server listens on (e.g. :8083). Empty disables profiling."`
	SyncPeriod               time.Duration `default:"1h"                                  help:"Controller manager sync period such as 300ms, 1.5h, or 2h45m"                                        short:"s"`
	WebhookPort              int           `default:"9443"                                env:"WEBHOOK_PORT"                                                                                         help:"The port the webhook listens on"`
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

	if !certsDirSet {
		cli.CertsDir = certsDirFromEnv(cli.CertsDir)
	}

	// The scheme must be fully populated before ctrl.NewManager runs: the CRD
	// cache transform below resolves its ByObject key's GVK synchronously at
	// manager-construction time, so apiextensionsv1 (and everything else we
	// need) has to be registered up front rather than via mgr.GetScheme()
	// afterwards. We only register the specific corev1 type we need (Secret,
	// for connection secrets and credential extraction) rather than the whole
	// corev1 group - this scheme exists solely to satisfy our own controllers
	// and the CRD cache transform, not to stand in for client-go's full scheme.
	scheme := runtime.NewScheme()
	scheme.AddKnownTypes(corev1.SchemeGroupVersion, &corev1.Secret{}, &corev1.SecretList{})
	metav1.AddToGroupVersion(scheme, corev1.SchemeGroupVersion)
	ctx.FatalIfErrorf(apisCluster.AddToScheme(scheme), "Cannot add cluster-scoped MongoDBAtlas APIs to scheme")
	ctx.FatalIfErrorf(apisNamespaced.AddToScheme(scheme), "Cannot add namespaced MongoDBAtlas APIs to scheme")
	ctx.FatalIfErrorf(apiextensionsv1.AddToScheme(scheme), "Cannot add api-extensions APIs to scheme")
	ctx.FatalIfErrorf(authv1.AddToScheme(scheme), "Cannot add k8s authorization APIs to scheme")

	var clientOpts client.Options
	if !cli.EnableSecretCache {
		clientOpts = client.Options{
			Cache: &client.CacheOptions{DisableFor: []client.Object{&corev1.Secret{}}},
		}
	}

	mgr, err := ctrl.NewManager(cfg, ctrl.Options{
		Scheme:           scheme,
		LeaderElection:   cli.LeaderElection,
		LeaderElectionID: "crossplane-leader-election-upjet-provider-mongodbatlas",
		Cache: cache.Options{
			SyncPeriod: &cli.SyncPeriod,
			// Strips the OpenAPI schema (and other non-critical fields) from
			// cached CRDs: the SafeStart gate below only needs GVK +
			// Established status, never the schema, and skipping it cuts
			// provider memory footprint substantially on clusters with many
			// CRDs. Only ever read this cached CRD, never full-object
			// Update() it - that would persist the stripped object back to
			// the cluster.
			ByObject: map[client.Object]cache.ByObject{
				&apiextensionsv1.CustomResourceDefinition{}: {
					Transform: customresourcesgate.TransformStripCRDSchema,
				},
			},
		},
		Client:           clientOpts,
		PprofBindAddress: cli.PprofBindAddress,
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

	metricRecorder := managed.NewMRMetricRecorder()
	stateMetrics := statemetrics.NewMRStateMetrics()

	metrics.Registry.MustRegister(metricRecorder)
	metrics.Registry.MustRegister(stateMetrics)

	sdkProvider := xpshim.GetSDKProvider()
	fwProvider := xpshim.GetFrameworkProvider()
	setupFn := clients.TerraformSetupBuilder(sdkProvider, fwProvider)
	otStore := tjcontroller.NewOperationStore(log)

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
		Provider:              config.GetidentifierFromProvider(),
		SetupFn:               setupFn,
		OperationTrackerStore: otStore,
		PollJitter:            pollJitter,
		StartWebhooks:         cli.CertsDir != "",
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
		Provider:              config.GetProviderNamespaced(),
		SetupFn:               setupFn,
		OperationTrackerStore: otStore,
		PollJitter:            pollJitter,
		StartWebhooks:         cli.CertsDir != "",
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

	if cli.CertsDir != "" {
		ctx.FatalIfErrorf(controllerCluster.SetupWebhookWithManager(mgr), "Cannot setup cluster-scoped conversion webhooks")
		ctx.FatalIfErrorf(controllerNamespaced.SetupWebhookWithManager(mgr), "Cannot setup namespaced conversion webhooks")
	}

	ctx.FatalIfErrorf(mgr.Start(ctrl.SetupSignalHandler()), "Cannot start controller manager")
}

// certsDirFromEnv returns the TLS certs directory from the environment
// variables set by Crossplane if they're available.
// In older XP versions we used WEBHOOK_TLS_CERT_DIR, in newer versions
// we use TLS_SERVER_CERTS_DIR. If an explicit certs dir is not supplied
// via the command-line options, then these environment variables are used
// instead.
func certsDirFromEnv(fallback certsDir) certsDir {
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
		return certsDir(xpCertsDir)
	}
	return fallback
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
			return false, errors.Wrapf(err, "unable to perform RBAC check for verb %s on CustomResourceDefinitions", verb)
		}
		if !sar.Status.Allowed {
			return false, nil
		}
	}
	return true, nil
}
