# Kafka message formats

- Topic keys:
  - request key: `chat:<id>` (example: `chat:123`)
  - response key: mirrors request key

## Request (topic: `${KAFKA_INPUT_TOPIC}`)

```json
{
  "chat_id": 123,
  "query": "string"
}
```

Fields:
- `chat_id` number (int64) — destination chat/user id
- `query` string — schedule query (free-form text)

## Response (topic: `${KAFKA_OUTPUT_TOPIC}`)

```json
{
  "chat_id": 123,
  "text": "Schedule (stub) for query: <query>",
  "ts": 1700000000
}
```

Fields:
- `chat_id` number (int64) — mirrors request chat id
- `text` string — rendered schedule text (stub content for now)
- `ts` number (unix seconds) — server-side timestamp

> Note: This is a stub. Replace with your real contract (e.g., JSON Schema/Avro/Protobuf) as needed.
