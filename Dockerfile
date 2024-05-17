FROM golang:1.22.3-alpine AS compiler
LABEL maintaier="karoon.siri@gmail.com"
WORKDIR /app
COPY . ./
RUN CGO_ENABLED=0 go build -o /stress

FROM alpine
COPY --from=compiler /stress /usr/bin
ENTRYPOINT ["/usr/bin/stress"]
CMD []
