FROM golang:1.21-alpine as build
RUN apk update && apk upgrade --available --no-interactive && apk add --update --no-cache bash make nodejs npm
COPY ./frontend/package.json /app/frontend/package.json
WORKDIR /app/frontend
RUN npm i
COPY . /app
WORKDIR /app
RUN make build
# COPY /app/sqlchk .


FROM alpine
COPY --from=build /app/sqlchk /sqlchk
COPY <<EOF /root/.sqlchk.yaml
---
connstr: 'postgresql://joee:password@postgres_local/joee?sslmode=require'
EOF
EXPOSE 3030/tcp
CMD ["/sqlchk","server"]
