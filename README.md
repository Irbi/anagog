### Project structure

- **client** -- test client which generates reports and send then to the server
- **api** -- HTTP server, sends received reports to NATS
- **worker** -- consumes reports from NATS and stores as files

### Project configuration

The project is configured in `docker-compose.yaml` by `ENV` variables

- **client**
  - `CLIENTS` - how many clients can run at the same time

- **worker**
  - `volumes` - mount path to host directory for storing files

### Running project

In the root directory run:

`docker-compose build`

`docker-compose up worker api client`
