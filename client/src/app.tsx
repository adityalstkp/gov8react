import React from "react";
import { StrictMode } from "react";
import { renderToString } from "react-dom/server";
import { matchRoutes } from "react-router-dom";
import { renderStylesToString } from "@emotion/server";
import { StaticRouter } from "react-router-dom/server";
import AppRoutes from "./routes";
import { StaticData } from "./model";
import "./styles/global";

interface AppProps {
  url: string;
  staticData: StaticData;
}

interface RenderOutput {
  markup: string;
}

const App = (props: AppProps) => {
  return (
    <StrictMode>
      <StaticRouter location={props.url}>
        <AppRoutes staticData={props.staticData} />
      </StaticRouter>
    </StrictMode>
  );
};

interface AppArgs {
  url: string;
  staticData: StaticData;
}

export async function renderApp(args: AppArgs): Promise<RenderOutput> {
  try {
    const markup = renderStylesToString(
      renderToString(<App url={args.url} staticData={args.staticData} />)
    );
    return { markup };
  } catch (err) {
    const error = err as Error;
    return {
      markup: error.stack || String(error),
    };
  }
}

export function getMatchRoutes(url: string) {
  const match = matchRoutes([{ path: "/" }, { path: "/about" }], url);
  return match;
}
