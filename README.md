# buildpack-application-config-environment-variable-labels

application.proerpties / application.yaml 上の環境変数で置換されるキーを抽出し、OCIイメージラベルとして出力するbuildpackです。

## Configration

| 環境編集                | 説明 |
|------------------------|--------------------------|
| `$BP_APPLICATION_CONFIG_ENVIRONMENT_VARIABLE_LABEL_NAME` | 挿入するラベルのキーを指定します。 <br> デフォルトでは `io.nncdevel.buildpacks.application-config.environment-variables` を利用します。|


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

イメージのpush

```bash
$ docker tag paketo-dd-java-agent oishikawa/paketo-dd-java-agent:0.0.1
$ docker push oishikawa/paketo-dd-java-agent:0.0.1
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

イメージのpush

```bash
$ ./scripts/push-image.sh
```
