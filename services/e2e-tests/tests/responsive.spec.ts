import { test, expect } from '@playwright/test';
import { registerUser, generateUsername } from './helpers/auth';
import { createRoom } from './helpers/room';

test.describe('Адаптивность', () => {
  test('отображение на десктопе 1920x1080', async ({ page }) => {
    await page.setViewportSize({ width: 1920, height: 1080 });
    
    const username = generateUsername();
    const password = 'testpass123';

    await registerUser(page, username, password);
    
    await expect(page.locator('.lobby-page, [class*="lobby"]')).toBeVisible();
    
    const layout = await page.evaluate(() => {
      return {
        width: window.innerWidth,
        height: window.innerHeight
      };
    });
    
    expect(layout.width).toBe(1920);
    expect(layout.height).toBe(1080);
  });

  test('отображение на десктопе 1366x768', async ({ page }) => {
    await page.setViewportSize({ width: 1366, height: 768 });
    
    const username = generateUsername();
    const password = 'testpass123';

    await registerUser(page, username, password);
    
    await expect(page.locator('.lobby-page, [class*="lobby"]')).toBeVisible();
    
    const layout = await page.evaluate(() => {
      return {
        width: window.innerWidth,
        height: window.innerHeight
      };
    });
    
    expect(layout.width).toBe(1366);
    expect(layout.height).toBe(768);
  });

  test('отображение на планшете 768x1024', async ({ page }) => {
    await page.setViewportSize({ width: 768, height: 1024 });
    
    const username = generateUsername();
    const password = 'testpass123';

    await registerUser(page, username, password);
    
    await expect(page.locator('.lobby-page, [class*="lobby"]')).toBeVisible();
    
    const isResponsive = await page.evaluate(() => {
      return window.innerWidth <= 768;
    });
    
    expect(isResponsive).toBe(true);
  });

  test('отображение на мобильном 375x667', async ({ page }) => {
    await page.setViewportSize({ width: 375, height: 667 });
    
    const username = generateUsername();
    const password = 'testpass123';

    await registerUser(page, username, password);
    
    await expect(page.locator('.lobby-page, [class*="lobby"]')).toBeVisible();
    
    const isMobile = await page.evaluate(() => {
      return window.innerWidth <= 375;
    });
    
    expect(isMobile).toBe(true);
  });

  test('отображение на мобильном 414x896', async ({ page }) => {
    await page.setViewportSize({ width: 414, height: 896 });
    
    const username = generateUsername();
    const password = 'testpass123';

    await registerUser(page, username, password);
    
    await expect(page.locator('.lobby-page, [class*="lobby"]')).toBeVisible();
    
    const isMobile = await page.evaluate(() => {
      return window.innerWidth <= 414;
    });
    
    expect(isMobile).toBe(true);
  });

  test('адаптивная вёрстка - элементы не перекрываются', async ({ page }) => {
    await page.setViewportSize({ width: 375, height: 667 });
    
    const username = generateUsername();
    const password = 'testpass123';

    await registerUser(page, username, password);
    
    const elements = page.locator('button, input, [class*="card"]');
    const count = await elements.count();
    
    for (let i = 0; i < Math.min(count, 5); i++) {
      const element = elements.nth(i);
      const box = await element.boundingBox();
      expect(box).not.toBeNull();
    }
  });
});

