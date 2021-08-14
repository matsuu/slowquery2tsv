# slowquery2tsv

## Description

MySQLの `performance_schema` やPostgreSQLの `pg_stat_statements` の内容を整形した上でTSV形式で出力するツール

## Install

```
# < go-1.16
go get -u github.com/matsuu/slowquery2tsv

# >= go-1.16
go install github.com/matsuu/slowquery2tsv@latest
```

## Prerequisite

MySQLの場合 `performance_schema` を有効にする。

```
performance_schema = ON
performance_schema_max_sql_text_length = 1024
```

PostgreSQLは `pg_stat_statements` を有効にする。

```
shared_preload_libraries = 'pg_stat_statements'
pg_stat_statements.max = 1000
pg_stat_statements.track = top
pg_stat_statements.save = on
track_activity_query_size = 1024
```

```
CREATE EXTENSION pg_stat_statements;
```

## Usage

```
slowquery2tsv -u <username> -p <password>
slowquery2tsv -h <host> -u <username> -p <password>
```

その他オプションは `--help` 参照。

## Tips

* MySQL
    * queryが途中で途切れる場合は [performance\_schema\_max\_sql\_text\_length](https://dev.mysql.com/doc/refman/8.0/en/performance-schema-system-variables.html#sysvar_performance_schema_max_sql_text_length) を引き上げてください
    * リセットは以下で可能です
        ```
        CALL sys.ps_truncate_all_tables(FALSE);
        ```
* PostgreSQL
    * queryが途中で途切れる場合は [track\_activity\_query\_size](https://www.postgresql.jp/document/current/html/runtime-config-statistics.html) を引き上げてください
    * リセットは以下で可能です
        ```
        SELECT pg_stat_statements_reset();
        ```
