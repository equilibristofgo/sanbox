.PHONY: run
run:
	go run .

.PHONY: docker
docker: 
	docker build -t echo-demoapp .
	docker run --rm -it -p 8080:8080 echo-demoapp:latest
