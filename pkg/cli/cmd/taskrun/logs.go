package taskrun

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/tektoncd/results/pkg/cli/flags"
	pb "github.com/tektoncd/results/proto/v1alpha2/results_go_proto"
)

type logsOptions struct {
	Namespace string
}

func logsCommand(params *flags.Params) *cobra.Command {
	opts := &logsOptions{Namespace: "default"}

	eg := `Show the logs of TaskRun named 'foo' from the namespace 'bar':

    tkn-results taskrun logs foo -n bar

Show the live logs of TaskRun named 'foo' from namespace 'bar':

    tkn-results taskrun logs -f foo -n bar
`
	cmd := &cobra.Command{
		Use:     "logs",
		Aliases: []string{"log"},
		Short:   "Show the logs of a TaskRun",
		Annotations: map[string]string{
			"commandType": "main",
		},
		Example: eg,
		RunE: func(cmd *cobra.Command, args []string) error {
			taskrunName := args[0]

			resp, err := params.ResultsClient.ListRecords(cmd.Context(), &pb.ListRecordsRequest{
				Parent:   fmt.Sprintf("%s/results/-", opts.Namespace),
				PageSize: 5,
				Filter:   fmt.Sprintf(`data_type==TASK_RUN && data.metadata.name=="%s" && data.metadata.namespace=="%s"`, taskrunName, opts.Namespace),
				OrderBy:  "create_time",
			})
			if err != nil {
				return fmt.Errorf("failed to list TaskRuns from namespace %s of name %s: %v", opts.Namespace, taskrunName, err)
			}

			if len(resp.Records) == 0 {
				return fmt.Errorf("no TaskRun found with name %s in namespace %s", taskrunName, opts.Namespace)
			}

			record := resp.Records[0]

			logName := strings.ReplaceAll(record.GetName(), "records", "logs")

			stream, err := params.LogsClient.GetLog(cmd.Context(), &pb.GetLogRequest{
				Name: logName,
			})
			if err != nil {
				return fmt.Errorf("failed to get logs for TaskRun %s in namespace %s: %v", taskrunName, opts.Namespace, err)
			}
			data, err := stream.Recv()
			if err != nil {
				return fmt.Errorf("failed to receive steaming logs for TaskRun %s in namespace %s: %v", taskrunName, opts.Namespace, err)
			}

			if data.ContentType != "text/plain" {
				return fmt.Errorf("unsupported content type: %s", data.ContentType)
			}
			_, err = fmt.Fprint(cmd.OutOrStdout(), string(data.Data))
			return err
		},
	}
	cmd.Flags().StringVarP(&opts.Namespace, "namespace", "n", "default", "Namespace to list TaskRuns in")
	return cmd
}
