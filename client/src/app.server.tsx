import { StrictMode } from "react";
import { renderToString } from "react-dom/server";
import {
    Routes,
    createRoutesFromElements,
    matchRoutes,
} from "react-router-dom";
import { StaticRouter } from "react-router-dom/server";
import { InitialData } from "./model";
import "./styles/global";
import { CacheProvider } from "@emotion/react";
import { extractCritical } from "@emotion/server";
import { EMOTION_CACHE_KEY, createEmotionCache } from "./emotionCache";
import { AppRoute } from "./routes";
import { AppProvider } from "./context/app";

interface AppProps {
    url: string;
    initialData: Record<string, unknown>;
}

interface EmotionOutput {
    key: string;
    css: string;
    ids: string[];
}

interface RenderOutput {
    html: string;
    emotion: EmotionOutput;
}

const cache = createEmotionCache();

const App = (props: AppProps) => {
    return (
        <StrictMode>
            <StaticRouter location={props.url}>
                <CacheProvider value={cache}>
                    <AppProvider initialData={props.initialData}>
                        <Routes>{AppRoute}</Routes>
                    </AppProvider>
                </CacheProvider>
            </StaticRouter>
        </StrictMode>
    );
};

interface AppArgs {
    url: string;
    initialData: InitialData;
}

export function renderApp(args: AppArgs): RenderOutput {
    try {
        const markup = renderToString(
            <App url={args.url} initialData={args.initialData} />
        );
        const { html, css, ids } = extractCritical(markup);
        return {
            html: html,
            emotion: {
                key: EMOTION_CACHE_KEY,
                css: css,
                ids: ids,
            },
        };
    } catch (err) {
        const error = err as Error;
        return {
            html: error.stack || String(error),
            emotion: {
                key: "",
                css: "",
                ids: [""],
            },
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
