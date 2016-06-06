## Workflow Manager

Workflow Manager is the core of ernest, it is responsible to manage all events on the service management workflow

## Build status

* master:  [![CircleCI Master](https://circleci.com/gh/r3labs/workflow-manager/tree/master.svg?style=svg&circle-token=627e89c447fe342aff9815ca146b081a37c075ad)](https://circleci.com/gh/r3labs/workflow-manager/tree/master)
* develop: [![CircleCI Develop](https://circleci.com/gh/r3labs/workflow-manager/tree/develop.svg?style=svg&circle-token=627e89c447fe342aff9815ca146b081a37c075ad)](https://circleci.com/gh/r3labs/workflow-manager/tree/develop)

## Code Example

Basically and by default Workflow Manager follows this workflow [workflow.json](workflow.json).
Additionally all messages coming in and out to the Workflow Manager are processed by the [subscriber.go](subscriber.go) and the [publisher.go](publisher.go).

In order to create new transitions you'll need to add new arcs to workflow.json, and also create the necessary subscriber and publisher message processor methods.

## Motivation

This library was conceived in order to clean and give flexibility to ernest platform, please keep this in mind if you're planning to improve it.

## Installation

This library asumes a NATS_URI environment variable is pointing to nats server. And you have the config-store service listening at the same nats.

Second, you'll need to install all dependencies
```
make deps
make install
```

## Running Tests

```
make test
```

## Contributing

Please read through our
[contributing guidelines](CONTRIBUTING.md).
Included are directions for opening issues, coding standards, and notes on
development.

Moreover, if your pull request contains patches or features, you must include
relevant unit tests.

## Versioning

For transparency into our release cycle and in striving to maintain backward
compatibility, this project is maintained under [the Semantic Versioning guidelines](http://semver.org/). 

## Copyright and License

Code and documentation copyright since 2015 r3labs.io authors.

Code released under
[the Mozilla Public License Version 2.0](LICENSE).

