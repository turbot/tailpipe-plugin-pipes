---
title: "Tailpipe Table: pipes_audit_log - Query Pipes audit logs"
description: "Pipes audit logs capture administrative actions and security events within your Pipes organization."
---

# Table: pipes_audit_log - Query Pipes audit logs

The `pipes_audit_log` table allows you to query data from Pipes audit logs. This table provides detailed information about API calls, resource modifications, security events, and administrative actions within your Pipes environment.

## Configure

Create a [partition](https://tailpipe.io/docs/manage/partition) for `pipes_audit_log` ([examples](https://hub.tailpipe.io/plugins/turbot/pipes/tables/pipes_audit_log#example-configurations)):

```sh
vi ~/.tailpipe/config/pipes.tpc
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

### Collect logs from audit logs API

Collect audit logs from the Pipes audit logs API.

```hcl
connection "pipes" "pipes_org" {
  token      = "tpt_pipestoken"
  org_handle = "org_handle_name"
}

partition "pipes_audit_log" "my_logs" {
  source "pipes_audit_log" {
    connection = connection.pipes.pipes_org
  }
}
```

### Exclude events

Use the filter argument in your partition to exclude specific events and and reduce log storage size.

```hcl
partition "pipes_audit_log" "my_logs_filtered" {
  # Avoid saving unnecessary events, which can drastically reduce local log size
  filter = "action_type != 'workspace.snapshot.create'"

  source "pipes_audit_log_api" {
    connection = connection.pipes.pipes_org
  }
}
```
