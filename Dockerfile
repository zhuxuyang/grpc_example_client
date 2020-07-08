FROM golang:1.13.9-stretch as build
ENV GO111MODULE on
ENV GOPROXY https://goproxy.io,direct
WORKDIR /home
COPY . .
RUN GOOS=linux GOARCH=amd64 go build -o grpc_example_client.bin .
CMD ["./grpc_example_client.bin", "-conf", "/config.yaml"]