############################
# STEP 1 Build App
############################
FROM golang:alpine AS build

RUN apk update && apk add --no-cache 'git'

WORKDIR /app
COPY . .

RUN go env -w GOPROXY=direct

RUN go mod download
RUN go build -o /app/go-gin

############################
# STEP 2 Build Image
############################
FROM alpine:latest

WORKDIR /app
COPY --from=build /app/go-gin .

USER user

CMD ["./go-gin"]