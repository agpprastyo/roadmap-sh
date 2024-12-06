import { Route, Routes } from "react-router-dom";
import { RoutePaths } from "./RoutePaths.jsx";
import { NotFound } from "./NotFound.jsx";
import { Layout } from "./Layout.jsx";
import Home from "../pages/home/Home.jsx";
import Detail from "../pages/home/Detail.jsx";
import {AuthProvider} from "../contexts/AuthContext.jsx";
import Login from "../pages/Login.jsx";
import Admin from "../pages/admin/Admin.jsx";
import PrivateRoute from "../components/PrivateRoute.jsx";


export const Router = () => (
    <AuthProvider>
        <Routes>
            <Route
                path={RoutePaths.HOME}
                element={
                    <Layout>
                        <Home />
                    </Layout>
                }
            />
            <Route
                path={RoutePaths.DETAIL}
                element={
                    <Layout>
                        <Detail />
                    </Layout>
                }
            />
            <Route
                path={RoutePaths.LOGIN}
                element={<Login />}
            />
            <Route
                path={RoutePaths.ADMIN}
                element={
                    <PrivateRoute element={<Layout><Admin /></Layout>} />
                }
            />
            <Route
                path="*"
                element={
                    <Layout>
                        <NotFound />
                    </Layout>
                }
            />
        </Routes>
    </AuthProvider>
);