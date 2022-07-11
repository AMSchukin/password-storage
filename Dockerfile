FROM golang:latest
WORKDIR /app
ADD . /app
COPY cmd/* /app
RUN go mod download
RUN go build -o /docker-ps
EXPOSE 8080
CMD [ "/docker-ps" ]