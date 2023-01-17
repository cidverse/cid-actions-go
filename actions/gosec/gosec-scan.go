package gosec

import (
	"os"
	"path"
	"strings"

	cidsdk "github.com/cidverse/cid-sdk-go"
)

type ScanAction struct {
	Sdk cidsdk.SDKClient
}

type ScanConfig struct {
}

func (a ScanAction) Execute() (err error) {
	cfg := ScanConfig{}
	ctx, err := a.Sdk.ModuleAction(&cfg)
	if err != nil {
		return err
	}
	sarifDir := path.Join(ctx.Config.ArtifactDir, ctx.Module.ModuleDir, "sarif")
	_ = os.MkdirAll(sarifDir, os.ModePerm)

	var opts []string
	opts = append(opts, "-no-fail")
	opts = append(opts, "-fmt sarif")
	opts = append(opts, "-out "+path.Join(sarifDir, "gosec.sarif"))
	opts = append(opts, "./...")
	_, err = a.Sdk.ExecuteCommand(cidsdk.ExecuteCommandRequest{
		Command: strings.TrimRight(`gosec `+strings.Join(opts, " "), " "),
		WorkDir: ctx.ProjectDir,
	})
	if err != nil {
		return err
	}

	return nil
}
