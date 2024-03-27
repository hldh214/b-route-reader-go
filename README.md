# b-route-reader-go

ROHM BP35A1を使用してBルート経由でスマートメーターから消費電力を取得する。

![overview](https://github.com/hldh214/b-route-reader-go/assets/5501843/e40dbe1f-6d63-4b41-80d0-9d1aed629f5c)

## 変更点

- 本プロダクトは、yakumo-saki 氏の[b-route-reader-go](https://github.com/yakumo-saki/b-route-reader-go)から派生しています。
- 変更点は以下の通りです。
  - 本プロダクトは、HomeAssistantとの連携を目的としています。
  - Atmark Techno社のSA-M0を使用しています。（2024年3月現在、ヤフオクにて、非常に手頃な価格でご購入いただけます。）
  - HomeAssistantのMQTT AutoDiscoveryに対応しています。
  - `godotenv`を使用して、環境変数の設定を簡略化しています。

## 動作環境

* SA-M0
* Golang 1.18

開発はSA-M0で行っていますが、Linuxが動いてBP35A1とシリアル通信できる環境であれば、
動作すると思います。

## ビルド方法

```
git pull https://github.com/hldh214/b-route-reader-go
cd b-route-reader-go
GOOS=linux GOARCH=arm GOARM=5 go build -ldflags="-s -w" .
```

b-route-reader-go というファイルが実行ファイルです。

## HomeAssistant との連携

HomeAssistant と連携するためには、以下のような設定を行います。

1. MQTTブローカーを立てる
2. `MQTT_BROKER` にブローカーのアドレスを設定する
3. AutoDiscovery に対応するため、 ほかの設定は不要です。

![energy](https://github.com/hldh214/b-route-reader-go/assets/5501843/0397b729-2362-480a-8152-6dfee3e9f956)

## SA-M0 の設定

### storageをつける

[Armadillo-Box WS1(旧 IIJ SA-M0)で遊ぶ準備をする](https://www2.hatenadiary.jp/entry/2023/04/26/202125) を参考にしてください。

### 起動時に自動実行

```shell
# cat /etc/config/rc.local

# ... 省略

while : ; do /mnt/sd/b-route-reader-go 2>&1 | tee -a /mnt/sd/b-route-reader-go.log ; done &
```

### コンフィグ領域の保存

```shell
flatfsd -s
```

## 動作

* 20秒ごとに瞬時電力(E7)と瞬時電流(E8)を取得します。
* 180秒ごとに積算電力(E0)を取得します。
* 取得したデータをJSONとして一時ファイルに書き込み、設定されたコマンドにパスを引き渡して実行します。

### 動作例

以下の2つの動作を行う。なお瞬時電力と積算電力の取得タイミングが同時に来たとしても、同時に二つの動作が行われることはない。
瞬時電力取得→積算電力取得 の順で処理される（はず）

### 瞬時電力と瞬時電流取得

1. 起動〜初期化
2. （20秒後）瞬時電力と瞬時電流を取得
3. 取得したデータを `/tmp/abcdefg` （ランダムなファイル名） に書き込む。
4. 内容はJSONである。 `{"E7":"1566","E8":"17","E8_Rphase":"8","E8_Tphase":"9","datetime":"2022-07-31T12:49:32.367Z"}`
5. EXEC_CMDで指定されたコマンドを実行する `./dist/zabbix.sh /tmp/abcdefg`
6. 繰り返し

### 積算電力取得

1. 起動〜初期化
2. （180秒後）積算電力を取得
3. 取得したデータを `/tmp/abcdefg` （ランダムなファイル名） に書き込む。
4. 内容はJSONである。 `{"E0":"10197.7","E3":"0.5","datetime":"2022-07-31T12:48:57.613Z"}`
5. EXEC_CMDで指定されたコマンドを実行する `./dist/zabbix.sh /tmp/abcdefg`
6. 繰り返し

## デバッグ・開発

BP35A1と接続された機器上で開発するのも手ですが、色々と面倒です。  
以下のコマンドを使うと、シリアルポートを自分のPCにリダイレクトすることができます。  

* シリアルポートがあるPCのIPアドレスを 10.1.0.191 として説明しています。
* 33444はポート番号です。何番でも構わないので好きな番号にしてください。

Run at serial port machine (eg: raspberry pi)
`socat -d -d /dev/ttyAMA0,echo=0 tcp-listen:33444`

running at your coding pc
`socat -d -d pty,link=$HOME/ttyAMA0,waitslave tcp-connect:10.1.0.191:33444`

### デバッグに関する蛇足

* socatは一度実行されると終了してしまう（なにかおかしい気がする）ので while で無限ループさせる(bash)

`while true; do <コマンド>; done`

### デバッグに関する蛇足２

* socatで転送する前に minicom でハングアップを行っておく必要があります。
* minicom を起動して ctrl-a -> h で `hangup line?` -> `yes` した後、ctrl-a -> q で終了します。（ctrl-a -> x ではない）
* このプログラムを実行すると、ローカルエコーがオフになるのでminicomには何も表示されませんが気にせずコマンドを入れることができます。

## 既知の問題

### 問い合わせ間隔が規約違反

瞬時電力取得→積算電力取得 を連続で取得するときに、ウェイトなしで問い合わせを2連続で投げるのは規約違反である。
最低20秒の間を空ける必要がある（が、普通に動いているので無視している）

### 積算電力量計測値（逆方向）の扱い

取り扱いがよくわからない。太陽光発電等で売電している場合に使用される？
