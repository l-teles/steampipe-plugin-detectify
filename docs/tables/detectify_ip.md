# Table: detectify_ip

The `detectify_ip` table contains information about IP addresses monitored and scanned for vulnerabilities by Detectify.

## Examples

### List all IP Addresses

```sql
select
  ip_address,
  active,
  domain_name,
  geolocation ->> 'country_name' as "country"
from
  detectify_ip;
```

### List the currently active IP addresses

```sql
select
  ip_address,
  active,
  domain_name,
  geolocation ->> 'country_name' as "country"
from
  detectify_ip;
where
  active = 'true';
```

### Group IP addresses by country

```sql
select
  count(*) as ip_count,
  geolocation ->> 'country_name' as "country"
from
  detectify_ip
group by
  "country";
```
