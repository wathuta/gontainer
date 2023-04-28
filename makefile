run:
	sudo go run main.go run /bin/bash

echo:
	sudo go run main.go run echo "The container has started"
container_Id:
	ps aux | grep sleep
hello:
	sudo go run main.go echo hello