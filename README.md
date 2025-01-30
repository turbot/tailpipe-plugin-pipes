# Turbot Pipes Plugin for Tailpipe

[Tailpipe](https://tailpipe.io) is an open-source CLI tool that allows you to collect logs and query them with SQL.

[Turbot Pipes](https://turbot.com/pipes/) is an intelligence, automation & security platform built specifically for DevOps.

The [Turbot Pipes Plugin for Tailpipe](https://hub.tailpipe.io/plugins/turbot/pipes) allows you to collect and query Pipes logs using SQL to track activity, monitor trends, detect anomalies, and more!

- **[Get started →](https://hub.tailpipe.io/plugins/turbot/pipes)**
- Documentation: [Table definitions & examples](https://hub.tailpipe.io/plugins/turbot/pipes/tables)
- Community: [Join #tailpipe on Slack →](https://turbot.com/community/join)
- Get involved: [Issues](https://github.com/turbot/tailpipe-plugin-pipes/issues)

Collect and query logs:
![image](docs/images/pipes_audit_log_terminal.png)

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

## Developing

Prerequisites:

- [Tailpipe](https://tailpipe.io/downloads)
- [Golang](https://golang.org/doc/install)

Clone:

```sh
git clone https://github.com/turbot/tailpipe-plugin-pipes.git
cd tailpipe-plugin-pipes
```

After making your local changes, build the plugin, which automatically installs the new version to your `~/.tailpipe/plugins` directory:

```sh
make
```

Re-collect your data:

```sh
tailpipe collect pipes_audit_log
```

Try it!

```sh
tailpipe query
> .inspect pipes_audit_log
```

## Open Source & Contributing

This repository is published under the [Apache 2.0](https://www.apache.org/licenses/LICENSE-2.0) (source code) and [CC BY-NC-ND](https://creativecommons.org/licenses/by-nc-nd/2.0/) (docs) licenses. Please see our [code of conduct](https://github.com/turbot/.github/blob/main/CODE_OF_CONDUCT.md). We look forward to collaborating with you!

[Tailpipe](https://tailpipe.io) is a product produced from this open source software, exclusively by [Turbot HQ, Inc](https://turbot.com). It is distributed under our commercial terms. Others are allowed to make their own distribution of the software, but cannot use any of the Turbot trademarks, cloud services, etc. You can learn more in our [Open Source FAQ](https://turbot.com/open-source).

## Get Involved

**[Join #tailpipe on Slack →](https://turbot.com/community/join)**

Want to help but don't know where to start? Pick up one of the `help wanted` issues:

- [Tailpipe](https://github.com/turbot/tailpipe/labels/help%20wanted)
- [Pipes Plugin](https://github.com/turbot/tailpipe-plugin-pipes/labels/help%20wanted)
