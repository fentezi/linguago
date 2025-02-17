# Приложение-переводчик

Приложение для перевода слов с использованием **Google Translate API** и озвучивания через **Elevenlabs.io**.

## Особенности
- **Автосохранение слов:** При каждом переводе слово автоматически добавляется в базу данных, формируя ваш персональный словарь.
- **Сохраненные слова:** Просматривайте и управляйте своими переводами в разделе "Сохраненные слова".
- **Добавление вручную:** Возможность самостоятельно добавлять слова и их переводы.

---

## Настройка переменных окружения

Создайте файл `.env` в корне проекта и добавьте в него следующие переменные:

```env
POSTGRES_HOST="ваш_адрес_хоста_PostgreSQL"
POSTGRES_PASSWORD="пароль_к_PostgreSQL"
API_KEY="ваш_API_ключ_Google_Translate"
REDIS_HOST="ваш_адрес_хоста_Redis"
REDIS_PASSWORD="пароль_к_Redis"
```

## Запуск проекта

### С использованием Docker

1. Убедитесь, что Docker установлен на вашем устройстве.
2. В терминале выполните следующую команду:
```bash
docker compose up
```

### Локальный запуск без Docker
1. Убедитесь, что все зависимости установлены.
2. Запустите проект командой:
```bash
make
```

### Стек технологий
- Go
- PostgreSQL
- Redis
- Google Translate API
- Elevenlabs.io API

## Лицензия
[LICENSE](./LICENSE) 