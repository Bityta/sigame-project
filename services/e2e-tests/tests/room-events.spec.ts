import { test, expect } from '@playwright/test';
import { registerUser, generateUsername } from './helpers/auth';
import { createRoom, joinRoom, setReady } from './helpers/room';

test.describe('SSE события комнаты', () => {
  test('обновление при присоединении нового игрока', async ({ page, context }) => {
    const hostUsername = generateUsername();
    const playerUsername = generateUsername();
    const password = 'testpass123';

    await registerUser(page, hostUsername, password);
    const roomId = await createRoom(page);
    
    await expect(page.getByText(hostUsername)).toBeVisible();
    
    const playerPage = await context.newPage();
    await registerUser(playerPage, playerUsername, password);
    await joinRoom(playerPage, roomId);
    
    await page.waitForTimeout(2000);
    
    await expect(page.getByText(playerUsername)).toBeVisible({ timeout: 5000 });
    
    await playerPage.close();
  });

  test('обновление при готовности игрока', async ({ page, context }) => {
    const hostUsername = generateUsername();
    const playerUsername = generateUsername();
    const password = 'testpass123';

    await registerUser(page, hostUsername, password);
    const roomId = await createRoom(page);
    
    const playerPage = await context.newPage();
    await registerUser(playerPage, playerUsername, password);
    await joinRoom(playerPage, roomId);
    
    await setReady(playerPage);
    
    await page.waitForTimeout(2000);
    
    await expect(page.getByText(/1 \/ 2/i).or(page.getByText(/готовы/i))).toBeVisible({ timeout: 5000 });
    
    await playerPage.close();
  });

  test('редирект при запуске игры через SSE', async ({ page, context }) => {
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
    await playerPage.waitForURL(/\/game\/.+/, { timeout: 30000 });
    
    await expect(page).toHaveURL(/\/game\/.+/);
    await expect(playerPage).toHaveURL(/\/game\/.+/);
    
    await playerPage.close();
  });

  test('редирект при закрытии комнаты', async ({ page, context }) => {
    const hostUsername = generateUsername();
    const playerUsername = generateUsername();
    const password = 'testpass123';

    await registerUser(page, hostUsername, password);
    const roomId = await createRoom(page);
    
    const playerPage = await context.newPage();
    await registerUser(playerPage, playerUsername, password);
    await joinRoom(playerPage, roomId);
    
    await page.getByRole('button', { name: /покинуть комнату/i }).click();
    
    await playerPage.waitForTimeout(2000);
    
    await expect(playerPage).toHaveURL(/\/lobby/, { timeout: 10000 });
    
    await playerPage.close();
  });
});

