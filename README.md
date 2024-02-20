# Документация к проекту YandexWork
![Uml Diagram](https://github.com/dadilll/YandexWork/assets/147308879/5fa19a2c-3206-4b22-b5ba-988ce6368eef)



Если что то не понятно пишите мне: https://t.me/Dadilesss

# Задачи
## Не начато

- [ ] Задача : Сделать настройку времени выражения через API
- [ ] Задача : Разработать веб-интерфейс для приложения

## В процессе

- [ ] Задача : Написать скрипт для решения задач формата 2+2+2 

## Завершено

- [x] Задача : Подготовить документацию по проекту
- [x] Задача : Написать скрипт для решения задач формата 2+2
- [x] Задача : Сделать агенты и воркеры
- [x] Задача : Сделать для разных выражений разное время выполнения
- [x] Задача : Сделать оркестратор

## Перед запуском
Перед запуском необходимо установить Redis и создать базу данных, также требуется настроить время выполнения в expression.go.

## Запуск 

Для запуска проекта требуется запустить два основных скрипта, расположенные в каталогах agentmain и orchestramain.

## EndPoint

### Получение списка задач
Для получения списка задач используйте следующий запрос:

```bash
curl -X GET http://localhost:8080/expressions 
```

### Добавление задач
Чтобы добавить новую задачу, выполните POST-запрос с указанием выражения в формате JSON. Пример запроса:

```bash
curl -X POST -H "Content-Type: application/json" -d '{"expression":"your_expression_here"}' http://localhost:8080/add
```

### Удаление всех задач
Чтобы удалить все задачи из системы, выполните DELETE-запрос на эндпоинт /delete-all. Ниже приведен пример:

```bash
curl -X DELETE http://localhost:8080/delete-all
```
