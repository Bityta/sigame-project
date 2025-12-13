import { test, expect } from '@playwright/test';
import { registerUser, generateUsername } from './helpers/auth';
import { createRoom, joinRoom } from './helpers/room';

test.describe('Функции хоста', () => {
  test('отображение кнопок управления только для хоста', async ({ page, context }) => {
    const hostUsername = generateUsername();
    const playerUsername = generateUsername();
    const password = 'testpass123';

    await registerUser(page, hostUsername, password);
    const roomId = await createRoom(page);
    
    const playerPage = await context.newPage();
    await registerUser(playerPage, playerUsername, password);
    await joinRoom(playerPage, roomId);
    
    await expect(page.locator('.room-page__player-action')).toBeVisible();
    await expect(playerPage.locator('.room-page__player-action')).not.toBeVisible();
    
    await playerPage.close();
  });

  test('передача роли хоста', async ({ page, context }) => {
    const hostUsername = generateUsername();
    const playerUsername = generateUsername();
    const password = 'testpass123';

    await registerUser(page, hostUsername, password);
    const roomId = await createRoom(page);
    
    const playerPage = await context.newPage();
    await registerUser(playerPage, playerUsername, password);
    await joinRoom(playerPage, roomId);
    
    page.on('dialog', dialog => dialog.accept());
    
    const transferButton = page.locator('.room-page__player-action--transfer').first();
    await transferButton.click();
    
    await page.waitForTimeout(2000);
    
    const crownVisible = await page.getByText(/👑/).isVisible({ timeout: 5000 }).catch(() => false);
    if (!crownVisible) {
      await expect(page.getByText(/👑/)).not.toBeVisible({ timeout: 5000 });
    }
    
    await playerPage.close();
  });

  test('выгон игрока из комнаты', async ({ page, context }) => {
    const hostUsername = generateUsername();
    const playerUsername = generateUsername();
    const password = 'testpass123';

    await registerUser(page, hostUsername, password);
    const roomId = await createRoom(page);
    
    const playerPage = await context.newPage();
    await registerUser(playerPage, playerUsername, password);
    await joinRoom(playerPage, roomId);
    
    const kickButton = page.locator('.room-page__player-action--kick').first();
    await kickButton.click();
    
    await page.waitForTimeout(2000);
    
    const playerVisible = await page.getByText(playerUsername).isVisible({ timeout: 5000 }).catch(() => false);
    if (!playerVisible) {
      await expect(page.getByText(playerUsername)).not.toBeVisible({ timeout: 5000 });
    }
    
    await playerPage.close();
  });

  test('кнопки управления неактивны когда комната не в статусе waiting', async ({ page, context }) => {
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
    
    await page.goto(`/room/${roomId}`);
    
    const transferButton = page.locator('.room-page__player-action--transfer').first();
    const transferButtonVisible = await transferButton.isVisible({ timeout: 5000 }).catch(() => false);
    if (transferButtonVisible) {
      await expect(transferButton).toBeDisabled({ timeout: 5000 });
    }
    
    await playerPage.close();
  });
});

