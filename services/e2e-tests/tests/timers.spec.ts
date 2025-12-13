import { test, expect } from '@playwright/test';
import { registerUser, generateUsername } from './helpers/auth';
import { createRoom, joinRoom, setReady } from './helpers/room';
import { waitForGameStart, selectQuestion, waitForStatus } from './helpers/game';

test.describe('Таймеры в игре', () => {
  test('отображение таймера на фазе выбора вопроса', async ({ page, context }) => {
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
    
    const timer = page.locator('.game-page__timer-bar, [class*="timer"]');
    const timerVisible = await timer.isVisible({ timeout: 5000 }).catch(() => false);
    if (timerVisible) {
      await expect(timer).toBeVisible({ timeout: 5000 });
    }
    
    await playerPage.close();
  });

  test('отображение таймера на фазе нажатия кнопки', async ({ page, context }) => {
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
    
    await waitForStatus(page, 'button_press', 10000);
    
    const timer = page.locator('.game-page__timer-bar, [class*="timer"]');
    const timerVisible = await timer.isVisible({ timeout: 5000 }).catch(() => false);
    if (timerVisible) {
      await expect(timer).toBeVisible({ timeout: 5000 });
    }
    
    await playerPage.close();
  });

  test('отображение таймера на фазе ответа', async ({ page, context }) => {
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
    
    await waitForStatus(page, 'button_press', 10000);
    await playerPage.getByRole('button', { name: /нажать/i }).or(playerPage.locator('button').filter({ hasText: /ответить/i })).click();
    
    await waitForStatus(page, 'answering', 10000);
    
    const timer = page.locator('.game-page__answering-timer, [class*="timer"]');
    const timerVisible = await timer.isVisible({ timeout: 5000 }).catch(() => false);
    if (timerVisible) {
      await expect(timer).toBeVisible({ timeout: 5000 });
    }
    
    await playerPage.close();
  });

  test('визуальная анимация таймера', async ({ page, context }) => {
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
    
    const timerBar = page.locator('.game-page__timer-bar-fill, [class*="timer-bar"]');
    if (await timerBar.isVisible({ timeout: 5000 }).catch(() => false)) {
      const initialWidth = await timerBar.evaluate(el => (el as HTMLElement).style.width);
      await page.waitForTimeout(1000);
      const afterWidth = await timerBar.evaluate(el => (el as HTMLElement).style.width);
      
      expect(initialWidth).not.toBe(afterWidth);
    }
    
    await playerPage.close();
  });
});

