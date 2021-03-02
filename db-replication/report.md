# Репликация
Поднимаем два инстанса mysql
на мастере запускаем команды:
```sql
GRANT REPLICATION SLAVE ON *.* TO "mydb_slave_user"@"%" IDENTIFIED BY "mydb_slave_pwd"; 
FLUSH PRIVILEGES;
```
проверяем статус
```sql
SHOW MASTER STATUS;
```
из статуса берём значение log файла и позицию для слейва. <br>
на слейве запускаем команды:
```sql
CHANGE MASTER TO MASTER_HOST ='172.27.0.2',
MASTER_USER ='mydb_slave_user',
MASTER_PASSWORD ='mydb_slave_pwd',
MASTER_LOG_FILE ='mysql-bin.000003',
MASTER_LOG_POS =638;
START SLAVE;

```

Запускаем приложение на чтение со слейва, и включаем нагрузку:
```bash
wrk -t1 -c1 -d1000s --timeout 10  --latency 'http://0.0.0.0:3000/account/search_user?firstname=Bobby&secondname=Chase' 
```
смотрим нагрузку на мастер (docker stats mysql_master) <br>

| Name      |  CPU % |  MEM USAGE / LIMIT   |
| ----------- | ----------- | ----------- |
| mysql_master     | 0.08%       |351.2MiB / 1.943GiB |

Мастер не нагружен. <br>

Добавляем еще инстанс mysql.<br>
В конфиге мастера добавляем ```binlog_format = ROW``` для row-based репликации.<br>
Всем инстансам в конфиг дописываем для включения GTID:
```gtid_mode=ON``` <br>
```enforce_gtid_consistency=ON```
Устанавливаем плагины и переменные для полусинхронной репликации на мастере:
```sql
Install plugin rpl_semi_sync_master soname 'semisync_master.so';
set global rpl_semi_sync_master_enabled = 1;
set global rpl_semi_sync_master_timeout = 2000;
```
Запускаем мастер
```sql
GRANT REPLICATION SLAVE ON *.* TO "mydb_slave_user"@"%" IDENTIFIED BY "mydb_slave_pwd"; 
```
проверяем статус мастера
```sql
SHOW MASTER STATUS;
```
Устанавливаем плагины и переменные для полусинхронной репликации на слейвах:
```sql
install plugin rpl_semi_sync_slave soname 'semisync_slave.so';
set global rpl_semi_sync_slave_enabled = 1;
```

Запускаем определяем мастер и запускаем слейвы:
```sql
CHANGE MASTER TO MASTER_HOST ='192.168.192.2',
    MASTER_USER ='mydb_slave_user',
    MASTER_PASSWORD ='mydb_slave_pwd',
    MASTER_LOG_FILE ='mysql-bin.000003',
    MASTER_LOG_POS = 484;
START SLAVE;
```
Создаем тестовую таблицу:
```sql
create table test
(
    id int null,
    dummy1 varchar(10000) null,
    dummy2 varchar(10000) null,
    dummy3 varchar(10000) null,
    dummy4 varchar(10000) null,
    dummy5 varchar(10000) null,
    dummy6 varchar(10000) null
);
```

Запускаем тестовое приложение (cmd/transaction/main.go) (заполняем строки в таблице test).<br>
Убиваем мастер - последняя строка на 1 слейве 1088, на втором 1188 -  второй слейв наиболее свежий после смерти мастера.<br>
Промоутим слейв 2 до мастера:
```sql
STOP SLAVE;
GRANT REPLICATION SLAVE ON *.* TO "mydb_slave_user"@"%" IDENTIFIED BY "mydb_slave_pwd";
stop slave;
set global read_only=OFF;
```
На 1 слейве проверяем позицию журнала и используем это значение при смене мастера:
```sql
STOP slave;
CHANGE MASTER TO MASTER_HOST ='192.168.192.4',
    MASTER_USER ='mydb_slave_user',
    MASTER_PASSWORD ='mydb_slave_pwd',
    MASTER_LOG_FILE ='mysql-bin.000003',
    MASTER_LOG_POS =62259612;
START SLAVE;
```

После запуска - данные в таблицах синхронизируются (1088 -> 1188 строк на первом слейве после подключения ко второму).
