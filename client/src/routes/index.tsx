import { Route } from "react-router-dom";
import Home from "./Home";
import About from "./About";
import NotFound from "./NotFound";

export const AppRoute = (
  <>
    <Route path="/" element={<Home />} />
    <Route path="/about" element={<About />} />
    <Route path="*" element={<NotFound />} />
  </>
);
