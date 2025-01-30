---
title: "Source: pipes_audit_log_api - Collect logs from Pipes audit log API"
description: "Allows users to collect logs from the Turbot Pipes audit og API."
---

# Source: pipes_audit_log_api - Collect logs from Pipes audit log API

The Pipes Audit log API provides access to audit logs for activities within Turbot Pipes. These logs help track administrative actions, security events, and system changes. The API enables users to query, monitor, and analyze audit logs for governance, compliance, and security investigations.

Using this source, you can collect, filter, and analyze logs retrieved from the Pipes audit log API to enhance visibility into operations and security monitoring within Turbot Pipes.

## Example Configurations

### Collect audit logs

Collect audit logs for a Turbot Pipes organization.

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

## Arguments

| Argument   | Required | Default                  | Description                                                                                                                 |
|------------|----------|--------------------------|-----------------------------------------------------------------------------------------------------------------------------|
| connection | No       | `connection.pipes.default` | The [Pipes connection](https://hub.tailpipe.io/plugins/turbot/pipes#connection-credentials) to use to connect to the Turbot Pipes workspace. |

 