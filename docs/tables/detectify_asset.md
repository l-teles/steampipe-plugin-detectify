# Table: detectify_asset

An asset in Detectify refers to any domain or subdomain that is monitored and scanned for vulnerabilities by the surface monitoring module, web application scanner or both. This table includes **only** the root domains. 

## Examples

### List all assets

```sql
select
  name,
  status,
  last_seen,
  monitored
from
  detectify_asset;
```

### List the non-monitored assets

```sql
select
  name,
  status,
  last_seen,
  monitored
from
  detectify_asset
where
  monitored = 'false';
```

### Group assets by monitored status

```sql
select
  count(*) as assets,
  monitored
from
  detectify_asset
group by
  monitored;
```
