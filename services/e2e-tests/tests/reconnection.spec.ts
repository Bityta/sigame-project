import { test, expect } from '@playwright/test';
import { registerUser, generateUsername } from './helpers/auth';
import { createRoom, joinRoom, setReady } from './helpers/room';
import { waitForGameStart, waitForStatus } from './helpers/game';

test.describe('Переподключение', () => {
  test('восстановление состояния игры после перезагрузки', async ({ page, browser }) => {
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
    
    await waitForGameStart(page);
    await waitForStatus(page, 'question_select');
    
    const gameUrl = page.url();
    
    await page.reload();
    
    await expect(page).toHaveURL(gameUrl);
    await page.waitForTimeout(2000);
    await expect(page.locator('.game-page')).toBeVisible({ timeout: 10000 });
    
    await playerContext.close();
  });

  test('переподключение WebSocket при разрыве соединения', async ({ page, browser }) => {
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
    
    await waitForGameStart(page);
    
    await page.evaluate(() => {
      if (window.WebSocket) {
        const ws = new WebSocket('ws://invalid');
        ws.close();
      }
    });
    
    await page.waitForTimeout(2000);
    
    await expect(page.locator('.game-page')).toBeVisible({ timeout: 10000 });
    
    await playerContext.close();
  });
});

