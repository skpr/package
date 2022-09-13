// Code generated by smithy-go-codegen DO NOT EDIT.

package types

import (
	smithydocument "github.com/aws/smithy-go/document"
	"time"
)

// Indicates that an Amazon S3 canned ACL should be set to control ownership of
// stored query results. When Athena stores query results in Amazon S3, the canned
// ACL is set with the x-amz-acl request header. For more information about S3
// Object Ownership, see Object Ownership settings
// (https://docs.aws.amazon.com/AmazonS3/latest/userguide/about-object-ownership.html#object-ownership-overview)
// in the Amazon S3 User Guide.
type AclConfiguration struct {

	// The Amazon S3 canned ACL that Athena should specify when storing query results.
	// Currently the only supported canned ACL is BUCKET_OWNER_FULL_CONTROL. If a query
	// runs in a workgroup and the workgroup overrides client-side settings, then the
	// Amazon S3 canned ACL specified in the workgroup's settings is used for all
	// queries that run in the workgroup. For more information about Amazon S3 canned
	// ACLs, see Canned ACL
	// (https://docs.aws.amazon.com/AmazonS3/latest/userguide/acl-overview.html#canned-acl)
	// in the Amazon S3 User Guide.
	//
	// This member is required.
	S3AclOption S3AclOption

	noSmithyDocumentSerde
}

// Provides information about an Athena query error. The AthenaError feature
// provides standardized error information to help you understand failed queries
// and take steps after a query failure occurs. AthenaError includes an
// ErrorCategory field that specifies whether the cause of the failed query is due
// to system error, user error, or other error.
type AthenaError struct {

	// An integer value that specifies the category of a query failure error. The
	// following list shows the category for each integer value. 1 - System 2 - User 3
	// - Other
	ErrorCategory *int32

	// Contains a short description of the error that occurred.
	ErrorMessage *string

	// An integer value that provides specific information about an Athena query error.
	// For the meaning of specific values, see the Error Type Reference
	// (https://docs.aws.amazon.com/athena/latest/ug/error-reference.html#error-reference-error-type-reference)
	// in the Amazon Athena User Guide.
	ErrorType *int32

	// True if the query might succeed if resubmitted.
	Retryable bool

	noSmithyDocumentSerde
}

// Contains metadata for a column in a table.
type Column struct {

	// The name of the column.
	//
	// This member is required.
	Name *string

	// Optional information about the column.
	Comment *string

	// The data type of the column.
	Type *string

	noSmithyDocumentSerde
}

// Information about the columns in a query execution result.
type ColumnInfo struct {

	// The name of the column.
	//
	// This member is required.
	Name *string

	// The data type of the column.
	//
	// This member is required.
	Type *string

	// Indicates whether values in the column are case-sensitive.
	CaseSensitive bool

	// The catalog to which the query results belong.
	CatalogName *string

	// A column label.
	Label *string

	// Indicates the column's nullable status.
	Nullable ColumnNullable

	// For DECIMAL data types, specifies the total number of digits, up to 38. For
	// performance reasons, we recommend up to 18 digits.
	Precision int32

	// For DECIMAL data types, specifies the total number of digits in the fractional
	// part of the value. Defaults to 0.
	Scale int32

	// The schema name (database name) to which the query results belong.
	SchemaName *string

	// The table name for the query results.
	TableName *string

	noSmithyDocumentSerde
}

// Contains metadata information for a database in a data catalog.
type Database struct {

	// The name of the database.
	//
	// This member is required.
	Name *string

	// An optional description of the database.
	Description *string

	// A set of custom key/value pairs.
	Parameters map[string]string

	noSmithyDocumentSerde
}

// Contains information about a data catalog in an Amazon Web Services account.
type DataCatalog struct {

	// The name of the data catalog. The catalog name must be unique for the Amazon Web
	// Services account and can use a maximum of 127 alphanumeric, underscore, at sign,
	// or hyphen characters. The remainder of the length constraint of 256 is reserved
	// for use by Athena.
	//
	// This member is required.
	Name *string

	// The type of data catalog to create: LAMBDA for a federated catalog, HIVE for an
	// external hive metastore, or GLUE for an Glue Data Catalog.
	//
	// This member is required.
	Type DataCatalogType

	// An optional description of the data catalog.
	Description *string

	// Specifies the Lambda function or functions to use for the data catalog. This is
	// a mapping whose values depend on the catalog type.
	//
	// * For the HIVE data catalog
	// type, use the following syntax. The metadata-function parameter is required. The
	// sdk-version parameter is optional and defaults to the currently supported
	// version. metadata-function=lambda_arn, sdk-version=version_number
	//
	// * For the
	// LAMBDA data catalog type, use one of the following sets of required parameters,
	// but not both.
	//
	// * If you have one Lambda function that processes metadata and
	// another for reading the actual data, use the following syntax. Both parameters
	// are required. metadata-function=lambda_arn, record-function=lambda_arn
	//
	// * If you
	// have a composite Lambda function that processes both metadata and data, use the
	// following syntax to specify your Lambda function. function=lambda_arn
	//
	// * The
	// GLUE type takes a catalog ID parameter and is required. The  catalog_id  is the
	// account ID of the Amazon Web Services account to which the Glue catalog belongs.
	// catalog-id=catalog_id
	//
	// * The GLUE data catalog type also applies to the default
	// AwsDataCatalog that already exists in your account, of which you can have only
	// one and cannot modify.
	//
	// * Queries that specify a Glue Data Catalog other than
	// the default AwsDataCatalog must be run on Athena engine version 2.
	Parameters map[string]string

	noSmithyDocumentSerde
}

// The summary information for the data catalog, which includes its name and type.
type DataCatalogSummary struct {

	// The name of the data catalog. The catalog name is unique for the Amazon Web
	// Services account and can use a maximum of 127 alphanumeric, underscore, at sign,
	// or hyphen characters. The remainder of the length constraint of 256 is reserved
	// for use by Athena.
	CatalogName *string

	// The data catalog type.
	Type DataCatalogType

	noSmithyDocumentSerde
}

// A piece of data (a field in the table).
type Datum struct {

	// The value of the datum.
	VarCharValue *string

	noSmithyDocumentSerde
}

// If query results are encrypted in Amazon S3, indicates the encryption option
// used (for example, SSE_KMS or CSE_KMS) and key information.
type EncryptionConfiguration struct {

	// Indicates whether Amazon S3 server-side encryption with Amazon S3-managed keys
	// (SSE_S3), server-side encryption with KMS-managed keys (SSE_KMS), or client-side
	// encryption with KMS-managed keys (CSE_KMS) is used. If a query runs in a
	// workgroup and the workgroup overrides client-side settings, then the workgroup's
	// setting for encryption is used. It specifies whether query results must be
	// encrypted, for all queries that run in this workgroup.
	//
	// This member is required.
	EncryptionOption EncryptionOption

	// For SSE_KMS and CSE_KMS, this is the KMS key ARN or ID.
	KmsKey *string

	noSmithyDocumentSerde
}

// The Athena engine version for running queries.
type EngineVersion struct {

	// Read only. The engine version on which the query runs. If the user requests a
	// valid engine version other than Auto, the effective engine version is the same
	// as the engine version that the user requested. If the user requests Auto, the
	// effective engine version is chosen by Athena. When a request to update the
	// engine version is made by a CreateWorkGroup or UpdateWorkGroup operation, the
	// EffectiveEngineVersion field is ignored.
	EffectiveEngineVersion *string

	// The engine version requested by the user. Possible values are determined by the
	// output of ListEngineVersions, including Auto. The default is Auto.
	SelectedEngineVersion *string

	noSmithyDocumentSerde
}

// A query, where QueryString contains the SQL statements that make up the query.
type NamedQuery struct {

	// The database to which the query belongs.
	//
	// This member is required.
	Database *string

	// The query name.
	//
	// This member is required.
	Name *string

	// The SQL statements that make up the query.
	//
	// This member is required.
	QueryString *string

	// The query description.
	Description *string

	// The unique identifier of the query.
	NamedQueryId *string

	// The name of the workgroup that contains the named query.
	WorkGroup *string

	noSmithyDocumentSerde
}

// A prepared SQL statement for use with Athena.
type PreparedStatement struct {

	// The description of the prepared statement.
	Description *string

	// The last modified time of the prepared statement.
	LastModifiedTime *time.Time

	// The query string for the prepared statement.
	QueryStatement *string

	// The name of the prepared statement.
	StatementName *string

	// The name of the workgroup to which the prepared statement belongs.
	WorkGroupName *string

	noSmithyDocumentSerde
}

// The name and last modified time of the prepared statement.
type PreparedStatementSummary struct {

	// The last modified time of the prepared statement.
	LastModifiedTime *time.Time

	// The name of the prepared statement.
	StatementName *string

	noSmithyDocumentSerde
}

// Information about a single instance of a query execution.
type QueryExecution struct {

	// The engine version that executed the query.
	EngineVersion *EngineVersion

	// A list of values for the parameters in a query. The values are applied
	// sequentially to the parameters in the query in the order in which the parameters
	// occur.
	ExecutionParameters []string

	// The SQL query statements which the query execution ran.
	Query *string

	// The database in which the query execution occurred.
	QueryExecutionContext *QueryExecutionContext

	// The unique identifier for each query execution.
	QueryExecutionId *string

	// The location in Amazon S3 where query results were stored and the encryption
	// option, if any, used for query results. These are known as "client-side
	// settings". If workgroup settings override client-side settings, then the query
	// uses the location for the query results and the encryption configuration that
	// are specified for the workgroup.
	ResultConfiguration *ResultConfiguration

	// The type of query statement that was run. DDL indicates DDL query statements.
	// DML indicates DML (Data Manipulation Language) query statements, such as CREATE
	// TABLE AS SELECT. UTILITY indicates query statements other than DDL and DML, such
	// as SHOW CREATE TABLE, or DESCRIBE TABLE.
	StatementType StatementType

	// Query execution statistics, such as the amount of data scanned, the amount of
	// time that the query took to process, and the type of statement that was run.
	Statistics *QueryExecutionStatistics

	// The completion date, current state, submission time, and state change reason (if
	// applicable) for the query execution.
	Status *QueryExecutionStatus

	// The name of the workgroup in which the query ran.
	WorkGroup *string

	noSmithyDocumentSerde
}

// The database and data catalog context in which the query execution occurs.
type QueryExecutionContext struct {

	// The name of the data catalog used in the query execution.
	Catalog *string

	// The name of the database used in the query execution. The database must exist in
	// the catalog.
	Database *string

	noSmithyDocumentSerde
}

// The amount of data scanned during the query execution and the amount of time
// that it took to execute, and the type of statement that was run.
type QueryExecutionStatistics struct {

	// The location and file name of a data manifest file. The manifest file is saved
	// to the Athena query results location in Amazon S3. The manifest file tracks
	// files that the query wrote to Amazon S3. If the query fails, the manifest file
	// also tracks files that the query intended to write. The manifest is useful for
	// identifying orphaned files resulting from a failed query. For more information,
	// see Working with Query Results, Output Files, and Query History
	// (https://docs.aws.amazon.com/athena/latest/ug/querying.html) in the Amazon
	// Athena User Guide.
	DataManifestLocation *string

	// The number of bytes in the data that was queried.
	DataScannedInBytes *int64

	// The number of milliseconds that the query took to execute.
	EngineExecutionTimeInMillis *int64

	// The number of milliseconds that Athena took to plan the query processing flow.
	// This includes the time spent retrieving table partitions from the data source.
	// Note that because the query engine performs the query planning, query planning
	// time is a subset of engine processing time.
	QueryPlanningTimeInMillis *int64

	// The number of milliseconds that the query was in your query queue waiting for
	// resources. Note that if transient errors occur, Athena might automatically add
	// the query back to the queue.
	QueryQueueTimeInMillis *int64

	// The number of milliseconds that Athena took to finalize and publish the query
	// results after the query engine finished running the query.
	ServiceProcessingTimeInMillis *int64

	// The number of milliseconds that Athena took to run the query.
	TotalExecutionTimeInMillis *int64

	noSmithyDocumentSerde
}

// The completion date, current state, submission time, and state change reason (if
// applicable) for the query execution.
type QueryExecutionStatus struct {

	// Provides information about an Athena query error.
	AthenaError *AthenaError

	// The date and time that the query completed.
	CompletionDateTime *time.Time

	// The state of query execution. QUEUED indicates that the query has been submitted
	// to the service, and Athena will execute the query as soon as resources are
	// available. RUNNING indicates that the query is in execution phase. SUCCEEDED
	// indicates that the query completed without errors. FAILED indicates that the
	// query experienced an error and did not complete processing. CANCELLED indicates
	// that a user input interrupted query execution. Athena automatically retries your
	// queries in cases of certain transient errors. As a result, you may see the query
	// state transition from RUNNING or FAILED to QUEUED.
	State QueryExecutionState

	// Further detail about the status of the query.
	StateChangeReason *string

	// The date and time that the query was submitted.
	SubmissionDateTime *time.Time

	noSmithyDocumentSerde
}

// The query execution timeline, statistics on input and output rows and bytes, and
// the different query stages that form the query execution plan.
type QueryRuntimeStatistics struct {

	// Stage statistics such as input and output rows and bytes, execution time, and
	// stage state. This information also includes substages and the query stage plan.
	OutputStage *QueryStage

	// Statistics such as input rows and bytes read by the query, rows and bytes output
	// by the query, and the number of rows written by the query.
	Rows *QueryRuntimeStatisticsRows

	// Timeline statistics such as query queue time, planning time, execution time,
	// service processing time, and total execution time.
	Timeline *QueryRuntimeStatisticsTimeline

	noSmithyDocumentSerde
}

// Statistics such as input rows and bytes read by the query, rows and bytes output
// by the query, and the number of rows written by the query.
type QueryRuntimeStatisticsRows struct {

	// The number of bytes read to execute the query.
	InputBytes *int64

	// The number of rows read to execute the query.
	InputRows *int64

	// The number of bytes returned by the query.
	OutputBytes *int64

	// The number of rows returned by the query.
	OutputRows *int64

	noSmithyDocumentSerde
}

// Timeline statistics such as query queue time, planning time, execution time,
// service processing time, and total execution time.
type QueryRuntimeStatisticsTimeline struct {

	// The number of milliseconds that the query took to execute.
	EngineExecutionTimeInMillis *int64

	// The number of milliseconds that Athena took to plan the query processing flow.
	// This includes the time spent retrieving table partitions from the data source.
	// Note that because the query engine performs the query planning, query planning
	// time is a subset of engine processing time.
	QueryPlanningTimeInMillis *int64

	// The number of milliseconds that the query was in your query queue waiting for
	// resources. Note that if transient errors occur, Athena might automatically add
	// the query back to the queue.
	QueryQueueTimeInMillis *int64

	// The number of milliseconds that Athena took to finalize and publish the query
	// results after the query engine finished running the query.
	ServiceProcessingTimeInMillis *int64

	// The number of milliseconds that Athena took to run the query.
	TotalExecutionTimeInMillis *int64

	noSmithyDocumentSerde
}

// Stage statistics such as input and output rows and bytes, execution time and
// stage state. This information also includes substages and the query stage plan.
type QueryStage struct {

	// Time taken to execute this stage.
	ExecutionTime *int64

	// The number of bytes input into the stage for execution.
	InputBytes *int64

	// The number of rows input into the stage for execution.
	InputRows *int64

	// The number of bytes output from the stage after execution.
	OutputBytes *int64

	// The number of rows output from the stage after execution.
	OutputRows *int64

	// Stage plan information such as name, identifier, sub plans, and source stages.
	QueryStagePlan *QueryStagePlanNode

	// The identifier for a stage.
	StageId *int64

	// State of the stage after query execution.
	State *string

	// List of sub query stages that form this stage execution plan.
	SubStages []QueryStage

	noSmithyDocumentSerde
}

// Stage plan information such as name, identifier, sub plans, and remote sources.
type QueryStagePlanNode struct {

	// Stage plan information such as name, identifier, sub plans, and remote sources
	// of child plan nodes/
	Children []QueryStagePlanNode

	// Information about the operation this query stage plan node is performing.
	Identifier *string

	// Name of the query stage plan that describes the operation this stage is
	// performing as part of query execution.
	Name *string

	// Source plan node IDs.
	RemoteSources []string

	noSmithyDocumentSerde
}

// The location in Amazon S3 where query results are stored and the encryption
// option, if any, used for query results. These are known as "client-side
// settings". If workgroup settings override client-side settings, then the query
// uses the workgroup settings.
type ResultConfiguration struct {

	// Indicates that an Amazon S3 canned ACL should be set to control ownership of
	// stored query results. Currently the only supported canned ACL is
	// BUCKET_OWNER_FULL_CONTROL. This is a client-side setting. If workgroup settings
	// override client-side settings, then the query uses the ACL configuration that is
	// specified for the workgroup, and also uses the location for storing query
	// results specified in the workgroup. For more information, see
	// WorkGroupConfiguration$EnforceWorkGroupConfiguration and Workgroup Settings
	// Override Client-Side Settings
	// (https://docs.aws.amazon.com/athena/latest/ug/workgroups-settings-override.html).
	AclConfiguration *AclConfiguration

	// If query results are encrypted in Amazon S3, indicates the encryption option
	// used (for example, SSE_KMS or CSE_KMS) and key information. This is a
	// client-side setting. If workgroup settings override client-side settings, then
	// the query uses the encryption configuration that is specified for the workgroup,
	// and also uses the location for storing query results specified in the workgroup.
	// See WorkGroupConfiguration$EnforceWorkGroupConfiguration and Workgroup Settings
	// Override Client-Side Settings
	// (https://docs.aws.amazon.com/athena/latest/ug/workgroups-settings-override.html).
	EncryptionConfiguration *EncryptionConfiguration

	// The Amazon Web Services account ID that you expect to be the owner of the Amazon
	// S3 bucket specified by ResultConfiguration$OutputLocation. If set, Athena uses
	// the value for ExpectedBucketOwner when it makes Amazon S3 calls to your
	// specified output location. If the ExpectedBucketOwner Amazon Web Services
	// account ID does not match the actual owner of the Amazon S3 bucket, the call
	// fails with a permissions error. This is a client-side setting. If workgroup
	// settings override client-side settings, then the query uses the
	// ExpectedBucketOwner setting that is specified for the workgroup, and also uses
	// the location for storing query results specified in the workgroup. See
	// WorkGroupConfiguration$EnforceWorkGroupConfiguration and Workgroup Settings
	// Override Client-Side Settings
	// (https://docs.aws.amazon.com/athena/latest/ug/workgroups-settings-override.html).
	ExpectedBucketOwner *string

	// The location in Amazon S3 where your query results are stored, such as
	// s3://path/to/query/bucket/. To run the query, you must specify the query results
	// location using one of the ways: either for individual queries using either this
	// setting (client-side), or in the workgroup, using WorkGroupConfiguration. If
	// none of them is set, Athena issues an error that no output location is provided.
	// For more information, see Query Results
	// (https://docs.aws.amazon.com/athena/latest/ug/querying.html). If workgroup
	// settings override client-side settings, then the query uses the settings
	// specified for the workgroup. See
	// WorkGroupConfiguration$EnforceWorkGroupConfiguration.
	OutputLocation *string

	noSmithyDocumentSerde
}

// The information about the updates in the query results, such as output location
// and encryption configuration for the query results.
type ResultConfigurationUpdates struct {

	// The ACL configuration for the query results.
	AclConfiguration *AclConfiguration

	// The encryption configuration for the query results.
	EncryptionConfiguration *EncryptionConfiguration

	// The Amazon Web Services account ID that you expect to be the owner of the Amazon
	// S3 bucket specified by ResultConfiguration$OutputLocation. If set, Athena uses
	// the value for ExpectedBucketOwner when it makes Amazon S3 calls to your
	// specified output location. If the ExpectedBucketOwner Amazon Web Services
	// account ID does not match the actual owner of the Amazon S3 bucket, the call
	// fails with a permissions error. If workgroup settings override client-side
	// settings, then the query uses the ExpectedBucketOwner setting that is specified
	// for the workgroup, and also uses the location for storing query results
	// specified in the workgroup. See
	// WorkGroupConfiguration$EnforceWorkGroupConfiguration and Workgroup Settings
	// Override Client-Side Settings
	// (https://docs.aws.amazon.com/athena/latest/ug/workgroups-settings-override.html).
	ExpectedBucketOwner *string

	// The location in Amazon S3 where your query results are stored, such as
	// s3://path/to/query/bucket/. For more information, see Query Results
	// (https://docs.aws.amazon.com/athena/latest/ug/querying.html) If workgroup
	// settings override client-side settings, then the query uses the location for the
	// query results and the encryption configuration that are specified for the
	// workgroup. The "workgroup settings override" is specified in
	// EnforceWorkGroupConfiguration (true/false) in the WorkGroupConfiguration. See
	// WorkGroupConfiguration$EnforceWorkGroupConfiguration.
	OutputLocation *string

	// If set to true, indicates that the previously-specified ACL configuration for
	// queries in this workgroup should be ignored and set to null. If set to false or
	// not set, and a value is present in the AclConfiguration of
	// ResultConfigurationUpdates, the AclConfiguration in the workgroup's
	// ResultConfiguration is updated with the new value. For more information, see
	// Workgroup Settings Override Client-Side Settings
	// (https://docs.aws.amazon.com/athena/latest/ug/workgroups-settings-override.html).
	RemoveAclConfiguration *bool

	// If set to "true", indicates that the previously-specified encryption
	// configuration (also known as the client-side setting) for queries in this
	// workgroup should be ignored and set to null. If set to "false" or not set, and a
	// value is present in the EncryptionConfiguration in ResultConfigurationUpdates
	// (the client-side setting), the EncryptionConfiguration in the workgroup's
	// ResultConfiguration will be updated with the new value. For more information,
	// see Workgroup Settings Override Client-Side Settings
	// (https://docs.aws.amazon.com/athena/latest/ug/workgroups-settings-override.html).
	RemoveEncryptionConfiguration *bool

	// If set to "true", removes the Amazon Web Services account ID previously
	// specified for ResultConfiguration$ExpectedBucketOwner. If set to "false" or not
	// set, and a value is present in the ExpectedBucketOwner in
	// ResultConfigurationUpdates (the client-side setting), the ExpectedBucketOwner in
	// the workgroup's ResultConfiguration is updated with the new value. For more
	// information, see Workgroup Settings Override Client-Side Settings
	// (https://docs.aws.amazon.com/athena/latest/ug/workgroups-settings-override.html).
	RemoveExpectedBucketOwner *bool

	// If set to "true", indicates that the previously-specified query results location
	// (also known as a client-side setting) for queries in this workgroup should be
	// ignored and set to null. If set to "false" or not set, and a value is present in
	// the OutputLocation in ResultConfigurationUpdates (the client-side setting), the
	// OutputLocation in the workgroup's ResultConfiguration will be updated with the
	// new value. For more information, see Workgroup Settings Override Client-Side
	// Settings
	// (https://docs.aws.amazon.com/athena/latest/ug/workgroups-settings-override.html).
	RemoveOutputLocation *bool

	noSmithyDocumentSerde
}

// The metadata and rows that make up a query result set. The metadata describes
// the column structure and data types. To return a ResultSet object, use
// GetQueryResults.
type ResultSet struct {

	// The metadata that describes the column structure and data types of a table of
	// query results.
	ResultSetMetadata *ResultSetMetadata

	// The rows in the table.
	Rows []Row

	noSmithyDocumentSerde
}

// The metadata that describes the column structure and data types of a table of
// query results. To return a ResultSetMetadata object, use GetQueryResults.
type ResultSetMetadata struct {

	// Information about the columns returned in a query result metadata.
	ColumnInfo []ColumnInfo

	noSmithyDocumentSerde
}

// The rows that make up a query result table.
type Row struct {

	// The data that populates a row in a query result table.
	Data []Datum

	noSmithyDocumentSerde
}

// Contains metadata for a table.
type TableMetadata struct {

	// The name of the table.
	//
	// This member is required.
	Name *string

	// A list of the columns in the table.
	Columns []Column

	// The time that the table was created.
	CreateTime *time.Time

	// The last time the table was accessed.
	LastAccessTime *time.Time

	// A set of custom key/value pairs for table properties.
	Parameters map[string]string

	// A list of the partition keys in the table.
	PartitionKeys []Column

	// The type of table. In Athena, only EXTERNAL_TABLE is supported.
	TableType *string

	noSmithyDocumentSerde
}

// A label that you assign to a resource. In Athena, a resource can be a workgroup
// or data catalog. Each tag consists of a key and an optional value, both of which
// you define. For example, you can use tags to categorize Athena workgroups or
// data catalogs by purpose, owner, or environment. Use a consistent set of tag
// keys to make it easier to search and filter workgroups or data catalogs in your
// account. For best practices, see Tagging Best Practices
// (https://aws.amazon.com/answers/account-management/aws-tagging-strategies/). Tag
// keys can be from 1 to 128 UTF-8 Unicode characters, and tag values can be from 0
// to 256 UTF-8 Unicode characters. Tags can use letters and numbers representable
// in UTF-8, and the following characters: + - = . _ : / @. Tag keys and values are
// case-sensitive. Tag keys must be unique per resource. If you specify more than
// one tag, separate them by commas.
type Tag struct {

	// A tag key. The tag key length is from 1 to 128 Unicode characters in UTF-8. You
	// can use letters and numbers representable in UTF-8, and the following
	// characters: + - = . _ : / @. Tag keys are case-sensitive and must be unique per
	// resource.
	Key *string

	// A tag value. The tag value length is from 0 to 256 Unicode characters in UTF-8.
	// You can use letters and numbers representable in UTF-8, and the following
	// characters: + - = . _ : / @. Tag values are case-sensitive.
	Value *string

	noSmithyDocumentSerde
}

// Information about a named query ID that could not be processed.
type UnprocessedNamedQueryId struct {

	// The error code returned when the processing request for the named query failed,
	// if applicable.
	ErrorCode *string

	// The error message returned when the processing request for the named query
	// failed, if applicable.
	ErrorMessage *string

	// The unique identifier of the named query.
	NamedQueryId *string

	noSmithyDocumentSerde
}

// The name of a prepared statement that could not be returned.
type UnprocessedPreparedStatementName struct {

	// The error code returned when the request for the prepared statement failed.
	ErrorCode *string

	// The error message containing the reason why the prepared statement could not be
	// returned. The following error messages are possible:
	//
	// * INVALID_INPUT - The name
	// of the prepared statement that was provided is not valid (for example, the name
	// is too long).
	//
	// * STATEMENT_NOT_FOUND - A prepared statement with the name
	// provided could not be found.
	//
	// * UNAUTHORIZED - The requester does not have
	// permission to access the workgroup that contains the prepared statement.
	ErrorMessage *string

	// The name of a prepared statement that could not be returned due to an error.
	StatementName *string

	noSmithyDocumentSerde
}

// Describes a query execution that failed to process.
type UnprocessedQueryExecutionId struct {

	// The error code returned when the query execution failed to process, if
	// applicable.
	ErrorCode *string

	// The error message returned when the query execution failed to process, if
	// applicable.
	ErrorMessage *string

	// The unique identifier of the query execution.
	QueryExecutionId *string

	noSmithyDocumentSerde
}

// A workgroup, which contains a name, description, creation time, state, and other
// configuration, listed under WorkGroup$Configuration. Each workgroup enables you
// to isolate queries for you or your group of users from other queries in the same
// account, to configure the query results location and the encryption
// configuration (known as workgroup settings), to enable sending query metrics to
// Amazon CloudWatch, and to establish per-query data usage control limits for all
// queries in a workgroup. The workgroup settings override is specified in
// EnforceWorkGroupConfiguration (true/false) in the WorkGroupConfiguration. See
// WorkGroupConfiguration$EnforceWorkGroupConfiguration.
type WorkGroup struct {

	// The workgroup name.
	//
	// This member is required.
	Name *string

	// The configuration of the workgroup, which includes the location in Amazon S3
	// where query results are stored, the encryption configuration, if any, used for
	// query results; whether the Amazon CloudWatch Metrics are enabled for the
	// workgroup; whether workgroup settings override client-side settings; and the
	// data usage limits for the amount of data scanned per query or per workgroup. The
	// workgroup settings override is specified in EnforceWorkGroupConfiguration
	// (true/false) in the WorkGroupConfiguration. See
	// WorkGroupConfiguration$EnforceWorkGroupConfiguration.
	Configuration *WorkGroupConfiguration

	// The date and time the workgroup was created.
	CreationTime *time.Time

	// The workgroup description.
	Description *string

	// The state of the workgroup: ENABLED or DISABLED.
	State WorkGroupState

	noSmithyDocumentSerde
}

// The configuration of the workgroup, which includes the location in Amazon S3
// where query results are stored, the encryption option, if any, used for query
// results, whether the Amazon CloudWatch Metrics are enabled for the workgroup and
// whether workgroup settings override query settings, and the data usage limits
// for the amount of data scanned per query or per workgroup. The workgroup
// settings override is specified in EnforceWorkGroupConfiguration (true/false) in
// the WorkGroupConfiguration. See
// WorkGroupConfiguration$EnforceWorkGroupConfiguration.
type WorkGroupConfiguration struct {

	// The upper data usage limit (cutoff) for the amount of bytes a single query in a
	// workgroup is allowed to scan.
	BytesScannedCutoffPerQuery *int64

	// If set to "true", the settings for the workgroup override client-side settings.
	// If set to "false", client-side settings are used. For more information, see
	// Workgroup Settings Override Client-Side Settings
	// (https://docs.aws.amazon.com/athena/latest/ug/workgroups-settings-override.html).
	EnforceWorkGroupConfiguration *bool

	// The engine version that all queries running on the workgroup use. Queries on the
	// AmazonAthenaPreviewFunctionality workgroup run on the preview engine regardless
	// of this setting.
	EngineVersion *EngineVersion

	// Indicates that the Amazon CloudWatch metrics are enabled for the workgroup.
	PublishCloudWatchMetricsEnabled *bool

	// If set to true, allows members assigned to a workgroup to reference Amazon S3
	// Requester Pays buckets in queries. If set to false, workgroup members cannot
	// query data from Requester Pays buckets, and queries that retrieve data from
	// Requester Pays buckets cause an error. The default is false. For more
	// information about Requester Pays buckets, see Requester Pays Buckets
	// (https://docs.aws.amazon.com/AmazonS3/latest/dev/RequesterPaysBuckets.html) in
	// the Amazon Simple Storage Service Developer Guide.
	RequesterPaysEnabled *bool

	// The configuration for the workgroup, which includes the location in Amazon S3
	// where query results are stored and the encryption option, if any, used for query
	// results. To run the query, you must specify the query results location using one
	// of the ways: either in the workgroup using this setting, or for individual
	// queries (client-side), using ResultConfiguration$OutputLocation. If none of them
	// is set, Athena issues an error that no output location is provided. For more
	// information, see Query Results
	// (https://docs.aws.amazon.com/athena/latest/ug/querying.html).
	ResultConfiguration *ResultConfiguration

	noSmithyDocumentSerde
}

// The configuration information that will be updated for this workgroup, which
// includes the location in Amazon S3 where query results are stored, the
// encryption option, if any, used for query results, whether the Amazon CloudWatch
// Metrics are enabled for the workgroup, whether the workgroup settings override
// the client-side settings, and the data usage limit for the amount of bytes
// scanned per query, if it is specified.
type WorkGroupConfigurationUpdates struct {

	// The upper limit (cutoff) for the amount of bytes a single query in a workgroup
	// is allowed to scan.
	BytesScannedCutoffPerQuery *int64

	// If set to "true", the settings for the workgroup override client-side settings.
	// If set to "false" client-side settings are used. For more information, see
	// Workgroup Settings Override Client-Side Settings
	// (https://docs.aws.amazon.com/athena/latest/ug/workgroups-settings-override.html).
	EnforceWorkGroupConfiguration *bool

	// The engine version requested when a workgroup is updated. After the update, all
	// queries on the workgroup run on the requested engine version. If no value was
	// previously set, the default is Auto. Queries on the
	// AmazonAthenaPreviewFunctionality workgroup run on the preview engine regardless
	// of this setting.
	EngineVersion *EngineVersion

	// Indicates whether this workgroup enables publishing metrics to Amazon
	// CloudWatch.
	PublishCloudWatchMetricsEnabled *bool

	// Indicates that the data usage control limit per query is removed.
	// WorkGroupConfiguration$BytesScannedCutoffPerQuery
	RemoveBytesScannedCutoffPerQuery *bool

	// If set to true, allows members assigned to a workgroup to specify Amazon S3
	// Requester Pays buckets in queries. If set to false, workgroup members cannot
	// query data from Requester Pays buckets, and queries that retrieve data from
	// Requester Pays buckets cause an error. The default is false. For more
	// information about Requester Pays buckets, see Requester Pays Buckets
	// (https://docs.aws.amazon.com/AmazonS3/latest/dev/RequesterPaysBuckets.html) in
	// the Amazon Simple Storage Service Developer Guide.
	RequesterPaysEnabled *bool

	// The result configuration information about the queries in this workgroup that
	// will be updated. Includes the updated results location and an updated option for
	// encrypting query results.
	ResultConfigurationUpdates *ResultConfigurationUpdates

	noSmithyDocumentSerde
}

// The summary information for the workgroup, which includes its name, state,
// description, and the date and time it was created.
type WorkGroupSummary struct {

	// The workgroup creation date and time.
	CreationTime *time.Time

	// The workgroup description.
	Description *string

	// The engine version setting for all queries on the workgroup. Queries on the
	// AmazonAthenaPreviewFunctionality workgroup run on the preview engine regardless
	// of this setting.
	EngineVersion *EngineVersion

	// The name of the workgroup.
	Name *string

	// The state of the workgroup.
	State WorkGroupState

	noSmithyDocumentSerde
}

type noSmithyDocumentSerde = smithydocument.NoSerde