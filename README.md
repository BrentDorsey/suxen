# Suxen

Suxen provides a clean way to search docker images that are hosted on Nexus.

**ALERT** this application is not production ready. There is a bug in Nexus itself (https://issues.sonatype.org/browse/NEXUS-15277) that causes high loads on a setup that has many resources. It is possible to re-architect the tool to work around this bug, as there hasn't been much progress from the Nexus team in resolving it yet.


## Getting Started

### Prerequisites

Most importantly, a Nexus container repository should be running.
If running in kubernetes, the best option is to use the Helm chart to install/configure Nexus: https://hub.helm.sh/charts/stable/sonatype-nexus

### Installing

The docker image (https://quay.io/repository/travelaudience/suxen) can be run on its own. Or Suxen can be installed using the helm chart.

### Helm

The helm chart configuration is kept in this repository. To perform an install, try:
```
helm install charts/suxen
```
Values that might be helpful to override would be:
```
envvar:
  nexusAddress:
    value: https://your-nexus.example.com
  nexusRegistryAddress:
    value: your-containers.example.com

ingress:
  enabled: true
  host: suxen.example.com
```
but more options can be found in the `charts/suxen/values.yaml` file.


## Developing

### setting up app to run locally

**go/backend**:

This is the most straight-forward. The `Makefile` contains most of the commands you'll need. If you run `make build`, the binary will be in the `bin/` folder, and can be run with:
```
./bin/suxend[_osx] -log.level debug -nexus.svc.address https://nexus.example.com -nexus.registry.address containers.example.com
```
Or try the `--help` to find which other parameters to configure.


**js/frontend**:

For this to work, the best option is to run the docker image with the backend in localhost, and then in the `ui/` dir,
```
yarn start
```
to have the code loaded.
You will also need to change the `uri` value in `ui/src/app.jsx` to: `http://localhost:8080/query`

To run the docker image, create an env variable file:
```
SUXEND_NEXUS_SVC_ADDRESS=https://nexus.example.com
SUXEND_NEXUS_REGISTRY_ADDRESS=containers.example.com
```
```
docker run --rm --env-file=[PATH TO ENV]/env-file -p 8080:8080 -it quay.io/travelaudience/suxen:v0.0.1
```

## Contributing

Contributions are welcomed! Read the [Contributing Guide](.github/CONTRIBUTING.md) for more information.

Things that would be the most interesting would be:

  - extend the ability to use ANY docker registry (not just Nexus)
  - new UI/UX improvements
  - create a backend datastore that can include more data about images (ie security scans)
    - maybe integrate with: https://grafeas.io/
  - with Nexus, extend ability to get other assets (jars/pip/etc)

## Licensing

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details
