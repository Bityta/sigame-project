/**
 * Auth Feature - Login Form
 * UI компонент формы входа
 */

import { useState } from 'react';
import type { FormEvent } from 'react';
import { useNavigate } from 'react-router-dom';
import { useLogin } from '../../model/mutations';
import { useAuthStore } from '../../model/authStore';
import { Button, Input, Card } from '@/shared/ui';
import { ROUTES, LIMITS } from '@/shared/config';
import './LoginForm.css';

export const LoginForm = () => {
  const navigate = useNavigate();
  const setAuthenticated = useAuthStore((state) => state.setAuthenticated);
  
  const [username, setUsername] = useState('');
  const [password, setPassword] = useState('');
  const [localError, setLocalError] = useState('');

  const loginMutation = useLogin({
    onSuccess: () => {
      setAuthenticated(true);
      navigate(ROUTES.LOBBY);
    },
    onError: () => {
      // Ошибка уже обработана глобально через error store
      // Можем добавить локальное сообщение если нужно
    },
  });

  const handleSubmit = (e: FormEvent) => {
    e.preventDefault();
    setLocalError('');

    // Валидация
    if (username.length < LIMITS.USERNAME.MIN) {
      setLocalError(`Имя пользователя должно быть не менее ${LIMITS.USERNAME.MIN} символов`);
      return;
    }
    if (password.length < LIMITS.PASSWORD.MIN) {
      setLocalError(`Пароль должен быть не менее ${LIMITS.PASSWORD.MIN} символов`);
      return;
    }

    loginMutation.mutate({ username, password });
  };

  return (
    <Card className="login-form">
      <h2 className="login-form__title">Вход</h2>
      
      <form onSubmit={handleSubmit} className="login-form__form">
        <Input
          label="Имя пользователя"
          type="text"
          value={username}
          onChange={(e) => setUsername(e.target.value)}
          placeholder="Введите имя пользователя"
          fullWidth
          disabled={loginMutation.isPending}
          autoComplete="username"
        />

        <Input
          label="Пароль"
          type="password"
          value={password}
          onChange={(e) => setPassword(e.target.value)}
          placeholder="Введите пароль"
          fullWidth
          disabled={loginMutation.isPending}
          autoComplete="current-password"
          error={localError}
        />

        <Button
          type="submit"
          variant="primary"
          size="large"
          fullWidth
          isLoading={loginMutation.isPending}
        >
          Войти
        </Button>
      </form>

      <div className="login-form__footer">
        <span className="login-form__footer-text">Нет аккаунта?</span>
        <button
          type="button"
          className="login-form__link"
          onClick={() => navigate(ROUTES.REGISTER)}
          disabled={loginMutation.isPending}
        >
          Зарегистрироваться
        </button>
      </div>
    </Card>
  );
};

