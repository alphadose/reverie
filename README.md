# Reverie

> Backend of EzFlo

## Development

You need to have [Golang 1.14.x](https://golang.org/dl/) or higher installed

Open your favourite terminal and perform the following tasks:-

1. Cross-check your golang version.

    ```bash
    $ go version
    go version go1.14.2 darwin/amd64
    ```

2. Clone this repository.

    ```bash
    $ git clone git@github.com:alphadose/reverie.git
    ```

3. Go inside the cloned directory and list available *makefile* commands.

    ```bash
    $ cd reverie && make help

    Reverie: The dark side(backend) of EzFlo

    install   Install missing dependencies
    build     Build the project binary
    tools     Install development tools
    start     Start in development mode with hot-reload enabled
    clean     Clean build files
    fmt       Format entire codebase
    vet       Vet entire codebase
    lint      Check codebase for style mistakes
    test      Run tests
    help      Display this help
    ```

4. Setup project configuration and make changes if required. The configuration file is well-documented so you
won't have a hard time looking around.

    ```bash
    $ cp config.sample.toml config.toml
    ```

5. Start the development server.

    ```bash
    $ make start
    ```

## Contributing

If you'd like to contribute to this project, refer to the [contributing documentation](./CONTRIBUTING.md).
