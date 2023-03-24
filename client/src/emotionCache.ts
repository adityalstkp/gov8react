import createCache from "@emotion/cache";

export const EMOTION_CACHE_KEY = "custom";

export const createEmotionCache = () => {
  return createCache({ key: EMOTION_CACHE_KEY });
};
