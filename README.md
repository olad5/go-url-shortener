# Go URL Shortener


## Description
The classic URL Shortener project. You feed it a long url and it shortens it to a short url. 

## Running the project

1. Install and start [Docker](https://docs.docker.com/compose/gettingstarted/) if you haven't already.
2. Copy the `.env` template file. Input the passwords and app secrets. 

```bash
cp .env.sample .env
```

```bash
  make run
```


## Run tests

```bash
  make test.integration
```
