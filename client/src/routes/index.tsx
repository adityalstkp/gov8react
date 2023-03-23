import React from "react";
import { Route, Routes } from "react-router-dom";
import Home from "./Home";
import About from "./About";
import { StaticData } from "../model";

interface AppRoutesProps {
  staticData: StaticData;
}

const AppRoutes = (props: AppRoutesProps) => {
  return (
    <Routes>
      <Route path="/" element={<Home data={props.staticData || {}} />} />
      <Route path="/about" element={<About />} />
    </Routes>
  );
};

export default AppRoutes;
