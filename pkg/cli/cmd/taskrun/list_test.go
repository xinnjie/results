package taskrun

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

func TestListTaskRuns_empty(t *testing.T) {
	results := []*pb.Result{}
	now := time.Now()
	cmd := command(results, now)

	output, err := test.ExecuteCommand(cmd, "list")
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	test.AssertOutput(t, "No TaskRuns found\n", output)
}

func TestListTaskRuns(t *testing.T) {
	clock := clockwork.NewFakeClock()
	createTime := clock.Now().Add(time.Duration(-3) * time.Minute)
	updateTime := clock.Now().Add(time.Duration(-2) * time.Minute)
	startTime := clock.Now().Add(time.Duration(-3) * time.Minute)
	endTime := clock.Now().Add(time.Duration(-1) * time.Minute)
	results := []*pb.Result{
		{
			Name:       "default/results/e6f35206-1778-4178-9b77-20bd26b8b789",
			Uid:        "16bed717-b8fc-44c8-b68c-ce7ffb2dde0e",
			CreateTime: timestamppb.New(createTime),
			UpdateTime: timestamppb.New(updateTime),
			Annotations: map[string]string{
				"object.metadata.name": "hello-task-run-bsft7",
			},
			Summary: &pb.RecordSummary{
				Record:    "default/results/e6f35206-1778-4178-9b77-20bd26b8b789/records/e6f35206-1778-4178-9b77-20bd26b8b789",
				Type:      "tekton.dev/v1.TaskRun",
				StartTime: timestamppb.New(startTime),
				EndTime:   timestamppb.New(endTime),
				Status:    pb.RecordSummary_SUCCESS,
			},
		},
		{
			Name:       "default/results/fc5de8f1-7071-4093-ba7f-4ad7f4ee993f",
			Uid:        "1c6d1d39-a5a5-4c21-a2cb-5b04befe0e77",
			CreateTime: timestamppb.New(createTime),
			UpdateTime: timestamppb.New(updateTime),
			Annotations: map[string]string{
				"object.metadata.name": "hello-task-run-9jfs5",
			},
			Summary: &pb.RecordSummary{
				Record:    "default/results/fc5de8f1-7071-4093-ba7f-4ad7f4ee993f/records/fc5de8f1-7071-4093-ba7f-4ad7f4ee993f",
				Type:      "tekton.dev/v1.TaskRun",
				StartTime: timestamppb.New(startTime),
				EndTime:   timestamppb.New(endTime),
				Status:    pb.RecordSummary_SUCCESS,
			},
		},
	}
	cmd := command(results, clock.Now())
	output, err := test.ExecuteCommand(cmd, "list")
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	test.AssertOutput(t, `NAMESPACE   UID                                    STARTED         DURATION   STATUS
default     16bed717-b8fc-44c8-b68c-ce7ffb2dde0e   3 minutes ago   2m0s       Succeeded
default     1c6d1d39-a5a5-4c21-a2cb-5b04befe0e77   3 minutes ago   2m0s       Succeeded`, output)
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
