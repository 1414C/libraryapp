# libraryapp

Libraryapp exposes a Library service and a Book service in order to simplify the test-drive of a jiffy generated application.  To make the experience as simple as possible, the services can be tested simply by loading the docker image (libraryapp.tar) contained in this repository.

Check out the tutorial in the [Jiffy Showcases](https://github.com/1414C/jiffy/showcases) section of the Jiffy documentation.

This repository has everything you need to create your own Dockerfile.  The *main* executable has been compiled for use on Alpine (musl) and can be added to whatever type of container you wish to use.  If you don't want to run the compiled application in a container, you can clone the repo to your local system and start the application as follows:

```code
go run main.go -dev
```

