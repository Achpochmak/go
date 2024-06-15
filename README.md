# Домашнее задание №3 «Рефакторинг слоя базы данных»

## Основное задание

### Цель

Модифицируйте приложение, написанное в "Домашнее задание №2", чтобы взаимодействие с хранением данных было через Postgres, а не через файл.

### Задание

- Переведите ваше приложение с хранения данных в файле на Postgres.
- Реализуйте миграцию для DDL операторов.
- Используйте транзакции.

## Дополнительное задание

### Анализ запросов в БД

Были составлены и выполнены следующие запросы:

#### Запрос 1: Удаление заказа

```sql
EXPLAIN ANALYZE DELETE FROM orders WHERE id = 1;
```
#### Результаты: 
```sql
QUERY PLAN
--------------------------------------------------------
Delete on orders  (cost=0.15..8.17 rows=0 width=0) (actual time=0.038..0.038 rows=0 loops=1)
   ->  Index Scan using orders_pkey on orders  (cost=0.15..8.17 rows=1 width=6) (actual time=0.014..0.015 rows=1 loops=1)
         Index Cond: (id = 1)
Planning Time: 0.048 ms
Execution Time: 0.052 ms
(5 rows)
```
#### Запрос 2: Выборка заказов с возвратом

```sql
EXPLAIN ANALYZE SELECT id, idReceiver, storageTime, delivered, refund, createdAt, deliveredAt FROM orders WHERE refund = true LIMIT 10 OFFSET 0;

```
#### Результаты: 
```sql
QUERY PLAN
--------------------------------------------------------
Limit  (cost=0.00..0.35 rows=10 width=34) (actual time=0.010..0.011 rows=1 loops=1)
   ->  Seq Scan on orders  (cost=0.00..23.10 rows=655 width=34) (actual time=0.009..0.010 rows=1 loops=1)
         Filter: refund
         Rows Removed by Filter: 1
Planning Time: 0.047 ms
Execution Time: 0.023 ms
(6 rows)
```

#### Запрос 3: Выборка всех заказов

```sql
EXPLAIN ANALYZE SELECT id, idReceiver, storageTime, delivered, refund, createdAt, deliveredAt FROM orders;
```
#### Результаты: 
```sql
QUERY PLAN
--------------------------------------------------------
Seq Scan on orders  (cost=0.00..23.10 rows=1310 width=34) (actual time=0.007..0.008 rows=2 loops=1)
Planning Time: 0.030 ms
Execution Time: 0.029 ms
(3 rows)
```
#### Запрос 4: Выборка заказов по получателю 

```sql
EXPLAIN ANALYZE SELECT id, idReceiver, storageTime, delivered, refund, createdAt, deliveredAt FROM orders WHERE idReceiver = 123 ORDER BY createdAt DESC LIMIT 10;
```
#### Результаты: 
```sql
QUERY PLAN
--------------------------------------------------------
Limit  (cost=26.47..26.49 rows=7 width=34) (actual time=0.036..0.036 rows=0 loops=1)
   ->  Sort  (cost=26.47..26.49 rows=7 width=34) (actual time=0.035..0.035 rows=0 loops=1)
         Sort Key: createdat DESC
         Sort Method: quicksort  Memory: 25kB
         ->  Seq Scan on orders  (cost=0.00..26.38 rows=7 width=34) (actual time=0.005..0.006 rows=0 loops=1)
               Filter: (idreceiver = 123)
               Rows Removed by Filter: 2
Planning Time: 0.089 ms
Execution Time: 0.051 ms
(9 rows)
```
### Вывод
Так как данных в базе довольно мало, затраты по времени довольно небольшие, но при необходимости можно добавить индексы на refund (возврат товара) и на createdAt idReceiver для более быстрой работы программы. В данном случае логично использовать B-tree, так как необходимо обеспечить сортировку и ограничение запроса.

#### Запрос 1: Выборка заказов с возвратом

```sql
EXPLAIN ANALYZE SELECT id, idReceiver, storageTime, delivered, refund, createdAt, deliveredAt FROM orders WHERE refund = true LIMIT 10 OFFSET 0;

```
#### Результаты: 
```sql
QUERY PLAN
--------------------------------------------------------
Limit  (cost=0.00..1.02 rows=1 width=34) (actual time=0.006..0.006 rows=1 loops=1)
   ->  Seq Scan on orders  (cost=0.00..1.02 rows=1 width=34) (actual time=0.005..0.005 rows=1 loops=1)
         Filter: refund
         Rows Removed by Filter: 1
Planning Time: 0.045 ms
Execution Time: 0.014 ms
(6 rows)
```
#### Запрос 2: Выборка заказов по получателю 

```sql
EXPLAIN ANALYZE SELECT id, idReceiver, storageTime, delivered, refund, createdAt, deliveredAt FROM orders WHERE idReceiver = 123 ORDER BY createdAt DESC LIMIT 10;

```
#### Результаты: 
```sql
QUERY PLAN
--------------------------------------------------------
Limit  (cost=1.03..1.04 rows=1 width=34) (actual time=0.027..0.027 rows=0 loops=1)
   ->  Sort  (cost=1.03..1.04 rows=1 width=34) (actual time=0.026..0.026 rows=0 loops=1)
         Sort Key: createdat DESC
         Sort Method: quicksort  Memory: 25kB
         ->  Seq Scan on orders  (cost=0.00..1.02 rows=1 width=34) (actual time=0.023..0.023 rows=0 loops=1)
               Filter: (idreceiver = 123)
               Rows Removed by Filter: 2
Planning Time: 0.072 ms
Execution Time: 0.042 ms
(9 rows)

```
### Итог
Как показал анализ, время на планирование запроса увеличилось, но время выполнения сократилось. Так как база данных маленькая, индексы вводить невыгодно, потому что временные затраты на планирование превышают профит от выполнения запроса. Однако, с увеличением количества данных выигрыш от уменьшения времени выполнения будет более существенным, и можно будет воспользоваться данной опцией.