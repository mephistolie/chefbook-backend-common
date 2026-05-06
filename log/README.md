# ChefBook Structured Logging

ChefBook backend logs should be JSON-only in production and structured-first in new code.

The logger backend is `zerolog`. The package keeps compatibility helpers such as `Infof`, `Warnf`, and `Errorf` during migration, but new log statements should prefer structured `Event` values.

## Environment

`environment` has only two valid values:

- `dev`
- `prod`

Do not introduce `stage`, `test`, `local`, or service-specific environment names into the log schema unless the runtime environment model changes.

## Base Schema

Every log entry should follow this shape:

```json
{
  "timestamp": "2026-05-05T12:00:00.000Z",
  "level": "info",
  "service": "auth",
  "environment": "prod",
  "component": "grpc",
  "event": "auth.session.created",
  "message": "session created",
  "trace_id": "trace-id",
  "request_id": "request-id",
  "user_id": "user-id",
  "duration_ms": 24,
  "payload": {}
}
```

Required fields:

- `timestamp`
- `level`
- `service`
- `environment`
- `event`
- `message`

Optional common fields:

- `component`: `grpc`, `http`, `amqp`, `postgres`, `s3`, `firebase`, `scheduler`, or another stable component name.
- `request_id`: incoming request correlation id.
- `trace_id`: distributed trace id when tracing is available.
- `span_id`: distributed trace span id when tracing is available.
- `user_id`: authenticated user id or domain user id.
- `message_id`: RabbitMQ message id.
- `operation`: stable operation name such as `SignIn`, `CreateRecipe`, or `GetShoppingList`.
- `duration_ms`: operation duration in milliseconds.
- `error`: error string for failed operations.
- `error_type`: stable domain or transport error type.
- `grpc_method`: full gRPC method name.
- `grpc_code`: gRPC status code.
- `http_method`: HTTP method.
- `http_path`: route pattern, not raw path with ids when a route pattern is available.
- `http_status`: HTTP response status.

## Events

`event` is a stable machine-readable event name. `message` is human-readable and can change without breaking dashboards or alerts.

Use dot-separated domain names:

- `auth.session.created`
- `auth.sign_in.failed`
- `profile.deletion.requested`
- `recipe.created`
- `recipe.picture.upload_link.generated`
- `shopping_list.purchases.added`
- `shopping_list.user.removed`
- `mq.message.received`
- `mq.message.requeued`
- `postgres.query.failed`
- `grpc.request.completed`

Prefer one event name per business outcome. Do not encode dynamic ids, usernames, emails, or error messages into `event`.

## Payload

Use `payload` for safe event-specific details that are useful for debugging but not common enough to become top-level indexed fields.

Example:

```json
{
  "level": "info",
  "service": "shopping-list",
  "environment": "prod",
  "component": "grpc",
  "event": "shopping_list.purchases.added",
  "message": "purchases added to shopping list",
  "user_id": "user-id",
  "payload": {
    "shopping_list_id": "list-id",
    "items_count": 3,
    "merged_count": 1,
    "version": 8
  }
}
```

Payload rules:

- Keep payload JSON-serializable.
- Keep payload small and shallow.
- Use snake_case keys.
- Prefer counts, booleans, enum values, ids, and safe state transitions.
- Do not put raw request bodies, raw SQL args, external API responses, or large arrays into payload.

## Sensitive Data

Never log:

- passwords or password hashes;
- access tokens or refresh tokens;
- OAuth codes;
- reset codes or activation codes;
- private keys, encrypted private keys, vault keys, recipe keys, or passphrase-derived data;
- raw Firebase, Google, Apple, or VK responses;
- raw request bodies;
- full config dumps with secrets;
- raw SQL arguments when they can contain user data.

Avoid logging these in production unless there is a documented reason:

- email addresses;
- raw IP addresses;
- user agent strings;
- external purchase tokens;
- public object upload URLs with signatures.

If a value is useful for correlation but sensitive, prefer a stable hash field such as `email_hash` or `client_ip_hash`.

## API

Initialize a service logger at process startup:

```go
log.InitWithService("auth", *cfg.LogsPath, *cfg.Environment == config.EnvDev)
```

The legacy initializer is still available for compatibility, but it writes `service: "unknown"`:

```go
log.Init(*cfg.LogsPath, *cfg.Environment == config.EnvDev)
```

New log statements should pass one `Event` value:

```go
log.Log(ctx, log.Event{
    Event: "auth.session.created",
    Message: "session created",
    Component: log.ComponentGRPC,
    UserID: userId.String(),
})

log.LogError(ctx, log.Event{
    Event: "postgres.query.failed",
    Message: "unable to get recipe",
    Component: log.ComponentPostgres,
    Payload: map[string]any{
        "recipe_id": recipeId.String(),
    },
}, err)

log.Log(ctx, log.Event{
    Event: "shopping_list.purchases.added",
    Message: "purchases added to shopping list",
    Component: log.ComponentGRPC,
    UserID: userId.String(),
    Payload: map[string]any{
        "shopping_list_id": shoppingListId.String(),
        "items_count": len(input.Purchases),
        "version": version,
    },
})
```

Compatibility helpers such as `Infof`, `Warnf`, and `Errorf` are acceptable for legacy call sites during migration, but new code should use structured fields.
