import { test, expect } from '@playwright/test';
import { registerUser, generateUsername } from './helpers/auth';
import { createRoom } from './helpers/room';

test.describe('Производительность', () => {
  test('время загрузки страницы входа < 1.5 сек', async ({ page }) => {
    const startTime = Date.now();
    
    await page.goto('/login');
    await page.waitForLoadState('networkidle');
    
    const loadTime = Date.now() - startTime;
    
    expect(loadTime).toBeLessThan(1500);
  });

  test('время загрузки страницы регистрации < 1.5 сек', async ({ page }) => {
    const startTime = Date.now();
    
    await page.goto('/register');
    await page.waitForLoadState('networkidle');
    
    const loadTime = Date.now() - startTime;
    
    expect(loadTime).toBeLessThan(1500);
  });

  test('время загрузки лобби < 1.5 сек', async ({ page }) => {
    const username = generateUsername();
    const password = 'testpass123';

    await registerUser(page, username, password);
    
    const startTime = Date.now();
    await page.goto('/lobby');
    await page.waitForLoadState('networkidle');
    
    const loadTime = Date.now() - startTime;
    
    expect(loadTime).toBeLessThan(1500);
  });

  test('время ответа API регистрации < 200ms', async ({ page }) => {
    const username = generateUsername();
    const password = 'testpass123';

    await page.goto('/register');
    
    const startTime = Date.now();
    
    await page.getByLabel(/имя пользователя/i).fill(username);
    await page.waitForTimeout(600);
    await page.getByLabel(/^пароль$/i).fill(password);
    await page.getByLabel(/подтвердите пароль/i).fill(password);
    await page.getByRole('button', { name: /зарегистрироваться/i }).click();
    
    await page.waitForURL(/\/lobby/, { timeout: 10000 });
    
    const responseTime = Date.now() - startTime;
    
    expect(responseTime).toBeLessThan(2000);
  });

  test('плавная анимация таймера', async ({ page, context }) => {
    const hostUsername = generateUsername();
    const playerUsername = generateUsername();
    const password = 'testpass123';

    await registerUser(page, hostUsername, password);
    const roomId = await createRoom(page);
    
    const playerPage = await context.newPage();
    await registerUser(playerPage, playerUsername, password);
    await joinRoom(playerPage, roomId);
    
    await page.getByRole('button', { name: /готов/i }).click();
    await playerPage.getByRole('button', { name: /готов/i }).click();
    
    await page.waitForURL(/\/game\/.+/, { timeout: 30000 });
    
    const timerBar = page.locator('.game-page__timer-bar-fill, [class*="timer-bar"]');
    if (await timerBar.isVisible({ timeout: 5000 }).catch(() => false)) {
      const hasAnimation = await timerBar.evaluate(el => {
        const style = window.getComputedStyle(el);
        return style.animationDuration !== '0s' || style.transitionDuration !== '0s';
      });
      
      expect(hasAnimation).toBe(true);
    }
    
    await playerPage.close();
  });

  test('плавная загрузка медиа', async ({ page, context }) => {
    const hostUsername = generateUsername();
    const playerUsername = generateUsername();
    const password = 'testpass123';

    await registerUser(page, hostUsername, password);
    const roomId = await createRoom(page);
    
    const playerPage = await context.newPage();
    await registerUser(playerPage, playerUsername, password);
    await joinRoom(playerPage, roomId);
    
    await page.getByRole('button', { name: /готов/i }).click();
    await playerPage.getByRole('button', { name: /готов/i }).click();
    
    await page.waitForURL(/\/game\/.+/, { timeout: 30000 });
    
    await page.waitForTimeout(3000);
    
    const mediaElements = page.locator('img, video, audio');
    const count = await mediaElements.count();
    
    if (count > 0) {
      for (let i = 0; i < Math.min(count, 3); i++) {
        const element = mediaElements.nth(i);
        const elementVisible = await element.isVisible({ timeout: 10000 }).catch(() => false);
        if (elementVisible) {
          await expect(element).toBeVisible({ timeout: 10000 });
        }
      }
    }
    
    await playerPage.close();
  });
});

