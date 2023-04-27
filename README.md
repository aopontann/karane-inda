### 事前準備
1. Google Cloudの[プロジェクトの作成](https://console.cloud.google.com/projectcreate)
2. Discordの[新規アプリケーションを作成](https://discord.com/developers/)
3. OpenAPIからアクセスキーを取得

### VMインスタンスの起動からSSH接続まで
1. Cloud Shellを[アクティブ](https://console.cloud.google.com/cloudshell/open?git_repo=https://github.com/aopontann/karane-inda)にする  
    セッション終了後にソースコードを削除したい場合は、Trust repoにチェックを入れないように
2. 作業するプロジェクトを設定
```
gcloud config set project <プロジェクトID>
```
<プロジェクトID>は先ほど作成したプロジェクトIDに書き換えてね  
3. 次のコマンドでCompute Engine VMインスタンスを作成
```
gcloud compute instances create discord-bot \
    --zone=us-west1-b \
    --machine-type=e2-micro \
    --metadata-from-file=startup-script=setup.sh
```
Would you like to enable and retry (this will take a few minutes)?と聞かれた場合、yを入力しエンターキーを押す  
Cloud Shell の承認画面が出てきた場合、承認をクリックする。

4. インスタンスにSSH接続
```
gcloud compute ssh --zone "us-west1-b" "discord-bot" 
```
Do you want to continue (Y/n)?　と聞かれた場合、yと入力しエンターキーを押す  
登録するパスワードを聞かれるので、パスワードを決めて入力する。パスワードを設定しない場合は何も入力せずエンターキーを押す  
もう一度パスワードを入力する  
少し時間が経ったあと、この文字が出力されるようになったら接続成功
```
<ユーザー名>@discord-bot:~$
```

### ソースコードの取得からBOTの起動まで
1. ソースコードダウンロード
```
wget https://github.com/aopontann/karane-inda/archive/main.tar.gz
```
2. 解凍
```
tar -xzf main.tar.gz
```
3. 作業ディレクトリ移動
```
cd karane-inda-main
```
4. .envを作成
```
echo -e 'OPENAI_API_KEY=<OpenAIのAPIキー> \nDISCORD_TOKEN=<Discordのトークン>' >> .env
```
<OpenAIのAPIキー>と<Discordのトークン>は書き換えてね </br>
5. BOT起動
```
go run main.go
```