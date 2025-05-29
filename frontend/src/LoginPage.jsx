import { useState } from 'react';

export default function LoginPage({ onLogin }) {
  const [username, setUsername] = useState('');
  const [password, setPassword] = useState('');
  const [isRegistering, setIsRegistering] = useState(false);
  const [error, setError] = useState('');

  const handleSubmit = async (e) => {
    e.preventDefault();
    const endpoint = isRegistering ? '/api/auth/register' : '/api/auth/login';
    try {
      const res = await fetch(endpoint, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ username, password }),
      });
      if (!res.ok) throw new Error('Invalid credentials');
      const data = await res.json();
      onLogin(data.token);
    } catch (err) {
      setError('Ошибка входа или регистрации');
    }
  };

  return (
    <form onSubmit={handleSubmit}>
      <h2>{isRegistering ? 'Регистрация' : 'Вход'}</h2>
      <input value={username} onChange={e => setUsername(e.target.value)} placeholder="Имя пользователя" />
      <input type="password" value={password} onChange={e => setPassword(e.target.value)} placeholder="Пароль" />
      <button type="submit">{isRegistering ? 'Зарегистрироваться' : 'Войти'}</button>
      <p onClick={() => setIsRegistering(!isRegistering)} style={{ cursor: 'pointer' }}>
        {isRegistering ? 'Уже есть аккаунт?' : 'Нет аккаунта?'}
      </p>
      {error && <p style={{ color: 'red' }}>{error}</p>}
    </form>
  );
}
