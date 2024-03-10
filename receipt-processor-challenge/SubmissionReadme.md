# Submission

## No DOCKER NEEDED
All commands can be run in terminal, and instructions will assume that GoLang is installed onto the computer.

## <ins>Instructions to Run</ins>
1. Download the source code
2. Open the code in VSCode (or similar IDE)
3. Open a terminal within the project (or navigate to the project within a terminal)
4. In the terminal, enter: 
go run main.go 
to start the application in localhost:8080 hosting the API
5. Split another terminal, and enter:
curl localhost:8080/receipts/process --include --header "Content-Type: application/json" -d {directory of .json} --request "POST"
- This is the command to read a .json file and save a record in memory. If you wish to use another .json file, replace {directory of .json} with the name of the json desired within the repo.
6. The API will return the ID, you can then run:
 curl localhost:8080/receipts/{id}/points 
 Where {id} is replaced by the ID returned by the API
7. The API will return the amount of points that the receipt earned.