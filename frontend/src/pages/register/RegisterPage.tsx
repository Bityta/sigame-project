import { RegisterForm } from '@/features/auth';
import { TEXTS } from '@/shared/config';
import './RegisterPage.css';

export const RegisterPage = () => {
  return (
    <div className="register-page">
      <div className="register-page__container">
        <div className="register-page__header">
          <h1 className="register-page__logo">{TEXTS.APP_NAME}</h1>
          <p className="register-page__subtitle">{TEXTS.AUTH.CREATE_ACCOUNT}</p>
        </div>
        <RegisterForm />
      </div>
    </div>
  );
};
