# Simple Chat
Simple Chat service using Hexagonal Architecture with Go

---


## How to run 

- Install [Golang](https://go.dev/)
- Install [tflocal](https://github.com/localstack/terraform-local)
- Install [Docker](https://www.docker.com/products/docker-desktop/) or [Podman](https://podman.io/)
- Install Docker [Compose](https://github.com/docker/compose)
- Install [AWS CLI](https://aws.amazon.com/cli/)
- Install [GolangCI Linter](https://golangci-lint.run/usage/install/)
- Install [gomock](https://thedevelopercafe.com/articles/mocking-interfaces-in-go-with-gomock-670b1a640b00) 


## Use Makefile targets

Create containers to run the service with docker compose

`make infra-up`

Run local infra setup - executes terraform code against [LocalStack](https://github.com/localstack/localstack)

`make infra-local` #

Executes tests and creates coverage report against local infra

`make cover`

--- 
# Architecture 


- Handler: http with WebSocket
- Repository: DynamoDB
- Notifier: RabbitMQ
- Moderator: Ignore (just passthrough the message, no moderation done)


```mermaid
sequenceDiagram 
    title: Connect Flow
    Client -->> WebSocket: Connect(uID, roomID)
    WebSocket -->> ChatService: Register(uID, roomID)
    ChatService -->> Notifier: Register(uID, roomID, receiver)
    Notifier ->> ChatService: Registered 
    ChatService --> WebSocket: Registered
    WebSocket -> Client: Connected 

```


```mermaid
sequenceDiagram 
    title: Send Message Flow
    Client -->> WebSocket: Send(message, roomID)
    WebSocket -->> ChatService: Send(message, roomID)    
    ChatService -->> Moderator: Check(message)
    Moderator ->> ChatService: ModeratedMessage
    ChatService -->> Repository: Store(ModeratedMessage)
    Repository ->> ChatService: ModeratedMessageStored
    ChatService -->> Notifier: Broadcast(ModeratedMessage, GroupId)
    Notifier ->> ChatService: Sent OK
    ChatService ->> WebSocket: Sent
    WebSocket ->> Client: Sent
```

--- 


# TODO
 
- [ ] Implement a Presence component that embeds the notification system
- [ ] Improve tests that cover errors
- [x] Integrate a mock framework (gomock)
- [ ] Separated UnitTest from Integration Tests (with local infra)
- [x] Integrate Golang Linters
- [ ] Implement WebSocket Handler 
- [ ] Add CICD with Github Actions
- [ ] Add Configuration (Consul, Vault)
- [ ] Add logging 
- [ ] Add metricts 
- [ ] Add rate limiting
- [ ] Add Load Tests Local
- [ ] Add Load Tests Cloud 
- [ ] Add Grafana, Prometheus to Docker Compose
- [ ] Add deployment on ECS 
- [ ] Add deployment on K8S
- [ ] Add deployment on Nomad

# FIX

- [ ] - Remove dependency of AWS CLI by running in a container
