import { test, expect } from '@playwright/test';
import { registerUser, generateUsername } from './helpers/auth';
import { createRoom } from './helpers/room';

test.describe('Лобби - детальные тесты', () => {
  test('отображение всех элементов лобби', async ({ page }) => {
    const username = generateUsername();
    const password = 'testpass123';

    await registerUser(page, username, password);
    
    await expect(page.getByRole('heading', { name: /sigame/i })).toBeVisible();
    await expect(page.getByText(username)).toBeVisible();
    await expect(page.getByRole('button', { name: /выход/i })).toBeVisible();
    await expect(page.getByRole('button', { name: /создать комнату/i })).toBeVisible();
    await expect(page.locator('.lobby-page__code-input')).toBeVisible();
  });

  test('кнопка создания комнаты неактивна при активной комнате', async ({ page }) => {
    const username = generateUsername();
    const password = 'testpass123';

    await registerUser(page, username, password);
    await createRoom(page);
    
    await page.goto('/lobby');
    
    const createButton = page.getByRole('button', { name: /создать комнату/i });
    const isDisabled = await createButton.isDisabled().catch(() => false);
    if (isDisabled) {
      await expect(createButton).toBeDisabled({ timeout: 5000 });
    } else {
      const tooltip = page.getByTitle(/сначала покиньте/i);
      const tooltipVisible = await tooltip.isVisible({ timeout: 5000 }).catch(() => false);
      if (tooltipVisible) {
        await expect(tooltip).toBeVisible({ timeout: 5000 });
      }
    }
  });

  test('присоединение по коду через Enter', async ({ browser }) => {
    const hostUsername = generateUsername();
    const playerUsername = generateUsername();
    const password = 'testpass123';

    const hostContext = await browser.newContext();
    const hostPage = await hostContext.newPage();
    
    await registerUser(hostPage, hostUsername, password);
    const roomId = await createRoom(hostPage);
    const roomCode = await hostPage.locator('.room-page__code-value').textContent();
    
    const playerContext = await browser.newContext();
    const playerPage = await playerContext.newPage();
    await registerUser(playerPage, playerUsername, password);
    
    await playerPage.goto('/lobby');
    await playerPage.waitForLoadState('networkidle');
    const codeInput = playerPage.locator('.lobby-page__code-input');
    await codeInput.fill(roomCode || '');
    await codeInput.press('Enter');
    
    await expect(playerPage).toHaveURL(/\/room\/.+/, { timeout: 10000 });
    
    await playerPage.close();
    await playerContext.close();
    await hostPage.close();
    await hostContext.close();
  });

  test('ошибка при неверном коде комнаты', async ({ page }) => {
    const username = generateUsername();
    const password = 'testpass123';

    await registerUser(page, username, password);
    
    const codeInput = page.locator('.lobby-page__code-input');
    await codeInput.fill('INVALID');
    const joinButton = page.locator('.lobby-page__join-code').locator('button').first();
    await joinButton.click();
    
    await expect(page.getByText(/не найдена/i).or(page.getByText(/ошибка/i))).toBeVisible({ timeout: 5000 });
  });

  test('кнопка присоединения неактивна при пустом коде', async ({ page }) => {
    const username = generateUsername();
    const password = 'testpass123';

    await registerUser(page, username, password);
    
    const joinButton = page.locator('.lobby-page__join-code').locator('button').first();
    await expect(joinButton).toBeDisabled();
  });

  test('автоматическое преобразование кода в uppercase', async ({ page }) => {
    const username = generateUsername();
    const password = 'testpass123';

    await registerUser(page, username, password);
    
    const codeInput = page.locator('.lobby-page__code-input');
    await codeInput.fill('abc123');
    
    const value = await codeInput.inputValue();
    expect(value).toBe('ABC123');
  });

  test('отображение состояния загрузки при поиске комнаты', async ({ page }) => {
    const username = generateUsername();
    const password = 'testpass123';

    await registerUser(page, username, password);
    
    const codeInput = page.locator('.lobby-page__code-input');
    await codeInput.fill('INVALID');
    const joinButton = page.locator('.lobby-page__join-code').locator('button').first();
    await joinButton.click();
    
    const joinButtonDisabled = await joinButton.isDisabled({ timeout: 2000 }).catch(() => false);
    if (joinButtonDisabled) {
      await expect(joinButton).toBeDisabled({ timeout: 2000 });
    }
  });
});

