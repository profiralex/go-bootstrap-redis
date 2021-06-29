swagger:
	swagger generate spec -m > swagger.json


swagger-serve:
	docker run --rm -p 80:8080 -e SWAGGER_JSON=/app/swagger.json -v $(shell pwd):/app swaggerapi/swagger-ui


lint:
	docker run --rm -v $(shell pwd):/app -w /app golangci/golangci-lint:v1.39.0 golangci-lint run -v

