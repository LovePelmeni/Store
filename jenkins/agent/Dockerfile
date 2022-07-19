FROM ubuntu:18.04
LABEL JenkinsAgent=StoreAgent Purpose="Agent for Managing Jenkins Pipeline for this Service."

RUN apt-get update -y && apt-get upgrade -y
RUN apt install apt-transport-https ca-certificates curl software-properties-common -y
RUN curl -fsSL https://download.docker.com/linux/ubuntu/gpg | gpg --dearmor -o /usr/share/keyrings/docker-archive-keyring.gpg
RUN echo \
  "deb [arch=$(dpkg --print-architecture) signed-by=/usr/share/keyrings/docker-archive-keyring.gpg] https://download.docker.com/linux/ubuntu \
  $(lsb_release -cs) stable" | tee /etc/apt/sources.list.d/docker.list > /dev/null
RUN apt-get update 
RUN apt-get install -y docker-ce-cli
USER jenkins 





