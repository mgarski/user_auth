# user_auth

## Set Up & Configuration

The database used in this service is Postgres, with the table schemas contained in the tables.sql file.

Update the database connection string in the conf.sql file, which needs to be in the working directory.

## User Management

User management is performed through the /user endpoint with the following operations available.

### Create User

To create a user perform an HTTP PUT to the /user endpoint with the following JSON

`{ "email":"<EMAIL>", "password":"<PASSWORD>", "name":"<NAME>" }`

The HTTP response code will reflect the status of the operation, and a JSON document will be returned with any appropriate error messages.

Success:

`{ "code":200, "message":"" }`

Input validation failed:

`{ "code":400, "message":"Fields missing from request: [name]" }`

### Edit User

To edit a user perform an HTTP POST to the /user endpoint with the following JSON

`{ "email":"<EMAIL>", "password":"<PASSWORD>", "name":"<NAME>", "id":"<USER_ID>" }`

The HTTP response code will reflect the status of the operation, and a JSON document will be returned with any appropriate error messages.

Success:

`{ "code":200, "message":"" }`

Input validation failed:

`{ "code":400, "message":"Fields missing from request: [name]" }`

### Delete User

To delete a user perform an HTTP DELETE to the /user endpoint with the following JSON

`{ "id":"<USER_ID>" }`

The user id is present in the claims portion of the JWT returned upon successful authentication. The HTTP response code will reflect the status of the operation, and a JSON document will be returned with any appropriate error messages.

Success:

`{ "code":200, "message":"" }`

User Not Found:

`{ "code":404, "message":"User not found" }`

## Authentication

The user's password is hashed with a salt that is randomly generated when the user is created and when the user is updated.
A signed JWT is used as the authentication token with contains the user's ID in the claims to allow for other services that use the authentication service to know the identity of the user. This id is used in calls to update, delete, or log out a user.

### Log In

To authenticate a user perform an HTTP POST to the /login endpoint with the following JSON

`{ "email":"<EMAIL>", "password":"<PASSWORD>" }`

The HTTP response code will reflect the status of the operation, and a JSON document will be returned with any appropriate error messages.

Success:

`{ "token":"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6MX0.1onc50C5uVIPMTj4iNNoullLlwvxOzBerbjlesgP0ok", "code":200, "message":"OK" }`

Failure:

`{ "code":401, "message":"Authentication Failed" }`

### Log Out

To log out a user perform an HTTP POST to the /logout endpoint with the following JSON

`{ "id":"<USER_ID>" }`

The HTTP response code will reflect the status of the operation, and a JSON document will be returned with any appropriate error messages.

Success:

Empty response with HTTP 200

Failure:

`{"code":400,"message":"ID field missing from request"}`

### Validate Token

To validate a token, a downstream service needs to pass the token into the /validate endpoint.

`{ "token":"<TOKEN>" }`

The HTTP response code will reflect the status of the operation, and a JSON document will be returned with any appropriate error messages.

Success:

Empty response with HTTP 200

Failure:

Empty response with HTTP 404