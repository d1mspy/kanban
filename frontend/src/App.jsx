import { useState, useEffect } from 'react';
import LoginPage from './LoginPage';
import BoardsPage from './BoardsPage';

function App() {
  const [token, setToken] = useState(localStorage.getItem('token'));

  useEffect(() => {
    if (token) {
      localStorage.setItem('token', token);
    }
  }, [token]);

  const handleLogout = () => {
    localStorage.removeItem('token');
    setToken(null);
  };

  return token
    ? <BoardsPage token={token} onLogout={handleLogout} />
    : <LoginPage onLogin={setToken} />;
}

export default App;
