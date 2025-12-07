/**
 * RoomSettings - Интеграционные тесты
 * 
 * Тестируют компонент настроек комнаты:
 * - Отображение настроек для хоста и обычного игрока
 * - Редактирование настроек
 * - Валидация кнопки сохранения (активна только при изменениях)
 */

import { describe, it, expect, vi, beforeEach } from 'vitest';
import { render, screen, fireEvent } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { RoomSettingsComponent } from './RoomSettings';
import type { GameRoom } from '@/shared/types';

// Мокаем мутацию обновления настроек
const mockUpdateSettingsMutate = vi.fn();
let mockIsPending = false;

vi.mock('@/entities/room', () => ({
  useUpdateRoomSettings: () => ({
    mutate: mockUpdateSettingsMutate,
    isPending: mockIsPending,
  }),
}));

const queryClient = new QueryClient({
  defaultOptions: { queries: { retry: false } },
});

// Базовые данные комнаты для тестов
const createMockRoom = (overrides?: Partial<GameRoom>): GameRoom => ({
  id: 'room-123',
  roomCode: 'ABC123',
  hostId: 'host-user-1',
  packId: 'pack-1',
  name: 'Тестовая комната',
  status: 'waiting',
  maxPlayers: 6,
  currentPlayers: 2,
  isPublic: true,
  hasPassword: false,
  settings: {
    timeForAnswer: 30,
    timeForChoice: 15,
    allowWrongAnswer: true,
    showRightAnswer: true,
  },
  players: [],
  createdAt: new Date().toISOString(),
  ...overrides,
});

const renderRoomSettings = (room: GameRoom, isHost: boolean) => {
  return render(
    <QueryClientProvider client={queryClient}>
      <RoomSettingsComponent room={room} isHost={isHost} />
    </QueryClientProvider>
  );
};

describe('RoomSettings', () => {
  beforeEach(() => {
    vi.clearAllMocks();
    mockIsPending = false;
  });

  /**
   * ТЕСТ: Отображение настроек для обычного игрока (не хоста)
   * 
   * Проверяет что обычный игрок видит настройки только для чтения:
   * - Показываются текущие значения настроек
   * - Нет слайдеров и переключателей для редактирования
   * - Нет кнопки сохранения
   */
  it('показывает настройки только для чтения для обычного игрока', () => {
    const room = createMockRoom();
    renderRoomSettings(room, false);

    // Проверяем заголовок
    expect(screen.getByText('Настройки игры')).toBeInTheDocument();

    // Проверяем что показаны значения настроек
    expect(screen.getByText('Время на ответ:')).toBeInTheDocument();
    expect(screen.getByText('30 сек')).toBeInTheDocument();
    expect(screen.getByText('Время на выбор:')).toBeInTheDocument();
    expect(screen.getByText('15 сек')).toBeInTheDocument();
    expect(screen.getByText('Неправильные ответы:')).toBeInTheDocument();
    expect(screen.getByText('Разрешены')).toBeInTheDocument();
    expect(screen.getByText('Показывать ответ:')).toBeInTheDocument();
    expect(screen.getByText('Да')).toBeInTheDocument();

    // Проверяем что нет кнопки сохранения
    expect(screen.queryByRole('button', { name: /сохранить/i })).not.toBeInTheDocument();

    // Проверяем что нет слайдеров
    expect(screen.queryByRole('slider')).not.toBeInTheDocument();
  });

  /**
   * ТЕСТ: Отображение формы редактирования для хоста
   * 
   * Проверяет что хост видит форму редактирования:
   * - Слайдеры для времени на ответ и выбор
   * - Переключатели для булевых настроек
   * - Кнопка сохранения
   */
  it('показывает форму редактирования для хоста', () => {
    const room = createMockRoom();
    renderRoomSettings(room, true);

    // Проверяем заголовок
    expect(screen.getByText('Настройки игры')).toBeInTheDocument();

    // Проверяем наличие слайдеров
    const sliders = screen.getAllByRole('slider');
    expect(sliders).toHaveLength(2);

    // Проверяем наличие переключателей
    expect(screen.getByText('Разрешить неправильные ответы')).toBeInTheDocument();
    expect(screen.getByText('Показывать правильный ответ')).toBeInTheDocument();

    // Проверяем наличие кнопки сохранения
    expect(screen.getByRole('button', { name: /сохранить/i })).toBeInTheDocument();
  });

  /**
   * ТЕСТ: Кнопка сохранения неактивна без изменений
   * 
   * Проверяет что кнопка "Сохранить настройки" disabled
   * когда настройки не были изменены
   */
  it('кнопка сохранения неактивна когда настройки не изменены', () => {
    const room = createMockRoom();
    renderRoomSettings(room, true);

    const saveButton = screen.getByRole('button', { name: /сохранить/i });
    expect(saveButton).toBeDisabled();
  });

  /**
   * ТЕСТ: Кнопка сохранения активируется после изменения времени на ответ
   * 
   * Проверяет что после изменения слайдера "Время на ответ"
   * кнопка сохранения становится активной
   */
  it('кнопка сохранения активна после изменения времени на ответ', async () => {
    const room = createMockRoom();
    renderRoomSettings(room, true);

    const saveButton = screen.getByRole('button', { name: /сохранить/i });
    expect(saveButton).toBeDisabled();

    // Находим первый слайдер (время на ответ) и меняем значение
    const sliders = screen.getAllByRole('slider');
    fireEvent.change(sliders[0], { target: { value: '45' } });

    expect(saveButton).not.toBeDisabled();
  });

  /**
   * ТЕСТ: Кнопка сохранения активируется после изменения времени на выбор
   * 
   * Проверяет что после изменения слайдера "Время на выбор"
   * кнопка сохранения становится активной
   */
  it('кнопка сохранения активна после изменения времени на выбор', async () => {
    const room = createMockRoom();
    renderRoomSettings(room, true);

    const saveButton = screen.getByRole('button', { name: /сохранить/i });
    const sliders = screen.getAllByRole('slider');

    // Меняем второй слайдер (время на выбор)
    fireEvent.change(sliders[1], { target: { value: '20' } });

    expect(saveButton).not.toBeDisabled();
  });

  /**
   * ТЕСТ: Кнопка сохранения активируется после переключения "allowWrongAnswer"
   * 
   * Проверяет что после переключения настройки "Разрешить неправильные ответы"
   * кнопка сохранения становится активной
   */
  it('кнопка сохранения активна после переключения allowWrongAnswer', async () => {
    const user = userEvent.setup();
    const room = createMockRoom();
    renderRoomSettings(room, true);

    const saveButton = screen.getByRole('button', { name: /сохранить/i });
    
    // Находим переключатель для allowWrongAnswer и кликаем
    const toggles = screen.getAllByRole('checkbox');
    await user.click(toggles[0]);

    expect(saveButton).not.toBeDisabled();
  });

  /**
   * ТЕСТ: Вызов мутации при сохранении настроек
   * 
   * Проверяет что при клике на кнопку "Сохранить"
   * вызывается мутация с правильными данными
   */
  it('вызывает мутацию обновления при сохранении', async () => {
    const user = userEvent.setup();
    const room = createMockRoom();
    renderRoomSettings(room, true);

    // Меняем настройки
    const sliders = screen.getAllByRole('slider');
    fireEvent.change(sliders[0], { target: { value: '45' } });

    // Сохраняем
    const saveButton = screen.getByRole('button', { name: /сохранить/i });
    await user.click(saveButton);

    expect(mockUpdateSettingsMutate).toHaveBeenCalledWith({
      id: 'room-123',
      settings: expect.objectContaining({
        timeForAnswer: 45,
        timeForChoice: 15,
        allowWrongAnswer: true,
        showRightAnswer: true,
      }),
    });
  });

  /**
   * ТЕСТ: Кнопка снова неактивна после возврата к исходным значениям
   * 
   * Проверяет что если пользователь изменил настройки,
   * а потом вернул их к исходным - кнопка снова становится неактивной
   */
  it('кнопка неактивна при возврате к исходным значениям', async () => {
    const room = createMockRoom();
    renderRoomSettings(room, true);

    const saveButton = screen.getByRole('button', { name: /сохранить/i });
    const sliders = screen.getAllByRole('slider');

    // Меняем значение
    fireEvent.change(sliders[0], { target: { value: '45' } });
    expect(saveButton).not.toBeDisabled();

    // Возвращаем к исходному
    fireEvent.change(sliders[0], { target: { value: '30' } });
    expect(saveButton).toBeDisabled();
  });

  /**
   * ТЕСТ: Отображение "Запрещены" для allowWrongAnswer = false
   * 
   * Проверяет что для обычного игрока корректно отображается
   * значение "Запрещены" когда неправильные ответы отключены
   */
  it('показывает "Запрещены" когда неправильные ответы отключены', () => {
    const room = createMockRoom({
      settings: {
        timeForAnswer: 30,
        timeForChoice: 15,
        allowWrongAnswer: false,
        showRightAnswer: true,
      },
    });
    renderRoomSettings(room, false);

    expect(screen.getByText('Запрещены')).toBeInTheDocument();
  });

  /**
   * ТЕСТ: Отображение "Нет" для showRightAnswer = false
   * 
   * Проверяет что для обычного игрока корректно отображается
   * значение "Нет" когда показ ответа отключен
   */
  it('показывает "Нет" когда показ ответа отключен', () => {
    const room = createMockRoom({
      settings: {
        timeForAnswer: 30,
        timeForChoice: 15,
        allowWrongAnswer: true,
        showRightAnswer: false,
      },
    });
    renderRoomSettings(room, false);

    expect(screen.getByText('Нет')).toBeInTheDocument();
  });

  /**
   * ТЕСТ: Корректное отображение при отсутствии настроек
   * 
   * Проверяет что компонент корректно работает
   * когда room.settings = undefined (использует дефолтные значения)
   */
  it('использует дефолтные настройки при отсутствии settings', () => {
    const room = createMockRoom({ settings: undefined });
    renderRoomSettings(room, false);

    // Проверяем дефолтные значения
    expect(screen.getByText('30 сек')).toBeInTheDocument();
    expect(screen.getByText('10 сек')).toBeInTheDocument();
    expect(screen.getByText('Разрешены')).toBeInTheDocument();
    expect(screen.getByText('Да')).toBeInTheDocument();
  });
});

