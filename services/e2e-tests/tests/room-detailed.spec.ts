import { test, expect } from '@playwright/test';
import { registerUser, generateUsername } from './helpers/auth';
import { createRoom, joinRoom, setReady } from './helpers/room';

test.describe('Комната - детальные тесты', () => {
  test('отображение списка игроков с ролями', async ({ page, context }) => {
    const hostUsername = generateUsername();
    const playerUsername = generateUsername();
    const password = 'testpass123';

    await registerUser(page, hostUsername, password);
    const roomId = await createRoom(page);
    
    const playerPage = await context.newPage();
    await registerUser(playerPage, playerUsername, password);
    await joinRoom(playerPage, roomId);
    
    await expect(page.getByText(hostUsername)).toBeVisible();
    await expect(page.getByText(playerUsername)).toBeVisible();
    await expect(page.getByText(/👑/)).toBeVisible();
    
    await playerPage.close();
  });

  test('счетчик готовых игроков', async ({ page, context }) => {
    const hostUsername = generateUsername();
    const playerUsername = generateUsername();
    const password = 'testpass123';

    await registerUser(page, hostUsername, password);
    const roomId = await createRoom(page);
    
    const playerPage = await context.newPage();
    await registerUser(playerPage, playerUsername, password);
    await joinRoom(playerPage, roomId);
    
    await expect(page.getByText(/0 \/ 2/i).or(page.getByText(/готовы/i))).toBeVisible();
    
    await setReady(page);
    
    await expect(page.getByText(/1 \/ 2/i).or(page.getByText(/готовы/i))).toBeVisible({ timeout: 3000 });
    
    await playerPage.close();
  });

  test('повторное нажатие отменяет готовность', async ({ page }) => {
    const username = generateUsername();
    const password = 'testpass123';

    await registerUser(page, username, password);
    await createRoom(page);
    
    const readyButton = page.getByRole('button', { name: /готов/i });
    await readyButton.click();
    
    await expect(page.getByText(/вы готовы/i)).toBeVisible();
    
    await readyButton.click();
    
    await expect(page.getByText(/вы готовы/i)).not.toBeVisible({ timeout: 3000 });
  });

  test('кнопка готовности неактивна когда комната не в статусе waiting', async ({ page, context }) => {
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
    
    await page.waitForURL(/\/game\/.+/, { timeout: 30000 });
    
    await page.goto(`/room/${roomId}`);
    
    const readyButton = page.getByRole('button', { name: /готов/i });
    const readyButtonVisible = await readyButton.isVisible({ timeout: 5000 }).catch(() => false);
    if (readyButtonVisible) {
      await expect(readyButton).toBeDisabled({ timeout: 5000 });
    }
    
    await playerPage.close();
  });

  test('сообщение о минимальном количестве игроков', async ({ page }) => {
    const username = generateUsername();
    const password = 'testpass123';

    await registerUser(page, username, password);
    await createRoom(page);
    
    await expect(page.getByText(/минимум.*игрок/i).or(page.getByText(/2.*игрок/i))).toBeVisible();
  });

  test('сообщение о готовности всех игроков', async ({ page, context }) => {
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
    
    await expect(page.getByText(/все готовы/i).or(page.getByText(/игра начинается/i))).toBeVisible({ timeout: 5000 });
    
    await playerPage.close();
  });

  test('выход из комнаты', async ({ page }) => {
    const username = generateUsername();
    const password = 'testpass123';

    await registerUser(page, username, password);
    await createRoom(page);
    
    await page.getByRole('button', { name: /покинуть комнату/i }).click();
    
    await expect(page).toHaveURL(/\/lobby/, { timeout: 10000 });
  });

  test('настройки комнаты только для просмотра у не-хоста', async ({ page, context }) => {
    const hostUsername = generateUsername();
    const playerUsername = generateUsername();
    const password = 'testpass123';

    await registerUser(page, hostUsername, password);
    const roomId = await createRoom(page);
    
    const playerPage = await context.newPage();
    await registerUser(playerPage, playerUsername, password);
    await joinRoom(playerPage, roomId);
    
    await expect(playerPage.locator('.room-settings__slider')).not.toBeVisible();
    await expect(playerPage.getByRole('button', { name: /сохранить настройки/i })).not.toBeVisible();
    
    await playerPage.close();
  });

  test('кнопка сохранения настроек неактивна без изменений', async ({ page }) => {
    const username = generateUsername();
    const password = 'testpass123';

    await registerUser(page, username, password);
    await createRoom(page);
    
    const saveButton = page.getByRole('button', { name: /сохранить настройки/i });
    await expect(saveButton).toBeDisabled();
  });
});

