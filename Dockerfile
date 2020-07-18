# A hello world exmple with Go
FROM golang:1.12

ENV GO111MODULE=on

RUN mkdir /leetcode-scheduler
ADD . /leetcode-scheduler

COPY  go.mod /leetcode-scheduler
COPY go.sum /leetcode-scheduler
RUN go mod download

WORKDIR /leetcode-scheduler
RUN go build -o main .
CMD ["/leetcode-scheduler/main"]