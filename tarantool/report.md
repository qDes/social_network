
Настроить репликацию не удалось (было Nное количество попыток собрать репликатор в докере центоса, но в итоге не завелось(скомпилировал из исходников репликатор, подмонитровал конфиги, при запуске replicatord ничего не происходило)). <br>
Написан скрипт для заливки данных аналогичных данных таблице users из MySQL в инстанс тарантула.
В качестве референсной таблицы использована таблица users.
В тарантуле создано пространство:
```json
c = box.schema.space.create('users')
c:format({
         {name = 'id', type = 'unsigned'},
         {name = 'username', type = 'string'},
         {name = 'first_name', type = 'string'},
         {name = 'second_name', type = 'string'}
         })
c:create_index('primary', {
         type = 'TREE',
         parts = {'id'}
         })
c:create_index('secondary', {
         type = 'TREE', unique=false,
         parts = {'first_name', 'second_name'}
         })
```
Для сравнения производительности выбран запрос из mysql вида  `firstName LIKE ? and secondName LIKE ?` и написана lua процедура для поиска записи по имени и фамилии:
```json
function name_search(first, second)
    local ret = {}
    for _, tuple in box.space.users.index.secondary:pairs({first, second}, {iterator='GE'}) do
        if (string.startswith(tuple[3], first, 1, -1) and string.startswith(tuple[4], second, 1, -1)) then
            table.insert(ret, tuple)
        end
    end
    return ret
end
```
Проведем нагрузочное тестирование с помощью wrk: <br>
wrk -t1 -c1 -d10s --timeout 10  --latency 'http://0.0.0.0:3000/account/search_user_tarantool?firstname=Bobby&secondname=Chase'
1 threads and 1 connections -  0.99 RPS <br>
1 threads and 10 connections - 1.00 RPS <br>
1 threads and 100 connections 0.4 RPS <br>
Сравним значения RPS для аналогичного запроса в MySQL:
wrk -t1 -c1 -d10s --timeout 10  --latency 'http://0.0.0.0:3000/account/search_us?firstname=Bobby&secondname=Chase'
1 threads and 1 connections -  405.82 RPS <br>
1 threads and 10 connections - 1480.50 RPS <br>
1 threads and 100 connections 1826.57 RPS <br>

Значения rps отличаются на 2-3 порядка в сторону MySQL. Это связано с тем, что запросы вида LIKE в тарантуле реализуются через итерацию по большому количеству данных и сравнениюстрок и имеет низкую эффективность.
Нагрузку такого характера бессмысленно переносить на тарантул.

# UPD
Сделаем вид что запрос может возвращать только 1 запись и немного изменим lua-процедуру
```json
function name_search(first, second)
    local ret = {}
    for _, tuple in box.space.users.index.secondary:pairs({second, first}, {iterator='GE'}) do
        if (string.startswith(tuple[3], first, 1, -1) and string.startswith(tuple[4], second, 1, -1)) then
            table.insert(ret, tuple)
            return ret
        end
    end
end
```
В этом случае данные нагрузочного тестирования
1 threads and 1 connections - 736 RPS
1 threads and 10 connections - 5716 RPS
1 threads and 100 connections - 15621 RPS
1 threads and 1000 connections - 23689 RPS
Значительное повышение производительности происходит за счет того что процедура выходит из цикла при первом совпадении без перебора всего набора данных.