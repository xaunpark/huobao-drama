# Locally-AI API Documentation

Locally-AI exposes an **OpenAI-compatible** REST API that routes prompts to real AI provider tabs (ChatGPT, Google AI Studio, …) running inside a Chrome browser via a browser extension.

**Base URL:** `http://127.0.0.1:1338`

---

## Table of Contents

- [Authentication](#authentication)
- [Endpoints](#endpoints)
  - [Chat Completions](#post-v1chatcompletions) ⭐ Main endpoint
  - [Health Check](#get-health)
  - [List Models](#get-v1models)
  - [List Workers](#get-apiworkers)
  - [List Tasks](#get-apitasks)
  - [Get Task Detail](#get-apitasksid)
  - [List Accounts](#get-apiaccounts)
  - [Get Account / Keys](#get-apiaccountsid)
  - [Create Account](#post-apiaccounts)
  - [Create API Key](#post-apiaccountsidkeys)
  - [Delete API Key](#delete-apiaccountsidkeyskeyid)
- [Providers](#providers)
- [Error Handling](#error-handling)
- [Integration Examples](#integration-examples)

---

## Authentication

All `/v1/*` endpoints require a **Bearer token** (API key).

```
Authorization: Bearer sk-acc-xxxxxxxx
```

API keys are generated when creating an account via `POST /api/accounts`.

> **Note:** Management endpoints (`/api/*`, `/health`) do **not** require authentication.

---

## Endpoints

### `POST /v1/chat/completions`

The main endpoint for sending prompts to AI providers. **OpenAI-compatible** format with Locally-AI extensions.

#### Request

```http
POST /v1/chat/completions
Content-Type: application/json
Authorization: Bearer sk-acc-xxxxxxxx
```

```json
{
  "model": "chatgpt/auto",
  "messages": [
    {
      "role": "system",
      "content": "You are a helpful assistant."
    },
    {
      "role": "user",
      "content": "Hello, how are you?"
    }
  ],
  "response_format": null,
  "urlContext": false
}
```

#### Request Fields

| Field | Type | Required | Default | Description |
|-------|------|----------|---------|-------------|
| `model` | string | **Yes** | `"chatgpt/auto"` | Model name prefixed by provider, e.g., `"chatgpt/gpt-4o"`, `"gemini/auto"`, `"aistudio/gemini-2.5-flash"`. |
| `messages` | array | **Yes** | — | Array of message objects with `role` and `content`. |
| `messages[].role` | string | **Yes** | — | One of `"system"`, `"user"`, `"assistant"`. |
| `messages[].content` | string \| array | **Yes** | — | The message text (string), or an array of content parts for multimodal input (see [Vision / Image Attachments](#vision--image-attachments)). |
| `response_format` | object | No | `null` | Set `{"type": "json_object"}` to request JSON output (enables Structured Output on AI Studio). |
| `urlContext` | boolean | No | `false` | Enable URL context browsing (AI Studio only). |
| `provider` | string | No | — | *Legacy:* Exists for backward compatibility but using the `provider/model` format in `model` is preferred. |

#### Vision / Image Attachments

To send images along with a text prompt, use the **multimodal content format** for `messages[].content`. Instead of a plain string, provide an array of content parts.

> **Supported providers:** `gemini/*` (via UI automation paste)
>
> **Supported image formats:** JPEG, PNG, WebP, GIF
>
> **Max image size:** 10 MB per image

##### Content Part Types

| Type | Description | Fields |
|------|-------------|--------|
| `text` | Text content | `text` (string) |
| `image_url` | Image reference | `image_url.url` (string) |

##### Image URL Formats

| Format | Example | Use Case |
|--------|---------|----------|
| **Local file** | `"file:///C:/photos/cat.jpg"` | Server reads local file, encodes to base64 |
| **Data URI** | `"data:image/png;base64,iVBOR..."` | Already-encoded image (e.g., from web UI) |

##### Example: Single Image

```json
{
  "model": "gemini/auto",
  "messages": [
    {
      "role": "user",
      "content": [
        {"type": "text", "text": "What is in this image?"},
        {"type": "image_url", "image_url": {"url": "file:///C:/photos/cat.jpg"}}
      ]
    }
  ]
}
```

##### Example: Multiple Images

```json
{
  "model": "gemini/auto",
  "messages": [
    {
      "role": "user",
      "content": [
        {"type": "text", "text": "Compare these two images"},
        {"type": "image_url", "image_url": {"url": "file:///C:/photos/before.jpg"}},
        {"type": "image_url", "image_url": {"url": "file:///C:/photos/after.jpg"}}
      ]
    }
  ]
}
```

##### Example: Text-Only (Backward Compatible)

Existing text-only requests continue to work exactly as before:

```json
{
  "model": "gemini/auto",
  "messages": [
    {"role": "user", "content": "Hello, how are you?"}
  ]
}
```

#### Response

```json
{
  "id": "chatcmpl-task_a1b2c3d4e5f6",
  "object": "chat.completion",
  "created": 1773500269,
  "model": "auto",
  "choices": [
    {
      "index": 0,
      "message": {
        "role": "assistant",
        "content": "Hello! I'm doing well, thank you!"
      },
      "finish_reason": "stop"
    }
  ],
  "usage": {
    "prompt_tokens": 0,
    "completion_tokens": 8,
    "total_tokens": 8
  },
  "conversation_id": "task_a1b2c3d4e5f6"
}
```

#### Response Fields

| Field | Type | Description |
|-------|------|-------------|
| `id` | string | Unique completion ID (prefixed with `chatcmpl-`). |
| `object` | string | Always `"chat.completion"`. |
| `created` | integer | Unix timestamp of creation. |
| `model` | string | The model name used. |
| `choices` | array | Array with one choice containing the assistant's response. |
| `choices[].message.role` | string | Always `"assistant"`. |
| `choices[].message.content` | string | The AI's response text. |
| `choices[].finish_reason` | string | Always `"stop"`. |
| `usage` | object | Token usage (approximate word count). |
| `conversation_id` | string | Internal task ID, can be used to query task details. |

---

### `GET /health`

Check server status and online worker count.

```http
GET /health
```

**Response:**

```json
{
  "status": "ok",
  "time": "2026-03-14T16:00:00Z",
  "online_workers": 1
}
```

---

### `GET /v1/models`

List available models (OpenAI-compatible).

```http
GET /v1/models
```

**Response:**

```json
{
  "object": "list",
  "data": [
    { "id": "chatgpt/auto", "object": "model", "owned_by": "Locally-AI" },
    { "id": "chatgpt/gpt-4o", "object": "model", "owned_by": "Locally-AI" },
    { "id": "aistudio/auto", "object": "model", "owned_by": "Locally-AI" },
    { "id": "aistudio/gemini-2.5-flash", "object": "model", "owned_by": "Locally-AI" },
    { "id": "gemini/auto", "object": "model", "owned_by": "Locally-AI" },
    { "id": "gemini/gemini-2.0-flash", "object": "model", "owned_by": "Locally-AI" },
    { "id": "gemini/gemini-2.5-pro", "object": "model", "owned_by": "Locally-AI" }
  ]
}
```

This API endpoints ensures that any 3rd party tool (like Cursor, Chatbox, Cline) can automatically discover and list all available model configurations natively without requiring custom code.

---

### `GET /api/workers`

List all connected browser extension workers.

```http
GET /api/workers
```

**Response:**

```json
{
  "data": [
    {
      "id": "worker-uuid",
      "account_id": "acc_xxx",
      "provider": "chatgpt",
      "status": "idle",
      "last_seen": "2026-03-14T16:00:00Z"
    }
  ]
}
```

| Status | Description |
|--------|-------------|
| `idle` | Worker is ready to accept tasks |
| `busy` | Worker is currently processing a task |
| `offline` | Worker has disconnected |

---

### `GET /api/tasks`

List recent tasks (up to 100).

```http
GET /api/tasks
```

**Response:**

```json
{
  "data": [
    {
      "id": "task_xxx",
      "status": "completed",
      "account_id": "acc_xxx",
      "worker_id": "worker-uuid",
      "model": "auto",
      "updated_at": "2026-03-14T16:00:00Z"
    }
  ]
}
```

---

### `GET /api/tasks/{id}`

Get detailed information about a specific task.

```http
GET /api/tasks/task_a1b2c3d4e5f6
```

---

### `GET /api/accounts`

List all accounts.

```http
GET /api/accounts
```

---

### `POST /api/accounts`

Create a new account and generate an API key.

```http
POST /api/accounts
Content-Type: application/json
```

```json
{
  "email": "user@example.com",
  "plan": "free"
}
```

**Response (201):**

```json
{
  "id": "acc_xxx",
  "email": "user@example.com",
  "plan": "free",
  "status": "active",
  "created_at": "2026-03-15T00:00:00Z",
  "keys": [
    {
      "id": "key_xxx",
      "account_id": "acc_xxx",
      "name": "Default Key",
      "preview": "sk-acc-xxxx...",
      "key": "sk-acc-xxxxxxxxxxxxxxxx",
      "created_at": "2026-03-15T00:00:00Z"
    }
  ]
}
```

> ⚠️ The `key` property inside the `keys` array is only returned in full **once** during creation. Save it immediately.

---

### `GET /api/accounts/{id}`

Get Account details, including a list of **API Key Previews** associated with the account.

```http
GET /api/accounts/acc_xxx
```

---

### `POST /api/accounts/{id}/keys`

Generate a new API key for an existing account.

```http
POST /api/accounts/acc_xxx/keys
Content-Type: application/json
```

```json
{
  "name": "My Tool Key"
}
```

**Response (201):**

```json
{
  "id": "key_xxx",
  "account_id": "acc_xxx",
  "name": "My Tool Key",
  "preview": "sk-acc-xxxx...",
  "key": "sk-acc-xxxxxxxxxxxxxxxx",
  "created_at": "2026-03-15T00:00:00Z"
}
```

> ⚠️ The `key` is only returned here once.

---

### `DELETE /api/accounts/{id}/keys/{keyId}`

Revoke (delete) a specific API Key. Any tools using it will immediately receive `401 Unauthorized` errors.

```http
DELETE /api/accounts/acc_xxx/keys/key_xxx
```

**Response:** `204 No Content`

---

## Providers

| Provider | Value | Description | Special Features |
|----------|-------|-------------|------------------|
| **ChatGPT** | `chatgpt` | OpenAI ChatGPT web interface | Default provider |
| **Gemini** | `gemini` | Google Gemini web interface (gemini.google.com) | **Image attachments** (vision) |
| **Google AI Studio** | `aistudio` | Google AI Studio web interface | `response_format`, `urlContext` |

### Provider-Specific Features

#### ChatGPT

Standard text completion. No special features beyond the base API. Use the `chatgpt/` prefix.

```json
{
  "model": "chatgpt/auto",
  "messages": [{"role": "user", "content": "Hello"}]
}
```

#### Gemini

Google Gemini web interface. Supports **image attachments** for vision/multimodal tasks. Use the `gemini/` prefix.

**Available models:** `gemini/auto`, `gemini/gemini-2.0-flash`, `gemini/gemini-2.5-pro`

```json
{
  "model": "gemini/auto",
  "messages": [
    {
      "role": "user",
      "content": [
        {"type": "text", "text": "Describe this image in detail"},
        {"type": "image_url", "image_url": {"url": "file:///C:/photos/landscape.jpg"}}
      ]
    }
  ]
}
```

> **How it works:** The server reads the local file, encodes it to base64, and sends it via WebSocket to the Chrome extension. The extension simulates a ClipboardEvent paste on the Gemini editor to attach the image, then types the text prompt and submits.

#### Google AI Studio

Supports additional features. Use the `aistudio/` prefix.

- **Structured Output (JSON):** Set `response_format` to get JSON responses.
- **URL Context:** Enable web browsing context.

```json
{
  "model": "aistudio/gemini-2.5-flash",
  "messages": [
    {"role": "system", "content": "Reply with valid JSON only."},
    {"role": "user", "content": "List 3 colors as JSON"}
  ],
  "response_format": {"type": "json_object"},
  "urlContext": false
}
```

---

## Error Handling

All errors follow this format:

```json
{
  "error": {
    "message": "error description",
    "type": "Error Type"
  }
}
```

### Common Errors

| HTTP Status | Type | Cause | Solution |
|-------------|------|-------|----------|
| `400` | Bad Request | Invalid request body | Check your JSON payload |
| `401` | Unauthorized | Missing or invalid API key | Use a valid `Bearer` token |
| `503` | Service Unavailable | No workers online | Ensure the Chrome extension is connected with popup windows open |
| `500` | Internal Server Error | Server-side error | Check server logs |

### Timeout

Tasks have a **30-minute** timeout. If the AI doesn't respond within 30 minutes, the task fails with a timeout error.

---

## Integration Examples

### Python

```python
import requests

BASE_URL = "http://127.0.0.1:1338"
API_KEY = "sk-acc-xxxxxxxx"

def chat(prompt, model="chatgpt/auto", use_json=False, url_context=False):
    payload = {
        "model": model,
        "messages": [{"role": "user", "content": prompt}]
    }

    if use_json:
        payload["response_format"] = {"type": "json_object"}
        payload["messages"].insert(0, {
            "role": "system",
            "content": "Reply with valid JSON only. No markdown."
        })

    if url_context:
        payload["urlContext"] = True

    response = requests.post(
        f"{BASE_URL}/v1/chat/completions",
        json=payload,
        headers={"Authorization": f"Bearer {API_KEY}"},
        timeout=1820  # slightly above 30-min server timeout
    )
    response.raise_for_status()
    data = response.json()
    return data["choices"][0]["message"]["content"]


# Example: ChatGPT
print(chat("Hello!", model="chatgpt/auto"))

# Example: AI Studio with JSON output
print(chat("List 5 colors", model="aistudio/gemini-2.5-flash", use_json=True))

# Example: AI Studio with URL context
print(chat("Summarize this article", model="aistudio/auto", url_context=True))
```

#### Python — Image Attachment (Gemini Vision)

```python
import requests

BASE_URL = "http://127.0.0.1:1338"
API_KEY = "sk-acc-xxxxxxxx"

def chat_with_images(prompt, image_paths, model="gemini/auto"):
    """Send a prompt with one or more local images."""
    content_parts = [{"type": "text", "text": prompt}]
    for path in image_paths:
        # Use file:// URL — server reads the file and encodes it
        file_url = "file:///" + path.replace("\\", "/")
        content_parts.append({
            "type": "image_url",
            "image_url": {"url": file_url}
        })

    payload = {
        "model": model,
        "messages": [{"role": "user", "content": content_parts}]
    }

    response = requests.post(
        f"{BASE_URL}/v1/chat/completions",
        json=payload,
        headers={"Authorization": f"Bearer {API_KEY}"},
        timeout=1820
    )
    response.raise_for_status()
    return response.json()["choices"][0]["message"]["content"]


# Single image
print(chat_with_images(
    "What is in this image?",
    ["C:/photos/cat.jpg"]
))

# Multiple images
print(chat_with_images(
    "Compare these two screenshots",
    ["C:/photos/before.png", "C:/photos/after.png"]
))
```

### cURL

```bash
# Simple ChatGPT request
curl -X POST http://127.0.0.1:1338/v1/chat/completions \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer sk-acc-xxxxxxxx" \
  -d '{
    "model": "chatgpt/auto",
    "messages": [{"role": "user", "content": "Hello!"}]
  }'

# AI Studio with JSON mode
curl -X POST http://127.0.0.1:1338/v1/chat/completions \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer sk-acc-xxxxxxxx" \
  -d '{
    "model": "aistudio/gemini-2.5-flash",
    "messages": [
      {"role": "system", "content": "Reply with valid JSON only."},
      {"role": "user", "content": "List 3 animals as JSON"}
    ],
    "response_format": {"type": "json_object"}
  }'

# Gemini with image attachment
curl -X POST http://127.0.0.1:1338/v1/chat/completions \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer sk-acc-xxxxxxxx" \
  -d '{
    "model": "gemini/auto",
    "messages": [{
      "role": "user",
      "content": [
        {"type": "text", "text": "Describe this image"},
        {"type": "image_url", "image_url": {"url": "file:///C:/photos/image.jpg"}}
      ]
    }]
  }'
```

### JavaScript / Node.js

```javascript
const BASE_URL = "http://127.0.0.1:1338";
const API_KEY = "sk-acc-xxxxxxxx";

async function chat(prompt, { model = "chatgpt/auto", useJson = false, urlContext = false } = {}) {
  const messages = [{ role: "user", content: prompt }];

  if (useJson) {
    messages.unshift({ role: "system", content: "Reply with valid JSON only. No markdown." });
  }

  const payload = {
    model,
    messages,
    ...(useJson && { response_format: { type: "json_object" } }),
    ...(urlContext && { urlContext: true })
  };

  const res = await fetch(`${BASE_URL}/v1/chat/completions`, {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
      "Authorization": `Bearer ${API_KEY}`
    },
    body: JSON.stringify(payload)
  });

  if (!res.ok) throw new Error(`HTTP ${res.status}: ${await res.text()}`);
  const data = await res.json();
  return data.choices[0].message.content;
}

// Usage
const reply = await chat("Hello!", { model: "aistudio/gemini-2.5-flash" });
console.log(reply);
```

### OpenAI Python SDK (Compatible)

```python
from openai import OpenAI

client = OpenAI(
    base_url="http://127.0.0.1:1338/v1",
    api_key="sk-acc-xxxxxxxx"
)

# Seamlessly integrates standard OpenAI SDK with any of Locally-AI's providers
response = client.chat.completions.create(
    model="chatgpt/auto",
    messages=[{"role": "user", "content": "Hello!"}]
)

print(response.choices[0].message.content)
```

#### OpenAI Python SDK — Image Attachment

```python
from openai import OpenAI

client = OpenAI(
    base_url="http://127.0.0.1:1338/v1",
    api_key="sk-acc-xxxxxxxx"
)

# Vision request with local image
response = client.chat.completions.create(
    model="gemini/auto",
    messages=[{
        "role": "user",
        "content": [
            {"type": "text", "text": "What animal is in this photo?"},
            {"type": "image_url", "image_url": {"url": "file:///C:/photos/pet.jpg"}}
        ]
    }]
)

print(response.choices[0].message.content)
```

---

## Architecture Overview

```
┌─────────────────────┐
│   Your Application  │  ← Calls POST /v1/chat/completions
│  (Python, JS, etc)  │     Supports text + image attachments
└──────────┬──────────┘
           │ HTTP REST
           ▼
┌─────────────────────┐
│  Locally-AI Server  │  ← Go server on port 1338
│  (API + WebSocket)  │     Reads local images, encodes base64
└──────────┬──────────┘
           │ WebSocket
           ▼
┌─────────────────────┐
│  Chrome Extension   │  ← Background service worker
│  (Task Router)      │
└──────┬────┬────┬────┘
       │    │    │
       ▼    ▼    ▼
┌──────┐ ┌──────┐ ┌──────────┐
│ChatGP│ │Gemini│ │AI Studio │  ← Popup windows (always active)
│  T   │ │      │ │          │
└──────┘ └──────┘ └──────────┘
                  ↑ images via
                    paste event
```

**Flow:**
1. Your app sends an HTTP request to the server.
2. Server finds an available worker (Chrome extension) via WebSocket.
3. If images are included, server reads local files, encodes to base64, and includes them in the task payload.
4. Extension routes the task to the correct provider popup window.
5. For image tasks (Gemini): content script decodes base64 → File, pastes via ClipboardEvent, waits for thumbnail, then types text prompt.
6. Content script submits the prompt and waits for the AI response.
7. Response flows back through WebSocket → HTTP response to your app.
