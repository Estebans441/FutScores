FROM golang:1.23

RUN mkdir -p /home/app

COPY . /home/app

WORKDIR /home/app

EXPOSE 8080

CMD [ "go", "run", "/home/app/main.go" ]