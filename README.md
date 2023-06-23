# Ultimate Service

As implemented by Ardan Labs in their course [Ultimate Go: Service with Kubernetes 4.0](https://github.com/ardanlabs/service/wiki/course-outline)

Copyright 2018, 2019, 2020, 2021, Ardan Labs  
info@ardanlabs.com

## Licensing

```
Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
```

## My notes

We start to setup the go project. Bill puts a big value into "Deploy First Mentality" - we should prepare our QA & Test environments in the first week before actually getting hands on development in order to make sure we develop something prod-suitable. It should be maintanable by more than one person (dev).

### Layers

We create the project structure that would contain five folders - a.k.a. upper layers. Each layer should not contain more than five layers respectively, because no person can keep in mind more than 5 things at once. The lesser - the better:

- **app** = main layer
    - services = services we are building
    - metrics - self-explanatory
    - tooling - other applications we are building that support our functionality (e.g., admin interface)
- **business** = application business logic. The insights. The problems the project is trying to solve.
    - core
    - cview
    - data
    - sys
    - web
- **foundation** = the standard library of the project. Packages not tied to business logic of the project. Eventually those packages might live in their own repos, assigned their owned repos and land here through vendor.
- **vendor** = third party dependencies.
- **zarf**  = holds everything related to configuration, docker, K8S, build & deployment.

**Our convention** - we only use *(import)* things from top to down: **app --> business --> foundation --> vendor --> zarf**. App can utilize everything below it. Business - everything **except of app**. Foundation can't make use of business and app. 

### Module

We use `go mod init` to initialize project as go module. Every module should have a name - it acts as modules namespace. This allows us to import code from the same project. 

The common convention is to give it a path to git repo without protocol - [github.com/kostjaigin/ultimategoservice](https://github.com/kostjaigin/ultimategoservice). 

So we initialize this project with
```
go mod init github.com/kostjaigin/ultimategoservice
```

and execute `go mod tidy` in order to add module requirements and sums. I will follow Bill's example and vendor my dependencies (= keep a local copy of them) `go mod vendor`. Those two commands together build my first make flow - `make tidy`.

| Group                     | Commands                  |
|---------------------------|---------------------------|
| Environment Setup         | `dev-docker`, `dev-gotooling`, `dev-brew-common` |
| Building                  | `all`, `service` |
| Local Kubernetes Management | `dev-up`, `dev-down`, `dev-load`, `dev-apply` |
| Monitoring and Inspection | `dev-status`, `dev-logs`             |
| Local Execution           | `run-local`              |
| Dependency Management     | `tidy`                   |

We can utilize the label from kustomize configuration to query the app carrying pod logs:
`kubectl logs --namespace=$(NAMESPACE) -l app=$(APP) ... where APP = sales`

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: sales
  namespace: sales-system

spec:
  selector:
    matchLabels:
      app: sales <---

  template:
    metadata:
      labels:
        app: sales <---
```

### Software Design Learnings

- We always switch between two hats. We are engineers (maintanance, usability) and programmers (algorithms)
- We do things that are easy to understand
- Do not add complexity until abs. necessary
- Rob Pike's approach - we don't design in interfaces and abstractions, we discover them in process
- When designing an API do not return abstract, decoupled types. Return concrete types. It's callers decision to decouple returned value. 

### Go Specific Design Learnings

- We only read configuration in main.go. Nowhere else.
- Go always tries to use the whole available processing power. This can be changed by setting GOMAXPROCS system variable (we are setting it to correspond to set k8s CPU limits) 
- We should always be able to type `--help` and `--version` in our services and be able to ovveride configuration with system variables. Their [conf](https://pkg.go.dev/github.com/ardanlabs/conf/) package helps us with this.
- Services should ALWAYS work on default settings.
- pointer/value semantics: go is balanced there. General rule: type represents data --> value semantics, type represents an API --> pointer semantics. 

Here comes a part about errors as signals (and as values)
```
(G_m) = main go routine
  |
  |---log
  |---conf
  |
  |        {👤} user
  |        ☁️☁️☁️
  |         |
  |-------(G_d) = debug go routine
  |      / | | \
  |     ⚪ 🟡 🔴 🟣 G_d serves spawns different go routines for each request
  |
  |        {👤} user
  |        ☁️☁️☁️
  |         |
  |-------(G_a) = API service routine
  |      / | | \
  |     ⚪ ⚫ 🟠 🔘 G_a serves spawns different go routines for each request
  |
  | 🔄 all while G_m is waiting for a signal to shutdown
  |
 ---
  -
```
Some of G_a spawned goroutines might execute write operations. If we signal shutdown to the G_m and do not let G_a spawned go routines finish what they were doing, we get data corruption. Parent routine should not terminate before children: if some of the spawning goroutines should be terminated, "orphan" goroutines should be adapted by the main routine. 

Debug go routine is a typicall "orphan" - we don't track it's state:

```go
log.Infow("startup", "status", "debug v1 router started", "host", cfg.Web.DebugHost)

go func() {
  if err := http.ListenAndServe(cfg.Web.DebugHost, debug.StandardLibraryMux()); err != nil {
    log.Errorw("shutdown", "status", "debug v1 router closed", "host", cfg.Web.DebugHost, "ERROR", err)
  }
}()
```

but we also don't care, cause it doesn't do any writes and can't corrupt the data.

- channels in go serve one purpose and that's --> horizontal signaling. With or without data. If the word 'signal' doesn't make sense for your application case you should not use channel. There is guaranteed signalling and non-guarantied signalling. A.k.a unbuffered & buffered channels. You get your garanties (unbuffered channels) in cost of latency - if the receiver is not there, the sender has to wait.

Using channels in API's is a bad practice - how do we define who is providing/deciding on garanty? 

- Bill himself doesn't know what those initial timeout values are supposed to be! We just set same values that are not too ridicously short or long. 

With current implementation of http.HandlerFunction we run into a problem. Our implemented Test function under testgrp is basically the outerlayer of the call that doesn't return anything - it is not allowed to, because http.HandlerFunc type is define strictly:
```go
type HandlerFunc func(ResponseWriter, *Request)
```
We can't return anything. But we said that Handler is supposed to do the following steps:

  - Validate the data
	- Call into the business layer
	- Return errors
	- Handle OK response

how are we going to return errors and responses if we can't return anything?...

We want to create an onion of the inside out of function: `(Router(Logger(ErrorHandler(PanicHandler(func T)))))`

