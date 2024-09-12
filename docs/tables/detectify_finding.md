# Table: detectify_finding

A Detectify vulnerability (or finding) is a security flaw or weakness identified by the Detectify web application security scanner or by the Detectify surface monitoring tool. 

Detectify performs automated scans of web applications to uncover potential vulnerabilities that could be exploited by attackers. These vulnerabilities can range from common issues such as Cross-Site Scripting (XSS) and SQL Injection to more complex and less common vulnerabilities. 

Each finding typically includes detailed information about the nature of the vulnerability, its location within the application, the potential impact, and recommendations for remediation. The severity of each vulnerability is also assessed to help prioritize the necessary actions to secure the application.

## Examples

### List all Detectify vulnerabilities

```sql
select
  uuid,
  title,
  severity,
  location,
  status,
  source ->>'value' as "source",
  updated_at
from
  detectify_finding;
```

### List all Detectify vulnerabilities that have been risk accepted

```sql
select
  uuid,
  title,
  severity,
  location,
  status,
  source ->>'value' as "source",
  updated_at
from
  detectify_finding
where
  status = 'accepted_risk';
```

### List all open Detectify findings for one specific asset

```sql
select
  uuid,
  title,
  severity,
  location,
  status,
  source ->>'value' as "source",
  updated_at
from
  detectify_finding
where
  status not in ('patched', 'accepted_risk', 'false_positive')
  and location like '%example.com%';
```

### Group open findings by severity

```sql
select
  count(*) as findings,
  severity
from
  detectify_finding
where
  status not in ('patched', 'accepted_risk', 'false_positive')
group by
  severity;
```
