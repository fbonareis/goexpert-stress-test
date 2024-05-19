# Go Expert - Stress Test

## Build image

### Generate the image
```shell
docker build -t fbonareis/goexpert-stress-test:latest . 
```
### Run application
```shell
docker run --rm fbonareis/goexpert-stress-test:latest --url=https://google.com --requests=100 --concurrency=10
```
