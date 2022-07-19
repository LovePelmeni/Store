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

RUN sudo service docker start 

RUN ls -la /var/run/docker.sock 
VOLUME /var/run/docker.sock

RUN adduser jenkins sudo
RUN echo “jenkins ALL=NOPASSWD: ALL” >> /etc/sudoers

RUN usermod -aG docker jenkins
RUN chmod 777 /var/run/docker.sock
RUN chown root:jenkins /var/run/docker.sock

USER jenkins