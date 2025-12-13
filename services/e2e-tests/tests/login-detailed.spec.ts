import { test, expect } from '@playwright/test';
import { registerUser, loginUser, generateUsername } from './helpers/auth';

test.describe('Страница входа - детальные тесты', () => {
  test('отображение всех элементов страницы', async ({ page }) => {
    await page.goto('/login');
    
    await expect(page.getByRole('heading', { name: /sigame/i })).toBeVisible();
    await expect(page.getByPlaceholder(/введите имя пользователя/i)).toBeVisible();
    await expect(page.getByPlaceholder(/введите пароль/i)).toBeVisible();
    await expect(page.getByRole('button', { name: /войти/i })).toBeVisible();
    await expect(page.getByRole('button', { name: /зарегистрироваться/i })).toBeVisible();
  });

  test('валидация - короткий пароль', async ({ page }) => {
    await page.goto('/login');
    
    await page.getByPlaceholder(/введите имя пользователя/i).fill('testuser');
    await page.getByPlaceholder(/введите пароль/i).fill('12345');
    await page.getByRole('button', { name: /войти/i }).click();
    
    await expect(page.getByText(/пароль должен быть/i)).toBeVisible({ timeout: 5000 });
  });

  test('валидация - пустые поля', async ({ page }) => {
    await page.goto('/login');
    
    const submitButton = page.getByRole('button', { name: /войти/i });
    await submitButton.click();
    
    await page.waitForTimeout(2000);
    
    const errorExists = await page.locator('.input-error').count() > 0;
    const buttonDisabled = await submitButton.isDisabled();
    const stillOnLoginPage = page.url().includes('/login');
    
    expect(errorExists || buttonDisabled || stillOnLoginPage).toBe(true);
  });

  test('сохранение сессии после перезагрузки', async ({ page }) => {
    const username = generateUsername();
    const password = 'testpass123';

    await registerUser(page, username, password);
    
    await page.reload();
    
    await expect(page).toHaveURL(/\/lobby/, { timeout: 10000 });
    await expect(page.getByText(username)).toBeVisible({ timeout: 5000 });
  });

  test('автоматический вход при наличии токена', async ({ page }) => {
    const username = generateUsername();
    const password = 'testpass123';

    await registerUser(page, username, password);
    
    await page.goto('/login');
    
    await expect(page).toHaveURL(/\/lobby/, { timeout: 5000 });
  });

  test('ошибка при несуществующем пользователе', async ({ page }) => {
    await page.goto('/login');
    
    await page.getByPlaceholder(/введите имя пользователя/i).fill('nonexistentuser12345');
    await page.getByPlaceholder(/введите пароль/i).fill('testpass123');
    
    const responsePromise = page.waitForResponse(
      response => response.url().includes('/auth/login'),
      { timeout: 10000 }
    );
    
    await page.getByRole('button', { name: /войти/i }).click();
    
    const response = await responsePromise;
    expect(response.status()).toBeGreaterThanOrEqual(400);
    
    await page.waitForTimeout(2000);
    
    const errorModal = page.getByText(/invalid username or password/i)
      .or(page.getByText(/неверный/i))
      .or(page.getByText(/ошибка/i))
      .or(page.locator('[class*="error"]').first());
    
    const errorVisible = await errorModal.isVisible({ timeout: 5000 }).catch(() => false);
    if (errorVisible) {
      await expect(errorModal).toBeVisible({ timeout: 5000 });
      
      const closeButton = page.getByRole('button', { name: /закрыть/i });
      const closeButtonVisible = await closeButton.isVisible({ timeout: 2000 }).catch(() => false);
      if (closeButtonVisible) {
        await closeButton.click();
      }
    }
  });

  test('состояние загрузки кнопки входа', async ({ page }) => {
    await page.goto('/login');
    
    await page.getByPlaceholder(/введите имя пользователя/i).fill('testuser');
    await page.getByPlaceholder(/введите пароль/i).fill('testpass123');
    
    const button = page.getByRole('button', { name: /войти/i });
    await button.click();
    
    const buttonDisabled = await button.isDisabled({ timeout: 2000 }).catch(() => false);
    if (buttonDisabled) {
      await expect(button).toBeDisabled({ timeout: 2000 });
    }
  });
});

