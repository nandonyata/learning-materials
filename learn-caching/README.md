# What is Caching?
In computing, a cache is a high-speed data storage layer which stores a subset of data, typically transient in nature, so that future requests for that data are served up faster than is possible by accessing the data’s primary storage location. Caching allows you to efficiently reuse previously retrieved or computed data.

# How does Caching work?
The data in a cache is generally stored in fast access hardware such as RAM (Random-access memory) and may also be used in correlation with a software component. A cache's primary purpose is to increase data retrieval performance by reducing the need to access the underlying slower storage layer. Trading off capacity for speed, a cache typically stores a subset of data transiently, in contrast to databases whose data is usually complete and durable.

source - https://aws.amazon.com/caching/

---

## Redis deployment
1. Deploy redis in docker `docker run -d --name redis -p 6379:6379 redis:latest`
2. connect to the db: `docker exec -it redis redis-cli`

source - https://redis.io/docs/latest/develop/setup/

---