package pipes_collection

import "github.com/turbot/tailpipe-plugin-pipes/helpers"

// AuditLogRow is the struct containing the enriched data for an AuditRecord
type AuditLogRow struct {

	// Metadata
	TpID              string             `json:"tp_id" parquet:"name=tp_id, type=BYTE_ARRAY, convertedtype=UTF8"`
	TpSourceType      string             `json:"tp_source_type" parquet:"name=tp_source_type, type=BYTE_ARRAY, convertedtype=UTF8, encoding=PLAIN_DICTIONARY"`
	TpSourceName      *string            `json:"tp_source_name" parquet:"name=tp_source_name, type=BYTE_ARRAY, convertedtype=UTF8, encoding=PLAIN_DICTIONARY"`
	TpSourceLocation  *string            `json:"tp_source_location" parquet:"name=tp_source_location, type=BYTE_ARRAY, convertedtype=UTF8, encoding=PLAIN_DICTIONARY"`
	TpIngestTimestamp helpers.UnixMillis `json:"tp_ingest_timestamp" parquet:"name=tp_ingest_timestamp, type=INT64, convertedtype=TIMESTAMP_MILLIS"`

	// Standardized
	TpTimestamp     helpers.UnixMillis `json:"tp_timestamp" parquet:"name=tp_timestamp, type=INT64, convertedtype=TIMESTAMP_MILLIS"`
	TpSourceIP      *string            `json:"tp_source_ip" parquet:"name=tp_source_ip, type=BYTE_ARRAY, convertedtype=UTF8"`
	TpDestinationIP *string            `json:"tp_destination_ip" parquet:"name=tp_destination_ip, type=BYTE_ARRAY, convertedtype=UTF8"`

	// Hive fields
	TpCollection string `json:"tp_collection" parquet:"name=tp_collection, type=BYTE_ARRAY, convertedtype=UTF8, encoding=PLAIN_DICTIONARY"`
	TpConnection string `json:"tp_connection" parquet:"name=tp_connection, type=BYTE_ARRAY, convertedtype=UTF8, encoding=PLAIN_DICTIONARY"`
	TpYear       int32  `json:"tp_year" parquet:"name=tp_year, type=INT32, convertedtype=INT32"`
	TpMonth      int32  `json:"tp_month" parquet:"name=tp_month, type=INT32, convertedtype=INT32"`
	TpDay        int32  `json:"tp_day" parquet:"name=tp_day, type=INT32, convertedtype=INT32"`

	// Searchable
	TpAkas      []string `json:"tp_akas,omitempty" parquet:"name=tp_akas, type=LIST, convertedtype=LIST, valuetype=BYTE_ARRAY, valueconvertedtype=UTF8, encoding=PLAIN_DICTIONARY"`
	TpIps       []string `json:"tp_ips,omitempty" parquet:"name=tp_ips, type=LIST, convertedtype=LIST, valuetype=BYTE_ARRAY, valueconvertedtype=UTF8"`
	TpTags      []string `json:"tp_tags,omitempty" parquet:"name=tp_tags, type=LIST, convertedtype=LIST, valuetype=BYTE_ARRAY, valueconvertedtype=UTF8"`
	TpDomains   []string `json:"tp_domains,omitempty" parquet:"name=tp_domains, type=LIST, convertedtype=LIST, valuetype=BYTE_ARRAY, valueconvertedtype=UTF8, encoding=PLAIN_DICTIONARY"`
	TpEmails    []string `json:"tp_emails,omitempty" parquet:"name=tp_emails, type=LIST, convertedtype=LIST, valuetype=BYTE_ARRAY, valueconvertedtype=UTF8, encoding=PLAIN_DICTIONARY"`
	TpUsernames []string `json:"tp_usernames,omitempty" parquet:"name=tp_usernames, type=LIST, convertedtype=LIST, valuetype=BYTE_ARRAY, valueconvertedtype=UTF8, encoding=PLAIN_DICTIONARY"`

	// The action performed on the resource.
	ActionType     string `json:"action_type" parquet:"name=action_type, type=BYTE_ARRAY, convertedtype=UTF8, encoding=PLAIN_DICTIONARY"`
	ActorAvatarUrl string `json:"actor_avatar_url" parquet:"name=actor_avatar_url, type=BYTE_ARRAY, convertedtype=UTF8, encoding=PLAIN_DICTIONARY"`
	// The display name of an actor.
	ActorDisplayName string `json:"actor_display_name" parquet:"name=actor_display_name, type=BYTE_ARRAY, convertedtype=UTF8, encoding=PLAIN_DICTIONARY"`
	// The handle name of an actor.
	ActorHandle string `json:"actor_handle" parquet:"name=actor_handle, type=BYTE_ARRAY, convertedtype=UTF8, encoding=PLAIN_DICTIONARY"`
	// The unique identifier of an actor.
	ActorId string `json:"actor_id" parquet:"name=actor_id, type=BYTE_ARRAY, convertedtype=UTF8, encoding=PLAIN_DICTIONARY"`
	ActorIp string `json:"actor_ip" parquet:"name=actor_ip, type=BYTE_ARRAY, convertedtype=UTF8"`
	// The time when the audit log was recorded.
	CreatedAt helpers.UnixMillis  `json:"created_at" parquet:"name=created_at, type=INT64, convertedtype=TIMESTAMP_MILLIS"`
	Data      *helpers.JSONString `json:"data" parquet:"name=data, type=BYTE_ARRAY, convertedtype=UTF8"`
	// The unique identifier for an audit log.
	Id string `json:"id" parquet:"name=id, type=BYTE_ARRAY, convertedtype=UTF8"`
	// The handle name for an identity where the action has been performed.
	IdentityHandle string `json:"identity_handle" parquet:"name=identity_handle, type=BYTE_ARRAY, convertedtype=UTF8, encoding=PLAIN_DICTIONARY"`
	// The unique identifier for an identity where the action has been performed.
	IdentityId string  `json:"identity_id" parquet:"name=identity_id, type=BYTE_ARRAY, convertedtype=UTF8, encoding=PLAIN_DICTIONARY"`
	ProcessId  *string `json:"process_id,omitempty" parquet:"name=process_id, type=BYTE_ARRAY, convertedtype=UTF8"`
	// The handle name of the entity on which the action has been performed.
	TargetHandle *string `json:"target_handle,omitempty" parquet:"name=target_handle, type=BYTE_ARRAY, convertedtype=UTF8, encoding=PLAIN_DICTIONARY"`
	// The unique identifier of the entity on which the action has been performed.
	TargetId *string `json:"target_id,omitempty" parquet:"name=target_id, type=BYTE_ARRAY, convertedtype=UTF8, encoding=PLAIN_DICTIONARY"`
	// The unique identifier for the where the action has been performed.
	TenantId string `json:"tenant_id" parquet:"name=tenant_id, type=BYTE_ARRAY, convertedtype=UTF8, encoding=PLAIN_DICTIONARY"`
}

func (r *AuditLogRow) GetTpID() string {
	return r.TpID
}

func (r *AuditLogRow) GetTpTimestamp() int64 {
	return int64(r.TpTimestamp)
}

func (r *AuditLogRow) GetConnection() string {
	return r.TpConnection
}

func (r *AuditLogRow) GetYear() int {
	return int(r.TpYear)
}

func (r *AuditLogRow) GetMonth() int {
	return int(r.TpMonth)
}

func (r *AuditLogRow) GetDay() int {
	return int(r.TpDay)
}
