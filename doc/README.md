# 
# Build
To build the project we need golang and nodejs with yarn installed.

First we need to download all dependencies for golang and react app.

In the root folder of your project type
```sh
go get -d ./...
```

Then change to frontend folder and type
```sh
yarn
```

Now we can build the project. There are some shell scripts in the bin directory.
To build a linux build type
```sh
bin/build_linux.sh
```

 



