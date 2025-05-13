# use Go image
FROM golang:1.23

# set the working directory
WORKDIR /app

# copy source code
COPY . .
RUN go mod download

# build app
RUN go build -o server ./cmd/server

# set entrypoint
CMD ["./server"]
