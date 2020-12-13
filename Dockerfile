FROM golang:1.15-alpine as build

WORKDIR /src
ADD . .
RUN go build -ldflags="-s -w" ./cmd/userdata

FROM scratch
COPY --from=build /src/userdata /bin/userdata
ENTRYPOINT [ "/bin/userdata" ]
