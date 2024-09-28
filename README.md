# todo-app

This is a simple ToDo application built with Go. It provides a REST API for managing ToDo tasks, allowing users to create, retrieve, update, and delete tasks. The application uses PostgreSQL as the database and integrates several key Go libraries to handle routing, database access, configuration, and logging.

## Features

- **Create ToDo**: Add new tasks to the ToDo list.
- **Retrieve ToDo**: Get a list of all tasks or retrieve specific tasks by ID.
- **Update ToDo**: Edit existing tasks.
- **Delete ToDo**: Remove tasks from the ToDo list.
- **Database Support**: PostgreSQL is used for persistent storage.
- **Logging**: Detailed application logs using Logrus.
- **Configuration Management**: Configurations are managed through Viper.
- **Routing**: HTTP routing is handled via the Chi router.

## Tech Stack

- **Language**: Go 1.23
- **Database**: PostgreSQL
- **Libraries**:
  - [Chi](https://github.com/go-chi/chi): HTTP router.
  - [PGX](https://github.com/jackc/pgx): PostgreSQL database driver.
  - [Logrus](https://github.com/sirupsen/logrus): Structured logging.
  - [Viper](https://github.com/spf13/viper): Configuration management.