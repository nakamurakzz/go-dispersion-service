# go-dispersion-service
## プロトコルバッファ
- Googleが設計したデータのシリアライズ形式
  - 一貫性のあるスキーマ
    - 分散サービスの色んなところで同じ型定義を使える
    - スキーマとコードが一致しているか、のテストができる
  - バージョン管理
    - フィールド追加、削除が容易
  - ボイラーテンプレートコードの削減
    - デコードとエンコードのコードはprotobufライブラリがやってくれる
  - 拡張性
    - 独自にコンパイルすることもできる
    - 共通的な独自メソッドを定義することもできる
  -  言語非依存
    -  様々な言語でも使える
      - Go
      - Java
      - Python
      - C++
      - JavaScript 等
  -  パフォーマンス
    - データが小さく、JSONよりも高速にシリアライズ可能
- gPRCではプロトコルバッファを使用してAPI定義を行い、メッセージをシリアライズしている

## プロトコルバッファの定義
インストール確認
```bash
protoc --version
libprotoc 22.2
```

プロトコルバッファの定義
```proto
syntax = "proto3";

package log.v1;

option go_package = "github.com/nakamurakzz/api/log_v1";

// 型 名前 フィールドID
message Record {
  bytes value = 1;
  uint64 offset = 2;
}
```

- コンパイルは各言語のランタイムに依存する
- コンパイルされると各言語のコードが生成される

