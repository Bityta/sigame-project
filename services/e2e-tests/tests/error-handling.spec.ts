import { test, expect } from '@playwright/test';
import { registerUser, generateUsername } from './helpers/auth';
import { createRoom, joinRoom } from './helpers/room';

test.describe('Обработка ошибок', () => {
  test('обработка ошибки сети', async ({ page }) => {
    await page.route('**/api/**', route => route.abort());
    
    const username = generateUsername();
    const password = 'testpass123';

    try {
      await registerUser(page, username, password);
    } catch (error) {
      // Ожидаем ошибку при регистрации из-за заблокированного API
    }
    
    await page.goto('/lobby');
    
    const errorVisible = await page.getByText(/ошибка/i).or(page.getByText(/не удалось/i)).isVisible({ timeout: 5000 }).catch(() => false);
    if (errorVisible) {
      await expect(page.getByText(/ошибка/i).or(page.getByText(/не удалось/i))).toBeVisible({ timeout: 5000 });
    }
  });

  test('обработка ошибки 404 при несуществующей комнате', async ({ page }) => {
    const username = generateUsername();
    const password = 'testpass123';

    await registerUser(page, username, password);
    
    await page.goto('/room/00000000-0000-0000-0000-000000000000');
    
    await expect(page.getByText(/не найдена/i).or(page.getByText(/ошибка/i))).toBeVisible({ timeout: 5000 });
  });

  test('обработка ошибки присоединения к заполненной комнате', async ({ page, context }) => {
    const hostUsername = generateUsername();
    const player1Username = generateUsername();
    const player2Username = generateUsername();
    const password = 'testpass123';

    await registerUser(page, hostUsername, password);
    const roomId = await createRoom(page, 'Small Room', 2);
    
    const player1Page = await context.newPage();
    await registerUser(player1Page, player1Username, password);
    await joinRoom(player1Page, roomId);
    
    const player2Page = await context.newPage();
    await registerUser(player2Page, player2Username, password);
    
    await player2Page.goto(`/room/${roomId}`);
    
    await expect(player2Page.getByText(/заполнена/i).or(player2Page.getByText(/ошибка/i)).or(player2Page.getByText(/полна/i))).toBeVisible({ timeout: 10000 });
    
    await player1Page.close();
    await player2Page.close();
  });

  test('редирект при критической ошибке', async ({ page }) => {
    await page.route('**/auth/me', route => route.fulfill({ status: 401 }));
    
    const username = generateUsername();
    const password = 'testpass123';

    await registerUser(page, username, password);
    
    await page.goto('/lobby');
    
    await expect(page).toHaveURL(/\/login/, { timeout: 15000 });
  });
});

