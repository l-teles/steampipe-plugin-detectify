# Table: detectify_member

This table contains information about users who have access to the Detectify platform.

## Examples

### List all users / members

```sql
select
  concat(first_name, ' ', last_name) AS name,  
  email,
  authentication,
  role,
  last_login
from
  detectify_member;
```

### List the users that have never logged in

```sql
select
  concat(first_name, ' ', last_name) AS name,  
  email,
  authentication,
  role,
  last_login
from
  detectify_member
where
  last_login is null;
```

### Group users by authentication method

```sql
select
  count(*) as users,
  authentication
from
  detectify_member
group by
 authentication;
```
