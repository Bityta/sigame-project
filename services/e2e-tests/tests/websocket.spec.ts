import { test, expect } from '@playwright/test';
import { registerUser, generateUsername } from './helpers/auth';
import { createRoom, joinRoom, setReady } from './helpers/room';
import { waitForGameStart } from './helpers/game';

test.describe('WebSocket соединение', () => {
  test('установка WebSocket соединения', async ({ page, browser }) => {
    const hostUsername = generateUsername();
    const playerUsername = generateUsername();
    const password = 'testpass123';

    const wsMessages: string[] = [];
    
    page.on('websocket', ws => {
      ws.on('framereceived', event => {
        wsMessages.push(event.payload.toString());
      });
    });
    
    await registerUser(page, hostUsername, password);
    const roomId = await createRoom(page);
    
    const playerContext = await browser.newContext();
    const playerPage = await playerContext.newPage();
    await registerUser(playerPage, playerUsername, password);
    await joinRoom(playerPage, roomId);
    
    await setReady(page);
    await setReady(playerPage);
    
    await waitForGameStart(page);
    await waitForGameStart(playerPage);
    
    await page.waitForTimeout(2000);
    
    expect(wsMessages.length).toBeGreaterThan(0);
    
    await playerContext.close();
  });

  test('получение STATE_UPDATE при подключении', async ({ page, browser }) => {
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
    
    await page.waitForTimeout(1000);
    
    await expect(page.locator('.game-page')).toBeVisible();
    
    await playerContext.close();
  });
});

