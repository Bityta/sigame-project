# Руководство по оптимизации E2E тестов

## Уже выполнено

1. ✅ **Увеличено количество workers в CI**: с 1 до 4 (параллелизация тестов)
2. ✅ **Оптимизирована запись видео**: `mode: 'on'` → `mode: 'retain-on-failure'` (видео только при ошибках)

## Ожидаемый эффект

- **Параллелизация**: тесты будут выполняться в 4 потока вместо 1 → ускорение в ~3-4 раза
- **Видео**: запись видео только при ошибках → экономия времени и места

## Дальнейшие оптимизации

### 1. Убрать фиксированные задержки (105 вызовов `waitForTimeout`)

**Проблема**: Фиксированные задержки замедляют тесты и делают их нестабильными.

**Решение**: Заменить на умные ожидания:

```typescript
// ❌ Плохо
await page.waitForTimeout(1000);

// ✅ Хорошо
await page.waitForLoadState('domcontentloaded');
await expect(element).toBeVisible();
```

**Примеры замены**:

```typescript
// Вместо:
await page.waitForTimeout(500);
await button.click();

// Использовать:
await button.waitFor({ state: 'visible' });
await button.click();

// Вместо:
await page.waitForTimeout(2000);
await expect(element).toBeVisible();

// Использовать:
await expect(element).toBeVisible({ timeout: 2000 });
```

### 2. Убрать искусственные задержки при действиях (13 вызовов `delay:`)

**Проблема**: `delay: 50` и `delay: 100` замедляют тесты без необходимости.

**Решение**: Убрать delay, Playwright сам управляет действиями:

```typescript
// ❌ Плохо
await input.fill(text, { delay: 50 });
await button.click({ delay: 100 });

// ✅ Хорошо
await input.fill(text);
await button.click();
```

### 3. Оптимизировать избыточные ожидания

**Проблема**: Комбинации `waitForLoadState` + `waitForTimeout` избыточны.

**Примеры оптимизации**:

```typescript
// ❌ Плохо
await page.waitForLoadState('networkidle');
await page.waitForTimeout(1000);

// ✅ Хорошо
await page.waitForLoadState('domcontentloaded');
await expect(element).toBeVisible();

// ❌ Плохо
await element.waitFor({ state: 'visible' });
await page.waitForTimeout(500);
await element.click({ delay: 100 });

// ✅ Хорошо
await element.waitFor({ state: 'visible' });
await element.click();
```

### 4. Оптимизировать таймауты

**Рекомендации**:
- Использовать разумные таймауты (5-10 секунд для большинства операций)
- Увеличивать таймауты только для действительно медленных операций (WebSocket, загрузка медиа)
- Использовать `waitForFunction` вместо длинных `waitForTimeout`

### 5. Использовать более эффективные стратегии ожидания

**Приоритет стратегий**:
1. `expect().toBeVisible()` - самое быстрое и надежное
2. `waitForSelector()` - для элементов, которые появляются
3. `waitForURL()` - для навигации
4. `waitForLoadState()` - для загрузки страницы
5. `waitForFunction()` - для сложных условий
6. `waitForTimeout()` - только в крайних случаях

## План действий

### Фаза 1: Быстрые победы (высокий эффект, низкая сложность)
1. ✅ Увеличить workers в CI
2. ✅ Оптимизировать запись видео
3. Убрать все `delay:` из действий (13 вызовов)
4. Убрать избыточные `waitForTimeout` после `waitForLoadState`

### Фаза 2: Оптимизация helpers (средний эффект, средняя сложность)
1. Оптимизировать `helpers/auth.ts` (много `waitForTimeout`)
2. Оптимизировать `helpers/room.ts` (много `waitForTimeout`)
3. Оптимизировать `helpers/game.ts` (длинные таймауты)

### Фаза 3: Оптимизация тестов (низкий эффект, высокая сложность)
1. Постепенно оптимизировать отдельные тесты
2. Заменять фиксированные задержки на умные ожидания

## Метрики для отслеживания

- Время выполнения всех тестов
- Количество `waitForTimeout` вызовов
- Количество `delay:` параметров
- Стабильность тестов (процент успешных прогонов)

## Ожидаемый результат

После полной оптимизации:
- **Ускорение в 3-5 раз** за счет параллелизации
- **Дополнительное ускорение на 20-30%** за счет оптимизации ожиданий
- **Повышение стабильности** за счет умных ожиданий вместо фиксированных задержек

