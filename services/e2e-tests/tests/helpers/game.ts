import { Page, expect } from '@playwright/test';

export async function waitForGameStart(page: Page): Promise<void> {
  await page.waitForURL(/\/game\/.+/, { timeout: 30000 });
  await expect(page).toHaveURL(/\/game\/.+/);
}

export async function selectQuestion(
  page: Page,
  themeName?: string,
  price?: number
): Promise<void> {
  // Ждем перехода в question_select - проверяем текст индикатора или статус
  try {
    await page.waitForFunction(
      () => {
        const turnIndicator = document.querySelector('.game-page__turn-indicator-text');
        const indicatorText = turnIndicator?.textContent || '';
        // Проверяем, что игра в состоянии выбора вопроса
        return indicatorText.includes('Выберите вопрос') || 
               indicatorText.includes('выбирает вопрос') ||
               document.querySelector('.game-board') !== null ||
               document.querySelector('.game-board-empty') !== null;
      },
      { timeout: 60000 }
    );
  } catch (error) {
    // Если не удалось дождаться по индикатору, просто ждем игровое поле
    console.log('Не удалось дождаться question_select по индикатору, ждем игровое поле напрямую');
  }
  
  // Ждем появления игрового поля или сообщения о загрузке
  await page.waitForSelector('.game-board, .game-board-empty', { timeout: 30000 });
  
  // Проверяем, что игровое поле загружено (не показывается "Загрузка игрового поля...")
  const loadingMessage = page.locator('.game-board-empty');
  const isVisible = await loadingMessage.isVisible().catch(() => false);
  if (isVisible) {
    // Ждем, пока загрузка завершится
    await page.waitForSelector('.game-board', { timeout: 30000 });
  }
  
  if (themeName && price) {
    const theme = page.locator('.game-board__theme').filter({ hasText: themeName }).first();
    const question = theme.locator('.game-board__question').filter({ hasText: price.toString() }).first();
    await question.waitFor({ state: 'visible' });
    await question.click();
  } else {
    const firstAvailableQuestion = page.locator('.game-board__question:not(.game-board__question--disabled)').first();
    await firstAvailableQuestion.waitFor({ state: 'visible' });
    await firstAvailableQuestion.click();
  }
}

export async function pressButton(page: Page): Promise<void> {
  const button = page.getByRole('button', { name: /нажать/i }).or(page.locator('button').filter({ hasText: /ответить/i }));
  await button.waitFor({ state: 'visible', timeout: 10000 });
  await button.click();
}

export async function waitForJudgeButtons(
  page: Page,
  timeout: number = 35000
): Promise<{ appeared: boolean; delay: number }> {
  const startTime = Date.now();
  
  try {
    await page.waitForSelector('.game-page__judging-buttons', { timeout });
    const endTime = Date.now();
    const delay = endTime - startTime;
    
    return { appeared: true, delay };
  } catch (error) {
    const endTime = Date.now();
    const delay = endTime - startTime;
    
    return { appeared: false, delay };
  }
}

export async function judgeAnswer(page: Page, correct: boolean): Promise<void> {
  await page.waitForSelector('.game-page__judging-buttons', { timeout: 5000 });
  
  if (correct) {
    await page.locator('.game-page__judge-btn--correct').click();
  } else {
    await page.locator('.game-page__judge-btn--wrong').click();
  }
}

export async function waitForStatus(page: Page, status: string, timeout: number = 10000): Promise<void> {
  const statusSelectors: Record<string, string> = {
    'question_select': '.game-board',
    'button_press': '.game-page__action-panel',
    'answering': '.game-page__answering',
    'answer_judging': '.game-page__judging',
    'secret_transfer': '.secret-transfer-panel',
    'stake_betting': '.stake-betting-panel',
    'for_all_answering': '.for-all-answer-input',
    'for_all_results': '.for-all-results',
    'game_end': '.game-end',
  };
  
  const selector = statusSelectors[status] || `.game-page`;
  await page.waitForSelector(selector, { timeout });
}

export async function waitForAnsweringPhase(page: Page): Promise<void> {
  await page.waitForSelector('.game-page__answering', { timeout: 15000 });
}

export async function transferSecret(page: Page, playerUsername: string): Promise<void> {
  await page.waitForSelector('.secret-transfer-panel', { timeout: 5000 });
  const playerButton = page.locator('.secret-transfer-panel__player-btn').filter({ hasText: playerUsername });
  await playerButton.click();
}

export async function placeStake(page: Page, amount: number): Promise<void> {
  await page.waitForSelector('.stake-betting-panel', { timeout: 5000 });
  const slider = page.locator('.stake-betting-panel__slider');
  await slider.fill(amount.toString());
  await page.getByRole('button', { name: /сделать ставку/i }).click();
}

export async function submitForAllAnswer(page: Page, answer: string): Promise<void> {
  await page.waitForSelector('.for-all-answer-input__input', { timeout: 5000 });
  await page.locator('.for-all-answer-input__input').fill(answer);
  await page.getByRole('button', { name: /отправить ответ/i }).click();
}

