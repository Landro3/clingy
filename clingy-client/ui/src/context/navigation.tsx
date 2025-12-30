import { createContext, useContext, useState, type PropsWithChildren } from 'react';

export enum Pages {
  Intro,
  Chat,
  Contacts,
  Config,
}

interface NavigationContextType {
  currentPage: Pages;
  navigate: (page: Pages) => void;
}

const NavigationContext = createContext<NavigationContextType | undefined>(undefined);

export function NavigationProvider({ children }: PropsWithChildren) {
  const [currentPage, setCurrentPage] = useState<Pages>(Pages.Chat);

  const navigate = (page: Pages) => {
    setCurrentPage(page);
  };

  const value: NavigationContextType = {
    currentPage,
    navigate,
  };

  return (
    <NavigationContext.Provider value={value}>
      {children}
    </NavigationContext.Provider>
  );
}

export function useNavigation() {
  const context = useContext(NavigationContext);
  if (context === undefined) {
    throw new Error('useNavigation must be used within a NavigationProvider');
  }
  return context;
}

