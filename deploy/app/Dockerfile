FROM golang:1.21
RUN apt update && apt upgrade -y && apt install -y curl make && curl -o- https://raw.githubusercontent.com/nvm-sh/nvm/v0.39.7/install.sh | bash
RUN export NVM_DIR="$([ -z "${XDG_CONFIG_HOME-}" ] && printf %s "${HOME}/.nvm" || printf %s "${XDG_CONFIG_HOME}/nvm")"; [ -s "$NVM_DIR/nvm.sh" ] && \. "$NVM_DIR/nvm.sh" && nvm install v20.10.0
COPY ./frontend/package.json /app/frontend/package.json
WORKDIR /app/frontend
RUN export NVM_DIR="$([ -z "${XDG_CONFIG_HOME-}" ] && printf %s "${HOME}/.nvm" || printf %s "${XDG_CONFIG_HOME}/nvm")"; [ -s "$NVM_DIR/nvm.sh" ] && \. "$NVM_DIR/nvm.sh" && [ -s "$NVM_DIR/bash_completion" ] && \. "$NVM_DIR/bash_completion" && npm i
COPY . /apps
WORKDIR /apps
RUN export NVM_DIR="$([ -z "${XDG_CONFIG_HOME-}" ] && printf %s "${HOME}/.nvm" || printf %s "${XDG_CONFIG_HOME}/nvm")"; [ -s "$NVM_DIR/nvm.sh" ] && \. "$NVM_DIR/nvm.sh" && [ -s "$NVM_DIR/bash_completion" ] && \. "$NVM_DIR/bash_completion" && make build


FROM scratch
#FROM debian
# RUN apt update && apt upgrade -y && apt install -y curl make && curl -o- https://raw.githubusercontent.com/nvm-sh/nvm/v0.39.7/install.sh | bash
COPY --from=0 /apps/sqlchk /bin/sqlchk
COPY <<EOF /root/.sqlchk.yaml
---
connstr: 'postgresql://joee:password@postgres_local/joee?sslmode=require'
EOF
EXPOSE 3030/tcp
CMD ["/bin/sqlchk","server"]
