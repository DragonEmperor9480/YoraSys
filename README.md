<p align="center">
  <img src="assets/YoraSys.png" alt="YoraSys Banner" width="100%" />
</p>

<h1 align="center">YoraSys</h1>
<p align="center"><strong>Storage Intelligence Platform for safe, transparent cleanup</strong></p>

<p align="center">
  Detect caches. Score confidence. Classify risk. Clean with control.
</p>

## Why YoraSys

Most cleanup tools optimize for deletion speed. YoraSys is built to optimize for trust: clear detection, explainable confidence, explicit safety levels, and controlled cleanup actions.

## Origin Of The Name

**YoraSys** is derived from:

- **Yora**: inspired by the early YoRHa-style internal naming used during the first system design drafts
- **Sys**: short for **System**, representing the shift to a practical, maintainable engineering platform

The name reflects the project philosophy: concept-driven beginnings, system-first execution.

## Vision

YoraSys is a data-driven storage intelligence system that:

- Detects system, developer, and GPU/toolchain caches
- Uses confidence-based classification
- Provides safe and transparent cleanup workflows
- Uses a community-driven knowledge registry

## Architecture

```text
[YAML Registry] -> [Go Engine (CLI)] -> [JSON Output] -> [Qt GUI]
```

### Components

- **Go Engine**: cache detection, confidence scoring, safety classification, cleanup execution
- **YAML Registry**: cache definitions, rules, and metadata
- **JSON Output Layer**: structured contract between engine and UI
- **Qt UI**: desktop interface via `QProcess` calls to the CLI

## CLI Direction

```bash
yorasys scan --json
yorasys clean --id <cache_id> --json
yorasys analyze --mode dev --json
```

## Confidence And Safety

Each detected cache includes:

- `confidence` score (`0.0` to `1.0`)
- `safety` classification:
  - `safe`
  - `rebuild_required`
  - `risky`

## Design Principles

- Optimize for correctness, not premature performance
- Keep engine logic separate from presentation logic
- Use YAML for knowledge and JSON for communication
- Build trust through transparent evidence and safe defaults

## Development Phases

1. **Phase 1: Detection Engine**
   - YAML registry loading
   - cache discovery
   - confidence scoring
   - no deletion
2. **Phase 2: Safe Cleanup**
   - high-confidence cleanup
   - dry-run mode
3. **Phase 3: Advanced Cleanup**
   - medium-confidence handling
   - explicit warnings
4. **Phase 4: UI Integration**
   - Qt desktop UI
   - JSON communication
5. **Phase 5: Intelligence Layer**
   - insights
   - persona-based recommendations

## Repository Layout

```text
engine/       # Go entrypoint
internal/     # Go backend packages
registry/     # YAML knowledge base
terminal/     # CLI-facing presentation area
ui/qt/        # Qt UI placeholder
web/          # Website placeholder
scripts/      # Build and development scripts
```

## Getting Started

```bash
go mod tidy
go run ./engine
```

## Contributing

Contributions are welcome in:

- registry definitions and validation
- detection/scoring accuracy
- cleanup safety improvements
- tests and documentation

## Status

Early-stage active development. The current focus is engine reliability, schema contracts, and safe cleanup behavior before advanced UI and intelligence features.
