/**
 * E2E тесты комнат
 * 
 * Тестируют полные пользовательские сценарии работы с комнатами:
 * - Создание комнаты
 * - Вход в комнату по коду
 * - Копирование кода комнаты
 * - Настройки комнаты
 * - Выход из комнаты
 */

import { test, expect, Page } from '@playwright/test';

// Генератор уникальных имен
const generateUsername = () => `player_${Date.now()}`;
const generateRoomName = () => `Room_${Date.now()}`;

// Хелпер для авторизации
async function login(page: Page, username?: string, password = 'testpassword123') {
  const user = username || generateUsername();
  
  await page.goto('/register');
  await page.getByLabel(/имя пользователя/i).fill(user);
  await page.getByLabel(/пароль/i).fill(password);
  await page.getByRole('button', { name: /зарегистрироваться/i }).click();
  
  await expect(page).toHaveURL('/lobby', { timeout: 10000 });
  
  return user;
}

test.describe('Комнаты', () => {
  /**
   * ТЕСТ: Отображение страницы создания комнаты
   * 
   * Проверяет что страница создания комнаты загружается
   * и содержит все необходимые элементы
   */
  test('страница создания комнаты загружается корректно', async ({ page }) => {
    await login(page);
    
    await page.getByRole('button', { name: /создать комнату/i }).click();
    
    await expect(page.getByRole('heading', { name: /создать комнату/i })).toBeVisible();
    await expect(page.getByLabel(/название комнаты/i)).toBeVisible();
    await expect(page.getByLabel(/пак вопросов/i)).toBeVisible();
    await expect(page.getByRole('button', { name: /создать/i })).toBeVisible();
    await expect(page.getByRole('button', { name: /отмена/i })).toBeVisible();
  });

  /**
   * ТЕСТ: Отмена создания комнаты возвращает в лобби
   * 
   * Проверяет что кнопка "Отмена" на странице создания
   * возвращает пользователя в лобби
   */
  test('отмена создания комнаты возвращает в лобби', async ({ page }) => {
    await login(page);
    
    await page.getByRole('button', { name: /создать комнату/i }).click();
    await expect(page).toHaveURL(/create/);
    
    await page.getByRole('button', { name: /отмена/i }).click();
    
    await expect(page).toHaveURL('/lobby');
  });

  /**
   * ТЕСТ: Валидация короткого названия комнаты
   * 
   * Проверяет что при вводе слишком короткого названия
   * показывается ошибка валидации
   */
  test('ошибка при коротком названии комнаты', async ({ page }) => {
    await login(page);
    
    await page.getByRole('button', { name: /создать комнату/i }).click();
    
    await page.getByLabel(/название комнаты/i).fill('ab');
    await page.getByLabel(/пак вопросов/i).selectOption({ index: 1 }); // Выбираем первый пак
    await page.getByRole('button', { name: /создать/i }).click();
    
    await expect(page.getByText(/название должно быть от 3/i)).toBeVisible();
  });

  /**
   * ТЕСТ: Валидация невыбранного пака
   * 
   * Проверяет что при попытке создать комнату без пака
   * показывается ошибка валидации
   */
  test('ошибка при невыбранном паке вопросов', async ({ page }) => {
    await login(page);
    
    await page.getByRole('button', { name: /создать комнату/i }).click();
    
    await page.getByLabel(/название комнаты/i).fill('Моя комната');
    // Не выбираем пак
    await page.getByRole('button', { name: /создать/i }).click();
    
    await expect(page.getByText(/выберите пак/i)).toBeVisible();
  });

  /**
   * ТЕСТ: Создание публичной комнаты
   * 
   * Полный сценарий создания комнаты:
   * - Заполнение формы
   * - Создание
   * - Редирект в комнату
   * - Отображение данных комнаты
   */
  test('успешное создание публичной комнаты', async ({ page }) => {
    await login(page);
    
    const roomName = generateRoomName();
    
    await page.getByRole('button', { name: /создать комнату/i }).click();
    
    await page.getByLabel(/название комнаты/i).fill(roomName);
    await page.getByLabel(/пак вопросов/i).selectOption({ index: 1 });
    await page.getByRole('button', { name: /создать/i }).click();
    
    // Ждем редирект в комнату
    await expect(page).toHaveURL(/\/room\//, { timeout: 10000 });
    
    // Проверяем что отображается название комнаты
    await expect(page.getByRole('heading', { name: roomName })).toBeVisible();
    
    // Проверяем что отображается код комнаты
    const codeElement = page.locator('.room-page__code-value');
    await expect(codeElement).toBeVisible();
  });

  /**
   * ТЕСТ: Отображение кода комнаты
   * 
   * Проверяет что после создания комнаты
   * отображается её код (6 символов)
   */
  test('отображение кода комнаты после создания', async ({ page }) => {
    await login(page);
    
    await page.getByRole('button', { name: /создать комнату/i }).click();
    await page.getByLabel(/название комнаты/i).fill(generateRoomName());
    await page.getByLabel(/пак вопросов/i).selectOption({ index: 1 });
    await page.getByRole('button', { name: /создать/i }).click();
    
    await expect(page).toHaveURL(/\/room\//, { timeout: 10000 });
    
    // Код должен быть из 6 символов
    const codeElement = page.locator('.room-page__code-value');
    await expect(codeElement).toBeVisible();
    const code = await codeElement.textContent();
    expect(code).toMatch(/^[A-Z0-9]{6}$/);
  });

  /**
   * ТЕСТ: Копирование кода комнаты по клику
   * 
   * Проверяет что при клике на блок с кодом
   * код копируется в буфер обмена
   */
  test('копирование кода комнаты по клику', async ({ page, context }) => {
    // Разрешаем доступ к clipboard
    await context.grantPermissions(['clipboard-read', 'clipboard-write']);
    
    await login(page);
    
    await page.getByRole('button', { name: /создать комнату/i }).click();
    await page.getByLabel(/название комнаты/i).fill(generateRoomName());
    await page.getByLabel(/пак вопросов/i).selectOption({ index: 1 });
    await page.getByRole('button', { name: /создать/i }).click();
    
    await expect(page).toHaveURL(/\/room\//, { timeout: 10000 });
    
    // Получаем код из элемента
    const codeElement = page.locator('.room-page__code-value');
    const code = await codeElement.textContent();
    
    // Кликаем на блок с кодом
    await page.locator('.room-page__code').click();
    
    // Проверяем что код скопирован
    const clipboardText = await page.evaluate(() => navigator.clipboard.readText());
    expect(clipboardText).toBe(code);
  });

  /**
   * ТЕСТ: Визуальная обратная связь при копировании
   * 
   * Проверяет что после клика на код
   * появляется визуальный индикатор копирования
   */
  test('визуальная обратная связь при копировании кода', async ({ page, context }) => {
    await context.grantPermissions(['clipboard-read', 'clipboard-write']);
    await login(page);
    
    await page.getByRole('button', { name: /создать комнату/i }).click();
    await page.getByLabel(/название комнаты/i).fill(generateRoomName());
    await page.getByLabel(/пак вопросов/i).selectOption({ index: 1 });
    await page.getByRole('button', { name: /создать/i }).click();
    
    await expect(page).toHaveURL(/\/room\//, { timeout: 10000 });
    
    const codeBlock = page.locator('.room-page__code');
    await codeBlock.click();
    
    // Проверяем что появился класс copied
    await expect(codeBlock).toHaveClass(/room-page__code--copied/);
  });

  /**
   * ТЕСТ: Хост видит себя в списке игроков
   * 
   * Проверяет что создатель комнаты видит себя
   * в списке игроков с короной
   */
  test('хост отображается в списке игроков', async ({ page }) => {
    const username = await login(page);
    
    await page.getByRole('button', { name: /создать комнату/i }).click();
    await page.getByLabel(/название комнаты/i).fill(generateRoomName());
    await page.getByLabel(/пак вопросов/i).selectOption({ index: 1 });
    await page.getByRole('button', { name: /создать/i }).click();
    
    await expect(page).toHaveURL(/\/room\//, { timeout: 10000 });
    
    // Проверяем что имя хоста отображается с короной
    await expect(page.getByText(new RegExp(`${username}.*👑`))).toBeVisible();
  });

  /**
   * ТЕСТ: Хост видит кнопку "Начать игру"
   * 
   * Проверяет что создатель комнаты видит
   * кнопку для запуска игры
   */
  test('хост видит кнопку "Начать игру"', async ({ page }) => {
    await login(page);
    
    await page.getByRole('button', { name: /создать комнату/i }).click();
    await page.getByLabel(/название комнаты/i).fill(generateRoomName());
    await page.getByLabel(/пак вопросов/i).selectOption({ index: 1 });
    await page.getByRole('button', { name: /создать/i }).click();
    
    await expect(page).toHaveURL(/\/room\//, { timeout: 10000 });
    
    await expect(page.getByRole('button', { name: /начать игру/i })).toBeVisible();
  });

  /**
   * ТЕСТ: Кнопка старта заблокирована при 1 игроке
   * 
   * Проверяет что нельзя начать игру
   * пока в комнате только один игрок
   */
  test('кнопка старта заблокирована при одном игроке', async ({ page }) => {
    await login(page);
    
    await page.getByRole('button', { name: /создать комнату/i }).click();
    await page.getByLabel(/название комнаты/i).fill(generateRoomName());
    await page.getByLabel(/пак вопросов/i).selectOption({ index: 1 });
    await page.getByRole('button', { name: /создать/i }).click();
    
    await expect(page).toHaveURL(/\/room\//, { timeout: 10000 });
    
    const startButton = page.getByRole('button', { name: /начать игру/i });
    await expect(startButton).toBeDisabled();
    
    // Проверяем подсказку
    await expect(page.getByText(/минимум 2 игрока/i)).toBeVisible();
  });

  /**
   * ТЕСТ: Выход из комнаты
   * 
   * Проверяет что при клике на "Покинуть комнату"
   * пользователь возвращается в лобби
   */
  test('выход из комнаты возвращает в лобби', async ({ page }) => {
    await login(page);
    
    await page.getByRole('button', { name: /создать комнату/i }).click();
    await page.getByLabel(/название комнаты/i).fill(generateRoomName());
    await page.getByLabel(/пак вопросов/i).selectOption({ index: 1 });
    await page.getByRole('button', { name: /создать/i }).click();
    
    await expect(page).toHaveURL(/\/room\//, { timeout: 10000 });
    
    await page.getByRole('button', { name: /покинуть/i }).click();
    
    await expect(page).toHaveURL('/lobby');
  });

  /**
   * ТЕСТ: Настройки комнаты для хоста
   * 
   * Проверяет что хост видит форму настроек
   * с возможностью редактирования
   */
  test('хост видит форму настроек комнаты', async ({ page }) => {
    await login(page);
    
    await page.getByRole('button', { name: /создать комнату/i }).click();
    await page.getByLabel(/название комнаты/i).fill(generateRoomName());
    await page.getByLabel(/пак вопросов/i).selectOption({ index: 1 });
    await page.getByRole('button', { name: /создать/i }).click();
    
    await expect(page).toHaveURL(/\/room\//, { timeout: 10000 });
    
    // Проверяем наличие секции настроек
    await expect(page.getByText(/настройки игры/i)).toBeVisible();
    
    // Проверяем наличие слайдеров
    const sliders = page.locator('input[type="range"]');
    await expect(sliders).toHaveCount(2);
    
    // Проверяем наличие кнопки сохранения
    await expect(page.getByRole('button', { name: /сохранить настройки/i })).toBeVisible();
  });

  /**
   * ТЕСТ: Кнопка сохранения настроек неактивна без изменений
   * 
   * Проверяет что кнопка "Сохранить настройки"
   * заблокирована пока настройки не изменены
   */
  test('кнопка сохранения настроек неактивна без изменений', async ({ page }) => {
    await login(page);
    
    await page.getByRole('button', { name: /создать комнату/i }).click();
    await page.getByLabel(/название комнаты/i).fill(generateRoomName());
    await page.getByLabel(/пак вопросов/i).selectOption({ index: 1 });
    await page.getByRole('button', { name: /создать/i }).click();
    
    await expect(page).toHaveURL(/\/room\//, { timeout: 10000 });
    
    const saveButton = page.getByRole('button', { name: /сохранить настройки/i });
    await expect(saveButton).toBeDisabled();
  });

  /**
   * ТЕСТ: Кнопка сохранения активируется после изменения настроек
   * 
   * Проверяет что после изменения любой настройки
   * кнопка сохранения становится активной
   */
  test('кнопка сохранения активна после изменения настроек', async ({ page }) => {
    await login(page);
    
    await page.getByRole('button', { name: /создать комнату/i }).click();
    await page.getByLabel(/название комнаты/i).fill(generateRoomName());
    await page.getByLabel(/пак вопросов/i).selectOption({ index: 1 });
    await page.getByRole('button', { name: /создать/i }).click();
    
    await expect(page).toHaveURL(/\/room\//, { timeout: 10000 });
    
    // Меняем слайдер
    const slider = page.locator('input[type="range"]').first();
    await slider.fill('45');
    
    // Кнопка должна стать активной
    const saveButton = page.getByRole('button', { name: /сохранить настройки/i });
    await expect(saveButton).not.toBeDisabled();
  });

  /**
   * ТЕСТ: Вход в комнату по коду
   * 
   * Проверяет что можно войти в существующую комнату
   * используя её код
   */
  test('вход в комнату по коду', async ({ browser }) => {
    // Создаем два контекста браузера для двух пользователей
    const context1 = await browser.newContext();
    const context2 = await browser.newContext();
    
    const page1 = await context1.newPage();
    const page2 = await context2.newPage();
    
    // Первый пользователь создает комнату
    await login(page1);
    await page1.getByRole('button', { name: /создать комнату/i }).click();
    await page1.getByLabel(/название комнаты/i).fill(generateRoomName());
    await page1.getByLabel(/пак вопросов/i).selectOption({ index: 1 });
    await page1.getByRole('button', { name: /создать/i }).click();
    
    await expect(page1).toHaveURL(/\/room\//, { timeout: 10000 });
    
    // Получаем код комнаты
    const codeElement = page1.locator('.room-page__code-value');
    const roomCode = await codeElement.textContent();
    
    // Второй пользователь входит по коду
    await login(page2);
    await page2.getByPlaceholder(/код комнаты/i).fill(roomCode!);
    await page2.getByRole('button', { name: /войти по коду/i }).click();
    
    // Должен попасть в комнату
    await expect(page2).toHaveURL(/\/room\//, { timeout: 10000 });
    
    // Очистка
    await context1.close();
    await context2.close();
  });

  /**
   * ТЕСТ: Ошибка при несуществующем коде комнаты
   * 
   * Проверяет что при вводе несуществующего кода
   * отображается сообщение об ошибке
   */
  test('ошибка при несуществующем коде комнаты', async ({ page }) => {
    await login(page);
    
    await page.getByPlaceholder(/код комнаты/i).fill('WRONG1');
    await page.getByRole('button', { name: /войти по коду/i }).click();
    
    // Должно появиться сообщение об ошибке
    // (зависит от реализации показа ошибок)
    await expect(page.getByText(/не найден/i)).toBeVisible({ timeout: 5000 });
  });

  /**
   * ТЕСТ: Автоматический uppercase для кода комнаты
   * 
   * Проверяет что введенный код автоматически
   * преобразуется в верхний регистр
   */
  test('код комнаты автоматически в верхнем регистре', async ({ page }) => {
    await login(page);
    
    const codeInput = page.getByPlaceholder(/код комнаты/i);
    await codeInput.fill('abc123');
    
    await expect(codeInput).toHaveValue('ABC123');
  });

  /**
   * ТЕСТ: Отображение счетчика игроков
   * 
   * Проверяет что отображается текущее количество
   * игроков в комнате
   */
  test('отображение счетчика игроков', async ({ page }) => {
    await login(page);
    
    await page.getByRole('button', { name: /создать комнату/i }).click();
    await page.getByLabel(/название комнаты/i).fill(generateRoomName());
    await page.getByLabel(/пак вопросов/i).selectOption({ index: 1 });
    await page.getByRole('button', { name: /создать/i }).click();
    
    await expect(page).toHaveURL(/\/room\//, { timeout: 10000 });
    
    // Должен показываться счетчик (1 из N)
    await expect(page.getByText(/игроки.*1/i)).toBeVisible();
  });
});


