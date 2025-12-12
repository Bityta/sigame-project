# Troubleshooting: 401 Unauthorized на /api/game/my-active

## Проблема
Эндпоинт `/api/game/my-active` возвращает 401 Unauthorized.

## Диагностика

### 1. Проверка логов Game Service

При старте сервиса проверьте логи:

```bash
# Должно быть:
Connected to Auth Service at <address>

# Если видите:
Failed to connect to Auth Service: <error>
Endpoints requiring authentication will only accept X-User-ID header
```

### 2. Проверка передачи токена из Frontend

В браузере (DevTools → Network):
1. Откройте запрос к `/api/game/my-active`
2. Проверьте заголовки запроса:
   - Должен быть: `Authorization: Bearer <token>`
   - Или fallback: `X-User-ID: <uuid>`

### 3. Проверка токена в localStorage

В консоли браузера:
```javascript
localStorage.getItem('access_token')
// или
sessionStorage.getItem('access_token')
```

### 4. Логи Game Service при запросе

При запросе к `/api/game/my-active` проверьте логи:

**Если токен передан:**
```
Token validated successfully, user_id=<uuid>
```

**Если токен не передан:**
```
No Authorization header found
Using X-User-ID header, user_id=<uuid>
```

**Если токен невалиден:**
```
Token is invalid: <error message>
```

**Если Auth client не инициализирован:**
```
Auth client not initialized, cannot validate token
```

## Решения

### Решение 1: Auth Service недоступен

Если Auth Service не запущен или недоступен:
- Запустите Auth Service
- Проверьте переменные окружения: `AUTH_SERVICE_HOST`, `AUTH_SERVICE_PORT`
- Используйте заголовок `X-User-ID` для внутренних вызовов

### Решение 2: Токен не передается из Frontend

Проверьте код frontend, который делает запрос:
```typescript
// Должно быть:
headers: {
  'Authorization': `Bearer ${token}`
}
```

### Решение 3: Токен невалиден

- Проверьте срок действия токена
- Проверьте, что токен не был отозван
- Попробуйте получить новый токен через Auth Service

## Конфигурация

Переменные окружения для Auth Service:
```bash
AUTH_SERVICE_HOST=localhost
AUTH_SERVICE_PORT=50051
```

