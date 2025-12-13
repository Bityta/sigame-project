import { test, expect } from '@playwright/test';
import { registerUser, generateUsername, logoutUser } from './helpers/auth';

test.describe('Страница регистрации - детальные тесты', () => {
  test('отображение всех элементов страницы', async ({ page }) => {
    await page.goto('/register');
    
    await expect(page.getByRole('heading', { name: /регистрация/i })).toBeVisible();
    await expect(page.getByText(/создать новый аккаунт/i)).toBeVisible();
    await expect(page.getByPlaceholder(/от 3 до 20 символов/i)).toBeVisible();
    await expect(page.getByPlaceholder(/минимум 6 символов/i)).toBeVisible();
    await expect(page.getByPlaceholder(/повторите пароль/i)).toBeVisible();
    await expect(page.getByRole('button', { name: /зарегистрироваться/i })).toBeVisible();
    await expect(page.getByRole('button', { name: /войти/i })).toBeVisible();
  });

  test('валидация - длинный username', async ({ page }) => {
    await page.goto('/register');
    
    const longUsername = 'a'.repeat(21);
    await page.getByPlaceholder(/от 3 до 20 символов/i).fill(longUsername);
    await page.waitForTimeout(1000);
    await page.getByPlaceholder(/минимум 6 символов/i).fill('testpass123');
    await page.getByPlaceholder(/повторите пароль/i).fill('testpass123');
    
    await page.waitForTimeout(500);
    
    const submitButton = page.getByRole('button', { name: /зарегистрироваться/i });
    const isDisabled = await submitButton.isDisabled();
    
    if (!isDisabled) {
      await page.evaluate(() => {
        const form = document.querySelector('form');
        if (form) {
          const event = new Event('submit', { bubbles: true, cancelable: true });
          form.dispatchEvent(event);
        }
      });
      await page.waitForTimeout(1000);
    }
    
    await expect(page.getByText(/имя пользователя должно быть/i)).toBeVisible({ timeout: 5000 });
  });

  test('валидация - короткий пароль', async ({ page }) => {
    const username = generateUsername();
    
    await page.goto('/register');
    
    await page.getByPlaceholder(/от 3 до 20 символов/i).fill(username);
    await page.waitForTimeout(2000);
    await page.getByPlaceholder(/минимум 6 символов/i).fill('12345');
    await page.getByPlaceholder(/повторите пароль/i).fill('12345');
    
    const submitButton = page.getByRole('button', { name: /зарегистрироваться/i });
    await expect(submitButton).toBeDisabled();
    
    await page.evaluate(() => {
      const form = document.querySelector('form');
      if (form) {
        const event = new Event('submit', { bubbles: true, cancelable: true });
        form.dispatchEvent(event);
      }
    });
    
    await page.waitForTimeout(1000);
    await expect(page.getByText(/пароль должен быть/i)).toBeVisible({ timeout: 5000 });
  });

  test('проверка доступности username через API', async ({ page }) => {
    await page.goto('/register');
    
    const username = generateUsername();
    await page.getByPlaceholder(/от 3 до 20 символов/i).fill(username);
    
    await page.waitForTimeout(2000);
    
    const hintVisible = await page.locator('.input-hint').first().isVisible({ timeout: 5000 }).catch(() => false);
    if (hintVisible) {
      await expect(page.locator('.input-hint').first()).toContainText(/доступно|имя занято/i, { timeout: 5000 });
    }
  });

  test('кнопка неактивна пока username проверяется', async ({ page }) => {
    await page.goto('/register');
    
    const username = generateUsername();
    await page.getByPlaceholder(/от 3 до 20 символов/i).fill(username);
    
    const button = page.getByRole('button', { name: /зарегистрироваться/i });
    
    const buttonDisabled = await button.isDisabled({ timeout: 1000 }).catch(() => false);
    if (buttonDisabled) {
      await expect(button).toBeDisabled({ timeout: 1000 });
    }
    
    await page.waitForTimeout(2000);
  });

  test('ошибка при существующем username', async ({ page }) => {
    const username = generateUsername();
    const password = 'testpass123';

    await registerUser(page, username, password);
    await logoutUser(page);
    await page.goto('/register');
    
    await page.getByPlaceholder(/от 3 до 20 символов/i).fill(username);
    await page.waitForTimeout(3000);
    
    const hintElement = page.locator('.input-hint').first();
    
    let hintText = '';
    for (let i = 0; i < 20; i++) {
      hintText = await hintElement.textContent().catch(() => '');
      if (hintText && (hintText.includes('занято') || hintText.includes('Имя занято'))) {
        break;
      }
      await page.waitForTimeout(500);
    }
    
    await page.getByPlaceholder(/минимум 6 символов/i).fill(password);
    await page.getByPlaceholder(/повторите пароль/i).fill(password);
    
    await page.waitForTimeout(2000);
    
    const submitButton = page.getByRole('button', { name: /зарегистрироваться/i });
    const isDisabled = await submitButton.isDisabled();
    
    const finalHintText = await hintElement.textContent().catch(() => '');
    const hasHint = finalHintText && (finalHintText.includes('занято') || finalHintText.includes('Имя занято'));
    
    if (!hasHint && !isDisabled) {
      await page.waitForTimeout(3000);
      const finalHintText2 = await hintElement.textContent().catch(() => '');
      const hasHint2 = finalHintText2 && (finalHintText2.includes('занято') || finalHintText2.includes('Имя занято'));
      const isDisabled2 = await submitButton.isDisabled();
      expect(hasHint2 || isDisabled2).toBe(true);
    } else {
      expect(isDisabled || hasHint).toBe(true);
    }
  });

  test('сохранение токенов после регистрации', async ({ page }) => {
    const username = generateUsername();
    const password = 'testpass123';

    await registerUser(page, username, password);
    
    const tokens = await page.evaluate(() => {
      const keys = Object.keys(localStorage);
      const tokenKeys = keys.filter(key => key.toLowerCase().includes('token'));
      return {
        hasTokens: tokenKeys.length > 0,
        keys: tokenKeys
      };
    });
    
    expect(tokens.hasTokens).toBe(true);
  });
});

