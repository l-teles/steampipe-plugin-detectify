![image](https://hub.steampipe.io/images/plugins/l-teles/detectify-social-graphic.png)

# Detectify Plugin for Steampipe

Use SQL to query your security vulnerabilities from [Detectify](https://detectify.com/)

- **[Get started →](https://hub.steampipe.io/plugins/l-teles/detectify)**
- Documentation: [Table definitions & examples](https://hub.steampipe.io/plugins/l-teles/steampipe-plugin-detectify/tables)
- Community: [Join #steampipe on Slack →](https://turbot.com/community/join)
- Get involved: [Issues](https://github.com/l-teles/steampipe-plugin-detectify/issues)

## Quick start

Install the plugin with [Steampipe](https://steampipe.io):

```shell
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
  # token = "abc123"

  # The access secret for API calls. Required.
  # This can also be set via the `DETECTIFY_API_SECRET` environment variable.
  # secret = "123"

  # The access secret for v3 API calls. Required.
  # This can also be set via the `DETECTIFY_API_SECRET_V3` environment variable.
  # secret = "123"
}
```

Or through environment variables:

```shell
export DETECTIFY_URL=https://api.detectify.com/rest
export DETECTIFY_API_TOKEN=abc123
export DETECTIFY_API_SECRET=123
export DETECTIFY_API_TOKEN_V3=abc123
```

Run a query:
Coming soon...