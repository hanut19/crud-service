FROM golang:1.18.2-alpine
RUN mkdir /app
ADD . /app
WORKDIR /app
RUN export GO111MODULE=on
RUN go build -o main .
EXPOSE 8080
CMD ["/app/main"]