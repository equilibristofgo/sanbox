## C->GO

https://stackoverflow.com/questions/62977352/using-go-in-c

This compiles to go file, generating libmylib.so and libmylib.h in the current directory.

    go build -o libmylib.so -buildmode=c-shared libmylib.go

The compiles the C++ program, linking it to the shared library above:

    g++ -L. main.cpp -lmylib -o hello_program

To run the program, LD_LIBRARY_PATH needs to be set to the current directory. That would be different if program was installed and the shared library put in a sensible place.

    LD_LIBRARY_PATH=. ./hello_program



## GO->C

https://dev.to/mattn/call-go-function-from-c-function-1n3

https://go.dev/blog/cgo

    cd callback
    go build
    ./c

## Socket

## Gorm Oracle

### Package Oracle


## Notes

### How to build Oracle via Docker
https://github.com/oracle/docker-images/tree/main/OracleDatabase/SingleInstance
    dockerfiles$ ./buildContainerImage.sh -x -v 18.4.0 -o '--build-arg SLIMMING=true'

### Free docker images
https://hub.docker.com/r/gvenzl/oracle-xe

### Okteto deploy
https://www.okteto.com/blog/connecting-to-database-with-port-forwarding/

### Helm
https://www.okteto.com/docs/cloud/deploy-from-helm/
https://github.com/oracle/docker-images/tree/main/OracleDatabase/SingleInstance/helm-charts/oracle-db