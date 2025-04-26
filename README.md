# Go Ent PostgreSQL Example

Example of a simple Go application using the Ent ORM framework with PostgreSQL.

## Running

### Starting PostgreSQL via Docker Compose

Launch PostgreSQL using Docker Compose:

```bash
docker-compose up -d
```

This will start a PostgreSQL service on port 5432 with the following settings:

- User: postgres
- Password: password
- Database: postgres

### Running the application

Start the application with the following command:

```bash
go run main.go
```

The application will perform the following actions:

1. Connect to PostgreSQL
2. Create a database schema based on Ent models
3. Create a user
4. Execute a query to retrieve users

## Project Structure

- `ent/` - generated Ent code
- `ent/schema/` - Ent schema definitions
- `main.go` - main application code
- `docker-compose.yml` - Docker Compose configuration for PostgreSQL

## Adding New Entities

To add new entities, use the command:

```bash
go run -mod=mod entgo.io/ent/cmd/ent new NameOfEntity
```

Then define fields and relationships in the schema file (`ent/schema/nameOfEntity.go`) and generate the code:

```bash
go generate ./ent
```
