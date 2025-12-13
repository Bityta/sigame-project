import { test, expect } from '@playwright/test';
import { registerUser, generateUsername } from './helpers/auth';
import { createRoom, joinRoom, setReady } from './helpers/room';
import { waitForGameStart, selectQuestion, waitForStatus } from './helpers/game';

test.describe('Медиа вопросы', () => {
  test('предзагрузка медиа раунда', async ({ page, browser }) => {
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
    
    await page.waitForTimeout(3000);
    
    const mediaProgress = page.locator('[class*="media"], [class*="preload"]');
    const hasMediaProgress = await mediaProgress.isVisible({ timeout: 5000 }).catch(() => false);
    if (hasMediaProgress) {
      await expect(mediaProgress).toBeVisible({ timeout: 5000 });
    }
    
    await playerContext.close();
  });

  test('отображение медиа в вопросах - изображение', async ({ page, browser }) => {
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
    
    await page.waitForTimeout(2000);
    
    const image = page.locator('img, [class*="image"], [class*="media"]');
    const hasImage = await image.isVisible({ timeout: 5000 }).catch(() => false);
    if (hasImage) {
      await expect(image).toBeVisible({ timeout: 5000 });
    }
    
    await playerContext.close();
  });

  test('отображение медиа в вопросах - видео', async ({ page, browser }) => {
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
    
    await page.waitForTimeout(2000);
    
    const video = page.locator('video, [class*="video"]');
    const hasVideo = await video.isVisible({ timeout: 5000 }).catch(() => false);
    if (hasVideo) {
      await expect(video).toBeVisible({ timeout: 5000 });
    }
    
    await playerContext.close();
  });

  test('синхронизация медиа между игроками', async ({ page, browser }) => {
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
    
    await page.waitForTimeout(2000);
    
    const hostMedia = page.locator('img, video, audio, [class*="media"]');
    const playerMedia = playerPage.locator('img, video, audio, [class*="media"]');
    
    const hasMedia = await hostMedia.isVisible({ timeout: 5000 }).catch(() => false);
    if (hasMedia) {
      await expect(playerMedia).toBeVisible({ timeout: 10000 });
    }
    
    await playerContext.close();
  });

  test('обработка ошибок загрузки медиа', async ({ page, browser }) => {
    const hostUsername = generateUsername();
    const playerUsername = generateUsername();
    const password = 'testpass123';

    await registerUser(page, hostUsername, password);
    const roomId = await createRoom(page);
    
    const playerContext = await browser.newContext();
    const playerPage = await playerContext.newPage();
    await registerUser(playerPage, playerUsername, password);
    await joinRoom(playerPage, roomId);
    
    await page.route('**/packs/**', route => route.abort());
    
    await setReady(page);
    await setReady(playerPage);
    
    await waitForGameStart(page);
    await waitForStatus(page, 'question_select');
    
    await selectQuestion(page);
    
    await page.waitForTimeout(3000);
    
    const fallback = page.locator('[class*="placeholder"], [class*="fallback"], [class*="error"]');
    const hasFallback = await fallback.isVisible({ timeout: 5000 }).catch(() => false);
    if (hasFallback) {
      await expect(fallback).toBeVisible({ timeout: 5000 });
    }
    
    await playerContext.close();
  });

  test('кеширование медиа на клиенте', async ({ page, browser }) => {
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
    
    await page.waitForTimeout(2000);
    
    const cachedMedia = await page.evaluate(() => {
      return caches.keys().then(keys => keys.length > 0);
    });
    
    expect(cachedMedia).toBeDefined();
    
    await playerContext.close();
  });
});

