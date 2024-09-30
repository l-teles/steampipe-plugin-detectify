---
organization: l-teles
category: ["security"]
icon_url: "/images/plugins/l-teles/detectify.svg"
brand_color: "#ECA87D"
display_name: "Detectify"
short_name: "detectify"
description: "Steampipe plugin to query vulnerabilities, assets, policies, IP addresses, technologies and members from Detectify."
og_description: "Query Detectify with SQL! Open source CLI. No DB required."
og_image: "/images/plugins/l-teles/detectify-social-graphic.png"
---

# Detectify + Steampipe

[Detectify](https://detectify.com/) offers surface monitoring and web app scanning to identify vulnerabilities using a continuously updated database of security tests, powered by a crowdsource community of ethical hackers.

[Steampipe](https://steampipe.io) is an open source CLI to instantly query cloud APIs using SQL.

Query your open security vulnerabilities and filter by status:

```sql
  select
    to_char(created_at, 'YYYY-MM-DD HH24:MI:SS') as "Creation Date",
    date_part('day', NOW() - created_at) AS "Days Open",
    status as "Status",
    cvss_scores -> 'cvss_3_1' ->> 'severity' as "Severity",
    host as "Asset",
    title as "Title",
    case
      when source ->> 'value' = 'surface-monitoring' then 'EASM'
      else 'WebApp Scan'
    end as "Source",
    location as "URL",
    definition ->> 'description' as "Description"
from
  detectify_finding
where
  status not in ('accepted_risk','patched','false_positive');
```

```
+--------+-------------+---------------------+-----------------------------------------------+-------------+----------------------------------+
| Status | Severity    | Asset               | Title                                         | Source      | URL                              |
+--------+-------------+---------------------+-----------------------------------------------+-------------+----------------------------------+
| active | medium      | gateway.example.com | Express Stack Trace                           | EASM        | https://gateway.example.com/%ff  |
| active | information | customer.example.com| Deprecated Security Header / X-XSS-Protection | WebApp Scan | https://customer.example.com/    |
+--------+-------------+---------------------+-----------------------------------------------+-------------+----------------------------------+
```

## Documentation

**[Table definitions & examples →](/plugins/l-teles/steampipe-plugin-detectify/tables)**

## Get started

### Install

Download and install the latest Detectify plugin:

```bash
steampipe plugin install l-teles/detectify
```

### Configuration

Installing the latest Detectify plugin will create a config file (`~/.steampipe/config/detectify.spc`) with a single connection named `detectify`:

```hcl
connection "detectify" {
  plugin = "l-teles/detectify"

  # The base URL of Detectify. Required.
  # This can be set via the `DETECTIFY_URL` environment variable.
  # base_url = "https://api.detectify.com/rest"

  # The API token for API calls. Required.
  # This can also be set via the `DETECTIFY_API_TOKEN` environment variable.
  # token = "abc123"

  # The access secret for API calls. Required.
  # This can also be set via the `DETECTIFY_API_SECRET` environment variable.
  # secret = "j1/1a2b3c="

  # The access secret for v3 API calls. Required.
  # This can also be set via the `DETECTIFY_API_TOKEN_V3` environment variable.
  # tokenv3 = "abc123"
}
```

- `token` - Required access token from Detectify - v2 of the API
- `secret` - Required secret token from Detectify - v2 of the API. This needs to be enabled manually on Detectify after the key is created. (more info [here](<https://support.detectify.com/support/solutions/articles/48001061878-how-to-create-and-manage-api-keys#:~:text=You%20can%20also%20enable%20if%20a%20message%20signature%20(based%20on%20secret%20key)%20should%20be%20required.>))
- `tokenv3` - Required access token from Detectify - v3 of the API

> ℹ️ Currently, one token per API version is required, since both API versions make different information available.

Alternatively, you can also use environment variables to obtain credentials only if other arguments (base_url, token and tokenv3) are not specified in the connection:

```
export DETECTIFY_URL="https://api.detectify.com/rest"
export DETECTIFY_API_TOKEN="45f86adc2nv54efd76151530rr629fc8953c2a111111fd74fa7d361d70e55759"
export DETECTIFY_API_SECRET="j1/1a2b3c="
export DETECTIFY_API_TOKEN_V3="abc123"
```

## Get involved

- Open source: https://github.com/l-teles/steampipe-plugin-detectify
- Community: [Join #steampipe on Slack →](https://turbot.com/community/join)
