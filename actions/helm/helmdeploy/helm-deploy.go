package helmdeploy

import (
	"fmt"
	"os"
	"path"

	"github.com/cidverse/cid-actions-go/actions/helm/helmcommon"
	"github.com/cidverse/cid-actions-go/util"
	cidsdk "github.com/cidverse/cid-sdk-go"
	"github.com/cidverse/cidverseutils/filesystem"
	cp "github.com/otiai10/copy"
)

type Action struct {
	Sdk cidsdk.SDKClient
}

type Config struct {
	DeploymentChart          string `json:"deployment_chart"        env:"DEPLOYMENT_CHART"`
	DeploymentChartVersion   string `json:"deployment_chart_version" env:"DEPLOYMENT_CHART_VERSION"`
	DeploymentChartLocalPath string `json:"deployment_chart_local_path"        env:"DEPLOYMENT_CHART_LOCAL_PATH"`
	DeploymentNamespace      string `json:"deployment_namespace" env:"DEPLOYMENT_NAMESPACE"`
	DeploymentID             string `json:"deployment_id"            env:"DEPLOYMENT_ID"`
	DeploymentEnvironment    string `json:"deployment_environment" env:"DEPLOYMENT_ENVIRONMENT"`
	HelmArgs                 string `json:"helm_args" env:"HELM_ARGS"`
}

func (a Action) Metadata() cidsdk.ActionMetadata {
	return cidsdk.ActionMetadata{
		Name:        "helm-deploy",
		Description: "The Helm Deploy action is used to deploy a Helm chart to a Kubernetes cluster.",
		Documentation: util.TrimLeftEachLine(`
			# Helm Deploy

			The Helm Deploy action is used to deploy a Helm chart to a Kubernetes cluster.
			...
		`),
		Category: "deploy",
		Scope:    cidsdk.ActionScopeModule,
		Rules: []cidsdk.ActionRule{
			{
				Type:       "cel",
				Expression: `MODULE_DEPLOYMENT_TYPE == "helm"`,
			},
		},
		Access: cidsdk.ActionAccess{
			Environment: []cidsdk.ActionAccessEnv{
				{
					Name:        "DEPLOYMENT_CHART",
					Description: "The Helm chart to deploy",
				},
				{
					Name:        "DEPLOYMENT_CHART_VERSION",
					Description: "The Helm chart version to deploy (deprecated)",
				},
				{
					Name:        "DEPLOYMENT_CHART_LOCAL_PATH",
					Description: "The path to a helm chart in the local filesystem. (cannot be used together with DEPLOYMENT_CHART/DEPLOYMENT_CHART_VERSION)",
					// Deprecated in favor of DEPLOYMENT_CHART,
				},
				{
					Name:        "DEPLOYMENT_NAMESPACE",
					Description: "The namespace the deployment should be created in",
				},
				{
					Name:        "DEPLOYMENT_ID",
					Description: "The unique identifier of the deployment",
				},
				{
					Name:        "DEPLOYMENT_ENVIRONMENT",
					Description: "The environment the deployment is targeting",
				},
				{
					Name:        "HELM_ARGS",
					Description: "Additional arguments to pass to the helm command",
				},
				{
					Name:        "KUBECONFIG_.*_BASE64",
					Description: "The base 64 encoded Kubernetes config file to use for installation",
					Pattern:     true,
				},
			},
			Executables: []cidsdk.ActionAccessExecutable{
				{
					Name:       "helm",
					Constraint: helmcommon.HelmVersionConstraint,
				},
			},
		},
	}
}

func (a Action) Execute() (err error) {
	// query action data
	d, err := a.Sdk.ModuleActionDataV1()
	if err != nil {
		return err
	}

	// parse config
	_ = a.Sdk.Log(cidsdk.LogMessageRequest{Level: "info", Message: "Starting Helm deployment..."})
	cfg := Config{}
	cidsdk.PopulateFromEnv(&cfg, d.Env)

	// verify required properties
	if cfg.DeploymentChart == "" && cfg.DeploymentChartLocalPath != "" {
		cfg.DeploymentChart = cfg.DeploymentChartLocalPath
	}
	if cfg.DeploymentChart == "" {
		return fmt.Errorf("missing required input: DEPLOYMENT_CHART")
	}
	if cfg.DeploymentNamespace == "" {
		// default to project id
		cfg.DeploymentNamespace = d.Env["NCI_PROJECT_ID"]
	}
	if cfg.DeploymentID == "" {
		return fmt.Errorf("missing required input: DEPLOYMENT_ID")
	}

	// prepare kubeconfig
	kubeConfigFile := cidsdk.JoinPath(d.Config.TempDir, "kube", "kubeconfig")
	_ = a.Sdk.Log(cidsdk.LogMessageRequest{Level: "info", Message: "Using kubeconfig...", Context: map[string]interface{}{"file": kubeConfigFile}})
	err = helmcommon.PrepareKubeConfig(kubeConfigFile, cfg.DeploymentEnvironment, d.Env)
	if err != nil {
		return err
	}

	// target cluster
	targetCluster, err := helmcommon.ParseKubeConfigCluster(kubeConfigFile)
	if err != nil {
		return err
	}
	_ = a.Sdk.Log(cidsdk.LogMessageRequest{Level: "info", Message: "Target cluster", Context: map[string]interface{}{"name": targetCluster.Name, "api": targetCluster.Cluster.Server}})

	// query chart information
	_ = a.Sdk.Log(cidsdk.LogMessageRequest{Level: "info", Message: "Querying Helm chart information", Context: map[string]interface{}{"chart": cfg.DeploymentChart, "chart-version": cfg.DeploymentChartVersion}})
	cmdResult, err := a.Sdk.ExecuteCommand(cidsdk.ExecuteCommandRequest{
		Command:       fmt.Sprintf(`helm show chart --version %q %q`, cfg.DeploymentChartVersion, cfg.DeploymentChart),
		WorkDir:       d.Module.ModuleDir,
		CaptureOutput: true,
	})
	if err != nil {
		return err
	} else if cmdResult.Code != 0 {
		return fmt.Errorf("command failed, exit code %d", cmdResult.Code)
	}
	chartMetadata := cmdResult.Stdout
	chart, err := helmcommon.ParseChart([]byte(chartMetadata))
	if err != nil {
		return err
	}
	_ = a.Sdk.Log(cidsdk.LogMessageRequest{Level: "info", Message: "Found Helm chart", Context: map[string]interface{}{"chart-version": chart.Version, "app-version": chart.AppVersion}})

	// properties
	chartsDir := cidsdk.JoinPath(d.Config.TempDir, "helm-charts")
	chartDir := path.Join(chartsDir, chart.Name)
	_ = os.MkdirAll(chartsDir, 0755)
	chartSource := helmcommon.GetChartSource(cfg.DeploymentChart)

	// local dir branch, copy dir and maybe pull requirements, if missing
	if chartSource == helmcommon.ChartSourceOCI || chartSource == helmcommon.ChartSourceRepository {
		// download chart
		_ = a.Sdk.Log(cidsdk.LogMessageRequest{Level: "info", Message: "Downloading Helm chart", Context: map[string]interface{}{"chart": cfg.DeploymentChart, "chart-version": cfg.DeploymentChartVersion}})
		cmdResult, err = a.Sdk.ExecuteCommand(cidsdk.ExecuteCommandRequest{
			Command: fmt.Sprintf(`helm pull --untar --destination %q --version %q %q`, chartsDir, cfg.DeploymentChartVersion, cfg.DeploymentChart),
			WorkDir: d.Module.ModuleDir,
		})
		if err != nil {
			return err
		} else if cmdResult.Code != 0 {
			return fmt.Errorf("command failed, exit code %d", cmdResult.Code)
		}
	} else if chartSource == helmcommon.ChartSourceLocal {
		// copy chart
		chartSourceDir := filesystem.ResolveAbsolutePath(cfg.DeploymentChart)
		_ = a.Sdk.Log(cidsdk.LogMessageRequest{Level: "info", Message: "Copying Helm chart", Context: map[string]interface{}{"chart-dir": chartSourceDir}})
		if chartSourceDir == "" {
			return fmt.Errorf("chart not found: %s", cfg.DeploymentChart)
		}

		err = cp.Copy(chartSourceDir, chartDir)
		if err != nil {
			return err
		}
	} else {
		return fmt.Errorf("unsupported chart source: %s", cfg.DeploymentChart)
	}

	// deploy
	_ = a.Sdk.Log(cidsdk.LogMessageRequest{Level: "info", Message: "installing helm chart onto cluster", Context: map[string]interface{}{"chart": cfg.DeploymentChart, "chart-version": cfg.DeploymentChartVersion, "namespace": cfg.DeploymentNamespace, "release": cfg.DeploymentID}})
	cmdResult, err = a.Sdk.ExecuteCommand(cidsdk.ExecuteCommandRequest{
		Command: fmt.Sprintf(`helm upgrade --namespace %q --install --disable-openapi-validation %s %q %q`, cfg.DeploymentNamespace, cfg.HelmArgs, cfg.DeploymentID, chartDir),
		WorkDir: d.Module.ModuleDir,
		Env: map[string]string{
			"KUBECONFIG": kubeConfigFile,
		},
	})
	if err != nil {
		return err
	} else if cmdResult.Code != 0 {
		return fmt.Errorf("command failed, exit code %d", cmdResult.Code)
	}

	return nil
}
