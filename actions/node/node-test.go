package node

import (
	"path"

	cidsdk "github.com/cidverse/cid-sdk-go"
)

type TestAction struct {
	Sdk cidsdk.SDKClient
}

type TestConfig struct {
}

func (a TestAction) Execute() (err error) {
	cfg := TestConfig{}
	ctx, err := a.Sdk.ModuleAction(&cfg)
	if err != nil {
		return err
	}

	// package.json
	content, err := a.Sdk.FileRead(path.Join(ctx.Module.ModuleDir, "package.json"))
	if err != nil {
		return err
	}
	pkg, err := ParsePackageJSON(content)
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
