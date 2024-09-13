# Table: detectify_port

This table contains information about open ports and their associated assets monitored by Detectify.

## Examples

### List all policies

```sql
select
  domain_name,  
  ip_address,
  port,
  status
from
  detectify_port;
```

### List the assets that currently have SSH port (22) open

```sql
select
  domain_name,  
  ip_address,
  port,
  status
from
  detectify_port
where
  port = 22;
```

### Count number of open ports per asset (domain name)

```sql
select
  count(*) as open_ports,
  domain_name
from
  detectify_port
where
  status not in('CLOSED', 'FILTERED')
group by
  domain_name;
```
