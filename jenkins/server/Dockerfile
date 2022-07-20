FROM jenkins/jenkins:lts

USER root
RUN apt-get update && apt-get upgrade -y

RUN apt-get update && \
apt-get -y install apt-transport-https \
ca-certificates \
curl \
gnupg2 \
software-properties-common && \
apt-get update && \
apt-get -y install docker.io && \
apt-get -y install sudo

RUN curl -fsSL https://download.docker.com/linux/debian/gpg | apt-key add -

RUN add-apt-repository \
   "deb [arch=amd64] https://download.docker.com/linux/debian \
   $(lsb_release -cs) \
   stable"

RUN curl -L "https://github.com/docker/compose/releases/download/v2.1.1/docker-compose-$(uname -s)-$(uname -m)" -o /usr/local/bin/docker-compose
RUN chmod +x /usr/local/bin/docker-compose

RUN curl https://get.docker.com/ > dockerinstall && chmod 777 dockerinstall && ./dockerinstall

# RUN if [$(getent group docker)]; then echo 'Docker Group Already Exists, Skipping...'; else groupadd docker; fi
RUN adduser jenkins sudo 
# Adding Jenkins User to the Docker Group and Making it as an owner of the File.
RUN usermod -aG docker jenkins 
USER jenkins