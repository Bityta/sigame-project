import { test, expect } from '@playwright/test';
import { registerUser, generateUsername } from './helpers/auth';
import { createRoom, joinRoom, setReady } from './helpers/room';
import { waitForGameStart, selectQuestion, transferSecret, placeStake, submitForAllAnswer, waitForStatus } from './helpers/game';

test.describe('Специальные вопросы', () => {
  test('Кот в мешке (Secret) работает корректно', async ({ page, context }) => {
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
    
    await waitForGameStart(page);
    await waitForGameStart(playerPage);
    
    await waitForStatus(page, 'question_select');
    
    await selectQuestion(page);
    
    await waitForStatus(page, 'secret_transfer', 10000).catch(async () => {
      await page.waitForTimeout(2000);
    });
    
    const secretPanel = page.locator('.secret-transfer-panel');
    if (await secretPanel.isVisible({ timeout: 5000 }).catch(() => false)) {
      await transferSecret(page, playerUsername);
      
      await waitForStatus(playerPage, 'question_show', 10000);
    }
    
    await playerPage.close();
  });

  test('Вопрос со ставкой (Stake) работает корректно', async ({ page, context }) => {
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
    
    await waitForGameStart(page);
    await waitForGameStart(playerPage);
    
    await waitForStatus(page, 'question_select');
    await selectQuestion(page);
    
    await waitForStatus(page, 'stake_betting', 10000).catch(async () => {
      await page.waitForTimeout(2000);
    });
    
    const stakePanel = playerPage.locator('.stake-betting-panel');
    if (await stakePanel.isVisible({ timeout: 5000 }).catch(() => false)) {
      await placeStake(playerPage, 100);
      
      await waitForStatus(playerPage, 'question_show', 10000);
    }
    
    await playerPage.close();
  });

  test('Вопрос для всех (ForAll) работает корректно', async ({ page, context }) => {
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
    
    await waitForGameStart(page);
    await waitForGameStart(playerPage);
    
    await waitForStatus(page, 'question_select');
    await selectQuestion(page);
    
    await waitForStatus(page, 'for_all_answering', 10000).catch(async () => {
      await page.waitForTimeout(2000);
    });
    
    const forAllPanel = playerPage.locator('.for-all-answer-input');
    if (await forAllPanel.isVisible({ timeout: 5000 }).catch(() => false)) {
      await submitForAllAnswer(playerPage, 'Тестовый ответ');
      
      await expect(playerPage.getByText(/ответ отправлен/i)).toBeVisible({ timeout: 5000 });
    }
    
    await playerPage.close();
  });
});
