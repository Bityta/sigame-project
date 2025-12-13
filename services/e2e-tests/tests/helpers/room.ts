import { Page, expect } from '@playwright/test';

export async function createRoom(
  page: Page,
  name: string = `Test Room ${Date.now()}`,
  maxPlayers: number = 4
): Promise<string> {
  await page.goto('/lobby');
  await page.waitForTimeout(500);
  
  await page.waitForTimeout(300);
  await page.getByRole('button', { name: /создать комнату/i }).click({ delay: 100 });
  await page.waitForURL(/\/lobby\/create/);
  await page.waitForTimeout(500);
  
  await page.getByPlaceholder(/введите название/i).fill(name, { delay: 50 });
  await page.waitForTimeout(300);
  
  const packSelect = page.locator('select').first();
  await packSelect.waitFor({ state: 'visible' });
  await page.waitForTimeout(1000);
  
  const packOptions = await packSelect.locator('option').all();
  if (packOptions.length > 1) {
    const firstRealOption = packOptions[1];
    const optionValue = await firstRealOption.getAttribute('value');
    if (optionValue) {
      await packSelect.selectOption(optionValue);
    } else {
      await packSelect.selectOption({ index: 1 });
    }
    await page.waitForTimeout(1000);
  }
  
  const createButton = page.getByRole('button', { name: /создать/i });
  await createButton.waitFor({ state: 'visible' });
  await page.waitForTimeout(500);
  await createButton.click({ delay: 100 });
  
  await page.waitForURL(/\/room\/.+/, { timeout: 20000 });
  await page.waitForTimeout(2000);
  const roomUrl = page.url();
  const roomId = roomUrl.split('/room/')[1]?.split('?')[0] || '';
  
  if (!roomId) {
    throw new Error(`Failed to create room. Current URL: ${roomUrl}`);
  }
  
  return roomId;
}

export async function joinRoomByCode(page: Page, code: string): Promise<void> {
  await page.goto('/lobby');
  await page.waitForLoadState('networkidle');
  await page.waitForTimeout(1000);
  
  const codeInput = page.locator('.lobby-page__code-input');
  await codeInput.waitFor({ state: 'visible', timeout: 10000 });
  await codeInput.fill(code.toUpperCase());
  
  await page.waitForTimeout(500);
  
  const joinCodeContainer = page.locator('.lobby-page__join-code');
  await joinCodeContainer.waitFor({ state: 'visible', timeout: 10000 });
  
  const joinButton = joinCodeContainer.locator('button').first();
  await joinButton.waitFor({ state: 'visible', timeout: 10000 });
  await joinButton.click();
  
  await page.waitForURL(/\/room\/.+/, { timeout: 10000 });
}

export async function joinRoom(page: Page, roomId: string): Promise<void> {
  await page.goto(`/room/${roomId}`);
  await page.waitForLoadState('networkidle');
  await page.waitForSelector('.room-page', { timeout: 15000 });
  await expect(page.locator('.room-page__title')).toBeVisible({ timeout: 15000 });
  await page.waitForTimeout(1000);
}

export async function setReady(page: Page): Promise<void> {
  const readyButton = page.getByRole('button', { name: /готов/i });
  await readyButton.waitFor({ state: 'visible' });
  await page.waitForTimeout(300);
  await readyButton.hover();
  await page.waitForTimeout(200);
  await readyButton.click({ delay: 100 });
  
  await page.waitForTimeout(500);
}

export async function waitForGameStart(page: Page): Promise<void> {
  await page.waitForURL(/\/game\/.+/, { timeout: 30000 });
  await expect(page).toHaveURL(/\/game\/.+/);
}

export async function getRoomCode(page: Page): Promise<string> {
  const codeElement = page.locator('.room-page__code-value');
  await codeElement.waitFor({ state: 'visible' });
  return await codeElement.textContent() || '';
}

