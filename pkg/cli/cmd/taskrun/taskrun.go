package taskrun

import (
	"github.com/spf13/cobra"
	"github.com/tektoncd/results/pkg/cli/flags"
)

// Command initializes a cobra command for `taskrun` sub commands
func Command(params *flags.Params) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "taskrun",
		Aliases: []string{"tr", "taskruns"},
		Short:   "Query TaskRuns",
		Annotations: map[string]string{
			"commandType": "main",
		},
	}

	cmd.AddCommand(
		listCommand(params),
		describeCommand(params),
		logsCommand(params),
	)

	return cmd
}
