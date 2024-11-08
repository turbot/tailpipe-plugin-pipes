# Audit Log Table

> [!NOTE]
> To run these against the generated test data:
>
>  cd ../../testdata
>  duckdb
>  CREATE VIEW pipes_audit_log AS SELECT * FROM read_parquet('audit_logs.parquet');

## Identify Suspicious Activity by Repeated Deletion Actions

```sql
select
  actor_handle,
  count(*) as delete_count 
from
  pipes_audit_log 
where
  action_type = 'workspace.delete' 
group by
  actor_handle 
having
  delete_count > 5;
```

## Detect IP Address Reuse Across Different Identities

```sql
select
  tp_ips,
  array_agg(distinct identity_id) as identities 
from
  pipes_audit_log 
group by
  tp_ips 
having
  count(distinct identity_id) > 1;
```

## Track Changes in State for Sensitive Resources

```sql
select
  target_handle,
  count(*) as change_count 
from
  pipes_audit_log 
where
  desired_state != state 
group by
  target_handle 
having
  change_count > 3;
```

## Identify Unusual Login Locations by Actor

```sql
select
  actor_handle,
  array_agg(distinct tp_source_location) as locations 
from
  pipes_audit_log 
group by
  actor_handle 
having
  count(distinct tp_source_location) > 3;
```

## Monitor for Large Number of Failed or Suspicious Login Attempts

```sql
select
  actor_handle,
  count(*) as suspicious_actions 
from
  pipes_audit_log 
where
  action_type in ('login.failed', 'login.suspicious') 
group by
  actor_handle 
having
  suspicious_actions > 5;
```

## Identify Accounts with Multiple Identifiers (Emails or IPs)

```sql
select
  identity_id,
  array_agg(distinct tp_usernames) as usernames,
  array_agg(distinct tp_ips) as ips 
from
  pipes_audit_log 
group by
  identity_id 
having
  cardinality(usernames) > 1 or cardinality(ips) > 1;
```

## Detect High Volume of Activity from Single Source IP

```sql
select
  tp_source_ip,
  count(*) as activity_count 
from
  pipes_audit_log 
group by
  tp_source_ip 
having
  activity_count > 100;
```

## Monitor for Access Outside Business Hours

```sql
select
  actor_handle,
  count(*) as off_hours_activity 
from
  pipes_audit_log 
where
  hour(cast(created_at as timestamp)) not between 9 and 18 
group by
  actor_handle 
having
  off_hours_activity > 5;
```

## Analyze Patterns of Rapid Successive Actions by Same Actor

```sql
with lagged_actions as (
  select 
    actor_handle,
    tp_timestamp,
    lag(tp_timestamp) over (partition by actor_handle order by tp_timestamp) as prev_timestamp 
  from 
    pipes_audit_log
)
select
  actor_handle, 
  count(*) as rapid_actions 
from
  lagged_actions 
where
  (tp_timestamp - prev_timestamp) < 10000 
group by
  actor_handle 
having
  rapid_actions > 5;
```

## Flag Potential Insider Threat by Monitoring for Changes on Sensitive Resources

```sql
select
  actor_handle,
  action_type,
  count(*) as sensitive_actions 
from
  pipes_audit_log 
where
  target_handle in ('database', 'workspace', 'sensitive_resource') 
group by
  actor_handle, 
  action_type 
having
  sensitive_actions > 3;
```

## Identify Accounts with Frequent Changes to Sensitive Settings

```sql
select
  actor_handle,
  count(*) as setting_changes 
from
  pipes_audit_log 
where
  action_type = 'settings.update' 
  and target_handle in ('database', 'security', 'access') 
group by
  actor_handle 
having
  setting_changes > 3;
```

## Detect Access from Non-Whitelisted Locations

```sql
select
  actor_handle,
  tp_source_location 
from
  pipes_audit_log 
where
  tp_source_location not in ('trusted_location_1', 'trusted_location_2') 
group by
  actor_handle, 
  tp_source_location;
```

## Track Multiple Failed Attempts to Perform Critical Actions

```sql
select
  actor_handle,
  count(*) as failed_actions 
from
  pipes_audit_log 
where
  action_type like '%failed' 
  and target_handle = 'critical_resource' 
group by
  actor_handle 
having
  failed_actions > 3;
```

## Monitor for Access to Restricted Database Instances

```sql
select
  actor_handle,
  count(*) as restricted_access 
from
  pipes_audit_log 
where
  database_name in ('restricted_db_1', 'restricted_db_2') 
group by
  actor_handle 
having
  restricted_access > 0;
```

## Identify Sudden Increase in Actions by an Actor

```sql
with recent_actions as (
  select
    actor_handle,
    count(*) as recent_count 
  from
    pipes_audit_log 
  where
    tp_date >= current_date - interval '1 day' 
  group by
    actor_handle
),
past_actions as (
  select
    actor_handle,
    count(*) as past_count 
  from
    pipes_audit_log 
  where
    tp_date < current_date - interval '1 day' 
    and tp_date >= current_date - interval '7 day' 
  group by
    actor_handle
)
select
  r.actor_handle,
  r.recent_count,
  p.past_count 
from
  recent_actions r 
  join past_actions p on r.actor_handle = p.actor_handle 
where
  r.recent_count > p.past_count * 2;
```

## Detect Use of Outdated API or CLI Versions

```sql
select
  actor_handle,
  data->'api_version' as api_version,
  data->'cli_version' as cli_version 
from
  pipes_audit_log 
where
  data->'api_version' < '1.9.0' 
  or data->'cli_version' < '0.20.0';
```

## Identify Actions from IP Addresses Known for Suspicious Activity

```sql
select
  actor_handle,
  tp_ips 
from
  pipes_audit_log 
where
  tp_ips in ('suspicious_ip_1', 'suspicious_ip_2') 
group by
  actor_handle, 
  tp_ips;
```

## Monitor for Changes Made Outside Standard Office Hours

```sql
select
  actor_handle,
  count(*) as off_hours_actions 
from
  pipes_audit_log 
where
  hour(cast(created_at as timestamp)) < 9 
  or hour(cast(created_at as timestamp)) > 18 
group by
  actor_handle 
having
  off_hours_actions > 3;
```

## Detect Access from Multiple Countries within Short Timeframe

```sql
with locations as (
  select
    actor_handle,
    tp_source_location,
    tp_timestamp,
    lag(tp_source_location) over (partition by actor_handle order by tp_timestamp) as prev_location,
    lag(tp_timestamp) over (partition by actor_handle order by tp_timestamp) as prev_timestamp 
  from
    pipes_audit_log
)
select
  actor_handle,
  count(*) as country_changes 
from
  locations 
where
  tp_source_location != prev_location 
  and (tp_timestamp - prev_timestamp) < 3600000  -- 1 hour in milliseconds
group by
  actor_handle 
having
  country_changes > 1;
```

## Identify High Volume of Requests to a Specific API Endpoint

```sql
select
  actor_handle,
  count(*) as request_count,
  data->>'api_version' as api_version 
from
  pipes_audit_log 
where
  tp_source_type = 'pipes_audit_log_api' 
group by
  actor_handle, 
  data->>'api_version' 
having
  request_count > 100;
```
