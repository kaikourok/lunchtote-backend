
# ディレクトリ構造
| ディレクトリ | 概要 |
| --- | --- |
| .vscode | VSCode用設定群 |
| cmd | 各インターフェース用実行ファイル実装 |
| config | 環境毎設定ファイル |
| domain | ドメイン層 |
| domain/service | データ構造のドメイン表現化実装 |
| infrastructure | インフラストラクチャー層 |
| infrastructure/config | 設定読み込み実装 |
| infrastructure/database | データベースによるリポジトリ実装 |
| infrastructure/logger | ロガー実装 |
| library | ユーティリティパッケージ群 |
| registry | DIコンテナ |
| repository | リポジトリ抽象実装 |
| usecase | アプリケーション層 |

# コミットメッセージ・ブランチ名等に関して
コントリビューター数は少ないことが予想されるため、現段階では厳密な規則は設けません。
だいたい分かればOKです。（もしコントリビューターが増えるのであればまた考えます）
英語/日本語についてもどちらを使用して頂いても構いません。

# ブランチ運用に関して
雰囲気としてはGitHub Flowにdevelopブランチを足したような感じです。

| ブランチ | 概要 |
| --- | --- |
| master | リリース用ブランチ。 |
| develop | デプロイ可能なブランチ。ある程度まとまったらdevelopからmasterにマージする。 |
| 作業用ブランチ | developから分岐させ、developに向けてマージさせる。ブランチ名はおおよその変更内容を表すようにする。 |

# 開発する上での目標
## APIレスポンスタイム
「ユーザーから見て操作がシンプルで機能も単純」なものについてはm4.large程度のサーバーで20ms以内の応答を目指して下さい。例えば他のユーザーのフォローなどです。操作が複雑だったり機能が複雑だったりするもの、例えばユーザーの検索やアカウント削除などは300ms以内の応答を目標とします。環境由来のばらつきを考慮し、応答時間が9割程度の確率で達成されていれば目標達成として扱います。

この目標は管理者用のAPIに対しては適用されません。また、操作ユーザーにとってのリアルタイム性がそこまで重要でないもの、例えばアクセスログや通知などはgo funcなどを用い操作ユーザーのAPIレスポンスタイムへの影響を少なくできる場合に限り、極度の効率化を行わなくても構いません。

## DB操作数の最小化
目標APIレスポンスタイム実現のため、APIアクセスごとのDB操作数を最小にすることに努めてください。直接SQLを記述し、JOINなども駆使し、可能な限りDB操作数を1に近づけます。IN句を使用したEager Loadingのようなことも避けてください。

DB操作数を減らすためであれば、エラーの粒度が荒くなることを許容します。例えば何らかのレコードをDELETEしようとした場合、DELETE対象が存在しなかった場合でも望んだ結果にはなっているので、わざわざレコードがあるかチェックし、最初からなかったらエラーを返す……などといったことをする必要はありません。

## 不要な永続化データ取得をしない
永続化データの取得は可能な限り省コストで行わなければなりません。

例えばユーザの取得を行う場合、Idと名前のみが必要なのであれば、それ以上の取得を行う関数を呼び出してはいけません。具体的な例としては、`SELECT id, name, nickname FROM users;`などを用いている関数は`nickname`を不要に取得しているので使用不可です。これは結果が単一行であることが予期される場合でも同じです。

上記実現のためであれば、`repository`、`infrastructure/database`の極度な細分化を許容します。

## API操作数の最小化
高速化のため、フロントエンド側でのユーザーの1操作に対してAPIの操作数は1を原則としています。そのため、例えばユーザー情報、フォロー情報などをそれぞれ取得し、フロントエンドなどで組み立てて表示……というユースケースは発生しません。ユーザーとフォロー情報をあわせた情報が必要なのであれば、そのような情報が1度で帰るAPIを使用することになります。

このとき、永続化データ取得時と同様、1つでも不要なフィールドが含まれるAPIを使用してはいけません。

# アーキテクチャ
クリーンアーキテクチャとなるよう設計しています。ただし速度の追求のため、一般的なクリーンアーキテクチャとは相違している点があります。以下にその相違点を記します。

## ユースケース駆動
このプロジェクトでは高速化のため「1操作、1APIアクセス」を原則としている都合、ユースケースに合わせて開発を進めるのが最も効率の良い方法となります。そのため、ユースケース駆動で開発を行っています。

よく採用されるドメイン駆動と違い、ドメインに対して深く洞察を行うということは行いません。今プロジェクトにおけるドメインの構造体はインターフェース層とユースケース層の橋渡しとなる存在として考えるとよいでしょう。

処理に向いたデータ構造からドメインに向いたデータ構造に変換する場合、多くは関数によって表現します。例えばドメインにおいて年齢というフィールドが必要とされがちであるにも関わらず、データ構造においては生年月日という保持形式であった場合、`entity/service`内に`CurrentAgeFromBirthday`というような関数を設けます。

# 命名規則
基本的に[Effective Go](https://go.dev/doc/effective_go#names)を参照しますが、Effective Goの命名規則に比べ、その箇所を見ただけで役割が明らかであることを優先します。そのため、全体を通し様々な命名は長いものとなります。

以下、Effective Goとは異なる箇所について大まかに記載しています。

## パッケージ名
単数形、小文字のみで1単語になるようにします。その際、1単語になる範囲内で可能な限り省略を避けてください。
混成語など、一般に1単語として扱われるような単語は許容します。

また、utilなど責任が明確でないパッケージは作成してはいけません。ユーティリティ用途として使用する関数についても適切な名前を持ったパッケージ内で管理してください。

なお、各パッケージ内でのみ使用されるユーティリティ関数はパッケージ内にutil.goなどとして配置しても構いません。また、パッケージ名が適切であればutilなどのディレクトリ内に配置しても構いません。

| ❌NG | ✅OK | NG理由・補足 |
| ---- | ---- | ---- |
| infra | infrastructure | 省略形であるため。 |
| db | database | 省略形であるため。 |
| entities | entity | 単数形でないため。 |
| protocol-buffers | protobuf | 2語であるため。protobufは1単語として用いられることが多いため許容される。　|
| library | secure | 責任が明確でないため。パッケージ名が適切であれば`library/secure`などに配置することは許容される。 |

## 変数・引数
一目見ただけでそれがどのような役割を持っているかを明確にします。
例えば、HTTPクランアントを作成する場合、用いる変数名は`c`ではなく`client`であるべきです。

ただし、その影響範囲が非常に短く役割の類推が非常に容易な場合のみ短い変数名が許容されます。
for文の`i`であったり、`concatStrings`関数内での`s`などです。

## レシーバ名
アルファベット1文字あるいは2文字程度の短い命名を行ってください。レシーバの場合、変数名で説明を行わずとも役割が明らかなためです。

同じ型に対してのレシーバ名は統一されている必要があります。

## 略語等の命名規則について
略語や固有名詞であっても強制的にキャメルケースを適用します。gRPCなど、小文字大文字が交じる場合などに対応しきれないためです。

以下は具体例です。

| ❌NG | ✅OK |
| ---- | ---- |
| URL | Url |
| ID | Id |
| gRPC | Grpc |

# ライセンス
MIT License
