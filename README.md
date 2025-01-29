Build docker image:

```
docker build --tag aggchain_prover_mock .
docker run --publish 50051:50051 aggchain_prover_mock
```