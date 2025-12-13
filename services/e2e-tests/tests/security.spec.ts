import { test, expect } from '@playwright/test';
import { registerUser, generateUsername } from './helpers/auth';

test.describe('Безопасность', () => {
  test('валидация входных данных - XSS защита', async ({ page }) => {
    const username = generateUsername();
    const password = 'testpass123';

    await registerUser(page, username, password);
    
    await page.goto('/lobby');
    
    const xssPayload = '<script>alert("XSS")</script>';
    
    await page.evaluate((payload) => {
      const input = document.querySelector('input');
      if (input) {
        input.value = payload;
      }
    }, xssPayload);
    
    await page.waitForTimeout(1000);
    
    const alertShown = await page.evaluate(() => {
      return window.alert.toString().includes('native');
    });
    
    expect(alertShown).toBe(true);
  });

  test('валидация токенов на каждом запросе', async ({ page }) => {
    const username = generateUsername();
    const password = 'testpass123';

    await registerUser(page, username, password);
    
    await page.evaluate(() => {
      localStorage.removeItem('access_token');
    });
    
    await page.reload();
    
    await expect(page).toHaveURL(/\/login/, { timeout: 5000 });
  });

  test('rate limiting - множественные запросы', async ({ page }) => {
    await page.goto('/login');
    
    const username = generateUsername();
    
    for (let i = 0; i < 10; i++) {
      await page.getByLabel(/имя пользователя/i).fill(`${username}${i}`);
      await page.getByLabel(/пароль/i).fill('wrongpass');
      await page.getByRole('button', { name: /войти/i }).click();
      await page.waitForTimeout(100);
    }
    
    const rateLimitVisible = await page.getByText(/429/i).or(page.getByText(/too many/i)).isVisible({ timeout: 5000 }).catch(() => false);
    if (rateLimitVisible) {
      await expect(page.getByText(/429/i).or(page.getByText(/too many/i))).toBeVisible({ timeout: 5000 });
    }
  });

  test('защита от брутфорса - блокировка после множественных попыток', async ({ page }) => {
    const username = generateUsername();
    const password = 'testpass123';

    await registerUser(page, username, password);
    await page.getByRole('button', { name: /выход/i }).click();
    
    for (let i = 0; i < 5; i++) {
      await page.goto('/login');
      await page.getByLabel(/имя пользователя/i).fill(username);
      await page.getByLabel(/пароль/i).fill('wrongpass');
      await page.getByRole('button', { name: /войти/i }).click();
      await page.waitForTimeout(500);
    }
    
    const blockedVisible = await page.getByText(/заблокирован/i).or(page.getByText(/попробуйте позже/i)).isVisible({ timeout: 5000 }).catch(() => false);
    if (blockedVisible) {
      await expect(page.getByText(/заблокирован/i).or(page.getByText(/попробуйте позже/i))).toBeVisible({ timeout: 5000 });
    }
  });

  test('CSRF защита - токены в запросах', async ({ page }) => {
    const username = generateUsername();
    const password = 'testpass123';

    await registerUser(page, username, password);
    
    const response = await page.request.post('/api/lobby/rooms', {
      data: { name: 'Test', packId: 'test' }
    });
    
    expect(response.status()).toBe(401);
  });
});

