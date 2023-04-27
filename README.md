### 事前準備
1. Google Cloudの[プロジェクトの作成](https://console.cloud.google.com/projectcreate)
2. Discordの[新規アプリケーションを作成](https://discord.com/developers/)
3. OpenAPIからアクセスキーを取得

### 開発
1. Cloud Shellを[アクティブ](https://console.cloud.google.com/cloudshell/open?git_repo=https://github.com/aopontann/karane-inda)にする </br>
    セッション終了後にソースコードを削除したい場合は、Trust repoにチェックを入れないように
2. 次のコマンドでCompute Engine VMインスタンスを作成
```
gcloud compute instances create discord-bot \
    --zone=us-west1-b \
    --machine-type=e2-micro \
    --metadata-from-file=startup-script=setup.sh
```
3. インスタンスにSSH接続
```
gcloud compute ssh --zone "us-west1-b" "discord-bot" 
```
4. ソースコードダウンロード
```
wget https://github.com/aopontann/karane-inda/archive/main.tar.gz
```
5. 解凍
```
tar -xzf main.tar.gz
```
6. 作業ディレクトリ移動
```
cd karane-inda-main
```
7. .envを作成
```
echo -e 'OPENAI_API_KEY=<OpenAIのAPIキー> \nDISCORD_TOKEN=<Discordのトークン>' >> .env
```
<OpenAIのAPIキー>と<Discordのトークン>は書き換えてね </br>
8. BOT起動
```
go run main.go
```