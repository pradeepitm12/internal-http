# internal-http

This is an internal http services where you have mentioned endpoints.
- **POST :/accounts** 
  - this takes a body with account id and initial amount
  - this creates an account with given id, and initial amount

  ```azure
  curl -X POST http://localhost:8080/accounts \
  -H "Content-Type: application/json" \
  -d '{"account_id": 1, "initial_balance": "1000"}'
  ```

- **GET : /accounts/{account_id}**
  - this get request gives you balance of account for a given account it
     ```azure
     curl http://localhost:8080/accounts/1                                                                                                                 
     {"account_id":1,"balance":"500"}
    ```

- **POST : /transactions**
  - this endpoint helps you to transfer amount from one account to other
  - this return with transaction id, 
    ```azure
     curl -X POST http://localhost:8080/transactions \
      -H "Content-Type: application/json" \
      -d '{"source_account_id": 1, "destination_account_id": 2, "amount": "100"}'
    {"transaction_id":101,"source_account_id":1,"destination_account_id":2,"amount":"100","created_at":"2025-06-13T20:57:16.871915184Z"}
    ```

- **POST : /deposit**
  - this endpoint helps you to deposit amount to an account
  - with response of account id and balance
  
    ```azure
      curl -X POST http://localhost:8080/deposit \
      -H "Content-Type: application/json" \
      -d '{"account_id": 2, "amount": "50"}'
      {"account_id":2,"balance":"650"}
      ```
- **POST : /withdraw**
  - this endpoint helps you to withdraw amount from an account
  - with response of account and amount
    
    ```azure
    curl -X POST http://localhost:8080/withdraw \
      -H "Content-Type: application/json" \
      -d '{"account_id": 1, "amount": "100"}'
      ```


#### Assumption: 
- I have assumed that account will be an integer always. within limit of go int.
- Used row locks `SELECT ... FOR UPDATE` to avoid any race condition.


#### ToDo
 - Middleware for RateLimiting.
 - Middleware for validations.
 - Fix lint.
 - 

#### To Run the server and postgres
- To build image of the server and run postgres. I have added a docker-compose file.
  - This brings db and service up.
  - you can now hit the curl requests mentioned above to test.
  ```azure
  docker-compose up --build
  ```
- To bring everying down 
    ```azure
    docker-compose down -v
    ```

- You can run unit test in transfer/transfer_service_test.go
  - This tried to mimic concurrent transfers and check atomicity. 
    ```azure
    docker-compose up --build
    cd internal-http
    make test
    ```