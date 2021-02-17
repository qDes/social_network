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
Запускаем мастер
```sql
GRANT REPLICATION SLAVE ON *.* TO "mydb_slave_user"@"%" IDENTIFIED BY "mydb_slave_pwd"; 
FLUSH PRIVILEGES;
```
проверяем статус
```sql
SHOW MASTER STATUS;
```
Запускаем 2 слейва:
```sql
CHANGE MASTER TO MASTER_HOST ='172.28.0.2',
MASTER_USER ='mydb_slave_user',
MASTER_PASSWORD ='mydb_slave_pwd',
MASTER_LOG_FILE ='mysql-bin.000003',
MASTER_LOG_POS =678;

SET @rpl_semi_sync_slave = 1;

START SLAVE;
```
Создаем тестовую таблицу:
```sql
create table mydb.test
(
    id bigint null
);
```

Запускаем тестовое приложение (cmd/transaction/main.go) - приложение пишет в таблицу - в 1 транзакции 100 строчек.<br>
Убиваем мастер - последняя закомиченная строчка 27600 - проверяем на слейвах - соотвествует.
Промоутим слейв 2 до мастера:
```sql
flush tables; flush logs;
stop slave;
set global read_only=OFF;
GRANT REPLICATION SLAVE ON *.* TO "mydb_slave_user"@"%" IDENTIFIED BY "mydb_slave_pwd";
```
Переключаем слейв 1 на новый мастер:
```sql
show slave status;
flush tables; flush logs;
stop slave;
CHANGE MASTER TO MASTER_HOST='172.28.0.4',
MASTER_USER='mydb_slave_user',
MASTER_PASSWORD='mydb_slave_pwd',
MASTER_LOG_FILE='mysql-bin.000004',
MASTER_LOG_POS=473148;
START SLAVE;
```

Запускаем приложение - пишем в новый мастер - читаем на слейве(потерь нету).