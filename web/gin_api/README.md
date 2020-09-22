## Web REST API built for Go training using:

- [Gin](https://github.com/gin-gonic/gin) web framework to handle HTTP requests 
- [Postgres](https://github.com/go-gorm/postgres) database with [GORM](https://github.com/go-gorm/gorm) object-relational mapping
- [JWT](https://github.com/dgrijalva/jwt-go) token-based authentication
- [Viper](https://github.com/spf13/viper) configuration management
- [Cobra](https://github.com/spf13/cobra) command line argument & flag handling
- [Logrus](https://github.com/sirupsen/logrus) for structured logging

Some tools may be a little overkill for the project dimensions, but consider this a proof of concept.

## Endpoints:
- POST ``/login``: User sign in, receives username and password. Return an access token (🔑)
- POST ``/logout`` 🔑: User sign out. Invalidates the token used on the authorization header by removing it from the database.
- GET ``/me`` 🔑: Returns the current user (using the access token)
- GET ``/users`` 🔑: Returns all the users from the database 
    - GET ``/users/{id}`` 🔑: Returns the user that corresponds to the specified id 
    - PUT ``/users/{id}`` 🔑: Updates an existing user
    - DELETE ``/users/{id}`` 🔑: Deletes and existing user
- GET ``/tasks`` 🔑: Returns all the tasks from the database 
    - GET ``/tasks/{id}`` 🔑: Returns the task that corresponds to the specified id
    - PUT ``/tasks/{id}`` 🔑: Updates an existing task
    - DELETE ``/tasks/{id}`` 🔑: Deletes and existing task
    - PUT ``/tasks/{id}/toggle`` 🔑: Toggles the "completed" field of an existing task
