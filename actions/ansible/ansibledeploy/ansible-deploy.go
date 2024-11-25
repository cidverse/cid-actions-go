package ansibledeploy

import (
	"fmt"
	"path"

	cidsdk "github.com/cidverse/cid-sdk-go"
)

type Action struct {
	Sdk cidsdk.SDKClient
}

type Config struct {
	PlaybookFile  string `json:"ansible_playbook"  env:"ANSIBLE_PLAYBOOK"`
	InventoryFile string `json:"ansible_inventory"  env:"ANSIBLE_INVENTORY"`
}

func (a Action) Execute() (err error) {
	cfg := Config{}
	ctx, err := a.Sdk.ModuleAction(&cfg)
	if err != nil {
		return err
	}

	// config
	playbookFile := cfg.PlaybookFile
	inventoryFile := cfg.InventoryFile
	if playbookFile == "" {
		playbookFile = ctx.Module.Discovery[0].File
	}
	if inventoryFile == "" {
		inventoryFile = path.Join(path.Dir(playbookFile), "inventory")
	}

	// role and collection requirements
	if a.Sdk.FileExists(path.Join(ctx.Module.ModuleDir, "requirements.yml")) {
		_, err = a.Sdk.ExecuteCommand(cidsdk.ExecuteCommandRequest{
			Command: `ansible-galaxy collection install -r requirements.yml`,
			WorkDir: ctx.Module.ModuleDir,
		})
		if err != nil {
			return err
		}
	}

	// deploy
	_, err = a.Sdk.ExecuteCommand(cidsdk.ExecuteCommandRequest{
		Command: fmt.Sprintf(`ansible-playbook %q -i %q`, playbookFile, inventoryFile),
		WorkDir: ctx.Module.ModuleDir,
	})
	if err != nil {
		return err
	}

	return nil
}
