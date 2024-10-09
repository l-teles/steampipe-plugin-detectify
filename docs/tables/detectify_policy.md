# Table: detectify_policy

This table contains information about policies and the assets that are breaching those policies.

> The `token_v3` argument is required to use this table, meaning you need to create an API key for v3 on Detectify.

## Examples

### List all policies

```sql
select
  policy_name,
  asset_name,
  severity,
  active,
  status
from
  detectify_policy;
```

### List the domains that are currently breaching policies

```sql
select
  policy_name,
  asset_name,
  severity,
  active,
  status
from
  detectify_policy
where
  active = 'true';
```

### Count assets by policy and status

```sql
select
  count(*) as assets,
  policy_name,
  status
from
  detectify_policy
group by
  policy_name,
  status;
```
