package deploy

import (
	"github.com/lovego/xiaomei/xiaomei/images"
	"github.com/lovego/xiaomei/xiaomei/release"
	"github.com/spf13/cobra"
)

// Run, Deploy, Ps, Logs commands
func Cmds(svcName string) (cmds []*cobra.Command) {
	if svcName == `access` {
		cmds = append(cmds, accessPrintCmd(), accessSetupCmd())
	}
	if svcName != `` {
		cmds = append(cmds, runCmdFor(svcName))
	}
	cmds = append(cmds,
		deployCmdFor(svcName),
		rmDeployCmdFor(svcName),
		psCmdFor(svcName),
		logsCmdFor(svcName),
	)
	return
}

func runCmdFor(svcName string) *cobra.Command {
	// var publish []string
	cmd := &cobra.Command{
		Use:   `run`,
		Short: `run    the ` + deployDesc(svcName) + `.`,
		RunE: release.NoArgCall(func() error {
			return run(svcName)
		}),
	}
	// cmd.Flags().StringSliceVarP(&publish, `publish`, `p`, nil, `publish ports for container.`)
	return cmd
}

func deployCmdFor(svcName string) *cobra.Command {
	var noBuild, noPush /*, rmCurrent*/ bool
	cmd := &cobra.Command{
		Use:   `deploy`,
		Short: `deploy the ` + deployDesc(svcName) + `.`,
		RunE: release.NoArgCall(func() error {
			if !noBuild {
				if err := images.Build(svcName, true); err != nil {
					return err
				}
			}
			if !noPush {
				if err := images.Push(svcName); err != nil {
					return err
				}
			}
			return deploy(svcName) // , rmCurrent)
		}),
	}
	cmd.Flags().BoolVarP(&noBuild, `no-build`, `B`, false, `do not build the images.`)
	cmd.Flags().BoolVarP(&noPush, `no-push`, `P`, false, `do not push the images.`)
	// cmd.Flags().BoolVar(&rmCurrent, `rm-current`, false, `remove the current running `+desc+`.`)
	return cmd
}

func rmDeployCmdFor(svcName string) *cobra.Command {
	cmd := &cobra.Command{
		Use:   `rm-deploy`,
		Short: `remove deployment of the ` + deployDesc(svcName) + `.`,
		RunE: release.NoArgCall(func() error {
			return rmDeploy(svcName)
		}),
	}
	return cmd
}

func psCmdFor(svcName string) *cobra.Command {
	var watch bool
	cmd := &cobra.Command{
		Use:   `ps`,
		Short: `list tasks of the ` + deployDesc(svcName) + `.`,
		RunE: func(c *cobra.Command, args []string) error {
			return ps(svcName, watch, args)
		},
	}
	cmd.Flags().BoolVarP(&watch, `watch`, `w`, false, `watch ps.`)
	return cmd
}

func logsCmdFor(svcName string) *cobra.Command {
	cmd := &cobra.Command{
		Use:   `logs`,
		Short: `list logs  of the ` + deployDesc(svcName) + `.`,
		RunE: func(c *cobra.Command, args []string) error {
			return logs(svcName)
		},
	}
	return cmd
}

func deployDesc(svcName string) string {
	if svcName == `` {
		return `project`
	} else {
		return svcName + ` service`
	}
}

func accessPrintCmd() *cobra.Command {
	return &cobra.Command{
		Use:   `print`,
		Short: `print access config for the project.`,
		RunE:  release.NoArgCall(accessPrint),
	}
}

func accessSetupCmd() *cobra.Command {
	return &cobra.Command{
		Use:   `setup`,
		Short: `setup access config for the project.`,
		RunE:  release.NoArgCall(accessSetup),
	}
}