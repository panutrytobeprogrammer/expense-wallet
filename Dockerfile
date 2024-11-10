FROM cgr.dev/chainguard/go:latest as build

ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

WORKDIR /app

COPY go.mod go.sum ./

COPY . .

RUN go build -o wallet-expense .

FROM cgr.dev/chainguard/static:latest

COPY --from=build /app/wallet-expense .

EXPOSE 8080

CMD ["/wallet-expense"]