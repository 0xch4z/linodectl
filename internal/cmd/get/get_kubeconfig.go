package get

import (
	"bytes"
	"context"
	"encoding/base64"
	"fmt"
	"io"
	"os"

	"github.com/Charliekenney23/linodectl/internal/cli/genericoptions"
	cmdutil "github.com/Charliekenney23/linodectl/internal/cmd/util"
	"github.com/Charliekenney23/linodectl/internal/resource/lkecluster"
	"github.com/Charliekenney23/linodectl/internal/resource/resourceref"
	"github.com/linode/linodego"
	"github.com/spf13/cobra"
)

type GetKubeconfigOptions struct {
	refs resourceref.List

	outFile string

	genericoptions.ProfileFlags
	cmdutil.IOStreams
}

func NewGetKubeconfigOptions(ioStreams cmdutil.IOStreams) *GetKubeconfigOptions {
	return &GetKubeconfigOptions{
		IOStreams: ioStreams,
	}
}

func NewCmdGetKubeconfig(f cmdutil.Factory, ioStreams cmdutil.IOStreams) *cobra.Command {
	o := NewGetKubeconfigOptions(ioStreams)

	cmd := &cobra.Command{
		Use:     "kubeconfig [NAME] [args...]",
		Aliases: []string{"kconfig", "kcfg"},
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := o.Complete(f, ioStreams, args); err != nil {
				return err
			}
			return o.Run(f, cmd)
		},
	}

	cmd.Flags().StringVarP(&o.outFile, "out", "o", "", "File to output kubeconfig")
	o.ProfileFlags.AddFlags(cmd)
	return cmd
}

func (o *GetKubeconfigOptions) Complete(f cmdutil.Factory, ioStreams cmdutil.IOStreams, args []string) (err error) {
	if o.refs, err = resourceref.ListFromArgs(args); err != nil {
		return err
	}
	return nil
}

func (o *GetKubeconfigOptions) Run(f cmdutil.Factory, cmd *cobra.Command) error {
	if len(o.refs) == 1 {
		return fmt.Errorf("need a reference to exactly one LKE Cluster")
	}

	client, err := f.Client(o.ProfileName())
	if err != nil {
		return err
	}

	ctx := context.Background()

	clusterID := o.refs.ID()
	if label := o.refs.Label(); label != "" {
		clusters, err := client.ListLKEClusters(ctx, &linodego.ListOptions{})
		if err != nil {
			return err
		}

		clusters = lkecluster.FilterByRefs(clusters, o.refs)
		if len(clusters) != 1 {
			return fmt.Errorf("could not find LKE Cluster %q", label)
		}
		clusterID = clusters[0].ID
	}

	kubeconfig, err := client.GetLKEClusterKubeconfig(ctx, clusterID)
	if err != nil {
		return err
	}

	decodedKubeconfigBytes, err := base64.StdEncoding.DecodeString(kubeconfig.KubeConfig)
	if err != nil {
		return err
	}

	var dst io.Writer = o.IOStreams.Out
	if o.outFile != "" {
		dst, err = os.Create(o.outFile)
		if err != nil {
			return err
		}
	}

	_, err = io.Copy(dst, bytes.NewReader(decodedKubeconfigBytes))
	return err
}
