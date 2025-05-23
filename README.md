# Kanban board

### Функционал
- Создание, редактирование и удаление задач
- Организация задач по колонкам (To Do, In Progress, Done)
- Перемещение задач между колонками
- Автоматическое добавление временных меток (создания и обновления задач)
### Технологии
- **Backend:** Go Gin
- **Frontend:** JS React
- **База данных:** PostgreSQL

### REST API

**POST   /auth/register**
*регистрация нового пользователя*
запрос:
```
{
  "username": "John Doe",
  "password": "qwerty123"
}
```
ответ:
```
{ "token": "<jwt>" }
```

**POST   /auth/login**
*вход в систему*
запрос:
```
{
  "username": "John Doe",
  "password": "qwerty123"
}
```
ответ:
```
{ "token": "<jwt>" }
```

**POST   /boards**
*создание доски*
запрос:
```
{ "name": "New Board" }
```

**GET    /boards**
*информация о всех досках пользователя*
ответ:
```
[
  { 
   "id": <uuid>, 
   "user_id": <uuid>, 
   "created_at": "...", 
   "updated_at": "...",
   "name": "Work"
  },
  { 
   "id": <uuid>, 
   "user_id": <uuid>,
   "name": "Personal", 
   ... 
  },
  ...
]
```

**GET    /boards/:id** 
*получение метаданных о конкретной доске*
ответ:
```
{
  "id": <uuid>,
  "user_id": <uuid>,
  "created_at": "...",
  "updated_at": "...",
  "name": "Work Board"
}

```

**PUT    /boards/:id**
*обновление названия доски*
запрос:
```
{ "name": "Renamed Board" }
```

**DELETE /boards/:id**
*удаление доски (и всего содержимого)*

**POST   /boards/:id/columns**
*создание колонки*
запрос:
```
{ "name": "Backlog" }
```

**GET    /boards/:id/columns**
*получение всех колонок конкретной доски*
ответ:
```
[
  { 
   "id": <uuid>, 
   "board_id": <uuid>,
   "created_at": "...",
   "updated_at": "...",
   "name": "To Do", 
   "position": 1 
  },
  {
   "id": <uuid>, 
   "board_id": <uuid>,
   "created_at": "...",
   "updated_at": "...",
   "name": "In Progress", 
   "position": 2 
  },
  ...
]
```

**GET /columns/:id**
*получение информации о конкретной колонке*
ответ:
```
{ 
 "id": <uuid>, 
 "board_id": <uuid>,
 "created_at": "...",
 "updated_at": "...",
 "name": "Backlog", 
 "position": 1 
}
```

**PATCH    /columns/:id**
*переименование и/или перемещение колонки*
запрос:
```
{
  "name": "Done",
  "position": "2"
}
```
*одно из полей может быть опущено, в таком случае оно просто не обновится
также могут быть опущены оба поля, но тогда единственное, что обновится - это поле `updated_at`*

**DELETE /columns/:id**
*удаление колонки и всех задач в ней*

**POST   /columns/:id/tasks**
*создание задачи*
запрос:
```
{ "name": "New task", "description": "..." }
```

**GET    /columns/:id/tasks**
*получение всех задач в конкретной колонке*
ответ:
```
[
  {
   "id": <uuid>,
   "column_id": <uuid>,
   "created_at": "...",
   "updated_at": "...",
   "name": "Fix bug", 
   "description": "Something",
   "position": "1",
   "done": false,
   "deadline": "2025-06-01T12:00:00Z" 
  },
  {
   "id": <uuid>,
   "column_id": <uuid>,
   "created_at": "...",
   "updated_at": "...",
   "name": "Fix bug", 
   "description": "Something",
   "position": "2",
   "done": false,
   "deadline": null
  },
  ...
]
```

**GET    /tasks/:id**
*получение информации о конкретной задаче*
ответ:
```
{
 "id": <uuid>,
 "column_id": <uuid>,
 "created_at": "...",
 "updated_at": "...",
 "name": "Do something", 
 "description": "Something",
 "position": 52,
 "done": false,
 "deadline": "2025-06-01T12:00:00Z" 
}
```

**PATCH    /tasks/:id**
*редактирование или перемещение задачи*
возможные запросы:
1. Обновление контента
*изменение имени*
```
{
  "name": "Updated name"
}
```
*изменение описания*
```
{
  "description": "Updated description"
}
```
*сделано/не сделано*
```
{
  "done": true,
}
```
*установка или изменение дедлайна*
```
{
  "deadline": "2025-06-01T12:00:00Z" 
}
```
2. Изменение позиции в колонке
```
{
  "position": 13
}
```
3. Изменение колонки
```
{
  "column_id": <uuid>,
  "position": 6
}
```

***!!! Все остальные комбинации полей в запросе вернут ошибку 400***

**DELETE /tasks/:id**
*удаление задачи*

### PostgreSQL
```
TABLE "user"(
  id uuid PRIMARY KEY,
  created_at timestamptz NOT NULL,
  username text NOT NULL UNIQUE,
  hashed_password text NOT NULL
);
```
```
TABLE board(
  id uuid PRIMARY KEY,
  user_id uuid REFERENCES "user"(id) ON DELETE CASCADE,
  created_at timestamptz NOT NULL,
  updated_at timestamptz NOT NULL,
  name text NOT NULL
);
```
```
TABLE "column"(
  id uuid PRIMARY KEY,
  board_id uuid REFERENCES "board"(id) ON DELETE CASCADE,
  created_at timestamptz NOT NULL,
  updated_at timestamptz NOT NULL,
  name text NOT NULL,
  position smallint NOT NULL
);
```
```
TABLE task(
  id uuid PRIMARY KEY,
  column_id uuid REFERENCES "column"(id) ON DELETE CASCADE,
  created_at timestamptz NOT NULL,
  updated_at timestamptz NOT NULL,
  name text NOT NULL,
  description text NOT NULL,
  position smallint NOT NULL,
  done boolean NOT NULL DEFAULT false,
  deadline timestamptz
);
```

### Структура проекта
```
/cmd
  main.go  # точка входа

/internal
 /config
   config.go  # конфигурация, env

 /server
   router.go  # конфигурация gin сервера, подключение всех роутов и middleware
   middleware.go  # глобальный middleware

 /auth
   handler.go  # ручки для авторизации
   service.go  # валидация, хэширование пароля
   middleware.go  # проверка авторизации
   jwt.go  # генерация, валидация, работа с JWT
   model.go  # модель для базы данных
   repo.go  # работа с базой данных

 /board
   handler.go  # ручки для досок
   model.go   # модель для базы данных
   repo.go  # работа с базой данных

 /column
   handler.go  # ручки для колонок
   model.go   # модель для базы данных
   repo.go  # работа с базой данных

 /task
   handler.go  # ручки для задач
   model.go   # модель для базы данных
   repo.go  # работа с базой данных

 /postgres
   postgres.go  # подключение к бд и инициализация
   queries.go  # сырые sql запросы

 /utils
   utils.go  # утилиты, используемые в разных пакетах
```
