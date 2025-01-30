---
title: "Tailpipe Table: pipes_audit_log - Query Pipes audit logs"
description: "Pipes audit logs capture administrative actions and security events within your Pipes organization."
---

# Table: pipes_audit_log - Query Pipes audit logs

The `pipes_audit_log` table allows you to query data from Pipes audit logs. This table provides detailed information about API calls, resource modifications, security events, and administrative actions within your Pipes environment.

## Configure

Create a [partition](https://tailpipe.io/docs/manage/partition) for `pipes_audit_log` ([examples](https://hub.tailpipe.io/plugins/turbot/azure/tables/pipes_audit_log#example-configurations)):

```sh
vi ~/.tailpipe/config/azure.tpc
```

```hcl
connection "pipes" "pipes_organization" {
  token      = "tpt_pipestoken"
  org_handle = "org_handle_name"
}

partition "pipes_audit_log" "my_logs" {
  source "pipes_audit_log" {
    connection = connection.pipes.pipes_organization
  }
}
```

## Collect

[Collect](https://tailpipe.io/docs/manage/collection) logs for all `pipes_audit_log` partitions:

```sh
tailpipe collect pipes_audit_log
```

Or for a single partition:

```sh
tailpipe collect pipes_audit_log.my_logs
```

## Query

**[Explore 40+ example queries for this table →](https://hub.tailpipe.io/plugins/turbot/pipes/queries/pipes_audit_log)**

### Role assigments

List role assignments to check for unexpected or suspicious role changes.

```sql
select
  created_at,
  actor_display_name,
  actor_handle,
  target_handle,
  action_type
from
  pipes_audit_log
where
  action_type = 'role_assignment'
order by
  created_at desc;
```

### Top 10 events

List the top 10 events and how many times they were called.

```sql
select
  action_type,
  count(*) as event_count
from
  pipes_audit_log
group by
  action_type
order by
  event_count desc
limit 10;
```

### High Volume Actions by User

Find users generating a high volume of audit log entries to identify potential anomalous activity.

```sql
select
  actor_handle,
  count(*) as event_count,
  date_trunc('minute', created_at) as event_minute
from
  pipes_audit_log
group by
  actor_handle,
  event_minute
having
  count(*) > 100
order by
  event_count desc;
```

## Example Configurations

### Collect logs from a storage account

Collect audit logs stored in a storage account that use the [default blob naming convention](https://learn.microsoft.com/en-us/azure/azure-monitor/essentials/audit-log?tabs=powershell#send-to-azure-storage).

```hcl
connection "azure" "my_logging_account" {
  tenant_id       = "00000000-0000-0000-0000-000000000000"
  subscription_id = "00000000-0000-0000-0000-000000000000"
  client_id       = "00000000-0000-0000-0000-000000000000"
  client_secret   = "my plaintext secret"
}

partition "pipes_audit_log" "my_logs" {
  source "azure_blob_storage" {
    connection   = connection.azure.my_logging_account
    account_name = "storage_account_name"
    container    = "container_name"
  }
}
```

### Collect logs from Monitor audit logs API

Collect audit logs from the Monitor audit logs API.

```hcl
connection "azure" "my_subscription" {
  tenant_id       = "00000000-0000-0000-0000-000000000000"
  subscription_id = "00000000-0000-0000-0000-000000000000"
  client_id       = "00000000-0000-0000-0000-000000000000"
  client_secret   = "my plaintext secret"
}

partition "pipes_audit_log" "my_logs" {
  source "pipes_audit_log_api" {
    connection = connection.azure.my_subscription
  }
}
```

### Exclude read-only events

Use the filter argument in your partition to exclude specific events and and reduce log storage size.

```hcl
partition "pipes_audit_log" "my_logs_filtered" {
  # Avoid saving unnecessary events, which can drastically reduce local log size
  filter = "operation_name != 'Microsoft.Storage/storageAccounts/listKeys/action'"

  source "pipes_audit_log_api" {
    connection = connection.azure.my_subscription
  }
}
```

### Collect logs for a single subscription

Collect logs for a specific subscription.

```hcl
partition "pipes_audit_log" "my_logs_subscription" {
  source "azure_blob_storage" {
    connection   = connection.azure.my_logging_account
    account_name = "storage_account_name"
    container    = "container_name"
    file_layout  = "/SUBSCRIPTIONS/12345678-1234-1234-1234-123456789012/y=%{YEAR:year}/m=%{MONTHNUM:month}/d=%{MONTHDAY:day}/h=%{HOUR:hour}/m=%{MINUTE:minute}/%{DATA:filename}.json"
  }
}
```

## Source Defaults

### azure_blob_storage

This table sets the following defaults for the [azure_blob_storage source](https://hub.tailpipe.io/plugins/turbot/azure/sources/azure_blob_storage#arguments):

| Argument    | Default |
|-------------|---------|
| file_layout | `/SUBSCRIPTIONS/%{DATA:subscription_id}/y=%{YEAR:year}/m=%{MONTHNUM:month}/d=%{MONTHDAY:day}/h=%{HOUR:hour}/m=%{MINUTE:minute}/%{DATA:filename}.json` |
