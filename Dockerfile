FROM golang:1.18.3-alpine 
LABEL Creator=Klimushin_Kirill, Email=kirklimushin@gmail.com 
RUN echo "Building Application... It is going to take some time..."

# Env vars... 
ENV GOOS=linux 
ENV GOARCH=amd64 
ENV GINMODE=release 

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
RUN go mod tidy
# RUN go test -v ./tests/...  
ENTRYPOINT ["go", "run", "./main/main.go"]


