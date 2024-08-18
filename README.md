## Сервис домов

### Запуск сервиса:

Для запуска необходимо
```bash
        git clone 
        cd ./
		make run
```  
после этого приложение доступно через 
http://localhost:8000

**make run**  запускает docker контейнеры с приложением и БД

**docker compose** лежат в папке ```./deploy```
 
Секрет токен авторизации сейчас лежит в файле ```.env```, данная реализация применима только к тестовому проекту, для реальных проектов понимаю, что такие данные не лежат в общем доступе

### Реализованный функционал сервиса:
Функционал отвечает на запросы описанные в [API](https://github.com/avito-tech/backend-bootcamp-assignment-2024/blob/main/api.yaml).
1. **Авторизация пользователей:**
	Регистрация и авторизация пользователей по почте и паролю:
	``` endpoint/register ```
	Принимает запрос в виде 
    ```json
    {
	"email": "test@gmail.com",
	"password": "Секретная строка",
	"user_type": "client"
	}
	```
    по результатам регистрации возвращается ``` uuid ```, формируется UUID V7.
	В базе создаётся и сохраняется новый пользователь желаемого типа: обычный пользователь (client) или модератор (moderator), который позволяет пройти аутентификацию и получить токен.

	У созданного пользователя появляется токен endpoint/login. Авторизация доступна в двух вариантах 
	UUID + пароль пример запроса
	```json
    {
	"id": "cae36e0f-69e5-4fa8-a179-a52d083c5549",
	"password": "Секретная строка",
	}
	```

	email+пароль
	 ```json
    {
	"email": "test@gmail.com",
	"password": "Секретная строка",
	}
	```

	При успешной авторизации по почте и паролю возвращается токен для пользователя с соответствующим ему уровнем доступа.

	**ручка /dummyLogin - не реализована** сервис работает через описанную выше регистрацию и авторизацию

	**БЕЗ РЕГИСТРАЦИИ И АУТЕНТИФИКАЦИИ ДАЛЬНЕЙШИЕ ДЕЙСТВИЯ НЕ ВОЗМОЖНЫ**

2. **Создание дома:**
```реализовано middleware для проверки авторизации```
	1. Только модератор имеет возможность создать дом используя endpoint /house/create.
	В случае успешного запроса возвращается полная информация о созданном доме
3. **Создание квартиры:**
	1. Создать квартиру может любой пользователь, используя endpoint /flat/create. При успешном запросе возвращается полная информация о квартире.
	2. Если жильё успешно создано через endpoint /flat/create, то объявление получает статус модерации created.
	3. У дома, в котором создали новую квартиру, обновляется дата последнего добавления жилья. 
4. **Модерация квартиры:**
	1. Статус модерации квартиры может принимать одно из четырёх значений: created, approved, declined, on moderation.
	2. Только модератор может изменить статус модерации квартиры с помощью endpoint /flat/update. При успешном запросе возвращается полная информация об обновленной квартире.
5. **Получение списка квартир по номеру дома:**
	1. Используя endpoint /house/{id}, обычный пользователь и модератор могут получить список квартир по номеру дома. Только обычный пользователь увидит все квартиры со статусом модерации approved, а модератор — жильё с любым статусом модерации.

**Связи между сущностями:**
1. Каждая квартира может иметь только одно соответствие с домом (один к одному).
2. Номер дома служит уникальным идентификатором самого дома.
``` По пункту 3 принципильно понятен смысл, но в сервис не передаются номера квартир, и как мне показалось в рамках одного дома всё таки номер квартиры должен быть уникален, поэтому номер квартиры формируется в БД как SERIAL PRIMARY KEY ```
3. Номер квартиры не является уникальным идентификатором. Например, квартира №1 может находиться как в доме №1, так и в доме №2, и в этом случае это будут разные квартиры.


- Модераторы — получают полный список всех объявлений в доме вне зависимости от статуса модерации.
- Пользователи — получают список только прошедших модерацию объявлений. 

### еще не реализовано
6. **Дополнительно.** Подписка на уведомления о новых квартирах в доме по его номеру. Обычный пользователь может подписаться на такие уведомления с помощью endpoint /house/{id}/subscribe.
