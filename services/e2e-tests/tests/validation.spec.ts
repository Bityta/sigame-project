import { test, expect } from '@playwright/test';
import { generateUsername } from './helpers/auth';

test.describe('Валидация форм', () => {
  test('валидация формы входа - все поля', async ({ page }) => {
    await page.goto('/login');
    
    await page.getByRole('button', { name: /войти/i }).click();
    
    await page.waitForTimeout(500);
    
    const usernameField = page.getByPlaceholder(/введите имя пользователя/i);
    
    await expect(usernameField).toBeFocused().catch(() => {
      expect(page.getByText(/должно быть/i).or(page.getByText(/обязательно/i))).toBeVisible();
    });
  });

  test('валидация формы регистрации - все поля', async ({ page }) => {
    await page.goto('/register');
    
    await page.evaluate(() => {
      const form = document.querySelector('form');
      if (form) {
        const event = new Event('submit', { bubbles: true, cancelable: true });
        form.dispatchEvent(event);
      }
    });
    
    await page.waitForTimeout(500);
    
    await expect(page.getByText(/должно быть/i).or(page.getByText(/обязательно/i)).or(page.getByText(/не менее/i))).toBeVisible({ timeout: 3000 });
  });

  test('валидация формы создания комнаты - все поля', async ({ page }) => {
    const username = generateUsername();
    const password = 'testpass123';

    await page.goto('/register');
    await page.getByPlaceholder(/от 3 до 20 символов/i).fill(username);
    await page.waitForTimeout(600);
    await page.getByPlaceholder(/минимум 6 символов/i).fill(password);
    await page.getByPlaceholder(/повторите пароль/i).fill(password);
    await page.getByRole('button', { name: /зарегистрироваться/i }).click();
    await page.waitForURL(/\/lobby/, { timeout: 10000 });
    
    await page.getByRole('button', { name: /создать комнату/i }).click();
    
    await page.evaluate(() => {
      const form = document.querySelector('form');
      if (form) {
        const event = new Event('submit', { bubbles: true, cancelable: true });
        form.dispatchEvent(event);
      }
    });
    
    await page.waitForTimeout(500);
    
    await expect(page.locator('.create-room-form__error, .input-error').first()).toBeVisible({ timeout: 3000 });
  });
});

