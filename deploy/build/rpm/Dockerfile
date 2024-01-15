FROM rockylinux:9
#ARG VERSION
#ENV VERSION=${VERSION}
#RUN echo "this version is ${VERSION}"
WORKDIR /root
RUN dnf update -y 
RUN dnf install -y rpmdevtools rpmlint golang && dnf clean all && useradd -d /home/builder -m -c "Build User" -s /bin/bash builder
WORKDIR /home/builder
USER builder
RUN rpmdev-setuptree
RUN rpmdev-newspec example
RUN curl -o- https://raw.githubusercontent.com/nvm-sh/nvm/v0.39.7/install.sh | bash && export NVM_DIR="$HOME/.nvm" && \. "$NVM_DIR/nvm.sh" && \. "$NVM_DIR/bash_completion" && nvm install v20.11.0


export NVM_DIR="$HOME/.nvm" && \. "$NVM_DIR/nvm.sh" && \. "$NVM_DIR/bash_completion"
