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
	"maps"
	"os"
	"slices"
	"strings"
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
	"github.com/crossplane/upjet/v2/pkg/diffserver"
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
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
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
	resolverapis "github.com/crossplane-contrib/provider-mongodbatlas/internal/apis"
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

	// The commands kong parses, as Context.Command reports them.
	cmdStart      = "start"
	cmdDiffServer = "internal diff-server"
)

type certsDir string

func (d certsDir) BeforeApply(certsDirSet *bool) error {
	// we record whether the command-line option "--certs-dir" was supplied
	*certsDirSet = true
	return nil
}

// startCmd runs the provider's controllers. It is the default command, so
// that the provider keeps behaving as it did before the internal commands
// were introduced when it is invoked with flags alone.
type startCmd struct {
	CertsDir                 certsDir      `default:"${defaultCertsDir}"                  env:"${defautCertsDirEnvVar}"                                                                              help:"The directory that contains the server key and certificate"`
	ChangelogsSocketPath     string        `default:"/var/run/changelogs/changelogs.sock" env:"CHANGELOGS_SOCKET_PATH"                                                                               help:"Path for changelogs socket (if enabled)"`
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

// diffServerCmd serves upjet's diff gRPC service, which the Crossplane CLI
// drives to simulate changes. It never contacts the MongoDB Atlas API or a
// cluster: plans are computed offline from the resources in each request.
type diffServerCmd struct {
	Network string `default:"tcp"   enum:"tcp,unix"           env:"DIFF_SERVER_NETWORK"                                                                   help:"The network the diff gRPC server listens on."`
	Address string `default:":9099" env:"DIFF_SERVER_ADDRESS" help:"The address the diff gRPC server listens on. A socket path when the network is unix."`
}

var cli struct {
	Debug bool `help:"Run with debug logging." short:"d"`

	Start    startCmd `cmd:"" default:"withargs" help:"Start the provider's controllers."`
	Internal struct {
		DiffServer diffServerCmd `cmd:"" help:"Start a gRPC server serving the provider's diff services. It never contacts the MongoDB Atlas API or a cluster: plans are computed offline from the resources in each request." name:"diff-server"`
	} `cmd:"" help:"Commands used internally by the Crossplane ecosystem. No compatibility guarantees are made for them."`
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

	switch ctx.Command() {
	case cmdDiffServer:
		// Dispatched before a cluster connection (ctrl.GetConfig and the
		// manager below) becomes a requirement for running this binary at
		// all: the diff server serves from the resources in each request and
		// never talks to a cluster or to the Atlas API.
		ctx.FatalIfErrorf(serveDiffServer(ctrl.SetupSignalHandler(), log, cli.Internal.DiffServer.Network, cli.Internal.DiffServer.Address), "Cannot run the diff gRPC server")
		return
	case cmdStart:
		// The provider's controllers are started below.
	default:
		ctx.Fatalf("unknown command %q", ctx.Command())
	}

	// currently, we configure the jitter to be the 5% of the poll interval
	pollJitter := time.Duration(float64(cli.Start.PollInterval) * 0.05)

	log.Debug("Starting", "sync-period", cli.Start.SyncPeriod.String(), "poll-interval", cli.Start.PollInterval.String(), "poll-jitter", pollJitter, "max-reconcile-rate", cli.Start.MaxReconcileRate)

	cfg, err := ctrl.GetConfig()
	ctx.FatalIfErrorf(err, "Cannot get API server rest config")

	if !certsDirSet {
		cli.Start.CertsDir = certsDirFromEnv(cli.Start.CertsDir)
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
	ctx.FatalIfErrorf(resolverapis.BuildScheme(apisCluster.AddToSchemes), "Cannot register the cluster-scoped MongoDBAtlas APIs with the API resolver's runtime scheme")
	ctx.FatalIfErrorf(resolverapis.BuildScheme(apisNamespaced.AddToSchemes), "Cannot register the namespaced MongoDBAtlas APIs with the API resolver's runtime scheme")
	ctx.FatalIfErrorf(apiextensionsv1.AddToScheme(scheme), "Cannot add api-extensions APIs to scheme")
	ctx.FatalIfErrorf(authv1.AddToScheme(scheme), "Cannot add k8s authorization APIs to scheme")

	var clientOpts client.Options
	if !cli.Start.EnableSecretCache {
		clientOpts = client.Options{
			Cache: &client.CacheOptions{DisableFor: []client.Object{&corev1.Secret{}}},
		}
	}

	mgr, err := ctrl.NewManager(cfg, ctrl.Options{
		Scheme:           scheme,
		LeaderElection:   cli.Start.LeaderElection,
		LeaderElectionID: "crossplane-leader-election-upjet-provider-mongodbatlas",
		Cache: cache.Options{
			SyncPeriod:       &cli.Start.SyncPeriod,
			DefaultTransform: cache.TransformStripManagedFields(),
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
		PprofBindAddress: cli.Start.PprofBindAddress,
		Metrics: metricsserver.Options{
			BindAddress: cli.Start.MetricsBindAddress,
		},
		WebhookServer: webhook.NewServer(
			webhook.Options{
				CertDir: string(cli.Start.CertsDir),
				Port:    cli.Start.WebhookPort,
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
		Logger:                  log,
		GlobalRateLimiter:       ratelimiter.NewGlobal(cli.Start.MaxReconcileRate),
		PollInterval:            cli.Start.PollInterval,
		MaxConcurrentReconciles: cli.Start.MaxReconcileRate,
		Features:                &feature.Flags{},
		MetricOptions: &xpcontroller.MetricOptions{
			PollStateMetricInterval: cli.Start.PollStateMetricInterval,
			MRMetrics:               metricRecorder,
			MRStateMetrics:          stateMetrics,
		},
		Provider:              config.GetidentifierFromProvider(),
		SetupFn:               setupFn,
		OperationTrackerStore: otStore,
		PollJitter:            pollJitter,
		StartWebhooks:         cli.Start.CertsDir != "",
	}

	namespacedOpts := tjcontroller.Options{
		Logger:                  log,
		GlobalRateLimiter:       ratelimiter.NewGlobal(cli.Start.MaxReconcileRate),
		PollInterval:            cli.Start.PollInterval,
		MaxConcurrentReconciles: cli.Start.MaxReconcileRate,
		Features:                &feature.Flags{},
		MetricOptions: &xpcontroller.MetricOptions{
			PollStateMetricInterval: cli.Start.PollStateMetricInterval,
			MRMetrics:               metricRecorder,
			MRStateMetrics:          stateMetrics,
		},
		Provider:              config.GetProviderNamespaced(),
		SetupFn:               setupFn,
		OperationTrackerStore: otStore,
		PollJitter:            pollJitter,
		StartWebhooks:         cli.Start.CertsDir != "",
	}

	if cli.Start.EnableManagementPolicies {
		clusterOpts.Features.Enable(features.EnableBetaManagementPolicies)
		namespacedOpts.Features.Enable(features.EnableBetaManagementPolicies)
		log.Info("Beta feature enabled", "flag", features.EnableBetaManagementPolicies)
	}

	if cli.Start.EnableChangeLogs {
		clusterOpts.Features.Enable(feature.EnableAlphaChangeLogs)
		namespacedOpts.Features.Enable(feature.EnableAlphaChangeLogs)
		log.Info("Alpha feature enabled", "flag", feature.EnableAlphaChangeLogs)

		conn, err := grpc.NewClient("unix://"+cli.Start.ChangelogsSocketPath, grpc.WithTransportCredentials(insecure.NewCredentials()))
		ctx.FatalIfErrorf(err, "failed to create change logs client connection at %s", cli.Start.ChangelogsSocketPath)

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
			Groups:                  providerGroups(scheme, clusterOpts.Provider.RootGroup, namespacedOpts.Provider.RootGroup),
		}), "Cannot setup CRD gate")
		ctx.FatalIfErrorf(controllerCluster.SetupGated(mgr, clusterOpts), "Cannot setup cluster-scoped MongoDBAtlas controllers")
		ctx.FatalIfErrorf(controllerNamespaced.SetupGated(mgr, namespacedOpts), "Cannot setup namespaced MongoDBAtlas controllers")
	} else {
		log.Info("Provider has missing RBAC permissions for watching CRDs, controller SafeStart capability will be disabled")
		ctx.FatalIfErrorf(controllerCluster.Setup(mgr, clusterOpts), "Cannot setup cluster-scoped MongoDBAtlas controllers")
		ctx.FatalIfErrorf(controllerNamespaced.Setup(mgr, namespacedOpts), "Cannot setup namespaced MongoDBAtlas controllers")
	}

	if cli.Start.CertsDir != "" {
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

// providerGroups returns the API groups registered in s that are a root group
// or one of its sub-groups, so the CRD gate ignores unrelated CRDs.
func providerGroups(s *runtime.Scheme, roots ...string) []string {
	set := map[string]struct{}{}
	for gvk := range s.AllKnownTypes() {
		for _, root := range roots {
			if gvk.Group == root || strings.HasSuffix(gvk.Group, "."+root) {
				set[gvk.Group] = struct{}{}
			}
		}
	}
	return slices.Sorted(maps.Keys(set))
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

// serveDiffServer runs the diff gRPC server on the given network and address
// until the context is cancelled.
func serveDiffServer(ctx context.Context, log logging.Logger, network, address string) error {
	s, scheme, err := newDiffServer(log)
	if err != nil {
		return err
	}
	return s.Serve(ctx, network, address, scheme)
}

// newDiffServer returns the diff gRPC server, wired with both provider
// configurations to compute plans offline, and the runtime scheme it decodes
// plan requests with. The scheme holds the provider's own APIs and the
// Kubernetes core API, which is all a request can carry: the server never
// talks to a cluster.
func newDiffServer(log logging.Logger) (*diffserver.Server, *runtime.Scheme, error) {
	scheme := runtime.NewScheme()
	if err := clientgoscheme.AddToScheme(scheme); err != nil {
		return nil, nil, errors.Wrap(err, "cannot add Kubernetes core APIs to the diff server scheme")
	}
	if err := apisCluster.AddToScheme(scheme); err != nil {
		return nil, nil, errors.Wrap(err, "cannot add cluster-scoped MongoDBAtlas APIs to the diff server scheme")
	}
	if err := apisNamespaced.AddToScheme(scheme); err != nil {
		return nil, nil, errors.Wrap(err, "cannot add namespaced MongoDBAtlas APIs to the diff server scheme")
	}
	s := diffserver.NewServer(
		diffserver.WithProviderConfigurations(config.GetidentifierFromProvider(), config.GetProviderNamespaced()),
		diffserver.WithLogger(log),
		diffserver.WithTerraformSetupFn(clients.OfflineTerraformSetupBuilder(xpshim.GetFrameworkProvider())),
	)
	return s, scheme, nil
}
