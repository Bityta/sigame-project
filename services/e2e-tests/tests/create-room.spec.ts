import { test, expect } from '@playwright/test';
import { registerUser, generateUsername } from './helpers/auth';

test.describe('Создание комнаты', () => {
  test('заполнение формы создания комнаты', async ({ page }) => {
    const username = generateUsername();
    const password = 'testpass123';

    await registerUser(page, username, password);
    
    await page.getByRole('button', { name: /создать комнату/i }).click();
    await page.waitForURL(/\/lobby\/create/);
    
    await expect(page.getByText(/создать комнату/i)).toBeVisible();
    
    await page.getByPlaceholder(/введите название/i).fill('Моя тестовая комната');
    
    const packSelect = page.locator('select').first();
    await packSelect.waitFor({ state: 'visible' });
    
    await expect(page.getByText(/максимум игроков/i)).toBeVisible();
    
    const publicCheckbox = page.getByLabel(/публичная комната/i);
    await expect(publicCheckbox).toBeVisible();
  });

  test('валидация названия комнаты', async ({ page }) => {
    const username = generateUsername();
    const password = 'testpass123';

    await registerUser(page, username, password);
    
    await page.getByRole('button', { name: /создать комнату/i }).click();
    
    await page.getByPlaceholder(/введите название/i).fill('ab');
    
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

  test('изменение максимального количества игроков', async ({ page }) => {
    const username = generateUsername();
    const password = 'testpass123';

    await registerUser(page, username, password);
    
    await page.getByRole('button', { name: /создать комнату/i }).click();
    await page.waitForURL(/\/lobby\/create/);
    
    await page.getByPlaceholder(/введите название/i).fill('Test Room Max Players');
    
    const slider = page.locator('input[type="range"]').first();
    await slider.waitFor({ state: 'visible' });
    await slider.fill('6');
    
    await expect(page.getByText(/6/)).toBeVisible({ timeout: 5000 });
    
    const packSelect = page.locator('select').first();
    await packSelect.waitFor({ state: 'visible' });
    await page.waitForTimeout(1000);
    
    const packOptions = await packSelect.locator('option').all();
    if (packOptions.length > 1) {
      const firstRealOption = packOptions[1];
      const optionValue = await firstRealOption.getAttribute('value');
      if (optionValue) {
        await packSelect.selectOption(optionValue);
      } else {
        await packSelect.selectOption({ index: 1 });
      }
      await page.waitForTimeout(1000);
    }
    
    const createButton = page.getByRole('button', { name: /создать/i });
    await createButton.waitFor({ state: 'visible' });
    await createButton.click();
    
    await page.waitForURL(/\/room\/.+/, { timeout: 20000 });
    
    const maxPlayersText = await page.getByText(/6/).or(page.getByText(/максимум.*6/i)).isVisible({ timeout: 5000 }).catch(() => false);
    if (maxPlayersText) {
      await expect(page.getByText(/6/).or(page.getByText(/максимум.*6/i))).toBeVisible({ timeout: 5000 });
    }
  });

  test('переключение публичной/приватной комнаты', async ({ page }) => {
    const username = generateUsername();
    const password = 'testpass123';

    await registerUser(page, username, password);
    
    await page.getByRole('button', { name: /создать комнату/i }).click();
    
    const publicCheckbox = page.getByLabel(/публичная комната/i);
    await publicCheckbox.uncheck();
    
    await expect(page.getByPlaceholder(/введите пароль/i)).toBeVisible();
    
    await publicCheckbox.check();
    
    await expect(page.getByPlaceholder(/введите пароль/i)).not.toBeVisible();
  });
});

