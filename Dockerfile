FROM golang:1.12-alpine AS builder
WORKDIR /go/src/github.com/donaldguy/artsrv
COPY . .
RUN go install .

FROM alpine
COPY --from=builder /go/bin/artsrv /bin/artsrv
VOLUME /data
EXPOSE 4444
CMD ["/bin/artsrv"]

