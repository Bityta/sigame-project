import { test, expect } from '@playwright/test';
import { registerUser, generateUsername } from './helpers/auth';
import { createRoom, joinRoom } from './helpers/room';

test.describe('Настройки комнаты - детальные тесты', () => {
  test('изменение времени на ответ', async ({ page }) => {
    const username = generateUsername();
    const password = 'testpass123';

    await registerUser(page, username, password);
    await createRoom(page);
    
    const slider = page.locator('.room-settings__slider').first();
    await slider.fill('45');
    
    await page.getByRole('button', { name: /сохранить настройки/i }).click();
    
    await expect(page.getByText(/45/i)).toBeVisible({ timeout: 5000 });
  });

  test('изменение времени на выбор вопроса', async ({ page }) => {
    const username = generateUsername();
    const password = 'testpass123';

    await registerUser(page, username, password);
    await createRoom(page);
    
    const sliders = page.locator('.room-settings__slider');
    const secondSlider = sliders.nth(1);
    await secondSlider.fill('90');
    
    await page.getByRole('button', { name: /сохранить настройки/i }).click();
    
    await expect(page.getByText(/90/i)).toBeVisible({ timeout: 5000 });
  });

  test('валидация диапазона слайдера времени на ответ', async ({ page }) => {
    const username = generateUsername();
    const password = 'testpass123';

    await registerUser(page, username, password);
    await createRoom(page);
    
    const slider = page.locator('.room-settings__slider').first();
    const min = await slider.getAttribute('min');
    const max = await slider.getAttribute('max');
    
    expect(min).toBe('10');
    expect(max).toBe('60');
  });

  test('валидация диапазона слайдера времени на выбор', async ({ page }) => {
    const username = generateUsername();
    const password = 'testpass123';

    await registerUser(page, username, password);
    await createRoom(page);
    
    const sliders = page.locator('.room-settings__slider');
    const secondSlider = sliders.nth(1);
    const min = await secondSlider.getAttribute('min');
    const max = await secondSlider.getAttribute('max');
    
    expect(min).toBe('5');
    expect(max).toBe('30');
  });

  test('отображение настроек для не-хоста', async ({ page, context }) => {
    const hostUsername = generateUsername();
    const playerUsername = generateUsername();
    const password = 'testpass123';

    await registerUser(page, hostUsername, password);
    const roomId = await createRoom(page);
    
    const playerPage = await context.newPage();
    await registerUser(playerPage, playerUsername, password);
    await joinRoom(playerPage, roomId);
    
    await expect(playerPage.locator('.room-settings__view')).toBeVisible();
    await expect(playerPage.locator('.room-settings__form')).not.toBeVisible();
    
    await playerPage.close();
  });
});

