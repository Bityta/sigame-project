import { test, expect } from '@playwright/test';
import { registerUser, generateUsername } from './helpers/auth';
import { createRoom, joinRoom, setReady, getRoomCode } from './helpers/room';

test.describe('Комната ожидания', () => {
  test('отображение информации о комнате', async ({ page }) => {
    const username = generateUsername();
    const password = 'testpass123';

    await registerUser(page, username, password);
    const roomId = await createRoom(page);
    
    await expect(page.locator('.room-page__title')).toBeVisible();
    await expect(page.locator('.room-page__code-value')).toBeVisible();
  });

  test('копирование кода комнаты', async ({ page }) => {
    const username = generateUsername();
    const password = 'testpass123';

    await registerUser(page, username, password);
    await createRoom(page);
    
    await expect(page).toHaveURL(/\/room\/.+/);
    await expect(page.locator('.room-page__title')).toBeVisible();
    await page.waitForTimeout(500);
    
    const codeElement = page.locator('.room-page__code');
    await expect(codeElement).toBeVisible();
    
    await page.context().grantPermissions(['clipboard-read', 'clipboard-write']);
    
    await codeElement.hover();
    await page.waitForTimeout(200);
    await codeElement.click({ delay: 100 });
    
    await page.waitForFunction(
      () => {
        const element = document.querySelector('.room-page__code');
        return element?.classList.contains('room-page__code--copied') || false;
      },
      { timeout: 10000 }
    );
    
    await expect(codeElement).toHaveClass(/room-page__code--copied/);
    
    await page.waitForTimeout(500);
  });

  test('готовность игрока', async ({ page, browser }) => {
    const hostUsername = generateUsername();
    const playerUsername = generateUsername();
    const password = 'testpass123';

    await registerUser(page, hostUsername, password);
    const roomId = await createRoom(page);
    
    await page.waitForTimeout(500);
    
    const playerContext = await browser.newContext();
    const playerPage = await playerContext.newPage();
    await registerUser(playerPage, playerUsername, password);
    await joinRoom(playerPage, roomId);
    
    await page.waitForTimeout(500);
    const readyButton = page.getByRole('button', { name: /готов/i });
    await readyButton.hover();
    await page.waitForTimeout(200);
    await readyButton.click({ delay: 100 });
    
    await expect(page.getByText(/вы готовы/i)).toBeVisible();
    
    await page.waitForTimeout(500);
    await playerContext.close();
  });

  test('автоматический запуск игры когда все готовы', async ({ page, browser }) => {
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
    
    await page.waitForURL(/\/game\/.+/, { timeout: 30000 });
    await expect(page).toHaveURL(/\/game\/.+/);
    
    await playerContext.close();
  });

  test('настройки комнаты для хоста', async ({ page }) => {
    const username = generateUsername();
    const password = 'testpass123';

    await registerUser(page, username, password);
    await createRoom(page);
    
    await page.waitForTimeout(500);
    
    const slider = page.locator('.room-settings__slider').first();
    await expect(slider).toBeVisible();
    
    await slider.hover();
    await page.waitForTimeout(200);
    await slider.fill('45', { delay: 100 });
    await page.waitForTimeout(300);
    
    const saveButton = page.getByRole('button', { name: /сохранить настройки/i });
    await saveButton.hover();
    await page.waitForTimeout(200);
    await saveButton.click({ delay: 100 });
    
    await expect(page.getByText(/45 сек/i)).toBeVisible({ timeout: 5000 });
    await page.waitForTimeout(500);
  });
});

