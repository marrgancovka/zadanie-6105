start:
	docker-compose up -d --build
stop:
	docker-compose down

build_docker:
	docker build -t main_service .

run_docker: build_docker
	docker run -d --name main_container main_service

stop_docker:
	docker stop main_container
