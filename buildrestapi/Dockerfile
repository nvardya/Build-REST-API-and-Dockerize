#Filename: Dockerfile
FROM golang:latest
RUN mkdir /app
WORKDIR /app
COPY . .
RUN go get -d -v ./...
RUN go install -v ./...
RUN go build -o /build
EXPOSE 10000
CMD [ "/build" ]
