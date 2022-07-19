FROM golang:1.18.3-alpine 
LABEL Creator=Klimushin_Kirill, Email=kirklimushin@gmail.com 
RUN echo "Building Application... It is going to take some time..."

USER jenkins 
# Initializing Project Directory...
CMD mkdir /project/dir/ 
WORKDIR /project/dir/ 

# Installing Dependencies...
RUN apk add git --no-cache 
RUN apk add build-base 

# Copying existed sources..
COPY ./go.mod ./ 
COPY ./go.sum ./
COPY . .

# Installing Dependencies + creating Vendor Directory.. + Running Tests...
RUN go mod tidy && go mod vendor 
RUN go test -v ./tests/...  

# Building Application...
RUN go build -o ./main/main.go 
ENTRYPOINT ["go", "run", "./main/main.go"]


