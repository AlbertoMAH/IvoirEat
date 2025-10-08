import { Refine } from "@refinedev/core";
import { RefineKbar, RefineKbarProvider } from "@refinedev/kbar";
import routerProvider from "@refinedev/react-router";
import dataProvider from "@refinedev/simple-rest";
import { BrowserRouter, Route, Routes } from "react-router-dom";

// We remove AntD related imports for this test
// import { ConfigProvider } from "antd";
// import "@refinedev/antd/dist/reset.css";

function App() {
  return (
    <BrowserRouter>
      <RefineKbarProvider>
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
                    <Route index element={<div>Le routeur fonctionne. Le problème vient d'AntD.</div>} />
                </Routes>
              <RefineKbar />
            </Refine>
      </RefineKbarProvider>
    </BrowserRouter>
  );
}

export default App;
