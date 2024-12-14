package taskrun

import (
	"github.com/spf13/cobra"
	"github.com/tektoncd/results/pkg/cli/flags"
)

type logsOptions struct {
	Namespace   string
	TaskRunName string
	StepName    string
}

func logsCommand(params *flags.Params) *cobra.Command {
	opts := &logsOptions{Namespace: "default"}

	eg := `Show the logs of TaskRun named 'foo' from the namespace 'bar':

    tkn-results taskrun logs foo -n bar

Show the live logs of TaskRun named 'foo' from namespace 'bar':

    tkn-results taskrun logs -f foo -n bar

Show the logs of TaskRun named 'microservice-1' for step 'build' only from namespace 'bar':

    tkn-results tr logs microservice-1 -s build -n bar
`
	cmd := &cobra.Command{
		Use:     "logs",
		Aliases: []string{"log"},
		Short:   "Show the logs of a TaskRun",
		Annotations: map[string]string{
			"commandType": "main",
		},
		Example: eg,
		RunE: func(cmd *cobra.Command, _ []string) error {
			// if opts.Limit < 0 {
			// 	return fmt.Errorf("limit was %d, but must be greater than 0", opts.Limit)
			// }

			// resp, err := params.ResultsClient.ListRecords(cmd.Context(), &pb.ListRecordsRequest{
			// 	Parent:   fmt.Sprintf("%s/results/-", opts.Namespace),
			// 	PageSize: int32(opts.Limit),
			// 	Filter:   `data_type==TASK_RUN`,
			// })
			// if err != nil {
			// 	return fmt.Errorf("failed to list TaskRuns from namespace %s: %v", opts.Namespace, err)
			// }
			// stream := &cli.Stream{
			// 	Out: cmd.OutOrStdout(),
			// 	Err: cmd.OutOrStderr(),
			// }
			// return printFormatted(stream, resp.Records, params.Clock)
			return nil
		},
	}
	cmd.Flags().StringVarP(&opts.Namespace, "namespace", "n", "default", "Namespace to list TaskRuns in")
	cmd.Flags().StringVarP(&opts.TaskRunName, "taskrun", "t", "", "Name of the TaskRun to show logs for")
	cmd.Flags().StringVarP(&opts.StepName, "step", "s", "", "Name of the step to show logs for")
	return cmd
}
