# Introduction
This backend API server is responsible for managing the wallets of the players of an online casino. It provides an API for getting and updating their account balances. 

## Dependencies

To correctly run this server, ensure the following dependencies are satisfied:

- [Go](https://go.dev/doc/install)
- [MySQL](https://www.digitalocean.com/community/tutorials/how-to-install-mysql-on-ubuntu-20-04)
- [Redis](https://redis.io/topics/quickstart)
- [Auth0](https://auth0.com/docs/quickstart/backend/golang/01-authorization#configure-auth0-apis)

## How to set up the project
1. Clone the repository
    ```bash
    serious@dev:~$ git clone git@github.com:ageeknamedslickback/wallet-API.git
    ```

2. Set up your MySQL user, password and database using your own preferred method

3. Create an Auth0 account, add a `Backend API` app and take note of it's credentials - `domain`, `audience`, `client ID` and `client secret`

4. Create `env.sh` and add the following environment variables
    ```bash
    export DB_USER=""
    export DB_PASS=""
    export DB_HOST=""
    export DB_PORT=""
    export DB_NAME=""
    export PORT=""
    export GIN_MODE=""
    export AUTH0_DOMAIN=""
    export AUTH0_AUDIENCE=""
    export AUTH0_CLIENT_ID=""
    export AUTH0_CLIENT_SECRET=""
    export AUTH0_GRANT_TYPE=""
    export REDIS_ADDR=""
    export REDIS_PASSWORD=""
    export REDIS_DB=""
    ```

5. Install Go dependencies
    ```bash
    serious@dev:~$ go mod tidy
    ```

6. Ensure your Redis service is properly configured and running

7. Run the server (performs the migrations to your database)
    ```bash
    serious@dev:~$ source env.sh
    serious@dev:~$ go run server.go
    ```

8. Pre-populate your database with a few dummy wallets
    ```bash
    serious@dev:~$ mysql -u <user> -p <password>
    mysql> INSERT INTO wallets(id,balance) VALUES(1,100);
    mysql> INSERT INTO wallets(id,balance) VALUES(2,10);
    mysql> INSERT INTO wallets(id,balance) VALUES(3,0);
    ```

9. Call the APIs on any client of your choice

## How to run the APIs
1. To run any API you have to first get an access token for authorization

    **Post:** `/access_token`

    **Response**
    ```json
    {
        "response": {
            "access_token": "",
            "token_type": "Bearer",
            "expires_in": 86400
        }
    }
    ```
2. For every request, pass the `access token` in the `Authorization Header`
    ```json
    {
        "Authorization": "Bearer <access token>"
    }
    ```

## How to run the tests

The server is covered by unit, integration and acceptance tests
```bash
serious@dev:~$ go test -v ./...
```

## API Spec

Export this collection to postman (if you are using it) to run the APIs:

[![Run in Postman](https://run.pstmn.io/button.svg)](https://app.getpostman.com/run-collection/a9495f127e246b807b17?action=collection%2Fimport)