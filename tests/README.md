# Test data

The `generator` folder contains a script, `generate.py`, which writes ~500,000 records to a parquet file. To run the sample queries in [docs/tables](../docs/tables), cd into `generator` and run:

```bash
$ python generate.py
Generated 499982 records
```

The output is `pipes_audit_log.parquet`.

In `generator`, run DuckDB.

```bash
$ duckdb
v1.1.3 19864453f7
Enter ".help" for usage hints.
Connected to a transient in-memory database.
Use ".open FILENAME" to reopen on a persistent database.
D CREATE VIEW pipes_audit_log AS SELECT * FROM read_parquet('pipes_audit_log.parquet');
```

You can copy queries from the table docs and paste them here.

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

```
┌──────────────┬────────────────────┐
│ actor_handle │ off_hours_activity │
│   varchar    │       int64        │
├──────────────┼────────────────────┤
│ dzhou        │              21899 │
│ jsmith       │             124174 │
│ awhite       │              21787 │
│ vhadianto    │             123745 │
└──────────────┴────────────────────┘
```
