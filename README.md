# gin-cognito-sample
Authentication/Authorization using Amazon Cognito via Gin http server.

## Setup

1. Copy `.env` file
    ```
    cp .env.sample .env
    ```

1. Put your cognito identify pool information
    - COGNITO_POOL_ID
    - COGNITO_CLIENT_ID
    - COGNITO_CLIENT_SECRET
    - COGNITO_REGION

1. Set AWS profiles. [ref](https://docs.aws.amazon.com/cli/latest/userguide/cli-configure-files.html)

1. Run go
    ```
    go run main.go
    ```
