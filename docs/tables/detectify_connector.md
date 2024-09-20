# Table: detectify_connector

This table contains information about connectors and their configurations in Detectify.

## Examples

### List all connectors

```sql
select
  id,
  name,
  team_token,
  provider,
  created_at,
  updated_at
from
  detectify_connector;
```

### List connectors with their last run status

```sql
select
  id,
  name,
  last_run ->> 'status' as "last_run_status",
  last_run ->> 'error' as "last_run_error",
  last_run ->> 'completed_at' as "last_run_completed_at"
from
  detectify_connector;
```

### Count connectors by provider

```sql
select
  count(*) as connectors,
  provider
from
  detectify_connector
group by
  provider;
```
