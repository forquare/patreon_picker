# patreon-picker

## Configuration

The configuration file is in the [YAML format](https://en.wikipedia.org/wiki/YAML).

Although not a requirement, it is advisable to use the `yamllint` command to check the syntax of the file to prevent any errors.  To do so, install the `yamllint` package, and then run the command:
```shell
yamllint /etc/picker.yaml
```

```yaml
---
credentials:
  id: "patreon_client_id"
  secret: "patreon_client_secret"
  redirect_url: "http://localhost:8080/auth"

# The levels are:
#   Trace - shows everything, will likely overwhelm you!
#   Debug - gives a good idea about what is going on
#   Info - shows noteworthy events
#   Warn - shows events that are not necessarily errors, but may be worth investigating
#   Error - shows errors that are not fatal
#   Fatal - shows fatal errors, this also crashes the service
#   Panic - similar to fatal, but worse
# It is recommended to use info or warn for normal operation
log_level: Info

session:
    cookie: "cookie"
    name: "myapp"

connection:
  port: "8080"
  address: "0.0.0.0"
```

