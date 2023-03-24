import { useEffect, useState } from "react";
import { useApp } from "../context/app";

interface IntroOutput {
  data: string;
  loading: boolean;
}

interface IntroData {
  message: string;
}

const matchRoute = "/";

export const useIntro = (): IntroOutput => {
  const { initialData } = useApp();
  const introInitialData = (initialData[matchRoute] as IntroData) || {};
  const greetData = introInitialData.message;

  const [data, setData] = useState(greetData);
  const [loading, setLoading] = useState(false);

  const handleGetIntroData = async () => {
    setLoading(true);
    try {
      const res = await fetch("/api/v1/intro");
      const payload = await res.json();
      setData(payload.message);
    } catch {
      // no-op
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    if (!greetData) {
      handleGetIntroData();
    }
  }, []);

  return { data, loading };
};
