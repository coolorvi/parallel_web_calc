# parallel_web_calc
Этот проект - распределённый вычислитель арифметических выражений. 

![image](https://github.com/user-attachments/assets/26249ab0-d21f-4823-bf80-ce22e163f15b)


## Порядок работы:
* Сначала пользователь должен зарегистрироваться, отправив запрос на эндпоинт /api/v1/register с логином и паролем. А после регистрации войти под своим логином и паролем по эндпоинту /api/v1/login и получить JWT токен для последующего использования при запросах на другие эндпоинты.
* При отправке выражения на эндпоинт api/v1/calculate выражение сохраняется в базу данных, парсится на подвыражения и создаются таски, а пользователю возвращается id выражения(calculate_handler). Далее хэндлер GetTask поочередно отправляет таски на эндпоинт internal/task(task_handler). Через этот эндпоинт агент получает по одному таску и добавляет их в канал jobs, после получения всех тасков агент создает определенное количество воркеров, которые получают таск через канал, вычисляют его и отправляют результат в канал results, агент достает результаты из этого канала и отправляет оркестратору по эндпоинту internal/task(agent). Хэндлер PostResult слушает этот эндпоинт, получает результат, добавляет в базу данных, и если все таски выполнены, меняет статус выражения и добавляет итоговый результат.
* При запросе на эндпоинт api/v1/expressions возвражается список всех выражений.
* При запроосе на эндпоинт api/v1/expressions/{id} возвращается выражение с конкретным id или же ошибка, если оно не найдено.

## Запуск:
1. Клонировать репозиторий с помощью команды:\
__git clone https://github.com/coolorvi/parallel_web_calc__
2. Перейти в папку проекта и запустить проект командой: __go run ./cmd/main.go__

## Примеры использования:
_Пользователь может отправить запрос на пять эндпоинтов: /api/v1/register, /api/v1/login, api/v1/calculate, api/v1/expressions и api/v1/expressions{id}_:
* Запросы на эндпоинт /api/v1/register:
 * ``` 
    curl --location 'http://localhost:8080/api/v1/register' \
    --header 'Content-Type: application/json' \
    --data '{"login": "1", "password": "1"}'
   ```
  Ответ(код 409): {"error":"User already exists"}
* ```
  curl --location 'http://localhost:8080/api/v1/register' \
  --header 'Content-Type: application/json' \
  --data '{"login": "13", "password": "13"}'
  ```
  Ответ(код 200): {"status":"registered"}

* Запросы на эндпоинт /api/v1/login:
 * ```
    curl --location 'http://localhost:8080/api/v1/login' \
    --header 'Content-Type:   application/json' \
    --data '{"login": "1", "password": "1"}'
    ```
  Ответ(код 200): {"token":"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoyLCJleHAiOjE3NDY5Njk5NTN9.gzKsaQkhXVATGZ1IWiys-MfdDqUzs8gSooFGWqbXe5k"}
 * ```
    curl --location 'http://localhost:8080/api/v1/login' \
   --header 'Content-Type: application/json' \
    --data '{"login": "22", "password": "1"}'
    ```
    Ответ(код 401): {"error":"Invalid login or password"}  

* Запросы на эндпоинт api/v1/calculate:
 * ```
    curl --location 'http://localhost:8080/api/v1/calculate' \
    --header 'Content-Type: application/json' \
    --data '{"expression": "45 + 3423"}'

   ```
   Ответ(код 201) {
    "id": "c7ed3459-9378-4ba7-be7c-b5cbbf87f221"
}
    
 *  ```
    curl --location 'http://localhost:8080/api/v1/calculate' \
    --header 'Content-Type: application/json' \
    --data '{"expression": "&&"}'
    ```

    Ответ(код 500) Fail to parse

* Запросы на api/v1/expressions
 * ```
   curl --location 'http://localhost:8080/api/v1/expressions'
   ```
   Ответ(код 200) {
    "expressions": [
        {
            "id": "c7ed3459-9378-4ba7-be7c-b5cbbf87f221",
            "status": "in_progress"
        }
    ]
  }

 * ```
   curl --location 'http://localhost:8080/api/v1/expressions'
   ```
   Ответ(код 500) {
    "expressions": null
  }   

* Запросы на api/v1/expressions/{id}:
 * ```
   curl --location 'http://localhost:8080/api/v1/expressions/3d576965-f5be-406a-84fc-982a1e0ff4de'
   ```
   Ответ(код 200) {
    "expression": {
        "id": "3d576965-f5be-406a-84fc-982a1e0ff4de",
        "status": "in_progress"
    }
    }
 * ```
   curl --location 'http://localhost:8080/api/v1/expressions/3dd'
   ```
   Ответ(код 404) Not Found
     
