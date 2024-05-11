## Сервис отслеживания посылок со следующими функциями:
регистрация посылки,  
получение списка посылок клиента,  
изменение статуса посылки,  
изменение адреса доставки,  
удаление посылки.  
План такой: информация о посылках хранится в БД. Посылка может быть зарегистрирована, отправлена или доставлена. При регистрации посылки создаётся новая запись в БД. У только что зарегистрированной должен быть статус «зарегистрирована». Трек-номер посылки равен её идентификатору в таблице. Если посылка в статусе «зарегистрирована», можно изменить адрес доставки или удалить посылку.  
В качестве СУБД используется SQLite. Файл с БД называется tracker.db. В БД всего одна таблица parcel со следующими колонками:  
_number_ — номер посылки, целое число, автоинкрементное поле.  
_client_ — идентификатор клиента, целое число.  
_status_ — статус посылки, строка.  
_address_ — адрес посылки, строка.  
_created_at_ — дата и время создания посылки, строка.  