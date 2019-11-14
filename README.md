# rest-service-mimic
A totally mocked REST service, driven by config files

## What does this do?
This project will let you define static, mocked REST endpoints with a JSON file. The intention is to use that as part of integration tests for things that call out to external services where we can't depend on those services serving consisting data. Somewhere between unit tests and full on integration tests.

## JSON config, you say...
So, the file itself is an array of objects with the following fields:
- `path` - The path of the request to handle. These can have parameters (`orgs/{orgId}/whatever`) or regexed parameters (`orgs/{orgId:[0-9]+}/whatever`)
- `methods` - An array of the HTTP methods to handle (`POST`, `GET`, etc.)
- `headers` - Key value pairs of the headers (and their values) that need to be part of the request. The values can be regexs, if you're in to that sort of thing.
- `query_params` - Key value paris of the query parameters that need to be part of the request. These adhere to the Gorilla Mux format (`"name": "{some name}"`)
- `response` - Information about how to respond to a valid request, specifically
  - `status_code` - The HTTP response code to respond with
  - `payload` - A JSON payload to include in the body of the response. This can also include values substituted from the request payload (see section below)

## Ok, but how do I use it?
This whole thing is Dockerized so, if you want to run locally all you need to do is:
- Create the json configuration file (or just use the example one at `examples/simple_example.json`)
- Run the most recent Docker image making sure to
  - Expose the port you want to hit
  - Volume mount the directory where your config file lives

So, something like this:
`docker run -p 8001:8001 -v $(pwd)/examples:/examples rest-service-mimic:0.0.1 -config=examples/simple_example.json -port=8001`

Now try hitting the end points:
`curl -v localhost:8001/hello -H "Is-Test: yeah"`
`curl -v localhost:8001/goodbye/123`
`curl -v localhost:8001/not_defined` <- This one will return a 404 because it's not defined in the json file

## How do I build this to run as a Docker image?
The easiest way to do this is just to run `./build.sh`. This will:
- Build the Golang Docker build environment
- Volume mount the `output` directory of your local machine to the build image
- Run the build image, creating an artifact (`output/rest-service-mimc`)
- Then make a Docker image with the artifact (`docker build -t rest-service-mimic:<version number> .`)

## Wait, my response payload can contain data derived from the request payload?
Yup. The response payload definition support the Go templating language. So, to include request payload information in your response payload, just add `{{ .Request.my_field }}`. So, for example, if you have the request payload:

```
{
    "id": "abc-123-def-456"
}
```

and you need your payload to look like

```
{
    "name": "Mock response",
    "provided_id": "abc-123-def-456"
}
```

you can define your response payload using the following Go template structure:

```
{
    "name": "Mock response",
    "provided_id": "{{ .Request.id }}"
}
```

then the `id` field of the request payload will be substituted into the `provided_id` field of the response payload. Additionally, this supports the sprig template functions (https://github.com/Masterminds/sprig) so, you can do things like this:

```
{
    "name": "Mock response",
    "provided_id": "{{ default "unknown" .Request.id | upper }}"
}
```

