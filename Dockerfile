FROM golang:1.13.5-alpine3.10

ADD app /build/exmonit

WORKDIR /build/exmonit
RUN go get -d .
RUN go install .

EXPOSE 8080

CMD ["exmonit", "-f/srv/config.yml"]
