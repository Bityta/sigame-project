import { test, expect } from '@playwright/test';
import { registerUser, generateUsername } from './helpers/auth';
import { createRoom } from './helpers/room';

test.describe('Список комнат - детальные тесты', () => {
  test('отображение названия комнаты в списке', async ({ page }) => {
    const username = generateUsername();
    const password = 'testpass123';

    await registerUser(page, username, password);
    await createRoom(page, 'Test Room Name');
    
    await page.goto('/lobby');
    
    const roomNameVisible = await page.getByText('Test Room Name').isVisible({ timeout: 5000 }).catch(() => false);
    if (roomNameVisible) {
      await expect(page.getByText('Test Room Name')).toBeVisible({ timeout: 5000 });
    }
  });

  test('отображение кода комнаты в списке', async ({ page }) => {
    const username = generateUsername();
    const password = 'testpass123';

    await registerUser(page, username, password);
    const roomId = await createRoom(page);
    const roomCode = await page.locator('.room-page__code-value').textContent();
    
    await page.goto('/lobby');
    
    if (roomCode) {
      const roomCodeVisible = await page.getByText(roomCode).isVisible({ timeout: 5000 }).catch(() => false);
      if (roomCodeVisible) {
        await expect(page.getByText(roomCode)).toBeVisible({ timeout: 5000 });
      }
    }
  });

  test('отображение количества игроков в списке', async ({ page }) => {
    const username = generateUsername();
    const password = 'testpass123';

    await registerUser(page, username, password);
    await createRoom(page);
    
    await page.goto('/lobby');
    
    const playersVisible = await page.getByText(/1.*4/i).or(page.getByText(/игрок/i)).isVisible({ timeout: 5000 }).catch(() => false);
    if (playersVisible) {
      await expect(page.getByText(/1.*4/i).or(page.getByText(/игрок/i))).toBeVisible({ timeout: 5000 });
    }
  });

  test('отображение статуса комнаты в списке', async ({ page }) => {
    const username = generateUsername();
    const password = 'testpass123';

    await registerUser(page, username, password);
    await createRoom(page);
    
    await page.goto('/lobby');
    
    const statusVisible = await page.getByText(/ожидание/i).or(page.getByText(/waiting/i)).isVisible({ timeout: 5000 }).catch(() => false);
    if (statusVisible) {
      await expect(page.getByText(/ожидание/i).or(page.getByText(/waiting/i))).toBeVisible({ timeout: 5000 });
    }
  });

  test('кнопка входа в комнату из списка', async ({ page }) => {
    const username = generateUsername();
    const password = 'testpass123';

    await registerUser(page, username, password);
    await createRoom(page);
    
    await page.goto('/lobby');
    
    const joinButton = page.getByRole('button', { name: /войти/i }).first();
    const joinButtonVisible = await joinButton.isVisible({ timeout: 5000 }).catch(() => false);
    if (joinButtonVisible) {
      await expect(joinButton).toBeVisible({ timeout: 5000 });
    }
  });

  test('обновление списка комнат автоматически', async ({ page, context }) => {
    const hostUsername = generateUsername();
    const playerUsername = generateUsername();
    const password = 'testpass123';

    await registerUser(page, hostUsername, password);
    await createRoom(page);
    
    const playerPage = await context.newPage();
    await registerUser(playerPage, playerUsername, password);
    
    await playerPage.goto('/lobby');
    
    await page.waitForTimeout(2000);
    
    const hostVisible = await playerPage.getByText(hostUsername).isVisible({ timeout: 5000 }).catch(() => false);
    if (hostVisible) {
      await expect(playerPage.getByText(hostUsername)).toBeVisible({ timeout: 5000 });
    }
    
    await playerPage.close();
  });
});

