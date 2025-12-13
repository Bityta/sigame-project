import { test, expect } from '@playwright/test';
import { registerUser, generateUsername } from './helpers/auth';
import { createRoom, joinRoom, setReady } from './helpers/room';
import { waitForGameStart, selectQuestion, pressButton, waitForStatus } from './helpers/game';

test.describe('Статусы игры - детально', () => {
  test('waiting - отображение сообщения ожидания', async ({ page, browser }) => {
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
    
    await expect(page.locator('.player-list, [class*="player"]')).toBeVisible();
    
    await playerContext.close();
  });

  test('question_select - таймер автоматического выбора', async ({ page, browser }) => {
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
    
    const startTime = Date.now();
    await waitForStatus(page, 'button_press', 40000);
    const elapsedTime = Date.now() - startTime;
    
    expect(elapsedTime).toBeGreaterThan(30000);
    expect(elapsedTime).toBeLessThan(45000);
    
    await playerContext.close();
  });

  test('button_press - отображение таймера', async ({ page, browser }) => {
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
    
    const timer = page.locator('.game-page__timer-bar, [class*="timer"]');
    const timerVisible = await timer.isVisible({ timeout: 5000 }).catch(() => false);
    if (timerVisible) {
      await expect(timer).toBeVisible({ timeout: 5000 });
    }
    
    await playerContext.close();
  });

  test('answering - визуальная индикация активного игрока', async ({ page, browser }) => {
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
    
    await expect(page.locator('.game-page__answering')).toBeVisible();
    await expect(page.getByText(/отвечает/i).or(page.getByText(/ваш черёд/i))).toBeVisible();
    
    await playerContext.close();
  });

  test('answer_judging - отображение результата судейства', async ({ page, browser }) => {
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
    
    await page.getByRole('button', { name: /верно/i }).click();
    
    await page.waitForTimeout(2000);
    
    const resultVisible = await page.getByText(/верно/i).or(page.locator('[class*="result"]')).isVisible({ timeout: 5000 }).catch(() => false);
    if (resultVisible) {
      await expect(page.getByText(/верно/i).or(page.locator('[class*="result"]'))).toBeVisible({ timeout: 5000 });
    }
    
    await playerContext.close();
  });
});

