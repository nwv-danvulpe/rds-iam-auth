FROM golang:1.15-alpine as builder

WORKDIR /workspace

COPY . .

RUN go build -o rds-iam-auth-test && \
    chmod +x rds-iam-auth-test

FROM alpine:latest

COPY --from=builder /workspace/rds-iam-auth-test /usr/bin/rds-iam-auth-test

RUN apk --no-cache add curl

ENTRYPOINT [ "/usr/bin/rds-iam-auth-test" ]