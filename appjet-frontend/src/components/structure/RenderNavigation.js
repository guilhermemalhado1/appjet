import { Route, Routes } from "react-router-dom";
import { AuthData } from "../../auth/AuthWrapper";
import Login from "../login/Login"
import Home from "../home/Home";


export const RenderRoutes = () => {

        const { user } = AuthData();
        
        return (
             <Routes>
                <Route path="/" element={<Login />} />
                <Route path="/login" element={<Login />} />
                <Route path="/home" element={<Home />} />
             </Routes>
        )
   }