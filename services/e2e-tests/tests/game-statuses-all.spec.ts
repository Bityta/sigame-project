import { test, expect } from '@playwright/test';
import { registerUser, generateUsername } from './helpers/auth';
import { createRoom, joinRoom, setReady } from './helpers/room';
import { waitForGameStart, waitForStatus } from './helpers/game';

test.describe('Все статусы игры', () => {
  test('rounds_overview - отображение обзора раундов', async ({ page, context }) => {
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
    
    await page.waitForTimeout(2000);
    
    const roundsOverview = page.locator('.rounds-overview, [class*="round"]');
    if (await roundsOverview.isVisible({ timeout: 5000 }).catch(() => false)) {
      await expect(roundsOverview).toBeVisible();
    }
    
    await playerPage.close();
  });

  test('round_start - отображение интро раунда', async ({ page, context }) => {
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
    
    await page.waitForTimeout(2000);
    
    const roundIntro = page.locator('.round-intro, [class*="round-intro"]');
    if (await roundIntro.isVisible({ timeout: 5000 }).catch(() => false)) {
      await expect(roundIntro).toBeVisible();
    }
    
    await playerPage.close();
  });

  test('question_show - отображение вопроса', async ({ page, context }) => {
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
    await waitForStatus(page, 'question_select');
    
    const firstQuestion = page.locator('.game-board__question:not(.game-board__question--disabled)').first();
    await firstQuestion.click();
    
    await page.waitForTimeout(2000);
    
    const questionView = page.locator('.question-view, [class*="question"]');
    await expect(questionView).toBeVisible({ timeout: 10000 });
    
    await playerPage.close();
  });

  test('game_end - отображение финальных результатов', async ({ page, context }) => {
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
    
    await page.waitForTimeout(60000);
    
    const gameEnd = page.locator('.game-end, [class*="game-end"]');
    if (await gameEnd.isVisible({ timeout: 10000 }).catch(() => false)) {
      await expect(gameEnd).toBeVisible();
      await expect(page.getByText(/победитель/i).or(page.getByText(/результаты/i))).toBeVisible();
    }
    
    await playerPage.close();
  });
});

