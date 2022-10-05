set -euo pipefail

pack buildpack package paketo-environment-variable-labels --config ./package.toml --format image
