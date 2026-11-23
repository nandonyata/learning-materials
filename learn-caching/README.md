# Redis backend examples (Go)

Redis is an in-memory data store. It is useful as a cache, index, rate limiter,
queue, counter, session store, and more. It usually complements the primary
database; it does not automatically understand SQL rows or API filters.

## Run

```bash
docker compose up -d
go run ./cmd/01-cache-aside
```

RedisInsight is available at <http://localhost:5540>. The examples use the
`example:` key prefix and can be run independently.

## 1. Does the value exist? Cache-aside

Run: `go run ./cmd/01-cache-aside`

The backend normally should not call `EXISTS` and then `GET`; that creates two
network calls and the key could change between them. Call `GET` once:

- `err == nil`: cache hit; use the value.
- `err == redis.Nil`: cache miss; query the database, then `SET` with a TTL.
- Any other error: Redis/network failure; log it and, where appropriate, fall
  back to the database rather than treating an outage as a cache miss.

This is the **cache-aside** pattern. On an update, write to the database first,
then delete the cache key. The next read rebuilds it. A TTL is still necessary
as a safety net. `EXISTS` is useful when only presence matters, but not before a
read.

## 2. Fetch a large list with pagination

Run: `go run ./cmd/02-paginated-list`

Do not cache a million products as one JSON value. Updating one item would
rewrite everything and every request would transfer and decode everything.
The example stores:

```text
example:product:42  -> product JSON
example:products    -> sorted set of product IDs, score controls order
```

`ZRANGE`/`ZREVRANGE` selects one page of IDs and `MGET` fetches those objects in
one round trip. For deep pagination, prefer a score/ID cursor over a very large
offset. Use pipelining or `MGET`, not one network request per item.

If Redis is only a cache, another valid design is to run the paginated SQL query
on a miss and cache its response under a key such as `products:page:2:size:20`.
Keep page TTLs short and consider invalidation when products change.

## 3. A user supplies filters

Run: `go run ./cmd/03-filtered-users`

The example creates secondary set indexes:

```text
country:eg -> {1, 2, 4}
plan:pro   -> {1, 3, 4}
SINTER     -> {1, 4}
```

It then uses `MGET` for the matching users. Exact-match filters work well with
sets. Sorted sets support numeric/date ranges. Geo indexes support location
queries. Redis will **not** dynamically filter arbitrary JSON like SQL unless
you maintain indexes or install/query Redis Search.

The example also demonstrates query-result caching. Filters are sorted and
hashed into a deterministic cache key, so filter order does not create duplicate
entries. In a typical backend with many optional or complex filters:

1. Build a canonical query key from filters, sort, page, tenant, and permission
   scope (never leak one user's authorized result to another).
2. `GET` that key.
3. On `redis.Nil`, query the primary database and cache the result briefly.
4. Invalidate affected query caches, use versioned key prefixes, or accept a
   bounded stale period via TTL.

Do not create every possible filter combination in advance. High-cardinality,
rarely repeated queries may not be worth caching at all.

## 4. Scan big collections safely

Run: `go run ./cmd/04-scan-large-data`

`SMEMBERS` returns an entire set and can use excessive server/client memory.
`SSCAN` returns cursor-based batches. Likewise use `SCAN`, `HSCAN`, and `ZSCAN`
for large keyspaces, hashes, and sorted sets. `COUNT` is only a batch-size hint;
continue until the returned cursor is `0` and tolerate duplicates while a
collection is changing.

Never use `KEYS *` in a production request path: it can block Redis while it
walks the whole keyspace. `SCAN` is for administrative/background iteration,
not normal application querying.

## 5. Atomic rate limiting

Run: `go run ./cmd/05-rate-limit`

The example atomically combines `INCR` and `EXPIRE` in Lua. Redis executes a Lua
script as one operation, avoiding a race that could leave a counter without a
TTL. This demonstrates counters, expiry, and server-side atomic operations.

## 6. Queue/event processing with Streams

Run: `go run ./cmd/06-stream-queue`

`XADD` appends an event. A consumer group divides work among workers, tracks
pending messages, and `XACK` acknowledges success. Production workers must also
retry/claim abandoned pending messages and send repeatedly failing events to a
dead-letter stream. Cap streams with `MAXLEN` when unlimited history is not
required.

Use Pub/Sub only for live, best-effort notifications: disconnected subscribers
miss messages. Use Streams when delivery tracking and replay matter.

## Redis structures to know

| Structure | Common backend use |
|---|---|
| String | JSON cache, session, counter (`GET`, `SET`, `INCR`) |
| Hash | Fields of one object (`HGET`, `HSET`, `HGETALL`) |
| List | Simple ordered list/queue (`LPUSH`, `BRPOP`) |
| Set | Unique membership and exact-filter intersections (`SADD`, `SINTER`) |
| Sorted set | Ranking, time ordering, ranges, pagination (`ZADD`, `ZRANGE`) |
| Stream | Durable event log and consumer groups (`XADD`, `XREADGROUP`) |
| Bitmap / HyperLogLog | Flags/cohorts and approximate unique counts |
| Geospatial | Radius/location lookup (`GEOADD`, `GEOSEARCH`) |

## Production checklist

- **Naming:** use predictable names such as `app:env:tenant:type:id`.
- **TTL:** cache keys should normally expire; add random TTL jitter to prevent
  many hot keys expiring simultaneously.
- **Stampedes:** for hot misses, use request coalescing/singleflight or a short
  lock so hundreds of requests do not query the database together.
- **Eviction:** configure `maxmemory` and an eviction policy deliberately.
  Monitor hits, misses, memory, evictions, latency, and rejected connections.
- **Large values:** keep values and collections bounded. Compression saves
  memory/bandwidth but costs CPU. Never assume Redis makes a huge payload cheap.
- **Round trips:** use `MGET` or pipelines for independent commands. A pipeline
  is not atomic; use transactions or Lua when atomicity is required.
- **Failure handling:** set short connection/read timeouts and decide whether an
  endpoint falls back to the database or fails closed (for example, security
  and strict rate-limit decisions need special care).
- **Consistency:** Redis and the database are separate systems. Define whether
  stale reads are acceptable and how writes invalidate/update keys.
- **Durability:** RDB snapshots and AOF improve recovery but are not the same as
  a relational database's guarantees. Replication improves availability, not
  correctness by itself.
- **Security:** do not expose Redis publicly. Use private networking, ACLs/TLS,
  secret management, and separate logical/physical instances where needed.
- **Hot keys:** one extremely popular key can bottleneck one Redis shard. Cache
  locally where safe or partition the workload.

Redis commands are atomic individually, but a multi-command application flow is
not automatically atomic. Always reason about concurrent requests, expiration,
partial failure, and stale indexes.
