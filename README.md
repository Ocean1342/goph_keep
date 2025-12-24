## Общая логика работы консольного приложения
* Приложение требует логин на каждый запрос
* синхронизация запускается при запуске консольного приложения

## Конфигурация
 берётся через билд тег configPath, где configPath - путь до json файла конфигурации
 
## Сборка 
для сборки консольного клиента под ос:
```shell
#windows 
GOOS=windows GOARCH=amd64 go build -o app.exe
#linux
GOOS=linux GOARCH=amd64 go build -o app
```

## Команды
запуск команды с флагом --help даёт расширенную информацию о используемых флагах.

### Общие флаги
```shell
  -l, --login string   login name
  -m, --meta string    metadata for stored data
  -n, --name string    data name
```
### buildData
выводит информацию о тегах при билде приложения 

### storeBin
ограничение на размер бинарного файла 5мб
```shell
storeBin --login=user_logo --name=unique_data_label \
--path=./path/to/bin_file --meta=data_description 
```

### getBin
```shell
getBin --login=login-data --name=test --binPathTo="./"  
```

### storeLogoPass
```shell
storeLogoPass --login=user_logo --name=secret \ 
--storeLogo=my-secret --storePass=my-secret-pass --meta=secret-meta
```
### getLogoPass
```shell
getLogoPass --login=login-data --name=secret 
```

## регистрация/авторизация
указываем пароль логин и фразу. всё храним на сервере

### gen docs
```bash
swag init -d ./server/cmd/,./internal/storage,./internal/api/,./cli -o ./docs --ot yaml
```

