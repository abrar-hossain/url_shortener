FROM golang:alpine

WORKDIR /myapp

COPY go.* ./

COPY . .

RUN go build -o /myapp .

EXPOSE 8080

ENTRYPOINT [ "/myapp/url_shortener" ]