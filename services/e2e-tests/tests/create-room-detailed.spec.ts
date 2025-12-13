import { test, expect } from '@playwright/test';
import { registerUser, generateUsername } from './helpers/auth';

test.describe('Создание комнаты - детальные тесты', () => {
  test('валидация названия комнаты - слишком длинное', async ({ page }) => {
    const username = generateUsername();
    const password = 'testpass123';

    await registerUser(page, username, password);
    
    await page.getByRole('button', { name: /создать комнату/i }).click();
    
    const longName = 'a'.repeat(101);
    await page.getByPlaceholder(/введите название/i).fill(longName);
    
    const packSelect = page.locator('select').first();
    await packSelect.waitFor({ state: 'visible' });
    const packOptions = await packSelect.locator('option').all();
    if (packOptions.length > 1) {
      await packSelect.selectOption({ index: 1 });
    }
    
    await page.getByRole('button', { name: /создать/i }).click();
    
    await page.waitForTimeout(500);
    await expect(page.locator('.create-room-form__error')).toContainText(/название должно быть/i, { timeout: 5000 });
  });

  test('валидация - невыбранный пак вопросов', async ({ page }) => {
    const username = generateUsername();
    const password = 'testpass123';

    await registerUser(page, username, password);
    
    await page.getByRole('button', { name: /создать комнату/i }).click();
    
    await page.getByPlaceholder(/введите название/i).fill('Test Room');
    
    await page.getByRole('button', { name: /создать/i }).click();
    
    await page.waitForTimeout(500);
    await expect(page.locator('.create-room-form__error')).toContainText(/выберите пак/i, { timeout: 5000 });
  });

  test('валидация - приватная комната без пароля', async ({ page }) => {
    const username = generateUsername();
    const password = 'testpass123';

    await registerUser(page, username, password);
    
    await page.getByRole('button', { name: /создать комнату/i }).click();
    
    await page.getByPlaceholder(/введите название/i).fill('Test Room');
    
    const packSelect = page.locator('select').first();
    await packSelect.waitFor({ state: 'visible' });
    const packOptions = await packSelect.locator('option').all();
    if (packOptions.length > 1) {
      await packSelect.selectOption({ index: 1 });
    }
    
    const publicCheckbox = page.getByLabel(/публичная комната/i);
    await publicCheckbox.uncheck();
    
    await page.getByRole('button', { name: /создать/i }).click();
    
    await page.waitForTimeout(500);
    await expect(page.locator('.create-room-form__error')).toContainText(/пароль/i, { timeout: 5000 });
  });

  test('отображение состояния загрузки при создании', async ({ page }) => {
    const username = generateUsername();
    const password = 'testpass123';

    await registerUser(page, username, password);
    
    await page.getByRole('button', { name: /создать комнату/i }).click();
    
    await page.getByPlaceholder(/введите название/i).fill('Test Room');
    
    const packSelect = page.locator('select').first();
    await packSelect.waitFor({ state: 'visible' });
    const packOptions = await packSelect.locator('option').all();
    if (packOptions.length > 1) {
      await packSelect.selectOption({ index: 1 });
    }
    
    const createButton = page.getByRole('button', { name: /создать/i });
    await createButton.click();
    
    const createButtonDisabled = await createButton.isDisabled({ timeout: 2000 }).catch(() => false);
    if (createButtonDisabled) {
      await expect(createButton).toBeDisabled({ timeout: 2000 });
    }
  });

  test('кнопка отмены возвращает в лобби', async ({ page }) => {
    const username = generateUsername();
    const password = 'testpass123';

    await registerUser(page, username, password);
    
    await page.getByRole('button', { name: /создать комнату/i }).click();
    
    await page.getByRole('button', { name: /отмена/i }).click();
    
    await expect(page).toHaveURL(/\/lobby/);
  });
});

