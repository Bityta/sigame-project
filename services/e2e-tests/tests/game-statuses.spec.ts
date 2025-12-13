import { test, expect } from '@playwright/test';
import { registerUser, generateUsername } from './helpers/auth';
import { createRoom, joinRoom, setReady } from './helpers/room';
import { waitForGameStart, waitForStatus, selectQuestion, pressButton } from './helpers/game';

test.describe('Статусы игры', () => {
  test('переход через все статусы игры', async ({ page, browser }) => {
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
    await waitForGameStart(playerPage);
    
    await waitForStatus(page, 'question_select', 15000);
    await expect(page.locator('.game-board')).toBeVisible();
    
    await selectQuestion(page);
    
    await waitForStatus(page, 'button_press', 10000);
    
    await pressButton(playerPage);
    
    await waitForStatus(page, 'answering', 10000);
    await expect(page.locator('.game-page__answering')).toBeVisible();
    
    await waitForStatus(page, 'answer_judging', 35000);
    await expect(page.locator('.game-page__judging')).toBeVisible();
    
    await playerContext.close();
  });

  test('отображение игрового поля', async ({ page, browser }) => {
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
    await waitForGameStart(playerPage);
    
    await waitForStatus(page, 'question_select');
    
    await expect(page.locator('.game-board')).toBeVisible();
    await expect(page.locator('.game-board__theme')).toHaveCount(5);
    
    await playerContext.close();
  });

  test('отображение списка игроков', async ({ page, browser }) => {
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
    await waitForGameStart(playerPage);
    
    await expect(page.locator('.player-list, [class*="player"]')).toBeVisible();
    
    await playerContext.close();
  });
});

