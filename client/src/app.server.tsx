import { StrictMode } from "react";
import { renderToString } from "react-dom/server";
import {
  Routes,
  createRoutesFromElements,
  matchRoutes,
} from "react-router-dom";
import { StaticRouter } from "react-router-dom/server";
import { StaticData } from "./model";
import "./styles/global";
import { CacheProvider } from "@emotion/react";
import { extractCritical } from "@emotion/server";
import { EMOTION_CACHE_KEY, createEmotionCache } from "./emotionCache";
import { AppRoute } from "./routes";

interface AppProps {
  url: string;
}

interface RenderOutput {
  markup: string;
  emotionKey: string;
  emotionCss: string;
  emotionIds: string[];
}

const cache = createEmotionCache();

const App = (props: AppProps) => {
  return (
    <StrictMode>
      <StaticRouter location={props.url}>
        <CacheProvider value={cache}>
          <Routes>{AppRoute}</Routes>
        </CacheProvider>
      </StaticRouter>
    </StrictMode>
  );
};

interface AppArgs {
  url: string;
  staticData: StaticData;
}

export function renderApp(args: AppArgs): RenderOutput {
  try {
    const markup = renderToString(<App url={args.url} />);
    const { html, css, ids } = extractCritical(markup);
    return {
      markup: html,
      emotionIds: ids,
      emotionKey: EMOTION_CACHE_KEY,
      emotionCss: css,
    };
  } catch (err) {
    const error = err as Error;
    return {
      markup: error.stack || String(error),
      emotionCss: "",
      emotionIds: [""],
      emotionKey: "",
    };
  }
}

export function getMatchRoutes(url: string) {
  const match = matchRoutes(createRoutesFromElements(AppRoute), url);
  return match;
}

// @ts-ignore - use for go communication
GO_APP.render = renderApp;
// @ts-ignore - use for go communication
GO_APP.getMatchRoutes = getMatchRoutes;
