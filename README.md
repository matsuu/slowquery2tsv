# mysql-ps2tsv

## Description

MySQLのperformance schemaの内容を整形した上でTSV形式で出力するツール

## Install

```
# < go-1.16
go get -u github.com/matsuu/mysql-ps2tsv

# >= go-1.16
go install github.com/matsuu/mysql-ps2tsv@latest
```

## Usage

```
mysql-ps2tsv -u <username> -p <password>
mysql-ps2tsv -h <host> -u <username> -p <password>
```

その他オプションは `--help` 参照。

## Tips

* MySQL 8.0環境でqueryが途中で途切れる場合は [performance\_schema\_max\_sql\_text\_length](https://dev.mysql.com/doc/refman/8.0/en/performance-schema-system-variables.html#sysvar_performance_schema_max_sql_text_length) を引き上げてください
