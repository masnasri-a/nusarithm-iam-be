# Database Migrations

This directory contains SQL migration files for the Nusarithm IAM backend database.

## Migration Files

- `001_create_domains_table.sql` - Creates the domains table with UUID primary key
- `002_create_users_table.sql` - Creates the users table with auto-incrementing ID

## Running Migrations

### Prerequisites

- PostgreSQL database server running
- `psql` command-line tool installed
- Database created (default: `nusarithm_iam`)

### Environment Variables

You can set the following environment variables to configure the database connection:

- `DB_HOST` - Database host (default: localhost)
- `DB_PORT` - Database port (default: 5432)
- `DB_NAME` - Database name (default: nusarithm_iam)
- `DB_USER` - Database user (default: postgres)
- `DB_PASSWORD` - Database password (leave empty if not required)

### Running the Migrations

From the project root directory:

```bash
./migrations/run_migrations.sh
```

Or manually with psql:

```bash
psql "postgresql://username:password@localhost:5432/nusarithm_iam?sslmode=disable" -f migrations/001_create_domains_table.sql
psql "postgresql://username:password@localhost:5432/nusarithm_iam?sslmode=disable" -f migrations/002_create_users_table.sql
```

## Table Schemas

### domains
- `domain_id` (UUID, Primary Key)
- `name` (VARCHAR(255), NOT NULL)
- `domain` (VARCHAR(255), NOT NULL, UNIQUE)
- `created_at` (TIMESTAMP WITH TIME ZONE)
- `updated_at` (TIMESTAMP WITH TIME ZONE)

### users
- `id` (SERIAL, Primary Key)
- `username` (VARCHAR(255), NOT NULL, UNIQUE)
- `email` (VARCHAR(255), NOT NULL, UNIQUE)
- `created_at` (TIMESTAMP WITH TIME ZONE)
- `updated_at` (TIMESTAMP WITH TIME ZONE)

## Adding New Migrations

When adding new migration files:

1. Use the format: `NNN_description.sql` where NNN is a zero-padded number
2. Include creation timestamp in comments
3. Make migrations idempotent using `IF NOT EXISTS` where appropriate
4. Add appropriate indexes for performance
5. Test migrations on a copy of production data before running in production
