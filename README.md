# Завдання

- N = 24
- V1 = (24 % 2) + 1 = 1
- V2 = (24 % 3) + 1 = 1
- V3 = (24 % 5) + 1 = 5

* Тематика веб застосунку: Notes Service 
* Спосіб конфігурації: аргументи командного рядка.
* Порт: 5000

## Notes Service 

Простий сервіс для зберігання текстових нотаток.

# Технічний стек
- Golang 
- Gin (web framework)
- Gorm (orm library)

# Запуск тестів

```bash
go test ./...
```

# Збірка додатку та запуск

```
go mod tidy
```

```
go build -o app ./cmd/app/main.go
```

```
go run main.go
```

# Документація
- [Ендпоінти](./docs/api.md)
- [Розгортання](./docs/deploy.md)
