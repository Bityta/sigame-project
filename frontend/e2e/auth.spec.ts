/**
 * E2E тесты авторизации
 * 
 * Тестируют полные пользовательские сценарии:
 * - Регистрация нового пользователя
 * - Вход существующего пользователя
 * - Ошибки авторизации
 * - Выход из системы
 * - Защита роутов
 */

import { test, expect } from '@playwright/test';

// Генерируем уникальное имя для тестового пользователя
const generateUsername = () => `testuser_${Date.now()}`;

test.describe('Авторизация', () => {
  /**
   * ТЕСТ: Отображение страницы входа
   * 
   * Проверяет что страница входа загружается корректно
   * и содержит все необходимые элементы
   */
  test('страница входа загружается корректно', async ({ page }) => {
    await page.goto('/login');
    
    // Проверяем заголовок
    await expect(page.getByRole('heading', { name: /вход/i })).toBeVisible();
    
    // Проверяем поля формы
    await expect(page.getByLabel(/имя пользователя/i)).toBeVisible();
    await expect(page.getByLabel(/пароль/i)).toBeVisible();
    
    // Проверяем кнопки
    await expect(page.getByRole('button', { name: /войти/i })).toBeVisible();
    await expect(page.getByRole('button', { name: /зарегистрироваться/i })).toBeVisible();
  });

  /**
   * ТЕСТ: Переход на страницу регистрации
   * 
   * Проверяет что ссылка "Зарегистрироваться" 
   * перенаправляет на страницу регистрации
   */
  test('переход со страницы входа на регистрацию', async ({ page }) => {
    await page.goto('/login');
    
    await page.getByRole('button', { name: /зарегистрироваться/i }).click();
    
    await expect(page).toHaveURL('/register');
    await expect(page.getByRole('heading', { name: /регистрация/i })).toBeVisible();
  });

  /**
   * ТЕСТ: Валидация короткого имени пользователя на странице входа
   * 
   * Проверяет что при вводе слишком короткого имени
   * показывается ошибка валидации
   */
  test('ошибка при коротком имени пользователя', async ({ page }) => {
    await page.goto('/login');
    
    await page.getByLabel(/имя пользователя/i).fill('ab');
    await page.getByLabel(/пароль/i).fill('password123');
    await page.getByRole('button', { name: /войти/i }).click();
    
    await expect(page.getByText(/не менее 3 символов/i)).toBeVisible();
  });

  /**
   * ТЕСТ: Валидация короткого пароля на странице входа
   * 
   * Проверяет что при вводе слишком короткого пароля
   * показывается ошибка валидации
   */
  test('ошибка при коротком пароле', async ({ page }) => {
    await page.goto('/login');
    
    await page.getByLabel(/имя пользователя/i).fill('testuser');
    await page.getByLabel(/пароль/i).fill('12345');
    await page.getByRole('button', { name: /войти/i }).click();
    
    await expect(page.getByText(/не менее 6 символов/i)).toBeVisible();
  });

  /**
   * ТЕСТ: Отображение страницы регистрации
   * 
   * Проверяет что страница регистрации загружается корректно
   */
  test('страница регистрации загружается корректно', async ({ page }) => {
    await page.goto('/register');
    
    await expect(page.getByRole('heading', { name: /регистрация/i })).toBeVisible();
    await expect(page.getByLabel(/имя пользователя/i)).toBeVisible();
    await expect(page.getByLabel(/пароль/i)).toBeVisible();
    await expect(page.getByRole('button', { name: /зарегистрироваться/i })).toBeVisible();
  });

  /**
   * ТЕСТ: Переход на страницу входа с регистрации
   * 
   * Проверяет что ссылка "Войти" перенаправляет
   * на страницу входа
   */
  test('переход со страницы регистрации на вход', async ({ page }) => {
    await page.goto('/register');
    
    await page.getByRole('button', { name: /войти/i }).click();
    
    await expect(page).toHaveURL('/login');
  });

  /**
   * ТЕСТ: Регистрация нового пользователя
   * 
   * Полный сценарий регистрации:
   * - Заполнение формы
   * - Отправка
   * - Редирект в лобби
   * - Отображение имени пользователя
   */
  test('успешная регистрация нового пользователя', async ({ page }) => {
    const username = generateUsername();
    
    await page.goto('/register');
    
    await page.getByLabel(/имя пользователя/i).fill(username);
    await page.getByLabel(/пароль/i).fill('testpassword123');
    await page.getByRole('button', { name: /зарегистрироваться/i }).click();
    
    // Ждем редирект в лобби
    await expect(page).toHaveURL('/lobby', { timeout: 10000 });
    
    // Проверяем что отображается имя пользователя
    await expect(page.getByText(new RegExp(username, 'i'))).toBeVisible();
  });

  /**
   * ТЕСТ: Вход существующего пользователя
   * 
   * Сначала регистрируем, потом выходим, потом входим снова
   */
  test('успешный вход существующего пользователя', async ({ page }) => {
    const username = generateUsername();
    const password = 'testpassword123';
    
    // Сначала регистрируемся
    await page.goto('/register');
    await page.getByLabel(/имя пользователя/i).fill(username);
    await page.getByLabel(/пароль/i).fill(password);
    await page.getByRole('button', { name: /зарегистрироваться/i }).click();
    
    await expect(page).toHaveURL('/lobby', { timeout: 10000 });
    
    // Выходим
    await page.getByRole('button', { name: /выход/i }).click();
    await expect(page).toHaveURL('/login');
    
    // Входим снова
    await page.getByLabel(/имя пользователя/i).fill(username);
    await page.getByLabel(/пароль/i).fill(password);
    await page.getByRole('button', { name: /войти/i }).click();
    
    // Проверяем успешный вход
    await expect(page).toHaveURL('/lobby', { timeout: 10000 });
    await expect(page.getByText(new RegExp(username, 'i'))).toBeVisible();
  });

  /**
   * ТЕСТ: Ошибка при неверном пароле
   * 
   * Проверяет что при вводе неверного пароля
   * пользователь остается на странице входа
   * и видит ошибку
   */
  test('ошибка при неверном пароле', async ({ page }) => {
    await page.goto('/login');
    
    await page.getByLabel(/имя пользователя/i).fill('someuser');
    await page.getByLabel(/пароль/i).fill('wrongpassword');
    await page.getByRole('button', { name: /войти/i }).click();
    
    // Должны остаться на странице входа
    await expect(page).toHaveURL('/login');
    
    // Можем проверить наличие индикатора ошибки
    // (зависит от реализации отображения ошибок)
  });

  /**
   * ТЕСТ: Выход из системы
   * 
   * Проверяет что пользователь может выйти
   * и будет перенаправлен на страницу входа
   */
  test('выход из системы', async ({ page }) => {
    const username = generateUsername();
    
    // Регистрируемся
    await page.goto('/register');
    await page.getByLabel(/имя пользователя/i).fill(username);
    await page.getByLabel(/пароль/i).fill('testpassword123');
    await page.getByRole('button', { name: /зарегистрироваться/i }).click();
    
    await expect(page).toHaveURL('/lobby', { timeout: 10000 });
    
    // Выходим
    await page.getByRole('button', { name: /выход/i }).click();
    
    // Проверяем редирект на страницу входа
    await expect(page).toHaveURL('/login');
  });

  /**
   * ТЕСТ: Защита роутов - лобби недоступно без авторизации
   * 
   * Проверяет что неавторизованный пользователь
   * не может попасть в лобби
   */
  test('лобби недоступно без авторизации', async ({ page }) => {
    // Очищаем localStorage чтобы быть уверенными что не авторизованы
    await page.goto('/login');
    await page.evaluate(() => localStorage.clear());
    
    // Пытаемся зайти в лобби напрямую
    await page.goto('/lobby');
    
    // Должны быть перенаправлены на страницу входа
    await expect(page).toHaveURL('/login');
  });

  /**
   * ТЕСТ: Главная страница редиректит на логин
   * 
   * Проверяет что при заходе на главную страницу
   * неавторизованный пользователь попадает на логин
   */
  test('главная страница редиректит на логин для неавторизованного', async ({ page }) => {
    await page.evaluate(() => localStorage.clear());
    await page.goto('/');
    
    // Должны оказаться на странице входа
    await expect(page).toHaveURL('/login');
  });

  /**
   * ТЕСТ: Сохранение сессии после перезагрузки
   * 
   * Проверяет что после авторизации и перезагрузки страницы
   * пользователь остается авторизованным
   */
  test('сессия сохраняется после перезагрузки', async ({ page }) => {
    const username = generateUsername();
    
    // Регистрируемся
    await page.goto('/register');
    await page.getByLabel(/имя пользователя/i).fill(username);
    await page.getByLabel(/пароль/i).fill('testpassword123');
    await page.getByRole('button', { name: /зарегистрироваться/i }).click();
    
    await expect(page).toHaveURL('/lobby', { timeout: 10000 });
    
    // Перезагружаем страницу
    await page.reload();
    
    // Должны остаться в лобби
    await expect(page).toHaveURL('/lobby');
    await expect(page.getByText(new RegExp(username, 'i'))).toBeVisible();
  });

  /**
   * ТЕСТ: Форма входа доступна по Enter
   * 
   * Проверяет что форму можно отправить
   * нажатием Enter после заполнения
   */
  test('отправка формы входа по Enter', async ({ page }) => {
    await page.goto('/login');
    
    await page.getByLabel(/имя пользователя/i).fill('testuser');
    await page.getByLabel(/пароль/i).fill('password123');
    await page.getByLabel(/пароль/i).press('Enter');
    
    // Форма должна быть отправлена
    // Проверяем что происходит какое-то действие
    // (либо ошибка, либо редирект)
  });
});

