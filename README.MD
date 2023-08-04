## GitHub Actoins Statistics collector
Collect statistics on your GitHub actions to get mor e insights on how your processes run

## How does it work?
You deploy server which accept requests with information about your GitHub actions. You can see usage in `example.yml` workflow in this repository.

## How to use it?

- Fork this repository
- Deploy backend server
Server uses mongodb as database so you need to deploy it before server. You can either use `docker-compose.yml`, which deploys both mongodb and server. If you have mongodb deployed somewhere already you can use `Dockerfile` to deploy just server, you can specify your mongodb credentials in `.env` file. Keep in mind that in any case you need to provide acces to server to GitHub actions.
- Setup `AUTH_TOKEN` (same as in `.env`) and `GHA_URL` (public url to your server)
- Run and checkout `example.yml` workflow to see how you can communicate with your new server

## Quickstart
There's a quick way to deploy this application if you don't have any public IP address.
1) Fork this repository
2) Deploy free instance of mongodb using mongodb [atlass](https://www.mongodb.com/atlas). Don't forget to allow all connections for you mongodb database!
3) Deploy server using Dockerfile via [Back4App](https://www.back4app.com/). Don't forget to setup your environment variables with your mongodb atlass url link and `.env` file
4) In your fork repository setup secrets `AUTH_TOKEN` and `GHA_URL` (there's `token-generator` utility under `cmd` which can generate `AUTH_TOKEN` for you)
5) Run `example.yml` to see that everything is working fine
