# GETIR-API-TASK

This repository represents a simple rest api which has 2 different endpoints.

If you would like to work on this repository, you need to follow the below steps.

# HOW TO RUN THE PROJECT
```
> git clone https://github.com/ocakhasan/getir-api-task.git
> cd getir-api-task
```

Then you need to create a `.env` file with following information
```
MONGO_URI=<your-mongo-uri>
DATABASE=<database-name>
COLLECTION=<collection-name>
PORT=3000
```
Then you can run the project with 
```
> go build
> ./getir-api-task
```
or 
```
go run main.go
```

Project automatically runs on port `3000` but you can change it via `.env` file.

Project consists of 3 parts.
- [GETIR-API-TASK](#getir-api-task)
- [HOW TO RUN THE PROJECT](#how-to-run-the-project)
		- [Get data from MongoDB](#get-data-from-mongodb)
		- [Post data to InMemory](#post-data-to-inmemory)
		- [Get data from Inmemory](#get-data-from-inmemory)
		- [HOW TO TEST](#how-to-test)


### Get data from MongoDB
What you need to do is to create a `POST` request to `/records` endpoint. Also, this request requires a request body in form of 
```json
{
    "startDate" : "2016-01-02",
    "endDate"   : "2021-01-02",
    "minCount"  : 1200,
    "maxCount"  : 4200
}

```

You can create the request with following command
```bash
curl -X POST "localhost:3000/records" -H 'Content-Type: application/json' -d '{ "startDate": "2016-01-02", "endDate": "2021-01-02", "minCount": 1200, "maxCount": 4200}'
```

### Post data to InMemory
What you need to do is to create a `POST` request to `inmemory` endpoint. Also this request requires a request body in form of
```json
{
	"key": "getir-task",
	"value": "api",
}
```

You can create the request with following command
```bash
curl -X POST "localhost:3000/inmemory" -H 'Content-Type: application/json' -d '{ "key": "getir-task", "value": "api"}'
```

### Get data from Inmemory
What you need to do is to create a `POST` request to `inmemory` endpoint with a query `key`. 

You can create the request with following command
```
curl -X GET "localhost:3000/inmemory?key=getir-task" 
```

It will return 
```json
{
	"key": "getir-task",
	"value": "api",
}
```

### HOW TO TEST
You can test the project with following command
```
go test -v ./...
```