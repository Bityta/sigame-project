import { test, expect } from '@playwright/test';
import { registerUser, loginUser, logoutUser, generateUsername } from './helpers/auth';

test.describe('Аутентификация', () => {
  test('регистрация нового пользователя', async ({ page }) => {
    const username = generateUsername();
    const password = 'testpass123';

    await registerUser(page, username, password);
    
    await expect(page).toHaveURL(/\/lobby/);
  });

  test('валидация формы регистрации - короткий username', async ({ page }) => {
    await page.goto('/register');
    
    await page.getByPlaceholder(/от 3 до 20 символов/i).fill('ab');
    await page.getByPlaceholder(/минимум 6 символов/i).fill('testpass123');
    await page.getByPlaceholder(/повторите пароль/i).fill('testpass123');
    
    const submitButton = page.getByRole('button', { name: /зарегистрироваться/i });
    
    await expect(submitButton).toBeDisabled();
    
    await page.evaluate(() => {
      const form = document.querySelector('form');
      if (form) {
        const event = new Event('submit', { bubbles: true, cancelable: true });
        form.dispatchEvent(event);
      }
    });
    
    await page.waitForTimeout(500);
    
    await expect(page.getByText(/имя пользователя должно быть от 3 до 20 символов/i)).toBeVisible({ timeout: 5000 });
  });

  test('валидация формы регистрации - несовпадение паролей', async ({ page }) => {
    const username = generateUsername();
    
    await page.goto('/register');
    
    await page.getByPlaceholder(/от 3 до 20 символов/i).fill(username);
    await page.waitForTimeout(2000);
    
    const submitButton = page.getByRole('button', { name: /зарегистрироваться/i });
    try {
      await page.waitForFunction(
        (button) => !button.disabled,
        await submitButton.elementHandle(),
        { timeout: 15000 }
      );
    } catch (error) {
      // Кнопка может остаться disabled если username занят
    }
    
    await page.getByPlaceholder(/минимум 6 символов/i).fill('testpass123');
    await page.getByPlaceholder(/повторите пароль/i).fill('differentpass');
    
    await page.waitForTimeout(1000);
    
    const isDisabled = await submitButton.isDisabled();
    expect(isDisabled).toBe(true);
    
    await page.evaluate(() => {
      const form = document.querySelector('form');
      if (form) {
        const event = new Event('submit', { bubbles: true, cancelable: true });
        form.dispatchEvent(event);
      }
    });
    
    await page.waitForTimeout(1000);
    
    await expect(page.getByText(/пароли не совпадают/i)).toBeVisible({ timeout: 5000 });
  });

  test('вход существующего пользователя', async ({ page }) => {
    const username = generateUsername();
    const password = 'testpass123';

    await registerUser(page, username, password);
    await logoutUser(page);
    
    await loginUser(page, username, password);
    
    await expect(page).toHaveURL(/\/lobby/);
  });

  test('ошибка при неверном пароле', async ({ page }) => {
    const username = generateUsername();
    const password = 'testpass123';

    await registerUser(page, username, password);
    await logoutUser(page);
    
    await page.goto('/login');
    await page.getByPlaceholder(/введите имя пользователя/i).fill(username);
    await page.getByPlaceholder(/введите пароль/i).fill('wrongpassword');
    await page.getByRole('button', { name: /войти/i }).click();
    
    await page.waitForTimeout(1000);
    
    const errorModal = page.getByText(/invalid username or password/i).or(page.getByText(/неверный/i)).or(page.getByText(/ошибка/i));
    await expect(errorModal).toBeVisible({ timeout: 5000 });
    
    const closeButton = page.getByRole('button', { name: /закрыть/i });
    const closeButtonVisible = await closeButton.isVisible({ timeout: 2000 }).catch(() => false);
    if (closeButtonVisible) {
      await closeButton.click();
    }
  });

  test('переход между страницами входа и регистрации', async ({ page }) => {
    await page.goto('/login');
    
    await page.getByRole('button', { name: /зарегистрироваться/i }).click();
    await expect(page).toHaveURL(/\/register/);
    
    await page.getByRole('button', { name: /войти/i }).click();
    await expect(page).toHaveURL(/\/login/);
  });
});

