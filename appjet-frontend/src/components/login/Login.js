import React, { useState } from 'react';
import './Login.css'; // Import the CSS file
import { AuthData } from "../../auth/AuthWrapper"
import { useNavigate } from "react-router-dom";

function Login() {
  const [username, setUsername] = useState('');
  const [password, setPassword] = useState('');
  const [loginStatus, setLoginStatus] = useState(null); // State for login status
  const [errorMessage, setErrorMessage] = useState(null);

  const { login } = AuthData();
  const navigate = useNavigate();

  const handleUsernameChange = (event) => {
    setUsername(event.target.value);
  };

  const handlePasswordChange = (event) => {
    setPassword(event.target.value);
  };

  const handleSubmit = async (event) => {
    event.preventDefault();
    console.log('Username:', username);
    console.log('Password:', password);
  
    try {
      await login(username, password);
      setLoginStatus('success'); // Set login status to success
      setErrorMessage(null); // Reset error message
      navigate("/home");
    } catch (error) {
      setLoginStatus('error'); // Set login status to error
      setErrorMessage(error);
    }

    console.log(errorMessage)
  };
  

  return (
    <div className="container">
      <form className="login-form" onSubmit={handleSubmit}>
        <h2>APPJET</h2>
        <div className="input-group">
          <label htmlFor="username">Username</label>
          <input 
            type="text" 
            id="username" 
            name="username" 
            value={username} // Connect value to username state
            onChange={handleUsernameChange} // Provide onChange handler
            required
          />
        </div>
        <div className="input-group">
          <label htmlFor="password">Password</label>
          <input 
            type="password" 
            id="password" 
            name="password" 
            value={password}
            onChange={handlePasswordChange}
            required
          />
        </div>
        <button type="submit">Login</button>
        {loginStatus === 'error' && <p className="error-message">{errorMessage}</p>} {/* Display error message if login status is error */}
      </form>
    </div>
  );
}

export default Login;
