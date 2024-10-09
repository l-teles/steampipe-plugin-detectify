# Table: detectify_technology

This table contains information about the technologies detected for each asset monitored by Detectify.

> The `token_v3` argument is required to use this table, meaning you need to create an API key for v3 on Detectify.

## Examples

### List all technologies

```sql
select
  domain_name,
  service_protocol,
  port,
  name,
  version,
  categories
from
  detectify_technology;
```

### List the currently active technologies

```sql
select
  domain_name,
  service_protocol,
  port,
  name,
  version,
  categories
from
  detectify_technology
where
  active = 'true';
```

### Count the number of assets using each technology

```sql
select
  count(*) as assets,
  name
from
  detectify_technology
group by
  name;
```

### Count the number of assets using each protocol

```sql
select
  count(*) as assets,
  service_protocol
from
  detectify_technology
group by
  service_protocol;
```
