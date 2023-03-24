import { ReactNode, createContext, useContext } from "react";

interface AppData {
  initialData: Record<string, unknown>;
}

const AppContext = createContext<AppData>({ initialData: {} });

interface AppProviderProps {
  children: ReactNode;
  initialData: Record<string, unknown>;
}

export const AppProvider = (props: AppProviderProps) => {
  const data: AppData = {
    initialData: props.initialData,
  };

  return (
    <AppContext.Provider value={data}>{props.children}</AppContext.Provider>
  );
};

export const useApp = () => {
  const context = useContext(AppContext);
  if (context === undefined) {
    throw new Error("useCount must be used within a CountProvider");
  }
  return context;
};
