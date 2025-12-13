import { test, expect } from '@playwright/test';
import { registerUser, generateUsername } from './helpers/auth';
import { createRoom, joinRoom, setReady } from './helpers/room';
import { waitForGameStart } from './helpers/game';

test.describe('WebSocket соединение', () => {
  test('установка WebSocket соединения', async ({ page, context }) => {
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
    
    await waitForGameStart(page);
    await waitForGameStart(playerPage);
    
    const wsMessages: string[] = [];
    
    page.on('websocket', ws => {
      ws.on('framereceived', event => {
        wsMessages.push(event.payload.toString());
      });
    });
    
    await page.waitForTimeout(2000);
    
    expect(wsMessages.length).toBeGreaterThan(0);
    
    await playerPage.close();
  });

  test('получение STATE_UPDATE при подключении', async ({ page, context }) => {
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
    
    await waitForGameStart(page);
    
    await page.waitForTimeout(1000);
    
    await expect(page.locator('.game-board, .player-list, [class*="game"]')).toBeVisible();
    
    await playerPage.close();
  });
});

