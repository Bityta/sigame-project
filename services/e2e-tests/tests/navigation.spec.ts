import { test, expect } from '@playwright/test';
import { registerUser, generateUsername } from './helpers/auth';
import { createRoom } from './helpers/room';

test.describe('Навигация между страницами', () => {
  test('навигация: лобби -> профиль -> лобби', async ({ page }) => {
    const username = generateUsername();
    const password = 'testpass123';

    await registerUser(page, username, password);
    
    await page.getByText(username).click();
    await expect(page).toHaveURL(/\/profile/);
    
    await page.getByRole('button', { name: /← в лобби/i }).click();
    await expect(page).toHaveURL(/\/lobby/);
  });

  test('навигация: лобби -> создание комнаты -> лобби', async ({ page }) => {
    const username = generateUsername();
    const password = 'testpass123';

    await registerUser(page, username, password);
    
    await page.getByRole('button', { name: /создать комнату/i }).click();
    await expect(page).toHaveURL(/\/lobby\/create/);
    
    await page.getByRole('button', { name: /отмена/i }).click();
    await expect(page).toHaveURL(/\/lobby/);
  });

  test('навигация: лобби -> комната -> игра -> лобби', async ({ page, context }) => {
    const hostUsername = generateUsername();
    const playerUsername = generateUsername();
    const password = 'testpass123';

    await registerUser(page, hostUsername, password);
    const roomId = await createRoom(page);
    
    const playerPage = await context.newPage();
    await registerUser(playerPage, playerUsername, password);
    await joinRoom(playerPage, roomId);
    
    await page.getByRole('button', { name: /готов/i }).click();
    await playerPage.getByRole('button', { name: /готов/i }).click();
    
    await page.waitForURL(/\/game\/.+/, { timeout: 30000 });
    
    await page.getByRole('button', { name: /выйти из игры/i }).click();
    
    await expect(page).toHaveURL(/\/lobby/);
    
    await playerPage.close();
  });

  test('редирект неавторизованного пользователя', async ({ page }) => {
    await page.goto('/lobby');
    
    await expect(page).toHaveURL(/\/login/, { timeout: 5000 });
  });

  test('редирект авторизованного пользователя со страницы входа', async ({ page }) => {
    const username = generateUsername();
    const password = 'testpass123';

    await registerUser(page, username, password);
    
    await page.goto('/login');
    
    await expect(page).toHaveURL(/\/lobby/, { timeout: 5000 });
  });
});

