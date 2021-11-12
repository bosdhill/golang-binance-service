PORT=4200

build:
	docker build --build-arg $(PORT) -t binance-server . 

run:
	docker run --rm -p $(PORT):$(PORT) binance-server:latest 

clean:
	docker rm -f binance-server:latest || true