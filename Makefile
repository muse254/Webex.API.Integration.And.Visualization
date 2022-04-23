# The {host} should be your publicly available server(this server) that can be called by the WebexAPI
run: build
	docker run -d -p 3000:3000 webex_app

build:
	docker build --build-arg HOST=${HOST} -t webex_app .
