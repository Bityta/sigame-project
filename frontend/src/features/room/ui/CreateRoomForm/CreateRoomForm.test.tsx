/**
 * CreateRoomForm - Интеграционные тесты
 * 
 * Тестируют форму создания комнаты:
 * - Валидацию полей (название, пак, пароль)
 * - Выбор пака вопросов
 * - Переключение публичности комнаты
 * - Успешное создание комнаты
 */

import { describe, it, expect, vi, beforeEach } from 'vitest';
import { render, screen, fireEvent, waitFor } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { BrowserRouter } from 'react-router-dom';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { CreateRoomForm } from './CreateRoomForm';

// Мокаем навигацию
const mockNavigate = vi.fn();
vi.mock('react-router-dom', async () => {
  const actual = await vi.importActual('react-router-dom');
  return {
    ...actual,
    useNavigate: () => mockNavigate,
  };
});

// Мокаем мутацию создания комнаты
const mockCreateRoomMutate = vi.fn();
let mockIsPending = false;

vi.mock('@/entities/room', () => ({
  useCreateRoom: (options: any) => ({
    mutate: (data: any) => {
      mockCreateRoomMutate(data);
      // Симулируем успешное создание
      options?.onSuccess?.({ id: 'new-room-123', roomCode: 'XYZ789' });
    },
    isPending: mockIsPending,
  }),
}));

// Мокаем получение паков
const mockPacks = [
  { id: 'pack-1', name: 'Общие знания', author: 'Admin' },
  { id: 'pack-2', name: 'Кино и музыка', author: 'User123' },
  { id: 'pack-3', name: 'Наука', author: 'ScienceGuy' },
];

vi.mock('@/entities/pack', () => ({
  usePacks: () => ({
    data: mockPacks,
    isLoading: false,
  }),
}));

const queryClient = new QueryClient({
  defaultOptions: { queries: { retry: false } },
});

const renderCreateRoomForm = () => {
  return render(
    <QueryClientProvider client={queryClient}>
      <BrowserRouter>
        <CreateRoomForm />
      </BrowserRouter>
    </QueryClientProvider>
  );
};

describe('CreateRoomForm', () => {
  beforeEach(() => {
    vi.clearAllMocks();
    mockIsPending = false;
  });

  /**
   * ТЕСТ: Корректный рендеринг формы
   * 
   * Проверяет что форма отображает все необходимые элементы:
   * - Заголовок "Создать комнату"
   * - Поле названия комнаты
   * - Выпадающий список паков
   * - Слайдер количества игроков
   * - Чекбокс публичности
   * - Кнопки "Создать" и "Отмена"
   */
  it('отображает все элементы формы', () => {
    renderCreateRoomForm();

    expect(screen.getByRole('heading', { name: /создать комнату/i })).toBeInTheDocument();
    expect(screen.getByPlaceholderText(/введите название/i)).toBeInTheDocument();
    expect(screen.getByRole('combobox')).toBeInTheDocument(); // select для паков
    expect(screen.getByText(/максимум игроков/i)).toBeInTheDocument();
    expect(screen.getByLabelText(/публичная комната/i)).toBeInTheDocument();
    expect(screen.getByRole('button', { name: /создать/i })).toBeInTheDocument();
    expect(screen.getByRole('button', { name: /отмена/i })).toBeInTheDocument();
  });

  /**
   * ТЕСТ: Валидация короткого названия комнаты
   * 
   * Проверяет что при вводе слишком короткого названия (менее 3 символов)
   * показывается ошибка валидации и запрос не отправляется
   */
  it('показывает ошибку при коротком названии (менее 3 символов)', async () => {
    const user = userEvent.setup();
    renderCreateRoomForm();

    const nameInput = screen.getByPlaceholderText(/введите название/i);
    const packSelect = screen.getByRole('combobox');
    const submitButton = screen.getByRole('button', { name: /создать/i });

    await user.type(nameInput, 'ab'); // Короткое название
    await user.selectOptions(packSelect, 'pack-1');
    await user.click(submitButton);

    expect(screen.getByText(/название должно быть от 3 до 50 символов/i)).toBeInTheDocument();
    expect(mockCreateRoomMutate).not.toHaveBeenCalled();
  });

  /**
   * ТЕСТ: Валидация слишком длинного названия комнаты
   * 
   * Проверяет что при вводе слишком длинного названия (более 50 символов)
   * показывается ошибка валидации
   */
  it('показывает ошибку при длинном названии (более 50 символов)', async () => {
    const user = userEvent.setup();
    renderCreateRoomForm();

    const nameInput = screen.getByPlaceholderText(/введите название/i);
    const packSelect = screen.getByRole('combobox');
    const submitButton = screen.getByRole('button', { name: /создать/i });

    const longName = 'a'.repeat(51);
    await user.type(nameInput, longName);
    await user.selectOptions(packSelect, 'pack-1');
    await user.click(submitButton);

    expect(screen.getByText(/название должно быть от 3 до 50 символов/i)).toBeInTheDocument();
    expect(mockCreateRoomMutate).not.toHaveBeenCalled();
  });

  /**
   * ТЕСТ: Валидация невыбранного пака вопросов
   * 
   * Проверяет что при попытке создать комнату без выбора пака
   * показывается ошибка валидации
   */
  it('показывает ошибку когда пак не выбран', async () => {
    const user = userEvent.setup();
    renderCreateRoomForm();

    const nameInput = screen.getByPlaceholderText(/введите название/i);
    const submitButton = screen.getByRole('button', { name: /создать/i });

    await user.type(nameInput, 'Моя комната');
    await user.click(submitButton);

    expect(screen.getByText(/выберите пак вопросов/i)).toBeInTheDocument();
    expect(mockCreateRoomMutate).not.toHaveBeenCalled();
  });

  /**
   * ТЕСТ: Отображение списка паков
   * 
   * Проверяет что все доступные паки отображаются
   * в выпадающем списке с названием и автором
   */
  it('отображает список доступных паков', () => {
    renderCreateRoomForm();

    const packSelect = screen.getByRole('combobox');

    expect(packSelect).toContainHTML('Общие знания (автор: Admin)');
    expect(packSelect).toContainHTML('Кино и музыка (автор: User123)');
    expect(packSelect).toContainHTML('Наука (автор: ScienceGuy)');
  });

  /**
   * ТЕСТ: Показ поля пароля при отключении публичности
   * 
   * Проверяет что при снятии галочки "Публичная комната"
   * появляется поле для ввода пароля
   */
  it('показывает поле пароля при отключении публичности', async () => {
    const user = userEvent.setup();
    renderCreateRoomForm();

    // Изначально поля пароля нет
    expect(screen.queryByPlaceholderText(/введите пароль/i)).not.toBeInTheDocument();

    // Снимаем галочку публичности
    const publicCheckbox = screen.getByLabelText(/публичная комната/i);
    await user.click(publicCheckbox);

    // Теперь поле пароля должно появиться
    expect(screen.getByPlaceholderText(/введите пароль/i)).toBeInTheDocument();
  });

  /**
   * ТЕСТ: Валидация пароля для приватной комнаты
   * 
   * Проверяет что для приватной комнаты обязательно
   * указать пароль, иначе будет ошибка валидации
   */
  it('показывает ошибку для приватной комнаты без пароля', async () => {
    const user = userEvent.setup();
    renderCreateRoomForm();

    const nameInput = screen.getByPlaceholderText(/введите название/i);
    const packSelect = screen.getByRole('combobox');
    const publicCheckbox = screen.getByLabelText(/публичная комната/i);
    const submitButton = screen.getByRole('button', { name: /создать/i });

    await user.type(nameInput, 'Приватная комната');
    await user.selectOptions(packSelect, 'pack-1');
    await user.click(publicCheckbox); // Делаем приватной
    await user.click(submitButton);

    expect(screen.getByText(/для приватной комнаты укажите пароль/i)).toBeInTheDocument();
    expect(mockCreateRoomMutate).not.toHaveBeenCalled();
  });

  /**
   * ТЕСТ: Успешное создание публичной комнаты
   * 
   * Проверяет полный сценарий создания комнаты:
   * - Заполнение всех полей
   * - Отправка формы
   * - Редирект на страницу комнаты
   */
  it('создает публичную комнату и перенаправляет на её страницу', async () => {
    const user = userEvent.setup();
    renderCreateRoomForm();

    const nameInput = screen.getByPlaceholderText(/введите название/i);
    const packSelect = screen.getByRole('combobox');
    const submitButton = screen.getByRole('button', { name: /создать/i });

    await user.type(nameInput, 'Моя игровая комната');
    await user.selectOptions(packSelect, 'pack-2');
    await user.click(submitButton);

    expect(mockCreateRoomMutate).toHaveBeenCalledWith(
      expect.objectContaining({
        name: 'Моя игровая комната',
        packId: 'pack-2',
        isPublic: true,
        password: undefined,
      })
    );
    expect(mockNavigate).toHaveBeenCalledWith('/room/new-room-123');
  });

  /**
   * ТЕСТ: Успешное создание приватной комнаты
   * 
   * Проверяет создание приватной комнаты с паролем
   */
  it('создает приватную комнату с паролем', async () => {
    const user = userEvent.setup();
    renderCreateRoomForm();

    const nameInput = screen.getByPlaceholderText(/введите название/i);
    const packSelect = screen.getByRole('combobox');
    const publicCheckbox = screen.getByLabelText(/публичная комната/i);
    const submitButton = screen.getByRole('button', { name: /создать/i });

    await user.type(nameInput, 'Секретная комната');
    await user.selectOptions(packSelect, 'pack-3');
    await user.click(publicCheckbox); // Делаем приватной

    // Появляется поле пароля
    const passwordInput = screen.getByPlaceholderText(/введите пароль/i);
    await user.type(passwordInput, 'secret123');
    
    await user.click(submitButton);

    expect(mockCreateRoomMutate).toHaveBeenCalledWith(
      expect.objectContaining({
        name: 'Секретная комната',
        packId: 'pack-3',
        isPublic: false,
        password: 'secret123',
      })
    );
  });

  /**
   * ТЕСТ: Изменение количества игроков через слайдер
   * 
   * Проверяет что слайдер максимального количества игроков
   * работает корректно и значение передается при создании
   */
  it('позволяет изменить максимальное количество игроков', async () => {
    const user = userEvent.setup();
    renderCreateRoomForm();

    const nameInput = screen.getByPlaceholderText(/введите название/i);
    const packSelect = screen.getByRole('combobox');
    const slider = screen.getByRole('slider');
    const submitButton = screen.getByRole('button', { name: /создать/i });

    await user.type(nameInput, 'Большая комната');
    await user.selectOptions(packSelect, 'pack-1');
    
    // Меняем слайдер на 6 игроков
    fireEvent.change(slider, { target: { value: '6' } });
    
    await user.click(submitButton);

    expect(mockCreateRoomMutate).toHaveBeenCalledWith(
      expect.objectContaining({
        maxPlayers: 6,
      })
    );
  });

  /**
   * ТЕСТ: Кнопка "Отмена" возвращает в лобби
   * 
   * Проверяет что клик на кнопку "Отмена"
   * перенаправляет пользователя обратно в лобби
   */
  it('кнопка Отмена возвращает в лобби', async () => {
    const user = userEvent.setup();
    renderCreateRoomForm();

    const cancelButton = screen.getByRole('button', { name: /отмена/i });
    await user.click(cancelButton);

    expect(mockNavigate).toHaveBeenCalledWith('/lobby');
  });

  /**
   * ТЕСТ: Очистка ошибки при повторной отправке
   * 
   * Проверяет что ошибка валидации очищается
   * при следующей попытке отправить форму
   */
  it('очищает ошибку при повторной попытке', async () => {
    const user = userEvent.setup();
    renderCreateRoomForm();

    const nameInput = screen.getByPlaceholderText(/введите название/i);
    const packSelect = screen.getByRole('combobox');
    const submitButton = screen.getByRole('button', { name: /создать/i });

    // Первая попытка с ошибкой
    await user.type(nameInput, 'ab');
    await user.selectOptions(packSelect, 'pack-1');
    await user.click(submitButton);
    expect(screen.getByText(/название должно быть от 3 до 50 символов/i)).toBeInTheDocument();

    // Исправляем и отправляем снова
    await user.clear(nameInput);
    await user.type(nameInput, 'Нормальное название');
    await user.click(submitButton);

    // Ошибка должна исчезнуть
    expect(screen.queryByText(/название должно быть от 3 до 50 символов/i)).not.toBeInTheDocument();
    expect(mockCreateRoomMutate).toHaveBeenCalled();
  });
});

