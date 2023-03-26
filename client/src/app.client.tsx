import { StrictMode } from "react";
import { CacheProvider } from "@emotion/react";
import { hydrateRoot } from "react-dom/client";
import { AppRoute } from "./routes";
import { createEmotionCache } from "./emotionCache";
import { BrowserRouter, Routes } from "react-router-dom";
import { AppProvider } from "./context/app";

const cache = createEmotionCache();

const App = () => {
  // @ts-ignore - too lazy
  const initialData = window.__GO_APP_STATE__ || {};

  return (
    <StrictMode>
      <CacheProvider value={cache}>
        <BrowserRouter>
          <AppProvider initialData={initialData}>
            <Routes>{AppRoute}</Routes>
          </AppProvider>
        </BrowserRouter>
      </CacheProvider>
    </StrictMode>
  );
};

hydrateRoot(document.getElementById("app") as HTMLElement, <App />);
