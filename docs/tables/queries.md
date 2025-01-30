## Activity Examples

### Daily activity trends

Count events per day to identify activity trends over time.

```sql
select
  strftime(created_at, '%Y-%m-%d') AS event_date,
  count(*) AS event_count
from
  pipes_audit_log
group by
  event_date
order by
  event_date asc;
```

### Top 10 events

List the 10 most frequently recorded events.

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

### Top events by actor

Identify the most frequently performed actions by an actor.

```sql
select
  actor_handle,
  action_type,
  count(*) as event_count
from
  pipes_audit_log
group by
  actor_handle,
  action_type
order by
  event_count desc;
```

## Detection Examples

### High privilege role assignments

Detect when high-privilege roles were assigned.

```sql
select
  created_at,
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

### Unusual login attempts

Identify failed login attempts and unusual authentication failures.

```sql
select
  created_at,
  actor_handle,
  action_type,
  actor_ip
from
  pipes_audit_log
where
  action_type in ('login_failed', 'unauthorized_access')
order by
  created_at desc;
```

### High volume API calls by user

Find users generating a high volume of API calls to detect potential abuse.

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

### Activity from unapproved IP addresses

Detect actions performed from outside approved network locations.

```sql
select
  created_at,
  action_type,
  actor_handle,
  actor_ip
from
  pipes_audit_log
where
  actor_ip not in ('192.168.1.1', '203.0.113.5')
order by
  created_at desc;
```

