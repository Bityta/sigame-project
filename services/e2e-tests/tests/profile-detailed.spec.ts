import { test, expect } from '@playwright/test';
import { registerUser, generateUsername, logoutUser } from './helpers/auth';

test.describe('Профиль - детальные тесты', () => {
  test('отображение всех элементов профиля', async ({ page }) => {
    const username = generateUsername();
    const password = 'testpass123';

    await registerUser(page, username, password);
    
    await page.getByText(username).click();
    
    await expect(page.getByRole('heading', { name: /профиль игрока/i })).toBeVisible();
    await expect(page.getByText(username)).toBeVisible();
    await expect(page.getByRole('button', { name: /← в лобби/i })).toBeVisible();
    await expect(page.getByRole('button', { name: /выход/i })).toBeVisible();
  });

  test('отображение аватара или дефолтного', async ({ page }) => {
    const username = generateUsername();
    const password = 'testpass123';

    await registerUser(page, username, password);
    await page.getByText(username).click();
    
    await expect(page.locator('.profile-card__avatar').first()).toBeVisible();
  });

  test('отображение даты регистрации', async ({ page }) => {
    const username = generateUsername();
    const password = 'testpass123';

    await registerUser(page, username, password);
    await page.getByText(username).click();
    
    await expect(page.getByText(/в игре с/i).or(page.getByText(/создан/i))).toBeVisible();
  });

  test('отображение заглушки в разделе статистика', async ({ page }) => {
    const username = generateUsername();
    const password = 'testpass123';

    await registerUser(page, username, password);
    await page.getByText(username).click();
    
    await page.getByRole('button', { name: /статистика/i }).click();
    
    await expect(page.getByText(/в разработке/i)).toBeVisible();
    await expect(page.getByText(/🚧/)).toBeVisible();
  });

  test('отображение заглушки в разделе достижения', async ({ page }) => {
    const username = generateUsername();
    const password = 'testpass123';

    await registerUser(page, username, password);
    await page.getByText(username).click();
    
    await page.getByRole('button', { name: /достижения/i }).click();
    
    await expect(page.getByText(/в разработке/i)).toBeVisible();
  });

  test('отображение заглушки в разделе история игр', async ({ page }) => {
    const username = generateUsername();
    const password = 'testpass123';

    await registerUser(page, username, password);
    await page.getByText(username).click();
    
    await page.getByRole('button', { name: /история игр/i }).click();
    
    await expect(page.getByText(/в разработке/i)).toBeVisible();
  });

  test('настройки - readonly поля', async ({ page }) => {
    const username = generateUsername();
    const password = 'testpass123';

    await registerUser(page, username, password);
    await page.getByText(username).click();
    
    await page.getByRole('button', { name: /настройки/i }).click();
    
    const usernameInput = page.getByLabel(/имя пользователя/i).or(page.locator('input[value*="' + username + '"]'));
    await expect(usernameInput).toBeDisabled().catch(() => {
      expect(usernameInput).toHaveAttribute('readonly');
    });
  });

  test('настройки - переключение уведомлений', async ({ page }) => {
    const username = generateUsername();
    const password = 'testpass123';

    await registerUser(page, username, password);
    await page.getByText(username).click();
    
    await page.getByRole('button', { name: /настройки/i }).click();
    
    const notificationCheckbox = page.getByLabel(/уведомления/i).first();
    const initialState = await notificationCheckbox.isChecked();
    await notificationCheckbox.click({ force: true });
    
    await page.waitForTimeout(500);
    
    const newState = await notificationCheckbox.isChecked();
    expect(newState).not.toBe(initialState);
  });

  test('настройки - кнопка удаления аккаунта disabled', async ({ page }) => {
    const username = generateUsername();
    const password = 'testpass123';

    await registerUser(page, username, password);
    await page.getByText(username).click();
    
    await page.getByRole('button', { name: /настройки/i }).click();
    
    const deleteButton = page.getByRole('button', { name: /удалить аккаунт/i });
    await expect(deleteButton).toBeDisabled();
  });

  test('выход из профиля', async ({ page }) => {
    const username = generateUsername();
    const password = 'testpass123';

    await registerUser(page, username, password);
    await page.getByText(username).click();
    
    await page.getByRole('button', { name: /выход/i }).click();
    
    await expect(page).toHaveURL(/\/login/);
  });
});

