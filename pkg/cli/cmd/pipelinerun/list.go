package pipelinerun

import (
	"fmt"
	"text/tabwriter"
	"text/template"

	"github.com/tektoncd/cli/pkg/cli"

	"github.com/jonboulle/clockwork"
	"github.com/spf13/cobra"
	"github.com/tektoncd/results/pkg/cli/flags"
	"github.com/tektoncd/results/pkg/cli/format"
	pb "github.com/tektoncd/results/proto/v1alpha2/results_go_proto"
)

const listTemplate = `{{- $size := len .Results -}}{{- if eq $size 0 -}}
No PipelineRuns found
{{ else -}}
NAMESPACE	UID	STARTED	DURATION	STATUS
{{- range $_, $result := .Results }}
{{ formatNamespace $result.Name }}	{{ $result.Uid }}	{{ formatAge $result.Summary.StartTime $.Time }}	{{ formatDuration $result.Summary.StartTime $result.Summary.EndTime }}	{{ formatStatus $result.Summary.Status }}
{{- end -}}
{{- end -}}`

type listOptions struct {
	Namespace string
	Limit     int
}

// listCommand initializes a cobra command to list PipelineRuns
func listCommand(params *flags.Params) *cobra.Command {
	opts := &listOptions{Limit: 0, Namespace: "default"}

	eg := `List all PipelineRuns in a namespace 'foo':
    tkn-results pipelinerun list -n foo

List all PipelineRuns in 'default' namespace:
    tkn-results pipelinerun list -n default
`
	cmd := &cobra.Command{
		Use:     "list",
		Aliases: []string{"ls"},
		Short:   "List PipelineRuns in a namespace",
		Annotations: map[string]string{
			"commandType": "main",
		},
		Example: eg,
		RunE: func(cmd *cobra.Command, _ []string) error {
			if opts.Limit < 0 {
				return fmt.Errorf("limit was %d, but must be greater than 0", opts.Limit)
			}

			resp, err := params.ResultsClient.ListResults(cmd.Context(), &pb.ListResultsRequest{
				Parent:   opts.Namespace,
				PageSize: int32(opts.Limit),
				Filter:   `summary.type==PIPELINE_RUN`,
			})
			if err != nil {
				return fmt.Errorf("failed to list PipelineRuns from namespace %s: %v", opts.Namespace, err)
			}
			stream := &cli.Stream{
				Out: cmd.OutOrStdout(),
				Err: cmd.OutOrStderr(),
			}
			return printFormatted(stream, resp.Results, params.Clock)
		},
	}
	cmd.Flags().StringVarP(&opts.Namespace, "namespace", "n", "default", "Namespace to list PipelineRuns in")
	cmd.Flags().IntVarP(&opts.Limit, "limit", "l", 0, "Limit the number of PipelineRuns to return")
	return cmd
}

func printFormatted(s *cli.Stream, results []*pb.Result, c clockwork.Clock) error {
	var data = struct {
		Results []*pb.Result
		Time    clockwork.Clock
	}{
		Results: results,
		Time:    c,
	}
	funcMap := template.FuncMap{
		"formatAge":       format.Age,
		"formatDuration":  format.Duration,
		"formatStatus":    format.Status,
		"formatNamespace": format.Namespace,
	}

	w := tabwriter.NewWriter(s.Out, 0, 5, 3, ' ', tabwriter.TabIndent)
	t := template.Must(template.New("List PipelineRuns").Funcs(funcMap).Parse(listTemplate))

	err := t.Execute(w, data)
	if err != nil {
		return err
	}
	return w.Flush()
}
