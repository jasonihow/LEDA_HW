FROM golang:1.23
WORKDIR /WORK

COPY go.mod go.mod
COPY go.sum go.sum 
COPY main.go main.go
COPY api_hw.go api_hw.go

RUN  echo "Before build" && sleep 2
RUN ls -al && sleep 2

RUN  go build -v
RUN  echo "After build" && sleep 2
RUN ls -al && sleep 2

ENTRYPOINT [ "./gogo" ]
