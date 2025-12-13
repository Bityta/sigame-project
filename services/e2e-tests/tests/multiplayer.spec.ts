import { test, expect } from '@playwright/test';
import { registerUser, generateUsername } from './helpers/auth';
import { createRoom, joinRoom, setReady } from './helpers/room';
import { waitForGameStart, selectQuestion, pressButton, waitForStatus } from './helpers/game';

test.describe('Множественные игроки', () => {
  test('работа с 3 игроками одновременно', async ({ page, browser }) => {
    const hostUsername = generateUsername();
    const player1Username = generateUsername();
    const player2Username = generateUsername();
    const password = 'testpass123';

    await registerUser(page, hostUsername, password);
    const roomId = await createRoom(page, 'Multiplayer Room', 4);
    
    const player1Context = await browser.newContext();
    const player1Page = await player1Context.newPage();
    await registerUser(player1Page, player1Username, password);
    await joinRoom(player1Page, roomId);
    
    const player2Context = await browser.newContext();
    const player2Page = await player2Context.newPage();
    await registerUser(player2Page, player2Username, password);
    await joinRoom(player2Page, roomId);
    
    await setReady(page);
    await setReady(player1Page);
    await setReady(player2Page);
    
    await waitForGameStart(page);
    await waitForGameStart(player1Page);
    await waitForGameStart(player2Page);
    
    await waitForStatus(page, 'question_select');
    
    await expect(page.getByText(hostUsername)).toBeVisible();
    await expect(page.getByText(player1Username)).toBeVisible();
    await expect(page.getByText(player2Username)).toBeVisible();
    
    await player1Page.close();
    await player1Context.close();
    await player2Page.close();
    await player2Context.close();
  });

  test('синхронизация состояния между игроками', async ({ page, browser }) => {
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
    await selectQuestion(page);
    
    await waitForStatus(playerPage, 'button_press', 10000);
    await expect(playerPage.locator('.game-board__question--disabled')).toBeVisible();
    
    await playerPage.close();
    await playerContext.close();
  });

  test('отображение действий других игроков в реальном времени', async ({ page, browser }) => {
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
    await selectQuestion(page);
    
    await waitForStatus(page, 'button_press', 10000);
    await pressButton(playerPage);
    
    await waitForStatus(page, 'answering', 10000);
    await expect(page.getByText(new RegExp(playerUsername, 'i')).or(page.locator('[class*="answering"]').filter({ hasText: new RegExp(playerUsername, 'i') }))).toBeVisible({ timeout: 10000 });
    
    await playerPage.close();
    await playerContext.close();
  });
});

