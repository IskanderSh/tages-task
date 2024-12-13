# Service for file uploading

## Information about service:
- Service is storing files on disk, foler ./data/
- For calling each function you should use gRPC
- Service is storing file's meta info in postgreSQL

## How to run?
- Clone the repository: ```git clone https://github.com/IskanderSh/tages-task.git```
- Start containers and check if they are running: ```docker-compose up -d```
- Make migrations for your postgres which are stored here: ./migrations/
- Run the service: ```make run```

## How to iterate with service:
- You could send requests from postman
- You could write your own client for sending requests, examples are in next point
- I wrote some typical client here: https://github.com/IskanderSh/file-uploader-client

## How to improve in the future:
- Add cache for fetching data faster, for example redis
- Implement migrations with goose or something like him, this will help to make migrations fast and simple