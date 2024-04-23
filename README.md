# Документация к проекту YandexWork
![Uml Diagram](https://github.com/dadilll/YandexWork/assets/147308879/5fa19a2c-3206-4b22-b5ba-988ce6368eef)



Если что то не понятно пишите мне: https://t.me/Dadilesss

# Задачи
## Не начато
- [ ] Задача : Разработать веб-интерфейс для приложения(Низкий Приоритет)
- [ ] Задача : Сделать настройку времени выражения через API(Низкий Приоритет)
- [ ] Задача : Переход на gRPC(Низкий Приоритет)

## В процессе
- [ ] Задача : Завернуть все в докер(Высокий Приоритет)
- [ ] Задача : Покрытие тестами проекта(Высокий Приоритет)

## Завершено
- [x] Задача : Покрытие тестами проекта(Высокий Приоритет)
- [x] Задача : Работа в конкретном пользователе
- [x] Задача : Переход с Redis на Postgres.
- [x] Задача : Возможность регистрировать и логинится в аккаунты.
- [x] Задача : Сделать для разных выражений разное время выполнения
- [x] Задача : Написать скрипт для решения задач формата 2+2+2 
- [x] Задача : Подготовить документацию по проекту
- [x] Задача : Написать скрипт для решения задач формата 2+2
- [x] Задача : Сделать агенты и воркеры
- [x] Задача : Сделать оркестратор

## Перед запуском
Перед запуском необходимо настроить время выполнения операторов и количество агентов и воркеров в configurations.go.

## Запуск без докера
Необхадимо установить PostgreSQL и разметить новые таблицы(Их можно будет посмотреть в файле init.sql). После чего необходимо будет настроить конфигурацию SQL в файле configurations.go. Для запуска проекта требуется запустить два основных скрипта, расположенные в каталогах agentmain и orchestramain.

## Запуск с докером
Пока что у меня нечего не работает. Если вам не в падлу можете попробовать сделать что то с моим говнокодом. 


## Тесты
описание тестов в файле TEST.md

## EndPoint
### Получение списка задач
```bash
curl -X GET http://localhost:8080/expressions \
-H "Authorization: Bearer YOUR_JWT_TOKEN"
```

### Добавление задач
```bash
curl -X POST http://localhost:8080/add \
-H "Content-Type: application/json" \
-H "Authorization: Bearer YOUR_JWT_TOKEN" \
-d '{"expression": "2 + 2"}'
```

### Удаление всех задач
```bash
curl -X DELETE http://localhost:8080/delete-tasks \
-H "Authorization: Bearer YOUR_JWT_TOKEN"
```
### Регистрация нового пользователя (/register)
```bash
curl -X POST -H "Content-Type: application/json" -d '{"login":"", "password":""}' http://localhost:8080/login
```

### Вход пользователя (/login)
```bash
curl -X POST -H "Content-Type: application/json" -d '{"login":"", "password":""}' http://localhost:8080/register
```
