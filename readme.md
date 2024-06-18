# PARTITIONED KAFKA TOPIC

- 3 producers
- 3 consumers

<img src="https://media.licdn.com/dms/image/D5612AQFkXNrpxs06MA/article-inline_image-shrink_400_744/0/1702956549332?e=1721865600&v=beta&t=PUmRBBIKd44srDJrNQmz2HOG_7a02tqZ8gayanWDWBs">
example

### docker run
```shell
    docker compose up --build
    docker exec -it kafka kafka-topics --create --topic test --bootstrap-server localhost:9092
    docker exec -it kafka kafka-topics --bootstrap-server localhost:9092 --alter --topic test --partitions 3
    docker exec -it kafka kafka-topics --describe --bootstrap-server localhost:9092 --topic test
```

### k8s run
```shell
    kubectl apply -f zookeeper-deployment.yaml
    kubectl apply -f kafka-deployment.yaml
    
    kubectl apply -f producer/k8s/ConfigMap.yaml
    kubectl apply -f producer/k8s/Deployment.yaml
    
    kubectl apply -f consumer/k8s/ConfigMap.yaml
    kubectl apply -f consumer/k8s/Deployment.yaml
```