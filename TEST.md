## Тесты для пакета `agent`

### TestAgent
- Проверяет, что агент корректно извлекает задачи из базы данных и обрабатывает их.
- Создает мок базы данных и настраивает ожидания для запроса `SELECT`.
- Запускает агента и ждет, чтобы он обработал задачи.
- Проверяет, что все ожидания выполнены.

### TestMarkTaskAsBeingProcessed
- Проверяет функцию `MarkTaskAsBeingProcessed`, которая должна помечать задачу как обрабатываемую.
- Создает тестовый агент.
- Вызывает функцию `MarkTaskAsBeingProcessed` с тестовым `taskID`.
- Проверяет, что задача помечена как обрабатываемая.

### TestMarkTaskAsFinished
- Проверяет функцию `MarkTaskAsFinished`, которая должна помечать задачу как завершенную.
- Создает тестовый агент.
- Помечает тестовую задачу как обрабатываемую.
- Вызывает функцию `MarkTaskAsFinished` для тестового `taskID`.
- Проверяет, что задача помечена как завершенная.

### TestNewAgent
- Проверяет функцию `NewAgent`, которая должна создавать экземпляр агента с заданными параметрами.
- Создает агента с заданными параметрами и проверяет их соответствие.

### TestGetQueueIndex
- Проверяет функцию `GetQueueIndex`, которая должна возвращать корректный индекс очереди для задачи.
- Создает тестовый агент с заданным количеством рабочих.
- Получает индекс очереди для задачи и проверяет его корректность.

### TestProcessTask
- Проверяет функцию `ProcessTask`, которая должна обрабатывать задачу.
- Создает мок базы данных и агента с этим моком.
- Устанавливает ожидания для запроса к базе данных.
- Обрабатывает тестовую задачу и проверяет выполнение всех ожиданий.

## Тесты для пакета `expression`

### TestParseExpression
- Проверяет функцию `ParseExpression`, которая должна правильно анализировать математические выражения.
- Задает различные тестовые выражения и ожидаемые результаты.
- Проверяет правильность вычислений и времени выполнения для каждого выражения.

### TestParseExpression_InvalidCharacters
- Проверяет обработку некорректных символов в выражении.
- Задает выражения с некорректными символами и проверяет, что они вызывают ошибку.

## Тесты для пакета `domain`

### TestAddTask
- Проверяет функцию `AddTask`, которая должна добавлять задачу в базу данных.
- Создает мок базы данных и оркестратор с этим моком.
- Устанавливает ожидания для запроса к базе данных.
- Добавляет задачу и проверяет возвращаемый идентификатор.

### TestGetTasks
- Проверяет функцию `GetTasks`, которая должна возвращать список задач из базы данных.
- Создает мок базы данных и оркестратор с этим моком.
- Устанавливает ожидания для запроса к базе данных.
- Получает список задач и проверяет их количество.

### TestAddTaskForUser
- Проверяет функцию `AddTaskForUser`, которая должна добавлять задачу для определенного пользователя в базу данных.
- Создает мок базы данных и оркестратор с этим моком.
- Устанавливает ожидания для запросов к базе данных.
- Добавляет задачу для пользователя и проверяет возвращаемый идентификатор.

### TestGetUserByLogin
- Проверяет функцию `GetUserByLogin`, которая должна возвращать пользователя по его логину из базы данных.
- Создает мок базы данных и оркестратор с этим моком.
- Устанавливает ожидания для запроса к базе данных.
- Получает пользователя и проверяет его наличие.

### TestDeleteAllTasksForUser
- Проверяет функцию `DeleteAllTasksForUser`, которая должна удалять все задачи пользователя из базы данных.
- Создает мок базы данных и оркестратор с этим моком.
- Устанавливает ожидания для запросов к базе данных.
- Удаляет все задачи пользователя и проверяет выполнение ожиданий.

### TestGetTasksForUser
- Проверяет функцию `GetTasksForUser`, которая должна возвращать список задач для определенного пользователя из базы данных.
- Создает мок базы данных и оркестратор с этим моком.
- Устанавливает ожидания для запроса к базе данных.
- Получает список задач для пользователя и проверяет его количество.