package dotnettest

import (
	"fmt"

	cidsdk "github.com/cidverse/cid-sdk-go"
)

type Action struct {
	Sdk cidsdk.SDKClient
}

type Config struct {
}

func (a Action) Execute() error {
	cfg := Config{}
	ctx, err := a.Sdk.ModuleAction(&cfg)
	if err != nil {
		return err
	}

	// restore
	restoreResult, err := a.Sdk.ExecuteCommand(cidsdk.ExecuteCommandRequest{
		Command: fmt.Sprintf(`dotnet restore --disable-parallel`),
		WorkDir: ctx.Module.ModuleDir,
	})
	if err != nil {
		return err
	} else if restoreResult.Code != 0 {
		return fmt.Errorf("dotnet restore failed, exit code %d", restoreResult.Code)
	}

	// test
	buildResult, err := a.Sdk.ExecuteCommand(cidsdk.ExecuteCommandRequest{
		Command: fmt.Sprintf(`dotnet test --configuration Release`),
		WorkDir: ctx.Module.ModuleDir,
	})
	if err != nil {
		return err
	} else if buildResult.Code != 0 {
		return fmt.Errorf("dotnet test failed, exit code %d", buildResult.Code)
	}

	return nil
}
