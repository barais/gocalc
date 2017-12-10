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

rununikdaemon:
	unik daemon --debug
rununiktarget:
	unik target --host localhost
builduniimage:
	unik build --name myImage --path ./ --base rump --language go --provider virtualbox
unistartinstance:
	unik run --instanceName myInstance --imageName myImage
unistopinstance:
	unik stop --instance --instanceName
unirestartinstance:
	unik stop --instance --instanceName
