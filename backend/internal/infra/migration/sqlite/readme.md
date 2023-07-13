# Migrator Package

The Migrator package provides functionality for handling the creation and rollback of tables in a database. It allows you to easily manage migrations, ensuring that your database schema stays up-to-date with your application's requirements.

## Usage

### Importing the Package

To use the Migrator package in your project, import it as follows:

```go
import "github.com/vanillazen/stl/backend/internal/infra/migration/sqlite"
```

### Creating Migrations

Migrations are used to define the changes you want to apply to your database schema. Each migration represents a specific set of changes, such as creating new tables or modifying existing ones.

To create a migration, follow these steps:

1. Define a new migration file, following the naming convention `00000001-table_name.sql`. The first part of the file name represents the migration index, and the second part represents the name of the table being created or modified.

2. Inside the migration file, use the following format to specify the SQL statements for the migration:

   ```sql
   --UP
   CREATE TABLE table_name (
       column1 datatype CONSTRAINT,
       column2 datatype CONSTRAINT,
       ...
   );

   --DOWN
   DROP TABLE table_name;
   ```

   The `--UP` section contains the SQL statements to apply the desired changes, while the `--DOWN` section contains the SQL statements to revert those changes.

   For example, a migration file named `00000001-create-table-users.sql` (or  `00000001-CreateTableUsers.sql`) would contain the following content:

   ```sql
   --UP
   CREATE TABLE users (
       id TEXT PRIMARY KEY,
       name TEXT NOT NULL,
       email TEXT NOT NULL,
       created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
   );

   --DOWN
   DROP TABLE users;
   ```

3. Save the migration file in the designated directory: `assets/migrations/sqlite/`.

### Applying Migrations

To apply the migrations and update your database schema, you can use the `Migrate()` function provided by the Migrator package. This function executes all pending migrations in the order they were created.

```go
import (
  "github.com/vanillazen/stl/backend/internal/infra/migration/sqlite"
)

//go:embed assets/seeding/sqlite/*.sql
var fs embed.FS

// ...
mig := migrator.NewMigrator(fs, db, opts)

err := mig.Migrate() // Apply the pending seeding
if err != nil {
    // Handle error
}
```

Here, `fs` represents the embedded filesystem that contains the migration files. The migrations are stored in the `assets/migrations/sqlite/` directory within the embedded filesystem. The `//go:embed` directive is used to include these migration files in the embedded filesystem.

The migrator uses the file names to infer the migration index and name. For example, the `Up` migration from file named `00000001-users.sql` would be represented in the database as when executed.

```
id: 5034c845-4e0f-43ae-ae73-325dc91d1f37
idx: 1
name: create-table-users
creaed_at: 2023-06-29T20:46:56+02:00
```

### Rolling Back Migrations

If you need to revert a migration, the Migrator package provides two options: `Rollback()` and `RollbackAll()`.

```go
err := migratorInstance.Rollback(2) // Rollback the last 2 seeding
if err != nil {
    // Handle error
}

```

The `Rollback()` function rolls back a specific number of migrations. In the example above, the last two migrations will be reverted. If no value is provided to `Rollback()`, it assumes that you want to rollback only the last migration. If n > total number of migrations then all migrations are rolled back (similar to RollbackAll()).

To rollback all migrations, you can use the `RollbackAll()` function:

```go
err := mig.RollbackAll() // Rollback all seeding
if err != nil {
    // Handle error
}
```

The `RollbackAll()` function reverts all applied migrations, effectively rolling back the entire database schema to its initial state.
