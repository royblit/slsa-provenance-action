package cli

import (
	"fmt"

	"github.com/docker/docker/client"
	"github.com/spf13/cobra"

	"github.com/royblit/slsa-provenance-action/cmd/slsa-provenance/cli/options"
	"github.com/royblit/slsa-provenance-action/lib/github"
	"github.com/royblit/slsa-provenance-action/lib/oci"
)

// OCI creates an instance of *cobra.Command to generate oci provenance
func OCI() *cobra.Command {
	o := &options.OCIOptions{}

	cmd := &cobra.Command{
		Use:   "container",
		Short: "Generate provenance on container assets",
		RunE: func(cmd *cobra.Command, args []string) error {
			outputPath, err := o.GetOutputPath()
			if err != nil {
				return err
			}

			gh, err := o.GetGitHubContext()
			if err != nil {
				return err
			}

			runner, err := o.GetRunnerContext()
			if err != nil {
				return err
			}

			materials, err := o.GetExtraMaterials()
			if err != nil {
				return err
			}

			repo, err := o.GetRepository()
			if err != nil {
				return err
			}

			digest, err := o.GetDigest()
			if err != nil {
				return err
			}

			tags, err := o.GetTags()
			if err != nil {
				return err
			}

			cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
			if err != nil {
				return err
			}
			subjecter := oci.NewContainerSubjecter(cli, repo, digest, tags...)

			env := &github.Environment{
				Context: gh,
				Runner:  runner,
			}
			stmt, err := env.GenerateProvenanceStatement(cmd.Context(), subjecter)
			if err != nil {
				return fmt.Errorf("failed to generate provenance: %w", err)
			}

			stmt.Predicate.Materials = append(stmt.Predicate.Materials, materials...)

			fmt.Fprintf(cmd.OutOrStdout(), "Saving provenance to %s\n", outputPath)

			return env.PersistProvenanceStatement(cmd.Context(), stmt, outputPath)
		},
	}

	o.AddFlags(cmd)

	return cmd
}
