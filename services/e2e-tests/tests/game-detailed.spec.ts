import { test, expect } from '@playwright/test';
import { registerUser, generateUsername } from './helpers/auth';
import { createRoom, joinRoom, setReady } from './helpers/room';
import { waitForGameStart, selectQuestion, pressButton, waitForStatus, waitForAnsweringPhase } from './helpers/game';

test.describe('Игра - детальные тесты', () => {
  test('отображение игрового поля с темами и вопросами', async ({ page, browser }) => {
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
    
    await waitForGameStart(page);
    await waitForStatus(page, 'question_select');
    
    await expect(page.locator('.game-board')).toBeVisible();
    await expect(page.locator('.game-board__theme')).toHaveCount(5);
    await expect(page.locator('.game-board__question')).toHaveCount(25);
    
    await playerPage.close();
    await playerContext.close();
  });

  test('вопрос становится недоступным после выбора', async ({ page, browser }) => {
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
    
    await waitForGameStart(page);
    await waitForStatus(page, 'question_select');
    
    const firstQuestion = page.locator('.game-board__question:not(.game-board__question--disabled)').first();
    await firstQuestion.click();
    
    await waitForStatus(page, 'button_press', 10000);
    
    await page.waitForTimeout(1000);
    
    await expect(firstQuestion).toHaveClass(/game-board__question--disabled/, { timeout: 5000 });
    
    await playerPage.close();
    await playerContext.close();
  });

  test('кнопка нажатия неактивна для ведущего', async ({ page, browser }) => {
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
    
    await waitForGameStart(page);
    await waitForStatus(page, 'question_select');
    await selectQuestion(page);
    await waitForStatus(page, 'button_press', 10000);
    
    const pressButtonHost = page.getByRole('button', { name: /нажать/i }).or(page.locator('button').filter({ hasText: /ответить/i }));
    const isVisible = await pressButtonHost.isVisible().catch(() => false);
    if (isVisible) {
      await expect(pressButtonHost).toBeDisabled({ timeout: 5000 });
    } else {
      await expect(pressButtonHost).not.toBeVisible({ timeout: 5000 });
    }
    
    await playerPage.close();
    await playerContext.close();
  });

  test('отображение микрофона для активного игрока', async ({ page, browser }) => {
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
    
    await waitForGameStart(page);
    await waitForStatus(page, 'question_select');
    await selectQuestion(page);
    await waitForStatus(page, 'button_press', 10000);
    await pressButton(playerPage);
    
    await waitForAnsweringPhase(playerPage);
    
    await expect(playerPage.locator('.game-page__answering')).toBeVisible();
    await expect(playerPage.getByText(/ваш черёд отвечать/i).or(playerPage.getByText(/говорите/i))).toBeVisible();
    
    await playerPage.close();
    await playerContext.close();
  });

  test('отображение сообщения ожидания для других игроков', async ({ page, browser }) => {
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
    
    await waitForGameStart(page);
    await waitForStatus(page, 'question_select');
    await selectQuestion(page);
    await waitForStatus(page, 'button_press', 10000);
    await pressButton(playerPage);
    
    await waitForAnsweringPhase(playerPage);
    
    await expect(page.getByText(new RegExp(playerUsername, 'i')).or(page.getByText(/отвечает/i)).or(page.getByText(/ждём/i))).toBeVisible({ timeout: 10000 });
    
    await playerPage.close();
    await playerContext.close();
  });

  test('отображение списка игроков с очками', async ({ page, browser }) => {
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
    
    await waitForGameStart(page);
    
    await expect(page.locator('.player-list, [class*="player"]')).toBeVisible();
    await expect(page.getByText(hostUsername)).toBeVisible();
    await expect(page.getByText(playerUsername)).toBeVisible();
    
    await playerPage.close();
    await playerContext.close();
  });

  test('выход из игры', async ({ page, browser }) => {
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
    
    await waitForGameStart(page);
    
    const exitButton = page.getByRole('button', { name: /выйти из игры/i });
    await exitButton.waitFor({ state: 'visible', timeout: 10000 });
    await exitButton.click();
    
    await expect(page).toHaveURL(/\/lobby|\/login/, { timeout: 10000 });
    
    await playerPage.close();
    await playerContext.close();
  });
});

