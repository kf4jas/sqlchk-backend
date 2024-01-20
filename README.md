# SQL Check 

This repo is a Golang Project that embeds a svelte app for a UI.

## Developing / Testing


## Configuration file

```
cat ~/.sqlchk.yaml
---
# in the docker setup below 
# connstr: 'postgresql://joee:password@localhost/joee?sslmode=require'
connstr: 'sqlite.db'
mq_mode: false
```

### Setting up the local repo

Naming the folder sqlchk helps with the docker nomenclature. Naming the frontend frontend helps with the deployment and build process.

```
# Testing
git clone https://github.com/kf4jas/sqlchk-backend sqlchk
cd sqlchk
git clone https://github.com/kf4jas/sqlchk-frontend frontend
npm i
cd ..
make dev

# Developers
git clone git@github.com:kf4jas/sqlchk-backend sqlchk
cd 
git clone git@github.com:kf4jas/sqlchk-frontend frontend
```

```bash
make dev
```

## Building

To create a production version of your app:

```bash
make
```

