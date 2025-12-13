import { test, expect } from '@playwright/test';
import { registerUser, generateUsername } from './helpers/auth';
import { createRoom, joinRoom, setReady } from './helpers/room';
import { waitForGameStart, selectQuestion, pressButton, judgeAnswer, waitForStatus, waitForAnsweringPhase } from './helpers/game';

test.describe('Полный цикл игры', () => {
  test('успешное прохождение раунда', async ({ page, browser }) => {
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
    
    await waitForStatus(page, 'button_press');
    await pressButton(playerPage);
    
    await waitForAnsweringPhase(playerPage);
    
    await waitForStatus(page, 'answer_judging', 35000);
    
    await judgeAnswer(page, true);
    
    await waitForStatus(page, 'question_select', 10000);
    
    await playerPage.close();
    await playerContext.close();
  });
});

