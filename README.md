# buildpack-application-config-environment-variable-labels

application.proerpties / application.yaml 上の環境変数で置換されるキーを抽出し、OCIイメージラベルとして出力するbuildpackです。

## How to use

```bash
pack build -b nncdevel/buildpack-application-config-environment-variable-labels:1.0.0
```

## Configration

| 環境変数                | 説明 |
|------------------------|--------------------------|
| `$BP_APPLICATION_CONFIG_ENVIRONMENT_VARIABLE_LABEL_NAME` | 挿入するラベルのキーを指定します。 <br> デフォルトでは `io.nncdevel.buildpacks.application-config.environment-variables` を利用します。|
| `$BP_APP_CONFIG_ENVIRONMENT_VARIABLE_TARGET_PATTERNS` | キーを抽出する対象ファイルパスをカンマ区切りで設定します。デフォルト値は後述。 |


## キーを抽出する対象ファイルパス

デフォルト値は以下の通り。  
各パスは github.com/mattn/go-zglob のパターンを設定できます。

```
BOOT-INF/classes/application.properties,BOOT-INF/classes/application.ya?ml,WEB-INF/classes/application.properties,WEB-INF/classes/application.ya?ml
```


## Output Label value

以下の3項目を持ったオブジェクト配列のJSON文字列が出力されます。

- `name` string : 環境変数名
- `required` boolean : 環境変数が必須（= プレースホルダのデフォルト値がない）の場合 `true` が設定されます。  
- `defaultValue` string : デフォルト値。設定されていない場合は空文字が設定されます。

例：

```json
[
    {"name": "JDBC_URL", "required": true, "defaultValue": ""},
    {"name": "API_ENDPOINT", "required": false, "defaultValue": "https://github.com"}
]

```

## ビルドに必要なもの

- docker

golang、pack CLIをローカルにインストールする場合は下記の2つが必要になる。

- [pack CLI](https://buildpacks.io/docs/install-pack/)
- [golang](https://golang.org/doc/install)

インストールしない場合は docker-compose をインストールしてください。

## docker-composeを利用したビルド手順

まずは `docker-compose.yml.example` を `docker-compose.yml` としてコピーします。

コピーしたファイルを開き、`http_proxy`、`https_proxy` を各人の環境に合わせて修正してください。

golangのビルド

```bash
$ docker-compose run build-golang
```

buildpackのパッケージング(Dockerイメージ化)

```bash
$ docker-compose run package-buildpack-image
```


## golang、pack-cliをインストールしている場合のビルド手順

golangのビルド

```bash
$ ./scripts/build-golang.sh
```

buildpackのパッケージング(Dockerイメージ化)

```bash
$ ./scripts/package-buildpack-image.sh
```
