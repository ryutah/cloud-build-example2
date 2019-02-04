# cloud-build-example
Cloud Buildを利用したGitHubの連携のサンプル


## Motivation
* GitHub Appsを利用する場合、選択可能なGCPプロジェクトがリポジトリごとではなく組織毎になってしまう等の問題があった
* CircleCIはGitHubの組織に紐づくため、状況によっては満足のいく環境でビルドをできないことがある


## 手順
### 1. Cloud Buildコンソールでトリガー追加する、Cloud Buildを作成する
GCPコンソールで、GitHubをソースとするトリガーを作成する

![create_trigger](https://user-images.githubusercontent.com/6662577/52204061-6aa70380-28b6-11e9-8600-f49569d10f95.png)


### 2. GitHubの設定画面からPersonal Access Tokenを生成する
`Settings > Developer settings` から生成

![generate_token](https://user-images.githubusercontent.com/6662577/52202980-5dd4e080-28b3-11e9-8554-217035d8700d.png)


### 3. Cloud Functionをデプロイする
```console
$ cd ./function
$ gcloud functions deploy hello \
    --entry-point HelloPubSub \
    --runtime go111 \
    --trigger-topic cloud-builds \
    --set-env-vars "AUTH_TOKEN=[PERSONAL_ACCESS_TOKEN]"
```

### 4. GitHubにコードをPushする
ビルドが実行され、コミットに対してビルド結果がCloud Funtionから設定される

![status](https://user-images.githubusercontent.com/6662577/52202090-ccfd0580-28b0-11e9-8780-bb07c40ca13a.png)
