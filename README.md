# Library of Open Source Hardware

This repository contains the source code of the Search Engine for Open Source Hardware that was developed as part of my master thesis: _Library of Open Source Hardware: Creating a Semantic Knowledge Base and Search Engine for Open _Source Hardware_.

The code is hosted on [GitHub](https://github.com/aisbergg/losh)

- [Development](#development)
    - [Get Started](#get-started)
    - [Useful Commands](#useful-commands)
- [License](#license)
- [Author Information](#author-information)

---

## Development

### Get Started

1. **Install Requirements:** Developing the application requires the following tools to be installed on your system:
    - [Go](https://golang.org/doc/install)
    - [Node.js](https://nodejs.org/en/download/)
    - [Yarn](https://classic.yarnpkg.com/en/docs/install/)
    - [Air](https://github.com/cosmtrek/air)
    - [Docker](https://docs.docker.com/get-docker/)
    - [Docker Compose](https://docs.docker.com/compose/install/)
2. **Clone Repository:**
    ```sh
    git clone https://github.com/aisbergg/losh.git
    cd losh
    ```
3. **Install Go Dependencies:**
    ```sh
    go install
    # or via proxy¹
    GOPROXY=https://proxy.golang.org,direct go install
    ```
4. **Setup Dev Tools:** This repository comes with tools, like code generators, that are required for development. You can install those by executing the following script:
    ```sh
    ./tools/setup-dev.sh
    ```
7. **Start Database:** This project uses [Dgraph](https://dgraph.io) as a database. Use the following command to spin up a new database using Docker:
    ```sh
    ./tools/start-db.sh
    ```
8. **Upload DB Schema:** You have to upload the database schema manually since this step is not included in the application yet.
    ```sh
    ./tools/upload-db-schema.sh
    ```
9. **Discover Products And Fill Database:** Let the crawler run for a few minutes to fill the database with product data. This is just so that you have some data to experiment with.
    ```sh
    go run ./crawler/main.go discover -c ./crawler/config-dev.yml
    ```
10. **Build Web Assets:** Web assets, such as images, CSS and JS, need to be prepared/compiled first before the web application can be run. This is done using Node.js and the Webpack Mix bundler.
    ```sh
    cd web
    yarn install
    yarn run dev
    ```
11. **Run Web Interface:** Run the web application with Air, for live reloading:
    ```sh
    air -c .air.toml
    ```
12. **Test Web Interface:** Open the web interface in your browser at http://localhost:3000 and test the service out for yourself.

¹ By setting the environment variable `GOPROXY`, Go will use a module proxy to download the dependencies instead of fetching them directly from the source repository. The simple reason for that is to ensure long-term availability. So even if a source repository is vandalized or deleted, an archived copy will be downloaded from the proxy.

### Useful Commands

Add License Header to all Go Files:

```sh
addlicense -c "André Lehmann" -l apache -v -ignore tools/tool.go -ignore gqlgenc_tool.go **/*.go
```

Build and run the web application manually:

```sh
cd web
go build -trimpath -o ./build/losh-web ./main.go && ./build/losh-web config show -c config-dev.yml
```

## License

[Apache-2.0](LICENSE)

## Author Information

André Lehmann (aisberg@posteo.de)
