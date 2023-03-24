import { StrictMode } from "react";
import { renderToString } from "react-dom/server";
import { matchRoutes } from "react-router-dom";
import { StaticRouter } from "react-router-dom/server";
import AppRoutes from "./routes";
import { StaticData } from "./model";
import "./styles/global";
import { CacheProvider } from "@emotion/react";
import createCache from "@emotion/cache";
import { extractCritical } from "@emotion/server";

interface AppProps {
  url: string;
  staticData: StaticData;
}

interface RenderOutput {
  markup: string;
  emotionKey: string;
  emotionCss: string;
  emotionIds: string[];
}

const key = "custom";
const cache = createCache({ key: key });

const App = (props: AppProps) => {
  return (
    <StrictMode>
      <StaticRouter location={props.url}>
        <CacheProvider value={cache}>
          <AppRoutes staticData={props.staticData} />
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
    const markup = renderToString(
      <App url={args.url} staticData={args.staticData} />
    );
    const { html, css, ids } = extractCritical(markup);
    return {
      markup: html,
      emotionIds: ids,
      emotionKey: key,
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
  const match = matchRoutes([{ path: "/" }, { path: "/about" }], url);
  return match;
}

// @ts-ignore - use for go communication
GO_APP.render = renderApp;
// @ts-ignore - use for go communication
GO_APP.getMatchRoutes = getMatchRoutes;
