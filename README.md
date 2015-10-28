### GOBERRY - Economist vanilla microservice template.

# SETTING UP YOUR EDITOR FOR GOLANG

Here is a good starting point for GoSublime:

```json
"on_save": [{
      "cmd": "gs9o_open", "args": {
        "run": ["sh",
          "errcheck && go build . errors && go test -i && go test && go vet && golint ."],
    "focus_view": false
    }}],
    "autocomplete_closures": true,
    "complete_builtins": true,
```

You'll need the relevant Go tools available for this to work. Errcheck can be found [here](https://github.com/kisielk/errcheck).

#### SETUP

1. Clone this repository into your GOPATH.
2. Run ```scripts/tools.sh``` to get the ```golint``` and ```errcheck``` tools if you don't have them.
3. Drop the .git folder.
4. Rename the goberry package with `scripts/rename.sh <package>`.
5. Design your [RAML](http://raml.org) API interface.
6. Run ```raml-gen``` which will generate HTTP handlers for your service.
7. Copy the generated ```handlers_gen.go``` file to ```handlers.go```.
8. Run ```git init```, ```git add .``` and ```git commit -m "initial commit"```.
9. Run ```make build``` to build a local version of the binary in the service folder.
10. Run ```make buildx``` to build a production (Linux) version of the binary.
11. The ramlapi package will wire up your endpoints to the handlers.
12. Now build out your service.
13. Run ```source dev_env``` to configure environment.

#### MAKEFILE

The comprehensive Makefile allows you to:

1. ```make build``` - build binary.
2. ```make buildx``` - build binary for deployment (Linux).
3. ```make tools``` - get ```golint``` and ```errcheck``` tools if needed.
4. ```make clean``` - remove binaries.
5. ```make lint``` - run linter recursively.
6. ```make vet``` - run ```go vet``` recursively.
7. ```make fmt``` - run ```go fmt``` recursively.
8. ```make test``` - run ```go test``` in vervbose mode recursively.
9. ```make race``` - run ```go test``` recursively with the race flag on.
10. ```make env``` - dump your Go enviroment variables.

### 12-FACTOR GOODNESS

We are aiming to make our microservices [12 factor](http://12factor.net/)

### BUILD INFORMATION

The make builders build the binary as follows:

```go build -ldflags "-X main.buildstamp `date -u '+%Y-%m-%d_%I:%M:%S%p'` -X main.githash `git rev-parse HEAD`"```

The build date and commit hash are then made available via the /version endpont.

### ASSET BUNDLING

We use gobundle to bundle assets (for now, just api.raml) with
the binary. To bundle an updated RAML file:

* Make sure BUNDLE_ASSETS=1 is included in your ```dev_env``` file.

Run ```go get github.com/alecthomas/gobundle/gobundle```

Run ```gobundle --compress --uncompress_on_init --package=main --target=bundle.go "api.raml"```

By default this is off in goberry so you can adapt your RAML file as necessary when you create a new service and then bundle it up.

### TESTS

Run ```make test``` for boring old black and white test output.

Run ```./pride``` to get nicely colorized test output.

### SERVICE DISCOVERY

The goconsul.json file is present to hook up to a package being built to plug into [consul](https://www.consul.io)

### CODE ANALYSIS

Run scripts/pythia to run the browser-based UI built on top of the Oracle code analysis tool and godoc.

### METRICS

This codebase provides low level metrics using the built-in expvar package. Simply navigate to /debug/vars to see basic
memory allocation and stack use information.