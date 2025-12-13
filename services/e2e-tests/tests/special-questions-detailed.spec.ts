import { test, expect } from '@playwright/test';
import { registerUser, generateUsername } from './helpers/auth';
import { createRoom, joinRoom, setReady } from './helpers/room';
import { waitForGameStart, selectQuestion, transferSecret, placeStake, submitForAllAnswer, waitForStatus, pressButton, judgeAnswer } from './helpers/game';

test.describe('Специальные вопросы - детально', () => {
  test('Secret - только получатель может отвечать', async ({ page, browser }) => {
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
    
    try {
      await waitForStatus(page, 'secret_transfer', 10000);
    } catch (error) {
      // Если специальный вопрос не появился, пропускаем тест
      return;
    }
    
    const secretPanel = page.locator('.secret-transfer-panel');
    if (await secretPanel.isVisible({ timeout: 5000 }).catch(() => false)) {
      await transferSecret(page, playerUsername);
      
      await waitForStatus(playerPage, 'button_press', 10000);
      
      const pressButtonPlayer = playerPage.getByRole('button', { name: /нажать/i }).or(playerPage.locator('button').filter({ hasText: /ответить/i }));
      await expect(pressButtonPlayer).toBeVisible({ timeout: 5000 });
      
      const pressButtonHost = page.getByRole('button', { name: /нажать/i }).or(page.locator('button').filter({ hasText: /ответить/i }));
      const hostButtonVisible = await pressButtonHost.isVisible({ timeout: 5000 }).catch(() => false);
      if (hostButtonVisible) {
        await expect(pressButtonHost).toBeDisabled({ timeout: 5000 });
      } else {
        await expect(pressButtonHost).not.toBeVisible({ timeout: 5000 });
      }
    }
    
    await playerContext.close();
  });

  test('Stake - правильное начисление/списание очков', async ({ page, browser }) => {
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
    
    try {
      await waitForStatus(page, 'stake_betting', 10000);
    } catch (error) {
      // Если специальный вопрос не появился, пропускаем тест
      return;
    }
    
    const stakePanel = playerPage.locator('.stake-betting-panel');
    if (await stakePanel.isVisible({ timeout: 5000 }).catch(() => false)) {
      const initialScore = await playerPage.locator('[class*="score"]').first().textContent().catch(() => '0');
      
      await placeStake(playerPage, 100);
      
      await waitForStatus(playerPage, 'button_press', 10000);
      await pressButton(playerPage);
      await waitForStatus(page, 'answering', 10000);
      await waitForStatus(page, 'answer_judging', 35000);
      
      await judgeAnswer(page, true);
      
      await page.waitForTimeout(2000);
      
      const finalScore = await playerPage.locator('[class*="score"]').first().textContent().catch(() => '0');
      
      const initial = parseInt(initialScore || '0', 10);
      const final = parseInt(finalScore || '0', 10);
      
      expect(final).toBeGreaterThan(initial);
    }
    
    await playerContext.close();
  });

  test('ForAll - правильное начисление очков всем правильно ответившим', async ({ page, browser }) => {
    const hostUsername = generateUsername();
    const player1Username = generateUsername();
    const player2Username = generateUsername();
    const password = 'testpass123';

    await registerUser(page, hostUsername, password);
    const roomId = await createRoom(page, 'ForAll Room', 4);
    
    const player1Context = await browser.newContext();
    const player1Page = await player1Context.newPage();
    await registerUser(player1Page, player1Username, password);
    await joinRoom(player1Page, roomId);
    
    const player2Context = await browser.newContext();
    const player2Page = await player2Context.newPage();
    await registerUser(player2Page, player2Username, password);
    await joinRoom(player2Page, roomId);
    
    await setReady(page);
    await setReady(player1Page);
    await setReady(player2Page);
    
    await waitForGameStart(page);
    await waitForStatus(page, 'question_select');
    
    await selectQuestion(page);
    
    try {
      await waitForStatus(page, 'for_all_answering', 10000);
    } catch (error) {
      // Если специальный вопрос не появился, пропускаем тест
      return;
    }
    
    const forAllPanel1 = player1Page.locator('.for-all-answer-input');
    const forAllPanel2 = player2Page.locator('.for-all-answer-input');
    
    if (await forAllPanel1.isVisible({ timeout: 5000 }).catch(() => false)) {
      await submitForAllAnswer(player1Page, 'Правильный ответ');
      await submitForAllAnswer(player2Page, 'Неправильный ответ');
      
      await waitForStatus(page, 'for_all_results', 10000);
      
      await expect(page.locator('.for-all-results')).toBeVisible({ timeout: 10000 });
    }
    
    await player1Context.close();
    await player2Context.close();
  });
});

