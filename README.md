# docapi

`docapi` is a tool to generate API documentation from comments in source code.

> **Warning**
> While it is the tool that I use in [Vertex](https://vertex.quentinguidee.dev), docapi is still in development. The tool is still subject to breaking changes before the v1.0 release.

## Usage

- Build

    ```bash
    go build -o docapi ./cmd
    ```

- Run

    ```bash
    ./docapi <path-to-project-source-code>
    ```

## Document the API

`docapi` uses comments in source code to generate the API documentation. The comments must be written in a specific format.

### Meta

You can add meta information to the API documentation by writing a comment in the following format:

```go
// docapi title The API Title
// docapi description Your API description.
// docapi version 0.0.0
```

The comment can be placed anywhere in the code.

### Types

Types are automatically documented. You don't need to write any comment for them.

### Status codes

You can declare status code one time and use them in multiple handlers.

```go
// docapi code 200 Success.
// docapi code 400 {YourErrorType} Bad request.
```

The optional `{YourErrorType}` allows you to specify the type of the error.

### Routes

To declare a route, you need to write a comment in the following format:

```go
// docapi route /your/path your_unique_identifier
```

The comment can be placed anywhere in the code, but I recommend to place it next to the route declaration.

### URLs

You can declare URL one time and use them in multiple routes via aliases. In the example below, v is the alias.

```go
// docapi url v http://{ip}:{port}/api
// docapi urlvar v ip localhost The IP address of the server.
// docapi urlvar v port 6130 The port of the server.
```

Then, you can use the alias in a route:

```go
// docapi:v route /your/path your_unique_identifier
```

### Handlers

To declare a handler, you need to write a comment in the following format:

```go
// docapi begin your_unique_identifier
// docapi method POST
// docapi summary Your handler summary
// docapi tags your-group
// docapi body {YourHandlerBodyStruct} Your handler body description.
// docapi query the-param-name {TheParamType} The param description.
// docapi response 200 {YourResponseType} The response description.
// docapi response 400
// docapi response 500
// docapi end
```

Again, the comment can be placed anywhere in the code, but I recommend to place it next to the handler declaration.

## License

`docapi` is released under the MIT License. See [LICENSE.md](./LICENSE.md).
