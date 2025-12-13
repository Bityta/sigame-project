import { test, expect } from '@playwright/test';
import { registerUser, generateUsername } from './helpers/auth';
import { createRoom, joinRoomByCode, getRoomCode } from './helpers/room';

test.describe('Лобби', () => {
  test('создание комнаты', async ({ page }) => {
    const username = generateUsername();
    const password = 'testpass123';

    await registerUser(page, username, password);
    
    await expect(page.getByRole('heading', { name: /sigame/i })).toBeVisible();
    await expect(page.getByText(username)).toBeVisible();
    
    await page.getByRole('button', { name: /создать комнату/i }).click();
    await page.waitForURL(/\/lobby\/create/);
    
    await expect(page.getByText(/создать комнату/i)).toBeVisible();
  });

  test('присоединение к комнате по коду', async ({ browser }) => {
    const hostUsername = generateUsername();
    const playerUsername = generateUsername();
    const password = 'testpass123';

    const hostContext = await browser.newContext();
    const hostPage = await hostContext.newPage();
    
    await registerUser(hostPage, hostUsername, password);
    const roomId = await createRoom(hostPage);
    const roomCode = await getRoomCode(hostPage);
    
    const playerContext = await browser.newContext();
    const playerPage = await playerContext.newPage();
    
    await registerUser(playerPage, playerUsername, password);
    
    await joinRoomByCode(playerPage, roomCode);
    
    await expect(playerPage).toHaveURL(/\/room\/.+/);
    
    await playerPage.close();
    await playerContext.close();
    await hostPage.close();
    await hostContext.close();
  });

  test('отображение списка комнат', async ({ page }) => {
    const username = generateUsername();
    const password = 'testpass123';

    await registerUser(page, username, password);
    
    await expect(page.locator('.room-list').first()).toBeVisible({ timeout: 5000 });
  });

  test('выход из системы', async ({ page }) => {
    const username = generateUsername();
    const password = 'testpass123';

    await registerUser(page, username, password);
    
    await page.getByRole('button', { name: /выход/i }).click();
    
    await expect(page).toHaveURL(/\/login/);
  });

  test('переход в профиль', async ({ page }) => {
    const username = generateUsername();
    const password = 'testpass123';

    await registerUser(page, username, password);
    
    await page.getByText(username).click();
    
    await expect(page).toHaveURL(/\/profile/);
  });
});

