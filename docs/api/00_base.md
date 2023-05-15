### API Doc
---

#### Introduction
- Environments
  - name: local, url: http://127.0.0.1:3091
- Authorization
  - ApiKey: xxxxxxxx
- Notation:
  - When the returning is normal, the HTTP status code is 200, and the code is 'OK'. In the case of an exception, data is {}.
  - json response format example:
```json
{
  "requestId": "44debe81-efa0-11ed-8c57-0242ac150002"
  "code": "ok",
  "msg": "ok",
  "data": {"OBJECT"}
}
```
