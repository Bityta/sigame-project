import { test, expect } from '@playwright/test';
import { registerUser, generateUsername } from './helpers/auth';
import { createRoom, joinRoom, setReady } from './helpers/room';
import { waitForGameStart, waitForStatus, selectQuestion, pressButton } from './helpers/game';

test.describe('Статусы игры', () => {
  test('переход через все статусы игры', async ({ page, context }) => {
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
    
    await waitForStatus(page, 'question_select', 15000);
    await expect(page.locator('.game-board')).toBeVisible();
    
    await selectQuestion(page);
    
    await waitForStatus(page, 'button_press', 10000);
    
    await pressButton(playerPage);
    
    await waitForStatus(page, 'answering', 10000);
    await expect(page.locator('.game-page__answering')).toBeVisible();
    
    await waitForStatus(page, 'answer_judging', 35000);
    await expect(page.locator('.game-page__judging')).toBeVisible();
    
    await playerPage.close();
  });

  test('отображение игрового поля', async ({ page, context }) => {
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
    
    await waitForStatus(page, 'question_select');
    
    await expect(page.locator('.game-board')).toBeVisible();
    await expect(page.locator('.game-board__theme')).toHaveCount(5);
    
    await playerPage.close();
  });

  test('отображение списка игроков', async ({ page, context }) => {
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
    
    await expect(page.locator('.player-list, [class*="player"]')).toBeVisible();
    
    await playerPage.close();
  });
});

