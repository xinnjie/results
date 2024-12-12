package pipelinerun

import (
	"testing"
	"time"

	"github.com/tektoncd/results/pkg/test"

	"github.com/tektoncd/results/pkg/cli/flags"
	"github.com/tektoncd/results/pkg/test/fake"

	"github.com/jonboulle/clockwork"

	"github.com/spf13/cobra"
	pb "github.com/tektoncd/results/proto/v1alpha2/results_go_proto"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func TestListPipelineRuns_empty(t *testing.T) {
	results := []*pb.Result{}
	now := time.Now()
	cmd := command(results, now)

	output, err := test.ExecuteCommand(cmd, "list")
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	test.AssertOutput(t, "No PipelineRuns found\n", output)
}

func TestListPipelineRuns(t *testing.T) {
	clock := clockwork.NewFakeClock()
	createTime := clock.Now().Add(time.Duration(-3) * time.Minute)
	updateTime := clock.Now().Add(time.Duration(-2) * time.Minute)
	startTime := clock.Now().Add(time.Duration(-3) * time.Minute)
	endTime := clock.Now().Add(time.Duration(-1) * time.Minute)
	results := []*pb.Result{
		{
			Name:       "default/results/e6b4b2e3-d876-4bbe-a927-95c691b6fdc7",
			Uid:        "949eebd9-1cf7-478f-a547-9ee313035f10",
			CreateTime: timestamppb.New(createTime),
			UpdateTime: timestamppb.New(updateTime),
			Annotations: map[string]string{
				"object.metadata.name": "hello-goodbye-run-vfsxn",
				"tekton.dev/pipeline":  "hello-goodbye",
			},
			Summary: &pb.RecordSummary{
				Record:    "default/results/e6b4b2e3-d876-4bbe-a927-95c691b6fdc7/records/e6b4b2e3-d876-4bbe-a927-95c691b6fdc7",
				Type:      "tekton.dev/v1.PipelineRun",
				StartTime: timestamppb.New(startTime),
				EndTime:   timestamppb.New(endTime),
				Status:    pb.RecordSummary_SUCCESS,
			},
		},
		{
			Name:       "default/results/3dacd30b-ce42-476c-be7e-84b0f664df55",
			Uid:        "c8d4cd50-06e8-4325-9ba2-044e6cc45235",
			CreateTime: timestamppb.New(createTime),
			UpdateTime: timestamppb.New(updateTime),
			Annotations: map[string]string{
				"object.metadata.name": "hello-goodbye-run-xtw2j",
			},
			Summary: &pb.RecordSummary{
				Record: "default/results/3dacd30b-ce42-476c-be7e-84b0f664df55/records/3dacd30b-ce42-476c-be7e-84b0f664df55",
				Type:   "tekton.dev/v1.PipelineRun",
			},
		},
	}
	cmd := command(results, clock.Now())
	output, err := test.ExecuteCommand(cmd, "list")
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	test.AssertOutput(t, `NAMESPACE   UID                                    STARTED         DURATION   STATUS
default     949eebd9-1cf7-478f-a547-9ee313035f10   3 minutes ago   2m0s       Succeeded
default     c8d4cd50-06e8-4325-9ba2-044e6cc45235   ---             ---        Unknown`, output)
}

func command(results []*pb.Result, now time.Time) *cobra.Command {
	clock := clockwork.NewFakeClockAt(now)

	param := &flags.Params{
		ResultsClient:    fake.NewResultsClient(results),
		LogsClient:       nil,
		PluginLogsClient: nil,
		Clock:            clock,
	}
	cmd := Command(param)
	return cmd
}
