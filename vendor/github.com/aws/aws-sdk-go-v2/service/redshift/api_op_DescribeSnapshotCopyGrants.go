// Code generated by smithy-go-codegen DO NOT EDIT.

package redshift

import (
	"context"
	"fmt"
	awsmiddleware "github.com/aws/aws-sdk-go-v2/aws/middleware"
	"github.com/aws/aws-sdk-go-v2/aws/signer/v4"
	"github.com/aws/aws-sdk-go-v2/service/redshift/types"
	"github.com/aws/smithy-go/middleware"
	smithyhttp "github.com/aws/smithy-go/transport/http"
)

// Returns a list of snapshot copy grants owned by the Amazon Web Services account
// in the destination region. For more information about managing snapshot copy
// grants, go to Amazon Redshift Database Encryption
// (https://docs.aws.amazon.com/redshift/latest/mgmt/working-with-db-encryption.html)
// in the Amazon Redshift Cluster Management Guide.
func (c *Client) DescribeSnapshotCopyGrants(ctx context.Context, params *DescribeSnapshotCopyGrantsInput, optFns ...func(*Options)) (*DescribeSnapshotCopyGrantsOutput, error) {
	if params == nil {
		params = &DescribeSnapshotCopyGrantsInput{}
	}

	result, metadata, err := c.invokeOperation(ctx, "DescribeSnapshotCopyGrants", params, optFns, c.addOperationDescribeSnapshotCopyGrantsMiddlewares)
	if err != nil {
		return nil, err
	}

	out := result.(*DescribeSnapshotCopyGrantsOutput)
	out.ResultMetadata = metadata
	return out, nil
}

// The result of the DescribeSnapshotCopyGrants action.
type DescribeSnapshotCopyGrantsInput struct {

	// An optional parameter that specifies the starting point to return a set of
	// response records. When the results of a DescribeSnapshotCopyGrant request exceed
	// the value specified in MaxRecords, Amazon Web Services returns a value in the
	// Marker field of the response. You can retrieve the next set of response records
	// by providing the returned marker value in the Marker parameter and retrying the
	// request. Constraints: You can specify either the SnapshotCopyGrantName parameter
	// or the Marker parameter, but not both.
	Marker *string

	// The maximum number of response records to return in each call. If the number of
	// remaining response records exceeds the specified MaxRecords value, a value is
	// returned in a marker field of the response. You can retrieve the next set of
	// records by retrying the command with the returned marker value. Default: 100
	// Constraints: minimum 20, maximum 100.
	MaxRecords *int32

	// The name of the snapshot copy grant.
	SnapshotCopyGrantName *string

	// A tag key or keys for which you want to return all matching resources that are
	// associated with the specified key or keys. For example, suppose that you have
	// resources tagged with keys called owner and environment. If you specify both of
	// these tag keys in the request, Amazon Redshift returns a response with all
	// resources that have either or both of these tag keys associated with them.
	TagKeys []string

	// A tag value or values for which you want to return all matching resources that
	// are associated with the specified value or values. For example, suppose that you
	// have resources tagged with values called admin and test. If you specify both of
	// these tag values in the request, Amazon Redshift returns a response with all
	// resources that have either or both of these tag values associated with them.
	TagValues []string

	noSmithyDocumentSerde
}

type DescribeSnapshotCopyGrantsOutput struct {

	// An optional parameter that specifies the starting point to return a set of
	// response records. When the results of a DescribeSnapshotCopyGrant request exceed
	// the value specified in MaxRecords, Amazon Web Services returns a value in the
	// Marker field of the response. You can retrieve the next set of response records
	// by providing the returned marker value in the Marker parameter and retrying the
	// request. Constraints: You can specify either the SnapshotCopyGrantName parameter
	// or the Marker parameter, but not both.
	Marker *string

	// The list of SnapshotCopyGrant objects.
	SnapshotCopyGrants []types.SnapshotCopyGrant

	// Metadata pertaining to the operation's result.
	ResultMetadata middleware.Metadata

	noSmithyDocumentSerde
}

func (c *Client) addOperationDescribeSnapshotCopyGrantsMiddlewares(stack *middleware.Stack, options Options) (err error) {
	err = stack.Serialize.Add(&awsAwsquery_serializeOpDescribeSnapshotCopyGrants{}, middleware.After)
	if err != nil {
		return err
	}
	err = stack.Deserialize.Add(&awsAwsquery_deserializeOpDescribeSnapshotCopyGrants{}, middleware.After)
	if err != nil {
		return err
	}
	if err = addSetLoggerMiddleware(stack, options); err != nil {
		return err
	}
	if err = awsmiddleware.AddClientRequestIDMiddleware(stack); err != nil {
		return err
	}
	if err = smithyhttp.AddComputeContentLengthMiddleware(stack); err != nil {
		return err
	}
	if err = addResolveEndpointMiddleware(stack, options); err != nil {
		return err
	}
	if err = v4.AddComputePayloadSHA256Middleware(stack); err != nil {
		return err
	}
	if err = addRetryMiddlewares(stack, options); err != nil {
		return err
	}
	if err = addHTTPSignerV4Middleware(stack, options); err != nil {
		return err
	}
	if err = awsmiddleware.AddRawResponseToMetadata(stack); err != nil {
		return err
	}
	if err = awsmiddleware.AddRecordResponseTiming(stack); err != nil {
		return err
	}
	if err = addClientUserAgent(stack); err != nil {
		return err
	}
	if err = smithyhttp.AddErrorCloseResponseBodyMiddleware(stack); err != nil {
		return err
	}
	if err = smithyhttp.AddCloseResponseBodyMiddleware(stack); err != nil {
		return err
	}
	if err = stack.Initialize.Add(newServiceMetadataMiddleware_opDescribeSnapshotCopyGrants(options.Region), middleware.Before); err != nil {
		return err
	}
	if err = addRequestIDRetrieverMiddleware(stack); err != nil {
		return err
	}
	if err = addResponseErrorMiddleware(stack); err != nil {
		return err
	}
	if err = addRequestResponseLogging(stack, options); err != nil {
		return err
	}
	return nil
}

// DescribeSnapshotCopyGrantsAPIClient is a client that implements the
// DescribeSnapshotCopyGrants operation.
type DescribeSnapshotCopyGrantsAPIClient interface {
	DescribeSnapshotCopyGrants(context.Context, *DescribeSnapshotCopyGrantsInput, ...func(*Options)) (*DescribeSnapshotCopyGrantsOutput, error)
}

var _ DescribeSnapshotCopyGrantsAPIClient = (*Client)(nil)

// DescribeSnapshotCopyGrantsPaginatorOptions is the paginator options for
// DescribeSnapshotCopyGrants
type DescribeSnapshotCopyGrantsPaginatorOptions struct {
	// The maximum number of response records to return in each call. If the number of
	// remaining response records exceeds the specified MaxRecords value, a value is
	// returned in a marker field of the response. You can retrieve the next set of
	// records by retrying the command with the returned marker value. Default: 100
	// Constraints: minimum 20, maximum 100.
	Limit int32

	// Set to true if pagination should stop if the service returns a pagination token
	// that matches the most recent token provided to the service.
	StopOnDuplicateToken bool
}

// DescribeSnapshotCopyGrantsPaginator is a paginator for
// DescribeSnapshotCopyGrants
type DescribeSnapshotCopyGrantsPaginator struct {
	options   DescribeSnapshotCopyGrantsPaginatorOptions
	client    DescribeSnapshotCopyGrantsAPIClient
	params    *DescribeSnapshotCopyGrantsInput
	nextToken *string
	firstPage bool
}

// NewDescribeSnapshotCopyGrantsPaginator returns a new
// DescribeSnapshotCopyGrantsPaginator
func NewDescribeSnapshotCopyGrantsPaginator(client DescribeSnapshotCopyGrantsAPIClient, params *DescribeSnapshotCopyGrantsInput, optFns ...func(*DescribeSnapshotCopyGrantsPaginatorOptions)) *DescribeSnapshotCopyGrantsPaginator {
	if params == nil {
		params = &DescribeSnapshotCopyGrantsInput{}
	}

	options := DescribeSnapshotCopyGrantsPaginatorOptions{}
	if params.MaxRecords != nil {
		options.Limit = *params.MaxRecords
	}

	for _, fn := range optFns {
		fn(&options)
	}

	return &DescribeSnapshotCopyGrantsPaginator{
		options:   options,
		client:    client,
		params:    params,
		firstPage: true,
		nextToken: params.Marker,
	}
}

// HasMorePages returns a boolean indicating whether more pages are available
func (p *DescribeSnapshotCopyGrantsPaginator) HasMorePages() bool {
	return p.firstPage || (p.nextToken != nil && len(*p.nextToken) != 0)
}

// NextPage retrieves the next DescribeSnapshotCopyGrants page.
func (p *DescribeSnapshotCopyGrantsPaginator) NextPage(ctx context.Context, optFns ...func(*Options)) (*DescribeSnapshotCopyGrantsOutput, error) {
	if !p.HasMorePages() {
		return nil, fmt.Errorf("no more pages available")
	}

	params := *p.params
	params.Marker = p.nextToken

	var limit *int32
	if p.options.Limit > 0 {
		limit = &p.options.Limit
	}
	params.MaxRecords = limit

	result, err := p.client.DescribeSnapshotCopyGrants(ctx, &params, optFns...)
	if err != nil {
		return nil, err
	}
	p.firstPage = false

	prevToken := p.nextToken
	p.nextToken = result.Marker

	if p.options.StopOnDuplicateToken &&
		prevToken != nil &&
		p.nextToken != nil &&
		*prevToken == *p.nextToken {
		p.nextToken = nil
	}

	return result, nil
}

func newServiceMetadataMiddleware_opDescribeSnapshotCopyGrants(region string) *awsmiddleware.RegisterServiceMetadata {
	return &awsmiddleware.RegisterServiceMetadata{
		Region:        region,
		ServiceID:     ServiceID,
		SigningName:   "redshift",
		OperationName: "DescribeSnapshotCopyGrants",
	}
}