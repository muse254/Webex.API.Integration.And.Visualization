# The {host} should be your publicly available server(this server) that can be called by the WebexAPI
run: build
	docker run -it -p 3000:3000 webex_app

build:
	docker build --build-arg HOST=http://3.222.86.122/ -t webex_app .
