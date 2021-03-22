# openapi-resp-convert

This program will convert anonymous responses (root schema has no ref to components) to ref with components; input is a yaml file and output as a JSON.

## Usage

```bash
openapi-resp-convert run [--in=./a.yaml] [--respcodes=200,400] [--mimetype='application/json'] [--out=a.json]
```

Help

```bash
openapi-resp-convert run -h
```

Install

```bash
go get github.com/sljeff/openapi-resp-convert
```

> How to convert JSON to yaml: https://editor.swagger.io/
