### Login Api Doc
---

#### 1. Request
- method: POST
- path: /api/v1/open/login
- headers:
  - Content-Type: application/json
  - Authorization: Bearer {{ .ApiKey }}

#### 2. Parameters
##### 2.1 fields

| field | type | required | default | note |
| ----  | ---- | ----     | ----    | ---- |
| page_no   | int | false | 1  | page number |
| page_size | int | false | 30 | page size   |

#### 3. Body
##### 3.1 content type
application/json

##### 3.2 fields

| field | type | required | default | note |
| ----  | ---- | ----     | ----    | ---- |
| model       | string    | false   | gpt-3.5-turbo | model name |
| messages    | []Message | true    | -             | message arrary |
| temperature | float     | false   | 1.0           | range 0.0~2.0 |
| max_tokens  | int       | false   | -             | default is unlimited |
| user        | string    | false   | -             | user identity |

##### 3.3 example
```json
{
  "model": "gpt-3.5-turbo",
  "messages": [{"role": "user", "content": "Hello!"}]
}
```

#### 4. Response
##### 4.1 content type
application/json

##### 4.2 fields

| field | type | required | example | note |
| ---- | ----  | ----     | ----    | ---- |
| id      | string   | true | gpt-3.5-turbo   | model name |
| object  | string   | true | chat.completion | chatgpt object |
| created | int64    | true | 1677649420      | created timestamp |
| model   | string   | true | gpt-3.5-turbo   | model name |
| usage   | Usage    | true | -               | usage of token |
| choices | []Choice | true | -               | response messages |

**Choice**

| field | type | required | example | note |
| ----  | ---- | ----     | ----    | ---- |
| role          | string | true | assistant | role |
| content       | string | true | -         | response message |
| finish_reason | string | true | stop      | - |
| index         | int    | true | 0         | index of choices |

##### 4.3 example
```json
{
 "id": "chatcmpl-6p9XYPYSTTRi0xEviKjjilqrWU2Ve",
 "object": "chat.completion",
 "created": 1677649420,
 "model": "gpt-3.5-turbo",
 "usage": {"prompt_tokens": 56, "completion_tokens": 31, "total_tokens": 87},
 "choices": [
   {
    "message": {
      "role": "assistant",
      "content": "The 2020 World Series was played in Arlington, Texas at the Globe Life Field, which was the new home stadium for the Texas Rangers."},
    "finish_reason": "stop",
    "index": 0
   }
  ]
}
```
