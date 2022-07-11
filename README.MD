Хранилище паролей с двумя API, telegram bot + rest api

Что позволяет:
1) Позволяет прочитать пароль по логину
2) Позволяет создать новую пару логин - пароль
3) Есть небольная аутентификация (реализована по простому base auth = login + pass) (чисто имитация, без всяких jwt токенов =) ) 
4) В параллели при записи в БД отбрасывается лог чере горутину

Информация о боте
Его ник DataDomovenok постучитесь к нему или воспользуйтесь api,
спека лежит в /api, рутовый адрес при старте /api/v1


Инфа о постгрюхи
Перед накатом поставьте docker и закиньте image postgre
docker pull postgres

далее run
docker run --name some-postgres -e POSTGRES_PASSWORD=12345 -d postgres

Накат БД
migrate -path resources/db/migration -database "postgresql://admin:12345@localhost:5432/postgres" -verbose up

Запуск
генерим модели oapi-codegen -package spec api/openapi.yml > internal/pkg/models/gen.go


Openapi лежит в password-storage/api/openapi.yml