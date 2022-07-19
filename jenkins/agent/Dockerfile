FROM ubuntu:20.04.3 

RUN apt update && apt-get install docker-ce-cli 
RUN usermod -aG docker jenkins 
USER jenkins 