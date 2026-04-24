# llama-run

[日本語版 README はこちら / Japanese README](./README_ja.md)

`llama-run` is a wrapper around [llama.cpp](https://github.com/ggml-org/llama.cpp)'s `llama-server` that lets you describe its configuration declaratively in YAML.

Instead of assembling a long list of command-line flags, you write a structured YAML file. `llama-run` reads the file, resolves models (downloading them on demand), generates any auxiliary files that `llama-server` needs, and then replaces itself with the `llama-server` process via `syscall.Exec(2)`. The wrapper does not linger: the PID is handed over to `llama-server`, which makes the tool well suited for running as PID 1 in a container or as a systemd service.

## Features

- **Declarative YAML configuration** — express every `llama-server` flag through a structured config file.
- **Process replacement via `syscall.Exec`** — no wrapper process is left behind; `llama-server` inherits the PID directly, which plays nicely with Docker, Kubernetes, and systemd.
- **Automatic model download** — fetch models from Hugging Face or S3 (including S3-compatible storage). Hugging Face downloads are verified against their published SHA-256.
- **Multi-model routing** — assign aliases to multiple models and let `llama-server`'s `--models-preset` feature switch between them. The `preset.ini` file is generated automatically.
- **cgroup-aware CPU thread detection** — when `threads: Auto` is set, `/sys/fs/cgroup/cpu.max` is consulted so the CPU share assigned to a container is respected.

## Installation

### Docker image

A Docker image that bundles `llama-server` itself is planned. Mount your config file and run:

```sh
docker run --rm -p 8080:8080 \
  -v $(pwd)/config.yaml:/etc/llama-run/config.yaml \
  -v llama-run-data:/var/lib/llama-server \
  ghcr.io/sileader/llama-run:latest
```

### Binary release

Prebuilt binaries will be distributed via GitHub Releases. Install `llama-server` separately and point `llamaServer.executable` in `config.yaml` at its path.

### Build from source

```sh
go build -o llama-run .
```

Go 1.26 or later is required.

## Usage

```sh
llama-run --config /etc/llama-run/config.yaml
```

If `--config` is omitted, `/etc/llama-run/config.yaml` is loaded by default.

## Configuration

The following example shows every top-level section. Any field that is not set falls back to its default.

```yaml
llamaServer:
  executable: /usr/bin/llama-server
  arguments: []
  directory:
    model: /var/lib/llama-server/model
    config: /var/lib/llama-server/config

downloader:
  s3:
    region: us-east-1
    endpoint: null           # set only when using S3-compatible storage
    accessKeyEnv: null       # env var to read the access key from (default: AWS_ACCESS_KEY_ID)
    secretKeyEnv: null       # env var to read the secret key from (default: AWS_SECRET_ACCESS_KEY)

server:
  host: 0.0.0.0
  port: 8080
  reusePort: false
  apiPrefix: null
  staticPath: null
  apiKey: []
  apiKeyFile: null
  tls: null                  # { certFile: ..., keyFile: ... }

features:
  embedding:
    enabled: false
    pooling: null            # None | Mean | Cls | Last | Rank
  rerank:
    enabled: false
  webui:
    enabled: false
    config: null
    configFile: null
  metrics:
    enabled: false
  properties:
    enabled: false
  jinja:
    enabled: true

log:
  enabled: true
  file: null
  level: Info                # Debug | Info | Warn | Error | Generic
  timestamp: true
  color: Auto                # Auto | On | Off

chat:
  template: null
  templateFile: null
  templateArguments: {}

reasoning:
  mode: Auto                 # Auto | On | Off
  format: None               # None | Deepseek | DeepseekLegacy
  budget: Unrestricted       # Unrestricted | Immediate | <number>
  budgetMessage: ""

device:
  cpu:
    threads: Auto            # Auto | <number>
  memory:
    mmap: true
  gpu:
    layers: Auto             # Auto | All | <number>
    mainIndex: 0

model:
  alias: null
  aliases: []
  docker: null               # passed to llama-server's --docker-repo
  huggingFace: null          # passed to llama-server's --hf-repo
  router:
    default:
      context: null
      gpuLayers: null
    models:
      - alias: gemma4-e4b
        huggingFace: ggml-org/gemma-4-E4B-it-GGUF:Q4_K_M

sampling:
  samplers: null
  seed: Random               # Random | <number>
  temperature: null
  topK: null
  topP: null
  minP: null
  repeatLastN: 64            # <number> | Disabled | Context
  repeatPenalty: Disabled    # Disabled | <number>
  presencePenalty: Disabled  # Disabled | <number>
  frequencyPenalty: Disabled # Disabled | <number>
```

### `llamaServer`

Points at the `llama-server` binary to exec into and the directories used for data. `directory.model` is where downloaded models are stored; `directory.config` is where the generated `preset.ini` is written.

### `model.router`

This is how you run multiple models behind a single `llama-server`. For every entry under `models`:

1. The model is downloaded from `huggingFace` or `s3` into `directory.model/<alias>.gguf`.
2. A `preset.ini` that maps each alias to its model path is written to `directory.config/preset.ini`.
3. `--models-preset <path>` is appended to the `llama-server` command line.

Each entry can override its own `context`. The `default` block sets router-wide defaults (`context`, `gpuLayers`).

## Model sources

### Hugging Face

Use the `<org>/<repo>:<quantization>` form, e.g. `ggml-org/gemma-4-E4B-it-GGUF:Q4_K_M`. For gated or private repositories, set the `HF_TOKEN` environment variable. Downloads are verified against the SHA-256 advertised by the Hugging Face API before being moved into place.

### S3 / S3-compatible storage

Use the `s3://<bucket>/<key>` form.

- **AWS S3** — leave `downloader.s3.endpoint` as `null` to use the default `aws-sdk-go-v2` credential chain (environment, shared config, IAM roles, etc.).
- **S3-compatible storage (MinIO and friends)** — set `endpoint` to switch to path-style requests. Credentials are read from the environment variables named in `accessKeyEnv` / `secretKeyEnv`, or, if those are not set, from `AWS_ACCESS_KEY_ID` / `AWS_SECRET_ACCESS_KEY`.

## Supported `llama-server` versions

`llama-run` relies on several relatively recent `llama-server` flags:

- `--models-preset` (multi-model routing)
- `--reasoning` / `--reasoning-format` / `--reasoning-budget`
- `--hf-repo` / `--docker-repo`
- `--webui` / `--no-webui` / `--webui-config*`

Use a version of `llama-server` that includes the features you plan to enable.

## License

MIT License. Copyright (c) SiLeader.
