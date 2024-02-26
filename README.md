# GoToll

Microservices to calculate toll cost basend on kms made by trucks / OBU.

![image](https://github.com/DenisBytes/GoToll/assets/130691305/7a962108-4a5b-4657-8a0e-bf928af75e33)

## Prerequisites


### Kafka

```
docker-compose up -d
```

### Install compiler for Protobuf

```
sudo snap install protobuf --classic
```


Install GRPC and protobuf plugins

```
go install google.golang.org/protobuf/cmd/protoc-gen-go
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc
```


Install GRPC and protobuf plugins

```
go get google.golang.org/protobuf
go get google.golang.org/grpc/
```


To compile the proto file

```
protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative types/ptypes.proto
```

### Prometheus


```
docker run --name prometheus -d -p 127.0.0.1:9090:9090 -v /local/path/to/GoToll/config:/etc/prometheus prom/prometheus --config.file=/etc/prometheus/prometheus.yml
```

### Grafana


```
docker run -d -p 3000:3000 --name=grafana grafana/grafana-oss
```

