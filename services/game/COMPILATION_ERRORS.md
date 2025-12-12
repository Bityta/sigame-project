# Ошибки компиляции требующие исправления

## Найдено проблем:

### 1. message/handler.go - отсутствующие типы
- `Hub`, `ClientMessageWrapper`, `sendErrorMessage`, `ErrorGameManagerNotFound`, `MaxRTTDuration`
- **Причина:** Эти типы и функции должны быть в пакете message или импортированы

### 2. hub/registry.go - неправильный type assertion
- Строка 46: `gameID.(interface{})` - нужно `gameID.(uuid.UUID)`

### 3. postgres repositories - undefined domain
- **Проблема:** использование `domain.X` вместо конкретных импортов
- **Решение:** добавить импорты `domainGame`, `domainPlayer` etc.

### 4. grpc/pack/pack_client.go - синтаксическая ошибка
- Строка 29: newline in string - незавершенная строка

### 5. ws/client/client.go - отсутствующие константы и типы
- `Hub`, `MaxMessageSize`, `PongWait`, `NewClientMessage`, `sendErrorMessage`, etc.
- **Причина:** Константы остались в старом пакете websocket

### 6. ws/handler - отсутствующие константы и функции
- `Hub`, `ErrorInvalidGameID`, `QueryParamUserID`, `NewClient`, etc.

### 7. redis repositories - unused imports + undefined domain
- Импорты game, pack не используются напрямую
- Использование `domain.X` вместо конкретных типов

## Общие проблемы:

1. **Константы не перенесены** - многие константы из ws/constants.go не используются
2. **Типы разбросаны** - Hub, Client, Message в разных пакетах
3. **Helper функции** - sendErrorMessage, NewClientMessage не найдены
4. **Domain импорты** - старый pattern `domain.X` вместо `domainGame.X`

## Рекомендации:

1. Проверить ws/constants.go - все ли константы экспортированы
2. Проверить ws/errors.go - все ли helper функции на месте
3. Обновить все domain импорты в postgres/redis repositories
4. Исправить синтаксис в pack_client.go
5. Добавить переэкспорт констант в ws/ws.go

## Статус:

❌ Код не компилируется
⏳ Требуется ручное исправление всех undefined references
✅ Структура правильная, нужно только связать все части

