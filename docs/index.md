---
organization: Turbot
category: ["saas"]
icon_url: "/images/plugins/turbot/pipes.svg"
brand_color: "#FABF1B"
display_name: "Turbot Pipes"
short_name: "pipes"
description: "Tailpipe plugin for collecting and querying logs from Turbot Pipes."
og_description: "Query Turbot Pipes logs with SQL! Open source CLI. No DB required."
og_image: "/images/plugins/turbot/pipes-social-graphic.png"
---

# Turbot Pipes + Tailpipe

[Tailpipe](https://tailpipe.io) is an open-source CLI tool that allows you to collect logs and query them with SQL.

[Turbot Pipes](https://turbot.com/pipes/) is an intelligence, automation & security platform built specifically for DevOps.

The [Turbot Pipes Plugin for Tailpipe](https://hub.tailpipe.io/plugins/turbot/pipes) allows you to collect and query Pipes logs using SQL to track activity, monitor trends, detect anomalies, and more!

- Documentation: [Table definitions & examples](https://hub.tailpipe.io/plugins/turbot/pipes/tables)
- Community: [Join #tailpipe on Slack →](https://turbot.com/community/join)
- Get involved: [Issues](https://github.com/turbot/tailpipe-plugin-pipes/issues)

<img src="https://raw.githubusercontent.com/turbot/tailpipe-plugin-pipes/main/docs/images/pipes_audit_log_terminal.png" width="50%" type="thumbnail"/>

## Getting Started

Install Tailpipe from the [downloads](https://tailpipe.io/downloads) page:

```sh
# MacOS
brew install turbot/tap/tailpipe
```

```sh
# Linux or Windows (WSL)
sudo /bin/sh -c "$(curl -fsSL https://tailpipe.io/install/tailpipe.sh)"
```

Install the plugin:

```sh
tailpipe plugin install pipes
```

Configure your [connection credentials](https://hub.tailpipe.io/plugins/turbot/pipes#connection-credentials), table partition, and data source ([examples](https://hub.tailpipe.io/plugins/turbot/pipes/tables/pipes_audit_log#example-configurations)):

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

Download, enrich, and save logs from your source ([examples](https://tailpipe.io/docs/reference/cli/collect)):

```sh
tailpipe collect pipes_audit_log
```

Enter interactive query mode:

```sh
tailpipe query
```

Run a query:

```sql
select
  action_type,
  target_handle,
  count(*) as event_count
from
  pipes_audit_log
group by
  action_type,
  target_handle
order by
  event_count desc;
```

```sh
+---------------------------------------+-----------------------+-------------+
| action_type                           | target_handle         | event_count |
+---------------------------------------+-----------------------+-------------+
| workspace.snapshot.create             | pipelingsscaletesting | 343         |
| workspace.schema.create               | smoketestv11925       | 8           |
| workspace.snapshot.create             | smoketestv11925       | 7           |
| pipeline.create                       | smoketestv11925       | 5           |
| pipeline.command.run                  | smoketestv11925       | 2           |
| workspace.mod.variable.setting.create | smoketestv11925       | 2           |
| workspace.mod.install                 | smoketestv11925       | 2           |
| workspace.snapshot.create             | smoketestv11923       | 2           |
| workspace.aggregator.create           | smoketestv11925       | 1           |
| workspace.delete                      | smoketestv11922       | 1           |
| workspace.snapshot.create             | smoketestv11922       | 1           |
| workspace.create                      | smoketestv11925       | 1           |
| datatank.table.create                 | smoketestv11925       | 1           |
| datatank.create                       | smoketestv11925       | 1           |
+---------------------------------------+-----------------------+-------------+
```

## Connection Credentials

### Arguments

| Name                   | Type          | Required | Description                                                                               |
|------------------------|---------------|----------|-------------------------------------------------------------------------------------------|
| `token`                | String        | Yes      | Authentication token for accessing the Turbot Pipes API.                                  |
| `org_handle`           | String        | Yes      | Unique identifier for the Turbot Pipes organization, used to scope API requests.          |


### Turbot Pipes Credentials

You need to specify the Pipes token along with the organization handle to authenticate to the pipes environment. A connection per organization, using tokens and organization handle is probably the most common configuration:

- `token`: A unique authentication [API token]((https://turbot.com/pipes/docs/da-settings#tokens)) used to securely connect to the Turbot Pipes API. This token grants access to query resources within the specified organization.
- `org_handle`: The unique identifier (handle) for your Turbot Pipes organization. This is required to scope API requests to the correct organization and ensure access to relevant resources.

#### pipes.tpc:

```hcl
connection "pipes" "pipes_org" {
  token      = "pipes_token"
  org_handle = "pipes_org_handle"
}
```

