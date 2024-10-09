![image](https://hub.steampipe.io/images/plugins/l-teles/detectify-social-graphic.png)

# Detectify Plugin for Steampipe

Use SQL to query your security vulnerabilities from [Detectify](https://detectify.com/)

- **[Get started →](https://hub.steampipe.io/plugins/l-teles/detectify)**
- Documentation: [Table definitions & examples](https://hub.steampipe.io/plugins/l-teles/steampipe-plugin-detectify/tables)
- Community: [Join #steampipe on Slack →](https://turbot.com/community/join)
- Get involved: [Issues](https://github.com/l-teles/steampipe-plugin-detectify/issues)

## Quick start

Install the plugin with [Steampipe](https://steampipe.io):

```sh
steampipe plugin install l-teles/detectify
```

Configure the API token in `~/.steampipe/config/detectify.spc`:

```hcl
connection "detectify" {
  plugin = "l-teles/detectify"

  # The base URL of Detectify. Required.
  # This can be set via the `DETECTIFY_URL` environment variable.
  # base_url = "https://api.detectify.com/rest"

  # The API token for API calls. Required.
  # This can also be set via the `DETECTIFY_API_TOKEN` environment variable.
  # token = "96d4y0631c31850v2g13e4rkqt50h1p8v"

  # The access secret for API calls. Required.
  # This can also be set via the `DETECTIFY_API_SECRET` environment variable.
  # secret = "zl/0kt4gvFsV43PQuhNJjZ-XSSIJKakoYY2pTax05zaY="

  # The access secret for v3 API calls. Required.
  # This can also be set via the `DETECTIFY_API_TOKEN_V3` environment variable.
  # token_v3 = "3cd16594-z302-4lgz-113e-b3a36xy2lt99"
}
```

Or through environment variables:

```sh
export DETECTIFY_URL="https://api.detectify.com/rest"
export DETECTIFY_API_TOKEN="96d4y0631c31850v2g13e4rkqt50h1p8v"
export DETECTIFY_API_SECRET="zl/0kt4gvFsV43PQuhNJjZ-XSSIJKakoYY2pTax05zaY="
export DETECTIFY_API_TOKEN_V3="3cd16594-z302-4lgz-113e-b3a36xy2lt99"
```

Run a query:

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

## Development

Prerequisites:

- [Steampipe](https://steampipe.io/downloads)
- [Golang](https://golang.org/doc/install)

Clone:

```sh
git clone https://github.com/l-teles/steampipe-plugin-detectify.git
cd steampipe-plugin-detectify
```

Build, which automatically installs the new version to your `~/.steampipe/plugins` directory:

```
make
```

Configure the plugin:

```
cp config/* ~/.steampipe/config
vi ~/.steampipe/config/detectify.spc
```

Try it!

```
steampipe query
> .inspect detectify
```

Further reading:

- [Writing plugins](https://steampipe.io/docs/develop/writing-plugins)
- [Writing your first table](https://steampipe.io/docs/develop/writing-your-first-table)

## Contributing

Please see the [contribution guidelines](https://github.com/turbot/steampipe/blob/main/CONTRIBUTING.md) and our [code of conduct](https://github.com/turbot/steampipe/blob/main/CODE_OF_CONDUCT.md). All contributions are subject to the [Apache 2.0 open source license](https://github.com/l-teles/steampipe-plugin-detectify/blob/main/LICENSE).

`help wanted` issues:

- [Steampipe](https://github.com/turbot/steampipe/labels/help%20wanted)
- [Detectify Plugin](https://github.com/l-teles/steampipe-plugin-detectify/labels/help%20wanted)
