import { Refine } from "@refinedev/core";
import { RefineKbar, RefineKbarProvider } from "@refinedev/kbar";
import { ThemedLayout } from "@refinedev/antd";
import routerProvider from "@refinedev/react-router-v6";
import dataProvider from "@refinedev/simple-rest";
import { BrowserRouter, Route, Routes, Outlet } from "react-router-dom";
import { ConfigProvider } from "antd";
import "@refinedev/antd/dist/reset.css";

import { LoginPage } from "./pages/Login";

function App() {
  return (
    <BrowserRouter>
      <RefineKbarProvider>
        <ConfigProvider>
            <Refine
              dataProvider={dataProvider("https://api.fake-rest.refine.dev")}
              routerProvider={routerProvider}
              resources={[]}
              options={{
                syncWithLocation: true,
                warnWhenUnsavedChanges: true,
              }}
            >
                <Routes>
                    <Route index element={<LoginPage />} />
                </Routes>
              <RefineKbar />
            </Refine>
        </ConfigProvider>
      </RefineKbarProvider>
    </BrowserRouter>
  );
}

export default App;
