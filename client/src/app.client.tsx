import { StrictMode } from "react";
import { CacheProvider } from "@emotion/react";
import { hydrate } from "react-dom";
import AppRoutes from "./routes";
import { createEmotionCache } from "./emotionCache";
import { BrowserRouter } from "react-router-dom";

const cache = createEmotionCache();

const App = () => {
  // @ts-ignore
  const initialStaticData = window.__INITIAL_STATE__ || {};
  return (
    <StrictMode>
      <CacheProvider value={cache}>
        <BrowserRouter>
          <AppRoutes staticData={initialStaticData} />
        </BrowserRouter>
      </CacheProvider>
    </StrictMode>
  );
};

hydrate(<App />, document.getElementById("app"));
