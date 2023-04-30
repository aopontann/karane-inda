### 事前準備
1. Discordの[新規アプリケーションを作成](https://discord.com/developers/)し、トークンを取得
2. OpenAPIからアクセスキーを取得
1. Google Cloudのプロジェクトを[作成](https://console.cloud.google.com/projectcreate)する

### 開発環境構築
1. [Cloud Shell](https://console.cloud.google.com/cloudshell/open?git_repo=https://github.com/aopontann/karane-inda)を起動
    - セッション終了後にソースコードを削除したい場合は、Trust repoにチェックを入れないように
    - 確認をクリックすると、Cloud Shell Editorが立ち上がります
2. 隠しファイルを表示するようにする。
    - Cloud Shell Editorの「View」タグの「Toggle Hidden Files」をクリックすると、隠しファイルが表示されるようになる。もう一度クリックすると、隠しファイルが表示されなくなる。
3. 次のコマンドを実行して.envファイルを作成
```
cp .env.sample .env
```
4. .envファイルをEditorで開き、事前準備で取得したアクセスキーとトークンを入力する。
```
OPENAI_API_KEY=sk-aaa...
DISCORD_TOKEN=aaa...
```
5. 次のコマンドを入力しBOTが起動できたら開発環境構築完了
```
go run main.go
```

### BOTを常に稼働状態にする（編集中）
Cloud Shell Editorを閉じてしまうと、BOTが停止してしまう。
常に起動しておくために、Google Cloud Compute Engine インスタンスを使用する。
1. Google Cloudのプロジェクトを[作成](https://console.cloud.google.com/projectcreate)する
    - プロジェクトIDを確認しておく
2. Cloud Consoleで[課金を有効](https://console.cloud.google.com/billing?hl=ja)にする
3. [Cloud Shell](https://console.cloud.google.com/cloudshell/open?git_repo=https://github.com/aopontann/karane-inda)を起動
4. 作業するプロジェクトを設定
```
gcloud config set project <プロジェクトID>
```
<プロジェクトID>は先ほど作成したプロジェクトIDに書き換えてね

5. setup.shの9行目のAPIキーとトークンを書き換える
```
echo -e 'OPENAI_API_KEY=<API_KEY> \nDISCORD_TOKEN=<TOKEN>' >> .env
```

6. 次のコマンドでCompute Engine VMインスタンスを作成
```
gcloud compute instances create discord-bot \
    --zone=us-west1-b \
    --machine-type=e2-micro \
    --metadata-from-file=startup-script=setup.sh
```
Would you like to enable and retry (this will take a few minutes)?と聞かれた場合、yを入力しエンターキーを押す  
Cloud Shell の承認画面が出てきた場合、承認をクリックする。

### 参考
- https://codelabs.developers.google.com/codelabs/cloud-compute-engine?hl=ja#0