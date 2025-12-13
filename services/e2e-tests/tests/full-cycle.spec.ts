import { test, expect } from '@playwright/test';
import { registerUser, generateUsername } from './helpers/auth';
import { createRoom, joinRoom, setReady } from './helpers/room';
import { waitForGameStart, selectQuestion, pressButton, judgeAnswer, waitForStatus } from './helpers/game';

test.describe('Полный цикл игры', () => {
  test('регистрация -> создание комнаты -> игра -> судейство -> следующий вопрос', async ({ page, browser }) => {
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
    
    await waitForStatus(page, 'answer_judging', 35000);
    await judgeAnswer(page, true);
    
    await waitForStatus(page, 'question_select', 10000);
    await expect(page.locator('.game-board')).toBeVisible();
    
    await playerPage.close();
    await playerContext.close();
  });

  test('несколько вопросов подряд', async ({ page, browser }) => {
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
    
    for (let i = 0; i < 3; i++) {
      await waitForStatus(page, 'question_select', 10000);
      await selectQuestion(page);
      
      await waitForStatus(page, 'button_press', 10000);
      await pressButton(playerPage);
      
      await waitForStatus(page, 'answering', 10000);
      await waitForStatus(page, 'answer_judging', 35000);
      
      await judgeAnswer(page, i % 2 === 0);
      
      await page.waitForTimeout(2000);
    }
    
    await playerPage.close();
    await playerContext.close();
  });
});

