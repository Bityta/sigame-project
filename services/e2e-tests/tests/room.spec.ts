import { test, expect } from '@playwright/test';
import { registerUser, generateUsername } from './helpers/auth';
import { createRoom, joinRoom, setReady, getRoomCode } from './helpers/room';

test.describe('Комната ожидания', () => {
  test('отображение информации о комнате', async ({ page }) => {
    const username = generateUsername();
    const password = 'testpass123';

    await registerUser(page, username, password);
    const roomId = await createRoom(page);
    
    await expect(page.locator('.room-page__title')).toBeVisible();
    await expect(page.locator('.room-page__code-value')).toBeVisible();
  });

  test('копирование кода комнаты', async ({ page }) => {
    const username = generateUsername();
    const password = 'testpass123';

    await registerUser(page, username, password);
    await createRoom(page);
    
    const codeElement = page.locator('.room-page__code');
    await codeElement.click();
    
    await expect(codeElement).toHaveClass(/room-page__code--copied/);
  });

  test('готовность игрока', async ({ page, context }) => {
    const hostUsername = generateUsername();
    const playerUsername = generateUsername();
    const password = 'testpass123';

    await registerUser(page, hostUsername, password);
    const roomId = await createRoom(page);
    
    const playerPage = await context.newPage();
    await registerUser(playerPage, playerUsername, password);
    await joinRoom(playerPage, roomId);
    
    const readyButton = page.getByRole('button', { name: /готов/i });
    await readyButton.click();
    
    await expect(page.getByText(/вы готовы/i)).toBeVisible();
    
    await playerPage.close();
  });

  test('автоматический запуск игры когда все готовы', async ({ page, context }) => {
    const hostUsername = generateUsername();
    const playerUsername = generateUsername();
    const password = 'testpass123';

    await registerUser(page, hostUsername, password);
    const roomId = await createRoom(page);
    
    const playerPage = await context.newPage();
    await registerUser(playerPage, playerUsername, password);
    await joinRoom(playerPage, roomId);
    
    await setReady(page);
    await setReady(playerPage);
    
    await page.waitForURL(/\/game\/.+/, { timeout: 30000 });
    await expect(page).toHaveURL(/\/game\/.+/);
    
    await playerPage.close();
  });

  test('настройки комнаты для хоста', async ({ page }) => {
    const username = generateUsername();
    const password = 'testpass123';

    await registerUser(page, username, password);
    await createRoom(page);
    
    await expect(page.locator('.room-settings__slider')).toBeVisible();
    
    const slider = page.locator('.room-settings__slider').first();
    await slider.fill('45');
    
    await page.getByRole('button', { name: /сохранить настройки/i }).click();
    
    await expect(page.getByText(/45 сек/i)).toBeVisible({ timeout: 5000 });
  });
});

