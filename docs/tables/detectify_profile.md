# Table: detectify_profile

This table contains information about web application scan profiles used by Detectify.

## Examples

### List all profiles

```sql
select
  name,  
  endpoint,
  status,
  latest_scan
from
  detectify_profile;
```

### List the profiles that are not verified

```sql
select
  name,  
  endpoint,
  status,
  latest_scan
from
  detectify_profile
where
  status != 'verified';
```

### Count the number of scans in each status

```sql
select
  count(*) as scans,
  latest_scan ->> 'status' as "status"
from
  detectify_profile
group by
  latest_scan ->> 'status';
```
