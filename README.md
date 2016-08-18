# go-pibot
Go library for a robot on raspberry pi

# UI development

Build the container if you haven't done so already, or if any changes have occured on the server side.
```
docker build -t pibot .
```

Then launch the container, this will share the `./public` directory, thus any changes you make inside that folder will appear immediately
```
docker run -i -t --rm -p 8080:8080 -v "$PWD/public"":/opt/src/app/public -w /opt/src/app pibot /go/src/app/app
```
