import { test, expect } from '@playwright/test';
import { registerUser, generateUsername } from './helpers/auth';
import { createRoom, joinRoom, setReady } from './helpers/room';
import { waitForGameStart, selectQuestion, pressButton, judgeAnswer, waitForStatus } from './helpers/game';

test.describe('Полный цикл игры - детально', () => {
  test('регистрация -> создание комнаты -> присоединение -> настройка -> игра -> завершение', async ({ page, context }) => {
    const hostUsername = generateUsername();
    const playerUsername = generateUsername();
    const password = 'testpass123';

    await registerUser(page, hostUsername, password);
    
    const roomId = await createRoom(page, 'Full Cycle Room', 4);
    
    const slider = page.locator('.room-settings__slider').first();
    await slider.fill('45');
    await page.getByRole('button', { name: /сохранить настройки/i }).click();
    
    const playerPage = await context.newPage();
    await registerUser(playerPage, playerUsername, password);
    await joinRoom(playerPage, roomId);
    
    await setReady(page);
    await setReady(playerPage);
    
    await waitForGameStart(page);
    await waitForGameStart(playerPage);
    
    await waitForStatus(page, 'question_select');
    await selectQuestion(page);
    
    await waitForStatus(page, 'button_press', 10000);
    await pressButton(playerPage);
    
    await waitForStatus(page, 'answering', 10000);
    await waitForStatus(page, 'answer_judging', 35000);
    
    await judgeAnswer(page, true);
    
    await waitForStatus(page, 'question_select', 10000);
    
    await playerPage.close();
  });

  test('завершение игры и возврат в лобби', async ({ page, context }) => {
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
    
    await page.getByRole('button', { name: /выйти из игры/i }).click();
    
    await expect(page).toHaveURL(/\/lobby/);
    
    await playerPage.close();
  });
});

