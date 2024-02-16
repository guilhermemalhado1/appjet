import { createContext, useContext, useState } from "react";
import { RenderRoutes } from "../components/structure/RenderNavigation";

const AuthContext = createContext();
export const AuthData = () => useContext(AuthContext);

export const AuthWrapper = () => {
    const [user, setUser] = useState({ name: "", token: "", isAuthenticated: false });
    const isAuthenticated = user.isAuthenticated;
    
    // Retrieve base URL from environment variable
    const baseURL = process.env.REACT_APP_BASE_URL || "http://localhost:9999";
    
    const login = (userName, password) => {
        return new Promise((resolve, reject) => {
            const formData = new URLSearchParams();
            formData.append("username", userName);
            formData.append("password", password);

            let url = `${baseURL}/appjet/login`;

            fetch(url, {
                method: "POST",
                body: formData,
            })
                .then((response) => {
                    if (!response.ok) {
                        throw new Error("Wrong credentials or inexistent account");
                    }
                    return response.json();
                })
                .then((data) => {
                    setUser({ name: userName, token: data.token, isAuthenticated: true });
                    resolve("success");
                })
                .catch((error) => {
                    reject(error.message);
                });
        });
    };

    const logout = async () => {
        return new Promise((resolve, reject) => {
            let url = `${baseURL}/appjet/logout/${user.token}`;

            fetch(url, {
                method: "GET"
            })
                .then((response) => {
                    setUser({ name: "", token: "", isAuthenticated: false });
                    resolve("success");
                })
                .catch((error) => {
                    reject(error.message);
                });
        });
    };

    return (
        <AuthContext.Provider value={{ user, login, logout, isAuthenticated: user.isAuthenticated }}>
            <>
                <RenderRoutes />
            </>
        </AuthContext.Provider>
    );
};
