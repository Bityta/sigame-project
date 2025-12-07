/**
 * Avatar Utilities
 * Утилиты для работы с аватарками пользователей
 */

/**
 * URL MinIO сервера для хранения файлов
 * В production должен быть настроен через переменную окружения
 */
const MINIO_URL = import.meta.env.VITE_MINIO_URL || 'http://localhost:9000';

/**
 * Бакет для хранения аватарок в MinIO
 */
const AVATARS_BUCKET = 'avatars';

/**
 * Получить URL аватарки пользователя
 * 
 * Поддерживает несколько форматов входных данных:
 * 1. Полный URL (https://...) - возвращается как есть
 * 2. UUID/ID аватарки - строится URL к MinIO
 * 3. null/undefined - возвращается null
 * 
 * @param avatarUrlOrId - полный URL, ID аватарки или null
 * @returns URL аватарки или null если аватарки нет
 * 
 * @example
 * // Полный URL
 * getAvatarUrl('https://minio.example.com/avatars/123.jpg')
 * // => 'https://minio.example.com/avatars/123.jpg'
 * 
 * @example
 * // UUID
 * getAvatarUrl('550e8400-e29b-41d4-a716-446655440000')
 * // => 'http://localhost:9000/avatars/550e8400-e29b-41d4-a716-446655440000'
 * 
 * @example
 * // Null
 * getAvatarUrl(null)
 * // => null
 */
export function getAvatarUrl(avatarUrlOrId: string | null | undefined): string | null {
  if (!avatarUrlOrId) {
    return null;
  }

  // Если это уже полный URL - возвращаем как есть
  if (avatarUrlOrId.startsWith('http://') || avatarUrlOrId.startsWith('https://')) {
    return avatarUrlOrId;
  }

  // Если это ID - строим URL к MinIO
  return `${MINIO_URL}/${AVATARS_BUCKET}/${avatarUrlOrId}`;
}

/**
 * Получить инициалы для дефолтного аватара
 * 
 * @param username - имя пользователя
 * @returns первая буква имени в верхнем регистре
 */
export function getAvatarInitial(username: string | null | undefined): string {
  if (!username) {
    return '?';
  }
  return username.charAt(0).toUpperCase();
}

/**
 * Проверить, является ли URL валидным URL изображения
 * 
 * @param url - URL для проверки
 * @returns true если URL похож на изображение
 */
export function isValidImageUrl(url: string | null | undefined): boolean {
  if (!url) {
    return false;
  }

  // Проверяем расширение или что это URL
  const imageExtensions = ['.jpg', '.jpeg', '.png', '.gif', '.webp', '.svg'];
  const lowercaseUrl = url.toLowerCase();
  
  return (
    lowercaseUrl.startsWith('http://') || 
    lowercaseUrl.startsWith('https://') ||
    imageExtensions.some(ext => lowercaseUrl.endsWith(ext))
  );
}

