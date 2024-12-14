package taskrun

import (
	"context"
	"fmt"
	"text/tabwriter"
	"text/template"

	"github.com/tektoncd/cli/pkg/cli"
	pipelinev1 "github.com/tektoncd/pipeline/pkg/apis/pipeline/v1"
	pipelinev1beta1 "github.com/tektoncd/pipeline/pkg/apis/pipeline/v1beta1"
	"k8s.io/apimachinery/pkg/util/json"

	"github.com/jonboulle/clockwork"
	"github.com/spf13/cobra"
	"github.com/tektoncd/cli/pkg/formatted"
	"github.com/tektoncd/results/pkg/cli/flags"

	pb "github.com/tektoncd/results/proto/v1alpha2/results_go_proto"
)

const listTemplate = `{{- $size := len .TaskRuns -}}{{- if eq $size 0 -}}
No TaskRuns found
{{ else -}}
NAMESPACE	UID	STARTED	DURATION	STATUS
{{- range $_, $tr := .TaskRuns }}
{{ $tr.ObjectMeta.Namespace }}	{{ $tr.ObjectMeta.Name }}	{{ formatAge $tr.Status.StartTime $.Time }}	{{ formatDuration $tr.Status.StartTime $tr.Status.CompletionTime }}	{{ formatCondition $tr.Status.Conditions }}
{{- end -}}
{{- end -}}`

type listOptions struct {
	Namespace string
	Limit     int
}

// listCommand initializes a cobra command to list PipelineRuns
func listCommand(params *flags.Params) *cobra.Command {
	opts := &listOptions{Limit: 0, Namespace: "default"}

	eg := `List all TaskRuns in a namespace 'foo':
    tkn-results taskrun list -n foo

List all TaskRuns in 'default' namespace:
    tkn-results taskrun list -n default
`
	cmd := &cobra.Command{
		Use:     "list",
		Aliases: []string{"ls"},
		Short:   "List TaskRuns in a namespace",
		Annotations: map[string]string{
			"commandType": "main",
		},
		Example: eg,
		RunE: func(cmd *cobra.Command, _ []string) error {
			if opts.Limit < 0 {
				return fmt.Errorf("limit was %d, but must be greater than 0", opts.Limit)
			}

			resp, err := params.ResultsClient.ListRecords(cmd.Context(), &pb.ListRecordsRequest{
				Parent:   fmt.Sprintf("%s/results/-", opts.Namespace),
				PageSize: int32(opts.Limit),
				Filter:   `data_type==TASK_RUN`,
			})
			if err != nil {
				return fmt.Errorf("failed to list TaskRuns from namespace %s: %v", opts.Namespace, err)
			}
			stream := &cli.Stream{
				Out: cmd.OutOrStdout(),
				Err: cmd.OutOrStderr(),
			}
			return printFormatted(stream, resp.Records, params.Clock)
		},
	}
	cmd.Flags().StringVarP(&opts.Namespace, "namespace", "n", "default", "Namespace to list TaskRuns in")
	cmd.Flags().IntVarP(&opts.Limit, "limit", "l", 0, "Limit the number of TaskRuns to return")
	return cmd
}

func taskRunFromRecord(record *pb.Record) (*pipelinev1.TaskRun, error) {
	tr := &pipelinev1.TaskRun{}
	if record.Data.GetType() == "tekton.dev/v1beta1.TaskRun" {
		trV1beta1 := &pipelinev1beta1.TaskRun{}
		if err := json.Unmarshal(record.Data.Value, trV1beta1); err != nil {
			return nil, fmt.Errorf("failed to unmarshal TaskRun data: %v", err)
		}
		if err := tr.ConvertFrom(context.TODO(), trV1beta1); err != nil {
			return nil, fmt.Errorf("failed to convert v1beta1 TaskRun to v1: %v", err)
		}
	} else {
		if err := json.Unmarshal(record.Data.Value, tr); err != nil {
			return nil, fmt.Errorf("failed to unmarshal TaskRun data: %v", err)
		}
	}
	return tr, nil
}

func printFormatted(s *cli.Stream, records []*pb.Record, c clockwork.Clock) error {
	var data = struct {
		TaskRuns []*pipelinev1.TaskRun
		Time     clockwork.Clock
	}{
		TaskRuns: []*pipelinev1.TaskRun{},
		Time:     c,
	}

	for _, record := range records {
		if tr, err := taskRunFromRecord(record); err != nil {
			continue
		} else {
			data.TaskRuns = append(data.TaskRuns, tr)
		}
	}

	funcMap := template.FuncMap{
		"formatAge":       formatted.Age,
		"formatDuration":  formatted.Duration,
		"formatCondition": formatted.Condition,
	}

	w := tabwriter.NewWriter(s.Out, 0, 5, 3, ' ', tabwriter.TabIndent)
	t := template.Must(template.New("List TaskRuns").Funcs(funcMap).Parse(listTemplate))

	err := t.Execute(w, data)
	if err != nil {
		return err
	}
	return w.Flush()
}
