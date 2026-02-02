# Database Configuration Guide

## Connection Pool Settings

The database connection pool can be configured via environment variables to optimize performance for different environments.

### Environment Variables

| Variable | Description | Default | Recommended Values |
|----------|-------------|---------|-------------------|
| `DB_MAX_OPEN_CONNS` | Maximum number of open connections to the database | `25` | Dev: 10-25<br>Staging: 25-50<br>Prod: 50-100 |
| `DB_MAX_IDLE_CONNS` | Maximum number of idle connections in the pool | `10` | Should be ≤ MaxOpenConns<br>Typically 25-50% of MaxOpenConns |
| `DB_CONN_MAX_LIFETIME` | Maximum lifetime of a connection in minutes | `15` | 5-30 minutes depending on DB timeout settings |

### Configuration Examples

#### Development (.env)
```bash
DB_MAX_OPEN_CONNS=10
DB_MAX_IDLE_CONNS=5
DB_CONN_MAX_LIFETIME=15
```

#### Staging (.env)
```bash
DB_MAX_OPEN_CONNS=50
DB_MAX_IDLE_CONNS=20
DB_CONN_MAX_LIFETIME=15
```

#### Production (.env)
```bash
DB_MAX_OPEN_CONNS=100
DB_MAX_IDLE_CONNS=50
DB_CONN_MAX_LIFETIME=10
```

### Tuning Guidelines

1. **MaxOpenConns**: 
   - Too low: Request queuing, slow response times
   - Too high: Database resource exhaustion
   - Start with 25 and increase based on load testing

2. **MaxIdleConns**:
   - Too low: Frequent connection creation overhead
   - Too high: Wasted resources
   - Typically 25-50% of MaxOpenConns

3. **ConnMaxLifetime**:
   - Should be less than database's connection timeout
   - Helps prevent stale connections
   - PostgreSQL default timeout is 10 hours, so 15 minutes is safe

### Monitoring

Monitor these metrics to tune your connection pool:

- **Connection wait time**: If high, increase MaxOpenConns
- **Connection creation rate**: If high, increase MaxIdleConns
- **Idle connection count**: If always at max, decrease MaxIdleConns
- **Database connection errors**: May indicate pool exhaustion

### References

- [Go database/sql documentation](https://pkg.go.dev/database/sql#DB.SetMaxOpenConns)
- [PostgreSQL connection limits](https://www.postgresql.org/docs/current/runtime-config-connection.html)
