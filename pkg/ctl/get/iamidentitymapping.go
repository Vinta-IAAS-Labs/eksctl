package get

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/pflag"

	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/authconfigmap"
	"github.com/weaveworks/eksctl/pkg/ctl/cmdutils"
	"github.com/weaveworks/eksctl/pkg/printers"
)

func getIAMIdentityMappingCmd(cmd *cmdutils.Cmd) {
	cfg := api.NewClusterConfig()
	cmd.ClusterConfig = cfg

	var arn authconfigmap.ARN

	params := &getCmdParams{}

	cmd.SetDescription("iamidentitymapping", "Get IAM identity mapping(s)", "")

	rc.SetRunFunc(func() error {
		return doGetIAMIdentityMapping(rc, params, arn)
	})

	rc.FlagSetGroup.InFlagSet("General", func(fs *pflag.FlagSet) {
		fs.Var(&arn, "arn", "ARN of the IAM role or user")
		cmdutils.AddNameFlag(fs, cfg.Metadata)
		cmdutils.AddRegionFlag(fs, cmd.ProviderConfig)
		cmdutils.AddCommonFlagsForGetCmd(fs, &params.chunkSize, &params.output)
		cmdutils.AddConfigFileFlag(fs, &cmd.ClusterConfigFile)
		cmdutils.AddTimeoutFlag(fs, &cmd.ProviderConfig.WaitTimeout)
	})

	cmdutils.AddCommonFlagsForAWS(cmd.FlagSetGroup, cmd.ProviderConfig, false)
}

<<<<<<< HEAD
func doGetIAMIdentityMapping(rc *cmdutils.Cmd, params *getCmdParams, arn string) error {
=======
func doGetIAMIdentityMapping(rc *cmdutils.ResourceCmd, params *getCmdParams, arn authconfigmap.ARN) error {
>>>>>>> Use dedicated ARN type instead of string
	if err := cmdutils.NewMetadataLoader(rc).Load(); err != nil {
		return err
	}

	cfg := cmd.ClusterConfig

	ctl, err := cmd.NewCtl()
	if err != nil {
		return err
	}

	if err := ctl.CheckAuth(); err != nil {
		return err
	}

	if cfg.Metadata.Name == "" {
		return cmdutils.ErrMustBeSet("--name")
	}

	if ok, err := ctl.CanOperate(cfg); !ok {
		return err
	}
	clientSet, err := ctl.NewStdClientSet(cfg)
	if err != nil {
		return err
	}
	acm, err := authconfigmap.NewFromClientSet(clientSet)
	if err != nil {
		return err
	}
	identities, err := acm.Identities()
	if err != nil {
		return err
	}

	if arn.Resource != "" {
		identities = identities.Get(arn)
		// If a filter was given, we error if none was found
		if len(identities) == 0 {
			return fmt.Errorf("no iamidentitymapping with arn %q found", arn)
		}
	}

	printer, err := printers.NewPrinter(params.output)
	if err != nil {
		return err
	}
	if params.output == "table" {
		addIAMIdentityMappingTableColumns(printer.(*printers.TablePrinter))
	}

	if err := printer.PrintObjWithKind("iamidentitymappings", identities, os.Stdout); err != nil {
		return err
	}

	return nil
}

func addIAMIdentityMappingTableColumns(printer *printers.TablePrinter) {
	printer.AddColumn("ARN", func(r authconfigmap.MapIdentity) string {
		return r.ARN.String()
	})
	printer.AddColumn("USERNAME", func(r authconfigmap.MapIdentity) string {
		return r.Username
	})
	printer.AddColumn("GROUPS", func(r authconfigmap.MapIdentity) string {
		return strings.Join(r.Groups, ",")
	})
}
