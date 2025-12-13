import { test, expect } from '@playwright/test';
import { registerUser, generateUsername } from './helpers/auth';
import { createRoom, joinRoom, setReady } from './helpers/room';
import { waitForGameStart, selectQuestion, pressButton, waitForJudgeButtons, waitForAnsweringPhase } from './helpers/game';

test.describe('Кнопки судейства', () => {
  test('появляются сразу после таймера ответа у ведущего', async ({ page, browser }) => {
    // Перехватываем все консольные логи и ошибки ДО загрузки страницы
    const consoleLogs: string[] = [];
    const consoleErrors: string[] = [];
    
    // Сохраняем логи в window для доступа из evaluate
    await page.addInitScript(() => {
      (window as any).__consoleLogs = [];
      (window as any).__consoleErrors = [];
      const originalLog = console.log;
      const originalError = console.error;
      const originalWarn = console.warn;
      console.log = (...args: any[]) => {
        const text = args.join(' ');
        (window as any).__consoleLogs.push(text);
        originalLog.apply(console, args);
      };
      console.error = (...args: any[]) => {
        const text = args.join(' ');
        (window as any).__consoleErrors.push(text);
        originalError.apply(console, args);
      };
      console.warn = (...args: any[]) => {
        const text = args.join(' ');
        (window as any).__consoleLogs.push(`[WARN] ${text}`);
        originalWarn.apply(console, args);
      };
    });
    
    page.on('console', msg => {
      const text = msg.text();
      consoleLogs.push(`[${msg.type()}] ${text}`);
      if (text.includes('[GamePage]') || text.includes('[useGameWebSocket]') || text.includes('ERROR') || text.includes('error') || text.includes('WebSocket')) {
        console.log(`[Browser Console] ${text}`);
      }
    });
    
    page.on('pageerror', error => {
      consoleErrors.push(error.message);
      console.log(`[Browser Error] ${error.message}`);
      console.log(`[Browser Error Stack] ${error.stack}`);
    });
    
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
    
    // Делаем скриншот перед выбором вопроса
    await page.screenshot({ path: 'test-results/before-select-question.png', fullPage: true });
    
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

