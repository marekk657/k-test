FROM golang:1.13 AS builder

RUN apt-get update && apt-get install -y zip
RUN wget https://github.com/protocolbuffers/protobuf/releases/download/v3.11.4/protoc-3.11.4-linux-x86_64.zip -O protoc.zip; unzip protoc.zip -d /protoc/; rm protoc.zip

RUN go get -u google.golang.org/grpc
RUN go get -u github.com/golang/protobuf/protoc-gen-go

WORKDIR /go/src/klaus
COPY . ./

ENV PATH=$PATH:/protoc/bin
# go generate steps defined in doc.go
RUN go generate  

RUN GOOS=linux go build -a -o /go/src/klaus/scoringserver /go/src/klaus/cmd/server/main.go

# Create release image
# FROM alpine:latest
# sqlit3 driver needs cgo, cgo need glibc.
# alternative: https://github.com/sgerrand/alpine-pkg-glibc
FROM frolvlad/alpine-glibc 
RUN apk --no-cache add ca-certificates
EXPOSE 8080
WORKDIR /app
COPY --from=builder /go/src/klaus/scoringserver .
COPY --from=builder /go/src/klaus/repository/database.db .
ENTRYPOINT ["scoringserver"]