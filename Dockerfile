FROM golang:1.18.3-alpine 
LABEL Creator=Klimushin_Kirill, Email=kirklimushin@gmail.com 

CMD mkdir /project/dir/ 
WORKDIR /project/dir/ 

RUN apk add git --no-cache 
RUN apk add build-base 

COPY ./go.mod ./ && COPY ./go.sum ./
COPY . .

RUN go mod tidy && go mod vendor 
RUN go test -v ./tests/...  

RUN go build -o ./main/main.go 
ENTRYPOINT ["go", "run", "./main/main.go"]


