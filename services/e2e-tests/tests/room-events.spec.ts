import { test, expect } from '@playwright/test';
import { registerUser, generateUsername } from './helpers/auth';
import { createRoom, joinRoom, setReady } from './helpers/room';

test.describe('SSE события комнаты', () => {
  test('обновление при присоединении нового игрока', async ({ page, browser }) => {
    const hostUsername = generateUsername();
    const playerUsername = generateUsername();
    const password = 'testpass123';

    await registerUser(page, hostUsername, password);
    const roomId = await createRoom(page);
    
    await expect(page.getByText(hostUsername)).toBeVisible();
    
    const playerContext = await browser.newContext();
    const playerPage = await playerContext.newPage();
    await registerUser(playerPage, playerUsername, password);
    await joinRoom(playerPage, roomId);
    
    await page.waitForTimeout(2000);
    
    await expect(page.getByText(playerUsername)).toBeVisible({ timeout: 5000 });
    
    await playerContext.close();
  });

  test('обновление при готовности игрока', async ({ page, browser }) => {
    const hostUsername = generateUsername();
    const playerUsername = generateUsername();
    const password = 'testpass123';

    await registerUser(page, hostUsername, password);
    const roomId = await createRoom(page);
    
    const playerContext = await browser.newContext();
    const playerPage = await playerContext.newPage();
    await registerUser(playerPage, playerUsername, password);
    await joinRoom(playerPage, roomId);
    
    await setReady(playerPage);
    
    await page.waitForTimeout(2000);
    
    await expect(page.locator('.room-page__ready-count')).toHaveText(/1 \/ 2/, { timeout: 5000 });
    
    await playerContext.close();
  });

  test('редирект при запуске игры через SSE', async ({ page, browser }) => {
    const hostUsername = generateUsername();
    const playerUsername = generateUsername();
    const password = 'testpass123';

    await registerUser(page, hostUsername, password);
    const roomId = await createRoom(page);
    
    const playerContext = await browser.newContext();
    const playerPage = await playerContext.newPage();
    await registerUser(playerPage, playerUsername, password);
    await joinRoom(playerPage, roomId);
    
    await setReady(page);
    await setReady(playerPage);
    
    await page.waitForURL(/\/game\/.+/, { timeout: 30000 });
    await playerPage.waitForURL(/\/game\/.+/, { timeout: 30000 });
    
    await expect(page).toHaveURL(/\/game\/.+/);
    await expect(playerPage).toHaveURL(/\/game\/.+/);
    
    await playerContext.close();
  });

  test('передача хоста при уходе хоста из комнаты', async ({ page, browser }) => {
    const hostUsername = generateUsername();
    const playerUsername = generateUsername();
    const password = 'testpass123';

    await registerUser(page, hostUsername, password);
    const roomId = await createRoom(page);
    
    const playerContext = await browser.newContext();
    const playerPage = await playerContext.newPage();
    await registerUser(playerPage, playerUsername, password);
    await joinRoom(playerPage, roomId);
    
    await playerPage.waitForTimeout(1000);
    
    await expect(playerPage.getByText(hostUsername)).toBeVisible();
    
    await page.getByRole('button', { name: /покинуть комнату/i }).click();
    
    await playerPage.waitForTimeout(3000);
    
    await expect(playerPage).toHaveURL(/\/room\/.+/);
    
    await expect(playerPage.getByText(playerUsername)).toBeVisible();
    await expect(playerPage.getByText('👑')).toBeVisible();
    
    await playerContext.close();
  });
});

