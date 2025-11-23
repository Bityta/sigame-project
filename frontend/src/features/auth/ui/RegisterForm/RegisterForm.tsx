/**
 * Auth Feature - Register Form
 * UI компонент формы регистрации
 */

import { useState, useEffect } from 'react';
import type { FormEvent } from 'react';
import { useNavigate } from 'react-router-dom';
import { useRegister } from '../../model/mutations';
import { useAuthStore } from '../../model/authStore';
import { useCheckUsername } from '@/entities/user';
import { Button, Input, Card } from '@/shared/ui';
import { ROUTES, LIMITS } from '@/shared/config';
import './RegisterForm.css';

export const RegisterForm = () => {
  const navigate = useNavigate();
  const setAuthenticated = useAuthStore((state) => state.setAuthenticated);
  
  const [username, setUsername] = useState('');
  const [debouncedUsername, setDebouncedUsername] = useState('');
  const [password, setPassword] = useState('');
  const [confirmPassword, setConfirmPassword] = useState('');
  const [localError, setLocalError] = useState('');

  // Debounce username для снижения нагрузки на API
  useEffect(() => {
    const timer = setTimeout(() => setDebouncedUsername(username), 500);
    return () => clearTimeout(timer);
  }, [username]);

  // Проверяем доступность username через API
  const isValidLength = debouncedUsername.length >= LIMITS.USERNAME.MIN && 
                        debouncedUsername.length <= LIMITS.USERNAME.MAX;
  
  const { data: isUsernameAvailable, isLoading: isCheckingUsername } = useCheckUsername(
    debouncedUsername,
    { 
      enabled: isValidLength,
      retry: false,
    }
  );

  // Определяем подсказку для поля username
  const getUsernameHint = () => {
    // Не показываем подсказку пока пользователь печатает
    if (username !== debouncedUsername || !isValidLength) {
      return { text: '', type: 'info' as const };
    }
    
    if (isUsernameAvailable === false) {
      return { text: 'Имя занято', type: 'error' as const };
    }
    
    if (isUsernameAvailable === true) {
      return { text: 'Доступно', type: 'success' as const };
    }
    
    return { text: '', type: 'info' as const };
  };

  const usernameHint = getUsernameHint();

  const registerMutation = useRegister({
    onSuccess: () => {
      setAuthenticated(true);
      navigate(ROUTES.LOBBY);
    },
    onError: () => {
      // Ошибка уже обработана глобально
    },
  });

  const handleSubmit = (e: FormEvent) => {
    e.preventDefault();
    setLocalError('');

    // Валидация длины username
    if (username.length < LIMITS.USERNAME.MIN || username.length > LIMITS.USERNAME.MAX) {
      setLocalError(`Имя пользователя должно быть от ${LIMITS.USERNAME.MIN} до ${LIMITS.USERNAME.MAX} символов`);
      return;
    }
    
    // Проверка доступности username
    if (isUsernameAvailable === false) {
      setLocalError('Это имя пользователя уже занято');
      return;
    }
    
    // Валидация длины пароля
    if (password.length < LIMITS.PASSWORD.MIN) {
      setLocalError(`Пароль должен быть не менее ${LIMITS.PASSWORD.MIN} символов`);
      return;
    }
    
    // Проверка совпадения паролей
    if (password !== confirmPassword) {
      setLocalError('Пароли не совпадают');
      return;
    }

    registerMutation.mutate({ username, password });
  };

  const isSubmitDisabled = 
    registerMutation.isPending || 
    isCheckingUsername || 
    isUsernameAvailable === false ||
    username.length < LIMITS.USERNAME.MIN ||
    password.length < LIMITS.PASSWORD.MIN ||
    password !== confirmPassword;

  return (
    <Card className="register-form">
      <h2 className="register-form__title">Регистрация</h2>
      
      <form onSubmit={handleSubmit} className="register-form__form">
        <Input
          label="Имя пользователя"
          type="text"
          value={username}
          onChange={(e) => setUsername(e.target.value)}
          placeholder="От 3 до 20 символов"
          fullWidth
          disabled={registerMutation.isPending}
          autoComplete="username"
          hint={usernameHint.text}
          hintType={usernameHint.type}
        />

        <Input
          label="Пароль"
          type="password"
          value={password}
          onChange={(e) => setPassword(e.target.value)}
          placeholder="Минимум 6 символов"
          fullWidth
          disabled={registerMutation.isPending}
          autoComplete="new-password"
        />

        <Input
          label="Подтвердите пароль"
          type="password"
          value={confirmPassword}
          onChange={(e) => setConfirmPassword(e.target.value)}
          placeholder="Повторите пароль"
          fullWidth
          disabled={registerMutation.isPending}
          autoComplete="new-password"
          error={localError}
        />

        <Button
          type="submit"
          variant="primary"
          size="large"
          fullWidth
          isLoading={registerMutation.isPending}
          disabled={isSubmitDisabled}
        >
          Зарегистрироваться
        </Button>
      </form>

      <div className="register-form__footer">
        <span className="register-form__footer-text">Уже есть аккаунт?</span>
        <button
          type="button"
          className="register-form__link"
          onClick={() => navigate(ROUTES.LOGIN)}
          disabled={registerMutation.isPending}
        >
          Войти
        </button>
      </div>
    </Card>
  );
};

