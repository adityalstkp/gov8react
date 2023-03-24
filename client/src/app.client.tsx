import { StrictMode } from "react";
import { CacheProvider } from "@emotion/react";
import { hydrate } from "react-dom";
import { AppRoute } from "./routes";
import { createEmotionCache } from "./emotionCache";
import { BrowserRouter, Routes } from "react-router-dom";

const cache = createEmotionCache();

const App = () => {
  // @ts-ignore
  // const initialStaticData = window.__INITIAL_STATE__ || {};
  return (
    <StrictMode>
      <CacheProvider value={cache}>
        <BrowserRouter>
          <Routes>{AppRoute}</Routes>
        </BrowserRouter>
      </CacheProvider>
    </StrictMode>
  );
};

hydrate(<App />, document.getElementById("app"));
