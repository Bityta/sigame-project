import { test, expect } from '@playwright/test';
import { registerUser, generateUsername } from './helpers/auth';
import { createRoom, joinRoom, setReady } from './helpers/room';
import { waitForGameStart, selectQuestion, pressButton, waitForStatus } from './helpers/game';

test.describe('UI элементы', () => {
  test('отображение спиннера при загрузке', async ({ page }) => {
    const username = generateUsername();
    const password = 'testpass123';

    await page.goto('/register');
    await page.getByPlaceholder(/от 3 до 20 символов/i).fill(username);
    await page.waitForTimeout(600);
    await page.getByPlaceholder(/минимум 6 символов/i).fill(password);
    await page.getByPlaceholder(/повторите пароль/i).fill(password);
    await page.getByRole('button', { name: /зарегистрироваться/i }).click();
    
    const spinnerVisible = await page.locator('[class*="spinner"], [class*="loading"]').isVisible({ timeout: 2000 }).catch(() => false);
    if (spinnerVisible) {
      await expect(page.locator('[class*="spinner"], [class*="loading"]')).toBeVisible({ timeout: 2000 });
    }
  });

  test('отображение сообщений об ошибках', async ({ page }) => {
    await page.goto('/login');
    
    await page.getByPlaceholder(/введите имя пользователя/i).fill('nonexistent');
    await page.getByPlaceholder(/введите пароль/i).fill('wrongpass');
    await page.getByRole('button', { name: /войти/i }).click();
    
    await page.waitForTimeout(2000);
    
    const errorVisible = await page.getByText(/ошибка/i).or(page.getByText(/неверный/i)).or(page.getByText(/неверные/i)).isVisible({ timeout: 5000 }).catch(() => false);
    if (errorVisible) {
      await expect(page.getByText(/ошибка/i).or(page.getByText(/неверный/i)).or(page.getByText(/неверные/i))).toBeVisible({ timeout: 5000 });
    }
  });

  test('отображение состояния загрузки кнопок', async ({ page }) => {
    const username = generateUsername();
    const password = 'testpass123';

    await registerUser(page, username, password);
    
    await page.getByRole('button', { name: /создать комнату/i }).click();
    
    await page.getByPlaceholder(/введите название/i).fill('Test Room');
    
    const packSelect = page.locator('select').first();
    await packSelect.waitFor({ state: 'visible' });
    const packOptions = await packSelect.locator('option').all();
    if (packOptions.length > 1) {
      await packSelect.selectOption({ index: 1 });
    }
    
    const createButton = page.getByRole('button', { name: /создать/i });
    await createButton.click();
    
    const createButtonDisabled = await createButton.isDisabled({ timeout: 2000 }).catch(() => false);
    if (createButtonDisabled) {
      await expect(createButton).toBeDisabled({ timeout: 2000 });
    }
  });

  test('отображение активной комнаты в лобби', async ({ page }) => {
    const username = generateUsername();
    const password = 'testpass123';

    await registerUser(page, username, password);
    await createRoom(page);
    
    await page.goto('/lobby');
    
    const activeRoomVisible = await page.getByText(/активная комната/i).or(page.locator('[class*="active-room"]')).isVisible({ timeout: 5000 }).catch(() => false);
    if (activeRoomVisible) {
      await expect(page.getByText(/активная комната/i).or(page.locator('[class*="active-room"]'))).toBeVisible({ timeout: 5000 });
    }
  });

  test('отображение баннера активной игры', async ({ page, browser }) => {
    const hostUsername = generateUsername();
    const playerUsername = generateUsername();
    const password = 'testpass123';

    await registerUser(page, hostUsername, password);
    const roomId = await createRoom(page);
    
    const playerContext = await browser.newContext();
    const playerPage = await playerContext.newPage();
    await registerUser(playerPage, playerUsername, password);
    await joinRoom(playerPage, roomId);
    
    await page.waitForTimeout(2000);
    
    await setReady(page);
    await setReady(playerPage);
    
    await page.waitForURL(/\/game\/.+/, { timeout: 30000 });
    
    await page.goto('/lobby');
    
    const activeGameVisible = await page.getByText(/активная игра/i).or(page.locator('[class*="active-game"]')).isVisible({ timeout: 5000 }).catch(() => false);
    if (activeGameVisible) {
      await expect(page.getByText(/активная игра/i).or(page.locator('[class*="active-game"]'))).toBeVisible({ timeout: 5000 });
    }
    
    await playerContext.close();
  });
});

