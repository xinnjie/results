package fake

import (
	"context"
	"fmt"

	pb "github.com/tektoncd/results/proto/v1alpha2/results_go_proto"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
)

// ResultsClient is a fake implementation of the ResultsClient interface
type ResultsClient struct {
	// Map of result name to Result for GetResult and ListResults
	results map[string]*pb.Result
}

// NewResultsClient creates a new fake ResultsClient
func NewResultsClient(testData []*pb.Result) *ResultsClient {
	r := &ResultsClient{
		results: make(map[string]*pb.Result),
	}
	for _, result := range testData {
		r.results[result.Name] = result
	}
	return r
}

// AddResult adds a Result to the fake client's data store
func (c *ResultsClient) AddResult(name string, result *pb.Result) {
	c.results[name] = result
}

// GetResult implements ResultsClient.GetResult
func (c *ResultsClient) GetResult(ctx context.Context, in *pb.GetResultRequest, opts ...grpc.CallOption) (*pb.Result, error) {
	result, exists := c.results[in.Name]
	if !exists {
		return nil, fmt.Errorf("result not found: %s", in.Name)
	}
	return result, nil
}

// ListResults implements ResultsClient.ListResults
func (c *ResultsClient) ListResults(ctx context.Context, in *pb.ListResultsRequest, opts ...grpc.CallOption) (*pb.ListResultsResponse, error) {
	results := make([]*pb.Result, 0, len(c.results))
	for _, result := range c.results {
		results = append(results, result)
	}

	return &pb.ListResultsResponse{
		Results: results,
	}, nil
}

// CreateResult is unimplemented
func (c *ResultsClient) CreateResult(ctx context.Context, in *pb.CreateResultRequest, opts ...grpc.CallOption) (*pb.Result, error) {
	return nil, fmt.Errorf("unimplemented")
}

// UpdateResult is unimplemented
func (c *ResultsClient) UpdateResult(ctx context.Context, in *pb.UpdateResultRequest, opts ...grpc.CallOption) (*pb.Result, error) {
	return nil, fmt.Errorf("unimplemented")
}

// DeleteResult is unimplemented
func (c *ResultsClient) DeleteResult(ctx context.Context, in *pb.DeleteResultRequest, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	return nil, fmt.Errorf("unimplemented")
}

// CreateRecord is unimplemented
func (c *ResultsClient) CreateRecord(ctx context.Context, in *pb.CreateRecordRequest, opts ...grpc.CallOption) (*pb.Record, error) {
	return nil, fmt.Errorf("unimplemented")
}

// UpdateRecord is unimplemented
func (c *ResultsClient) UpdateRecord(ctx context.Context, in *pb.UpdateRecordRequest, opts ...grpc.CallOption) (*pb.Record, error) {
	return nil, fmt.Errorf("unimplemented")
}

// GetRecord is unimplemented
func (c *ResultsClient) GetRecord(ctx context.Context, in *pb.GetRecordRequest, opts ...grpc.CallOption) (*pb.Record, error) {
	return nil, fmt.Errorf("unimplemented")
}

// ListRecords is unimplemented
func (c *ResultsClient) ListRecords(ctx context.Context, in *pb.ListRecordsRequest, opts ...grpc.CallOption) (*pb.ListRecordsResponse, error) {
	return nil, fmt.Errorf("unimplemented")
}

// DeleteRecord is unimplemented
func (c *ResultsClient) DeleteRecord(ctx context.Context, in *pb.DeleteRecordRequest, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	return nil, fmt.Errorf("unimplemented")
}

// GetRecordListSummary is unimplemented
func (c *ResultsClient) GetRecordListSummary(ctx context.Context, in *pb.RecordListSummaryRequest, opts ...grpc.CallOption) (*pb.RecordListSummary, error) {
	return nil, fmt.Errorf("unimplemented")
}
