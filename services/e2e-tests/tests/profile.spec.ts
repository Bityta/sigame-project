import { test, expect } from '@playwright/test';
import { registerUser, generateUsername, logoutUser } from './helpers/auth';

test.describe('Профиль', () => {
  test('отображение страницы профиля', async ({ page }) => {
    const username = generateUsername();
    const password = 'testpass123';

    await registerUser(page, username, password);
    
    await page.getByText(username).click();
    await expect(page).toHaveURL(/\/profile/);
    
    await expect(page.getByRole('heading', { name: /профиль игрока/i })).toBeVisible();
    await expect(page.getByText(username)).toBeVisible();
  });

  test('переключение табов', async ({ page }) => {
    const username = generateUsername();
    const password = 'testpass123';

    await registerUser(page, username, password);
    await page.getByText(username).click();
    
    await page.getByRole('button', { name: /статистика/i }).click();
    await expect(page.getByText(/в разработке/i)).toBeVisible();
    
    await page.getByRole('button', { name: /достижения/i }).click();
    await expect(page.getByText(/в разработке/i)).toBeVisible();
    
    await page.getByRole('button', { name: /история игр/i }).click();
    await expect(page.getByText(/в разработке/i)).toBeVisible();
    
    await page.getByRole('button', { name: /настройки/i }).click();
    await expect(page.getByRole('heading', { name: /уведомления/i })).toBeVisible();
  });

  test('выход из профиля', async ({ page }) => {
    const username = generateUsername();
    const password = 'testpass123';

    await registerUser(page, username, password);
    await page.getByText(username).click();
    
    await page.getByRole('button', { name: /выход/i }).click();
    
    await expect(page).toHaveURL(/\/login/);
  });

  test('возврат в лобби', async ({ page }) => {
    const username = generateUsername();
    const password = 'testpass123';

    await registerUser(page, username, password);
    await page.getByText(username).click();
    
    await page.getByRole('button', { name: /← в лобби/i }).click();
    
    await expect(page).toHaveURL(/\/lobby/);
  });
});

