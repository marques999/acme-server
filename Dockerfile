FROM golang:latest

WORKDIR /go/src/github.com/marques999/acme-server

COPY . .

RUN go get github.com/Masterminds/squirrel
RUN go get github.com/lib/pq
RUN go get github.com/jmoiron/sqlx
RUN go get github.com/gin-gonic/gin
RUN go get github.com/joho/godotenv
RUN go get github.com/appleboy/gin-jwt
RUN go get golang.org/x/crypto/bcrypt
RUN go get github.com/speps/go-hashids

CMD ["go", "run", "server.go"]