## Activity Examples

### Daily Activity Trends

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

```
folder: Organization
```

### Top 10 Events

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

```
folder: Organization
```

### Top Events by Actor

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

```
folder: Organization
```

## Detection Examples

### Unusual Workspace Deletions

Detect when multiple workspaces are deleted in a short period.

```sql
select
  created_at,
  actor_handle,
  target_handle,
  action_type
from
  pipes_audit_log
where
  action_type = 'workspace.delete'
group by
  created_at, actor_handle, target_handle, action_type
having
  count(*) > 2
order by
  created_at desc;
```

```
folder: Workspace
```

### High Privilege Role Changes

Identify when members of an organization or tenant are updated or removed.

```sql
select
  created_at,
  actor_handle,
  target_handle,
  action_type
from
  pipes_audit_log
where
  action_type in ('org.member.update', 'org.member.delete', 'tenant.member.update', 'tenant.member.delete')
order by
  created_at desc;
```

```
folder: Member
```

### Unauthorized Token Activity

Detect unusual token creation, updates, or deletions.

```sql
select
  created_at,
  actor_handle,
  target_handle,
  action_type
from
  pipes_audit_log
where
  action_type in ('token.create', 'token.update', 'token.delete')
order by
  created_at desc;
```

```
folder: Token
```

### Organization Subscription Cancellations

Monitor if organization or user subscriptions are being canceled.

```sql
select
  created_at,
  actor_handle,
  target_handle,
  action_type
from
  pipes_audit_log
where
  action_type in ('org.subscription.canceled', 'user.subscription.canceled')
order by
  created_at desc;
```

```
folder: Subscription
```

### Workspace Schema Changes

Track modifications in workspace schemas.

```sql
select
  created_at,
  actor_handle,
  target_handle,
  action_type
from
  pipes_audit_log
where
  action_type in ('workspace.schema.create', 'workspace.schema.update', 'workspace.schema.delete', 'workspace.schema.attach', 'workspace.schema.detach')
order by
  created_at desc;
```

```
folder: Schema
```

### High Volume API Calls by User

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

```
folder: Organization
```

### Activity from Unapproved IP Addresses

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

```
folder: Organization
```
