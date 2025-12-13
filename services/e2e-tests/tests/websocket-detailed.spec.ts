import { test, expect } from '@playwright/test';
import { registerUser, generateUsername } from './helpers/auth';
import { createRoom, joinRoom, setReady } from './helpers/room';
import { waitForGameStart, selectQuestion, pressButton, judgeAnswer, waitForStatus } from './helpers/game';

test.describe('WebSocket - детальные тесты', () => {
  test('отправка SELECT_QUESTION через WebSocket', async ({ page, browser }) => {
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
    
    const wsMessages: string[] = [];
    page.on('websocket', ws => {
      ws.on('framesent', event => {
        if (event.payload.toString().includes('SELECT_QUESTION')) {
          wsMessages.push(event.payload.toString());
        }
      });
    });
    
    await selectQuestion(page);
    
    await page.waitForTimeout(1000);
    
    expect(wsMessages.length).toBeGreaterThan(0);
    
    await playerContext.close();
  });

  test('отправка PRESS_BUTTON через WebSocket', async ({ page, browser }) => {
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
    await selectQuestion(page);
    await waitForStatus(page, 'button_press', 10000);
    
    const wsMessages: string[] = [];
    playerPage.on('websocket', ws => {
      ws.on('framesent', event => {
        if (event.payload.toString().includes('PRESS_BUTTON')) {
          wsMessages.push(event.payload.toString());
        }
      });
    });
    
    await pressButton(playerPage);
    
    await page.waitForTimeout(1000);
    
    expect(wsMessages.length).toBeGreaterThan(0);
    
    await playerContext.close();
  });

  test('отправка JUDGE_ANSWER через WebSocket', async ({ page, browser }) => {
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
    await selectQuestion(page);
    await waitForStatus(page, 'button_press', 10000);
    await pressButton(playerPage);
    await waitForStatus(page, 'answering', 10000);
    await waitForStatus(page, 'answer_judging', 35000);
    
    const wsMessages: string[] = [];
    page.on('websocket', ws => {
      ws.on('framesent', event => {
        if (event.payload.toString().includes('JUDGE_ANSWER')) {
          wsMessages.push(event.payload.toString());
        }
      });
    });
    
    await judgeAnswer(page, true);
    
    await page.waitForTimeout(1000);
    
    expect(wsMessages.length).toBeGreaterThan(0);
    
    await playerContext.close();
  });

  test('получение PING и отправка PONG', async ({ page, browser }) => {
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
    
    const pingReceived: boolean[] = [];
    const pongSent: boolean[] = [];
    
    page.on('websocket', ws => {
      ws.on('framereceived', event => {
        if (event.payload.toString().includes('PING')) {
          pingReceived.push(true);
        }
      });
      ws.on('framesent', event => {
        if (event.payload.toString().includes('PONG')) {
          pongSent.push(true);
        }
      });
    });
    
    await page.waitForTimeout(6000);
    
    expect(pingReceived.length).toBeGreaterThan(0);
    expect(pongSent.length).toBeGreaterThan(0);
    
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
    
    const stateUpdates: string[] = [];
    page.on('websocket', ws => {
      ws.on('framereceived', event => {
        if (event.payload.toString().includes('STATE_UPDATE')) {
          stateUpdates.push(event.payload.toString());
        }
      });
    });
    
    await waitForGameStart(page);
    
    await page.waitForTimeout(2000);
    
    expect(stateUpdates.length).toBeGreaterThan(0);
    
    await playerContext.close();
  });

  test('обработка ошибок WebSocket', async ({ page, browser }) => {
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
    
    const wsErrors: Error[] = [];
    const wsClosed: boolean[] = [];
    
    page.on('websocket', ws => {
      ws.on('socketerror', error => {
        wsErrors.push(error);
      });
      ws.on('close', () => {
        wsClosed.push(true);
      });
    });
    
    await page.waitForTimeout(5000);
    
    await expect(page.locator('.game-board, [class*="game"]')).toBeVisible({ timeout: 10000 });
    
    await playerContext.close();
  });
});

