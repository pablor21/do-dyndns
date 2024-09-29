# DigitalOcean Dynamic Dns

App to update digitalocean DNS entries with your public IP

## Config

Please rename [config.example.yml](./src/config.example.yml) to config.yml and modify the values according to your needs.

The app will also read configuration values from environment variables (os variables and .env file)

## Run in container

Build the docker container 

```
docker build . -t dyndns
```

Run the container (please make sure you pass the config file)

```
docker run -v ./config.yml:/config.yml -v .env:/.env dydns
```

## License

MIT License
