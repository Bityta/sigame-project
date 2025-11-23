import { LoginForm } from '@/features/auth';
import { TEXTS } from '@/shared/config';
import './LoginPage.css';

export const LoginPage = () => {
  return (
    <div className="login-page">
      <div className="login-page__container">
        <div className="login-page__header">
          <h1 className="login-page__logo">{TEXTS.APP_NAME}</h1>
          <p className="login-page__subtitle">{TEXTS.APP_SUBTITLE}</p>
        </div>
        <LoginForm />
      </div>
    </div>
  );
};
