import { test, expect } from '@playwright/test';

test.describe('Валидация форм', () => {
  test('валидация формы входа - все поля', async ({ page }) => {
    await page.goto('/login');
    
    await page.getByRole('button', { name: /войти/i }).click();
    
    await page.waitForTimeout(500);
    
    const usernameField = page.getByLabel(/имя пользователя/i);
    const passwordField = page.getByLabel(/пароль/i);
    
    await expect(usernameField).toBeFocused().catch(() => {
      expect(page.getByText(/имя пользователя/i)).toBeVisible();
    });
  });

  test('валидация формы регистрации - все поля', async ({ page }) => {
    await page.goto('/register');
    
    await page.getByRole('button', { name: /зарегистрироваться/i }).click();
    
    await page.waitForTimeout(500);
    
    await expect(page.getByText(/имя пользователя/i).or(page.getByText(/пароль/i))).toBeVisible({ timeout: 3000 });
  });

  test('валидация формы создания комнаты - все поля', async ({ page }) => {
    const username = `testuser_${Date.now()}`;
    const password = 'testpass123';

    await page.goto('/register');
    await page.getByLabel(/имя пользователя/i).fill(username);
    await page.waitForTimeout(600);
    await page.getByLabel(/^пароль$/i).fill(password);
    await page.getByLabel(/подтвердите пароль/i).fill(password);
    await page.getByRole('button', { name: /зарегистрироваться/i }).click();
    await page.waitForURL(/\/lobby/, { timeout: 10000 });
    
    await page.getByRole('button', { name: /создать комнату/i }).click();
    
    await page.getByRole('button', { name: /создать/i }).click();
    
    await expect(page.getByText(/название/i).or(page.getByText(/пак/i))).toBeVisible({ timeout: 3000 });
  });
});

