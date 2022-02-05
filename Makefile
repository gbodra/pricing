default:
	@echo "Building Go app..."
	go build
	@echo "Running test coverage report..."
	go test -v -cover ./...

dockerbuild:
	@echo "Building Dockerfile..."
	docker build -f ./Dockerfile -t docker-pricing-api .

up: dockerbuild
	@echo "Starting API Locally..."
	docker-compose up -d

logs:
	docker-compose logs -f

down:
	docker-compose down

test:
	go test -v -cover ./...

clean: down
	@echo "Cleaning Up..."
	rm -f app
	docker system prune -f
	docker volume prune -f

github:
	@echo "Committing changes to Github..."
	git add -A
	git commit -m "$m"
	git push

heroku:
	@echo "Publishing on Heroku..."
	heroku login
	git push heroku master