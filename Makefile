all:
	go get
	go build -o calc
docker:
	CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .
	docker build -t barais/calc .
clean:
	rm calc
	rm main
run:
	./calc
rundocker:
	docker run -p8080:8080 barais/calc
