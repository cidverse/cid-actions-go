package nodetest

import (
	"github.com/cidverse/cid-actions-go/actions/node/nodecommon"
	cidsdk "github.com/cidverse/cid-sdk-go"
)

type Action struct {
	Sdk cidsdk.SDKClient
}

type Config struct {
}

func (a Action) Execute() (err error) {
	cfg := Config{}
	ctx, err := a.Sdk.ModuleAction(&cfg)
	if err != nil {
		return err
	}

	// package.json
	content, err := a.Sdk.FileRead(cidsdk.JoinPath(ctx.Module.ModuleDir, "package.json"))
	if err != nil {
		return err
	}
	pkg, err := nodecommon.ParsePackageJSON(content)
	if err != nil {
		return err
	}

	// test, if script is present
	_, buildScriptPresent := pkg.Scripts[`test`]
	if buildScriptPresent {
		// install
		_, err = a.Sdk.ExecuteCommand(cidsdk.ExecuteCommandRequest{
			Command: `yarn install`,
			WorkDir: ctx.Module.ModuleDir,
		})
		if err != nil {
			return err
		}

		// test
		_, err := a.Sdk.ExecuteCommand(cidsdk.ExecuteCommandRequest{
			Command: `yarn test`,
			WorkDir: ctx.Module.ModuleDir,
		})
		if err != nil {
			return err
		}
	}

	return nil
}
