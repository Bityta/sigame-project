import { test, expect } from '@playwright/test';
import { registerUser, generateUsername } from './helpers/auth';
import { createRoom, joinRoom, setReady } from './helpers/room';
import { waitForGameStart, selectQuestion, pressButton, waitForJudgeButtons, waitForAnsweringPhase } from './helpers/game';

test.describe('Кнопки судейства', () => {
  test('появляются сразу после таймера ответа у ведущего', async ({ page, browser }) => {
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
    
    await selectQuestion(page);
    
    await pressButton(playerPage);
    
    await waitForAnsweringPhase(playerPage);
    
    const result = await waitForJudgeButtons(page, 35000);
    
    expect(result.appeared).toBe(true);
    expect(result.delay).toBeLessThan(35000);
    
    const judgeButtons = page.locator('.game-page__judging-buttons');
    await expect(judgeButtons).toBeVisible();
    
    const correctButton = page.locator('.game-page__judge-btn--correct');
    const wrongButton = page.locator('.game-page__judge-btn--wrong');
    
    await expect(correctButton).toBeVisible();
    await expect(wrongButton).toBeVisible();
    
    await playerPage.close();
    await playerContext.close();
  });
});

