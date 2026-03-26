## AI Secrets Validation – Flow Overview

This project exposes an HTTP API that accepts raw text, sends it to an AI model that detects secrets (including potential API keys), and returns the detected values plus some metadata.


### 0. How to run ?

Install docker command and docker-compose if you are in linux or windows or even MAC 

https://docs.docker.com/engine/install/

Then run `docker-compose up --build -d` in the root directory of project

### 1. Model: StarPII (Hugging Face)

- **Model used**: `bigcode/starpii`  
- **Role**: Given a text string, it performs Named Entity Recognition (NER) and returns a list of entities, including secrets / API-like tokens when present.  
- **Code**: the model pipeline is created in `app/ner.py` via `load_model()`.  
- **Config and auth**:
  - Model name and Hugging Face token are read from `config/config.yaml` (`huggingface` section) and/or the `HUGGINGFACEHUB_API_TOKEN` environment variable.

### 2. API server startup

1. `python main.py` is the container entrypoint.
2. `main.py`:
   - Ensures a config file exists via `ensure_default_config()`.
   - Loads config with `load_config()` (host/port and Hugging Face settings).
   - Starts Uvicorn with `app.app` on the configured host/port (default `0.0.0.0:8000`).
3. `app/app.py`:
   - Calls `load_model()` from `app.ner` at startup.
   - If the model loads, it prints a success message; otherwise it keeps running but endpoints will return errors.
   - Builds the FastAPI app and attaches routes from `app/routes.py`.

### 3. Request flow: `/detect-api-keys`

1. A client sends a `POST` request to `/detect-api-keys` with JSON body:
   - `{"text": "<any text to scan for secrets>"}`.
2. FastAPI validates the body against the `TextRequest` Pydantic model.
3. In `app/routes.py`:
   - If the model is not loaded, it returns HTTP 500.
   - Otherwise, it calls `result = model(request.text)`:
     - This passes the input text to the StarPII NER pipeline.
     - The pipeline returns a list of detected entities.
   - It extracts potential secrets/API keys with:
     - `api_keys = [item["word"] for item in result if "word" in item]`.
   - It computes metadata (`input_length`, `processing_time_seconds`, `valid_keys_count`) and returns an `ApiKeyResponse`:
     - `api_keys`: list of strings detected as secrets/API keys.
     - `total_count`: number of detected values.
     - `details`: extra info about the request and processing time.

### 4. Health checks and Docker

- **Health endpoint**: `/health` in `app/routes.py` returns:
  - `{"status": "healthy", "model_loaded": <true/false>}`.
- **Docker**:
  - `Dockerfile` builds a Python 3.11 image, installs requirements, copies `main.py` and the `app/` package, and runs `python main.py`.
  - `docker-compose.yml` maps container port `8000` to host port `9001`, sets the Hugging Face token env var, and defines a healthcheck against `http://localhost:8000/health`.

### 5. How automation should use this service

1. Wait for the container to be healthy (`docker ps` or hitting `/health`).
2. For each text (file content, log fragment, etc.) you want to scan:
   - Send a `POST` to `/detect-api-keys` with the text.
   - Read the returned `api_keys` list and `total_count`.
3. Use this information in your automation (e.g., fail a CI job, raise alerts, or block deployments when secrets are detected).

