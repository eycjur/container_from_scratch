# Container from Scratch

Golangを使ってDockerのようなコンテナランタイムを自作します。Linux Namespace、cgroups、chrootといったコンテナ技術の基礎を実装しています。

このリポジトリは以下のLiz Rice氏による講演動画に沿って作成しています。

[![Container From Scratch](https://i.ytimg.com/vi/8fi7uSYlOdc/hq720.jpg?sqp=-oaymwEnCNAFEJQDSFryq4qpAxkIARUAAIhCGAHYAQHiAQoIGBACGAY4AUAB&rs=AOn4CLCHquVb4pt4jfJqWO799-WwXGgp6A)](https://www.youtube.com/watch?v=8fi7uSYlOdc)

https://www.youtube.com/watch?v=8fi7uSYlOdc

## 🚀 Quick Start

### 1. Lima環境のセットアップ（macOSの場合）

Linuxカーネルの機能が必要なため、Limaを使用してLinux仮想環境を用意します。

```bash
# Limaのインストール
brew install lima

# Ubuntu LTS環境の起動
limactl start --name=default template://ubuntu-lts

# Limaシェルに入る
limactl shell default

# 環境の確認
grep VERSION_ID /etc/os-release
# VERSION_ID="24.04"
```

### 2. コンテナの実行
```bash
# Lima環境内で以下を実行

# root権限を取得
sudo -s

# コンテナを起動（docker run相当）
go run main.go run /bin/bash
# Running [/bin/bash] as 8150
# Running [/bin/bash] as 1
```

### 3. 動作確認

コンテナ内で以下を実行：

```bash
# ホスト名の分離を確認
hostname
# container

# PID namespaceの分離を確認（PID 1から始まる）
ps
# PID TTY          TIME CMD
#   1 pts/1    00:00:00 exe
#   6 pts/1    00:00:00 bash
#  13 pts/1    00:00:00 ps
```

別のターミナルからホスト側のプロセスツリーを確認：

```bash
# ホスト側（Lima内のroot）で実行
ps fa
#  PID TTY      STAT   TIME COMMAND
# 8072 pts/0    Ss     0:00 /bin/bash --login
# 8084 pts/0    S+     0:00  \_ sudo -s
# 8085 pts/1    Ss     0:00      \_ sudo -s
# 8086 pts/1    S      0:00          \_ /bin/bash
# 8191 pts/1    Sl     0:00              \_ go run main.go run /bin/bash
# 8248 pts/1    Sl     0:00                  \_ /tmp/go-build2483503065/b001/exe/main ru
# 8253 pts/1    Sl     0:00                      \_ /proc/self/exe child /bin/bash
# 8258 pts/1    S+     0:00                          \_ /bin/bash
```
