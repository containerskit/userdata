FROM golang:1.15-alpine as build

WORKDIR /src
ADD . .
RUN CGO_ENABLED=0 go build -ldflags="-s -w" ./cmd/userdata
RUN apk add upx
RUN upx userdata

FROM scratch
COPY --from=build /src/userdata /bin/userdata
ENTRYPOINT [ "/bin/userdata" ]
