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