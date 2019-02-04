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


## Cloud Functionで受け取るPub/Subからのメッセージ(Data)のサンプル
```json
{
  "id": "[BUILD_iD]",
  "projectId": "[PROJECT_ID]",
  "status": "WORKING",
  "source": {
    "repoSource": {
      "projectId": "[PROJECT_ID]",
      "repoName": "github_[REPO_OWNER]_[REPO_NAME]",
      "branchName": "[BRANCH]"
    }
  },
  "steps": [
    {
      "name": "gcr.io/cloud-builders/go",
      "entrypoint": "pwd"
    },
    {
      "name": "gcr.io/cloud-builders/go",
      "args": [
        "-c",
        "pwd\nls ./\necho \"END...............\"\n"
      ],
      "entrypoint": "/bin/sh"
    }
  ],
  "createTime": "2019-02-04T08:48:37.196669771Z",
  "startTime": "2019-02-04T08:48:38.777833995Z",
  "timeout": "600s",
  "logsBucket": "gs://[PROJECT_NUMBER].cloudbuild-logs.googleusercontent.com",
  "sourceProvenance": {
    "resolvedRepoSource": {
      "projectId": "[PROJECT_ID]",
      "repoName": "github_[REPO_OWNER]_[REPO_NAME]",
      "commitSha": "[COMMIT_SHA]"
    }
  },
  "buildTriggerId": "[BUILD_TRIGGER_ID]",
  "options": {
    "substitutionOption": "ALLOW_LOOSE",
    "logging": "LEGACY"
  },
  "logUrl": "https://console.cloud.google.com/gcr/builds/[BUILD_ID]?project=[PROJECT_ID]",
  "tags": [
    "event-cb00ff11-ff51-4df1-ac47-6d4c90f54e5c",
    "trigger-09537d59-8e4a-4251-a3af-8bd05150c3eb"
  ]
}
```
