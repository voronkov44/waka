
import { createRoot } from "react-dom/client";
import App from "./app/App.tsx";
import { initializeTheme } from "./app/hooks/useTheme";
import "./styles/index.css";

initializeTheme();
createRoot(document.getElementById("root")!).render(<App />);
  
