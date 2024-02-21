# Require
Table plus

Docker Desktop

Visual Studio Code
# Set up
Run the below commend in VS code terminal
```Make
    Make postgres
    
# Then you can see a Postgres image in your Docker Desktop
    Make createdb
    
# Now you have a db call simple_bank, you can view it in table plus
    Make migrateup
    
# Now you can see you tables created in your simple_bank db
```
# Run the application
Run the below commend in VS code terminal
```go
    go run main.go
```
# Run the unit test
You can go to Testing tab in left hand side of VS Code, then you can run the by clicking Run test button

# Step for integrating GraphQL
Install GraphQl library's command
```go
    go get github.com/graphql-go/graphql
```



